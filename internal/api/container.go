package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/kraken-hpc/go-fork"
	"github.com/kraken-hpc/imageapi/models"
	"github.com/kraken-hpc/uinit"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type container struct {
	ctn    *models.Container
	log    *log.Logger
	mnt    string
	cancel context.CancelFunc
}

type ContainersType struct {
	next  models.ID
	ctns  map[models.ID]*container
	names map[models.Name]models.ID
	mutex *sync.Mutex
	log   *logrus.Entry
}

func (c *ContainersType) Init() {
	c.next = 1
	c.ctns = make(map[models.ID]*container)
	c.names = make(map[models.Name]models.ID)
	c.mutex = &sync.Mutex{}
	c.log = Log.WithField("subsys", "container")
	c.log.Trace("initialized")
}

func (c *ContainersType) List() (r []*models.Container) {
	l := c.log.WithField("operation", "list")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, ctn := range c.ctns {
		r = append(r, ctn.ctn)
	}
	l.WithField("entries", len(r)).Trace("listing entries")
	return
}

func (c *ContainersType) Get(id models.ID) (*models.Container, error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "get",
		"id":        id,
	})
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if r, ok := c.ctns[id]; ok {
		l.Trace("found")
		return r.ctn, nil
	}
	l.Trace("not found")
	return nil, fmt.Errorf("no container by id %d", id)
}

func (c *ContainersType) Create(ctn *models.Container) (r *models.Container, err error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "create",
		"name":      ctn.Name,
	})
	// This creates a container in our list, and activates its initial state
	// find the mount
	n := &container{
		ctn: ctn,
	}
	if ctn.Mount.MountID == 0 {
		if ctn.Mount, err = Mount(ctn.Mount); err != nil {
			l.WithError(err).Error("failed to mount")
			return nil, fmt.Errorf("mount failed: %v", err)
		}
		if n.mnt, err = MountGetMountpoint(ctn.Mount); err != nil {
			l.WithError(err).Error("failed to get mountpoint")
			MountRefAdd(ctn.Mount, -1)
			return nil, fmt.Errorf("could not get mountpoint: %v", err)
		}
	} else {
		if n.mnt, err = MountGetMountpoint(ctn.Mount); err != nil {
			return nil, fmt.Errorf("could not get mountpoint: %v", err)
		}
		MountRefAdd(ctn.Mount, 1)
	}

	defer func() {
		if err != nil {
			MountRefAdd(ctn.Mount, -1)
		}
	}()
	l = l.WithField("mount_point", n.mnt)

	// ok, we've got a valid mountpoint
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// fail early on non-unique name
	if _, ok := c.names[ctn.Name]; ctn.Name != "" && ok {
		err := fmt.Errorf("container with name %s already exists", ctn.Name)
		l.Debug("non-unique name")
		return nil, err
	}
	ctn.ID = c.next

	// set up logger
	if err = os.MkdirAll(logDir, 0700); err != nil {
		l.WithError(err).Error("could not make log directory")
		return nil, fmt.Errorf("could not create log directory: %v", err)
	}
	ctn.Logfile = path.Join(logDir, fmt.Sprintf("%d-%d.log", ctn.ID, time.Now().Unix()))
	f, err := os.Create(ctn.Logfile)
	if err != nil {
		l.WithError(err).Error("failed to creat elog file")
		return nil, fmt.Errorf("could not create log file: %v", err)
	}
	n.log = log.New(f, fmt.Sprintf("container(%d): ", ctn.ID), log.Ldate|log.Ltime|log.Lmsgprefix)
	n.log.Printf("container created")

	// handle initial state
	switch ctn.State {
	case models.ContainerStateRunning:
		if err := c.run(n); err != nil {
			l.WithError(err).Error("failed to start container")
			return nil, fmt.Errorf("failed to start container: %v", err)
		}
		ctn.State = models.ContainerStateRunning
	case models.ContainerStateStopping,
		models.ContainerStateExited,
		models.ContainerStateDead:
		err = fmt.Errorf("requested invalid initial container state: %s.  valid initial states: [ %s, %s ]", ctn.State, models.ContainerStateCreated, models.ContainerStateRunning)
		l.WithError(err).Debug("invalid initial state")
		return nil, err
	case models.ContainerStateCreated:
		fallthrough
	default: // wasn't specified
		ctn.State = models.ContainerStateCreated
	}

	// container is ready to be entered
	c.ctns[ctn.ID] = n
	c.next++
	c.names[ctn.Name] = ctn.ID

	l.Info("successfully created")
	return ctn, nil
}

func (c *ContainersType) SetState(id models.ID, state models.ContainerState) (err error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "setstate",
		"id":        id,
	})
	var ctn *container
	var ok bool
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if ctn, ok = c.ctns[id]; !ok {
		l.Error("not found")
		return fmt.Errorf("invalid container id: %d", id)
	}
	// handle state request
	switch state {
	case models.ContainerStateRunning:
		if ctn.ctn.State == state {
			return
		}
		ctn.ctn.State = models.ContainerStateRunning
		l.Info("starting container")
		if err = c.run(ctn); err != nil {
			l.WithError(err).Error("failed to start")
			ctn.ctn.State = models.ContainerStateDead
			return
		}
	case models.ContainerStateExited:
		if ctn.ctn.State == state {
			return
		}
		l.Info("stopping container")
		c.stop(ctn)
	default: // something not valid
		err := fmt.Errorf("can't set state to: %s.  valid states to request: [ %s, %s ]", state,
			models.ContainerStateRunning,
			models.ContainerStateExited)
		l.WithError(err).Error("failed")
		return err
	}
	return
}

func (c *ContainersType) Delete(id models.ID) (ret *models.Container, err error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "setstate",
		"id":        id,
	})
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var ctn *container
	var ok bool
	if ctn, ok = c.ctns[id]; !ok {
		l.Error("container not found")
		return nil, fmt.Errorf("invalid container id: %d", id)
	}
	switch ctn.ctn.State {
	//case models.ContainerStatePaused:
	//case models.ContainerStateRestarting:
	case models.ContainerStateRunning:
		l.Trace("attempt to delete running container")
		return nil, fmt.Errorf("cannot delete a running container")
	case models.ContainerStateStopping:
		l.Trace("attempt to delete stopping container")
		return nil, fmt.Errorf("cannot delete a container that hasn't fully stopped")
	}
	ctn.log.Printf("container deleted")
	ctn.log.Writer().(io.WriteCloser).Close()
	delete(c.ctns, id)
	if ctn.ctn.Name != "" {
		delete(c.names, ctn.ctn.Name)
	}
	MountRefAdd(ctn.ctn.Mount, -1)
	l.Info("container deleted")
	// garbage collection should take care of our mount if it's now unused
	return
}

// NameGetID will return the ID for a given name
// This is used to implement the `byname` calls
// If the name is not found,  it will return -1
func (c *ContainersType) NameGetID(name models.Name) models.ID {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if id, ok := c.names[name]; ok {
		return id
	}
	return -1
}

func (c *ContainersType) stop(ctn *container) error {
	// this will need to be updated if we implement pause and/or restart
	if ctn.ctn.State != models.ContainerStateRunning {
		// it's not an error to stop an already stopped container
		return nil
	}
	// we trust the watcher to take care of everything for us
	ctn.cancel()
	ctn.ctn.State = models.ContainerStateStopping
	return nil
}

// this is the workhorse
func (c *ContainersType) run(ctn *container) (err error) {
	log := ctn.log
	log.Print("starting container")

	// 0. For a container, we want to be able to launch bare directory containers
	//    We implement this by bind mounting newRoot on itself.
	// NOTE: this should always be true, but it's a good check anyway
	if !isMountpoint(ctn.mnt) {
		if err = bindMountSelf(ctn.mnt); err != nil {
			return fmt.Errorf("clone: could not self-bind mount bare directory: %v", err)
		}
	}

	// 1. Is our image valid?
	log.Print("validating image")
	if err = validateImage(ctn.mnt); err != nil {
		return fmt.Errorf("clone: image validation failed: %v", err)
	}

	// 2 parse command into args
	args := uinit.SplitCommandLine(*ctn.ctn.Command)
	if len(args) < 1 {
		return fmt.Errorf("command appears to be invalid: %s", *ctn.ctn.Command)
	}

	// 3. Is our init valid?
	log.Print("validating init")
	if err = validateInit(ctn.mnt, args[0]); err != nil {
		return fmt.Errorf("clone: init validationfailed: %v", err)
	}

	// 3. Launch new process
	f := fork.NewFork("containerInit", containerInit)
	f.Stdout = log.Writer().(*os.File)
	f.Stderr = log.Writer().(*os.File)
	f.Stdin = nil
	f.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS,
	}
	if err := f.Fork(ctn.mnt, ctn.ctn.Systemd, args); err != nil {
		return fmt.Errorf("clone: failed to start pid_init: %v", err)
	}

	// 4. Launch the process watcher
	ctx := context.Background()
	ctx, ctn.cancel = context.WithCancel(ctx)
	go c.watcher(ctx, ctn, f)
	return
}

type mountType struct {
	dev    string
	path   string
	fstype string
	flags  uintptr
}

// see libcontainer SPEC.md
var specialMounts = []mountType{
	{"proc", "/proc", "proc", unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV},
	{"tmpfs", "/dev", "tmpfs", unix.MS_NOEXEC | unix.MS_STRICTATIME},
	{"tmpfs", "/dev/shm", "tmpfs", unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV},
	{"mqueue", "/dev/mqueue", "mqueue", unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV},
	{"devpts", "/dev/pts", "devpts", unix.MS_NOEXEC | unix.MS_NOSUID},
	{"sysfs", "/sys", "sysfs", unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV | unix.MS_RDONLY},
}

var systemdMounts = []mountType{
	{"tmpfs", "/run", "tmpfs", unix.MS_NOSUID | unix.MS_NODEV},
	{"tmpfs", "/tmp", "tmpfs", unix.MS_NOSUID | unix.MS_NODEV},
	{"tmpfs", "/sys/fs/cgroup", "tmpfs", unix.MS_NOSUID | unix.MS_NODEV},
	{"tmpfs", "/var/lib/journal", "tmpfs", unix.MS_NOSUID | unix.MS_NODEV},
}

type deviceFileType struct {
	path string
	mode uint32
	dev  uint64
}

var specialDevices = []deviceFileType{
	{"/dev/null", syscall.S_IFCHR | uint32(os.FileMode(0666)), unix.Mkdev(1, 3)},
	{"/dev/zero", syscall.S_IFCHR | uint32(os.FileMode(0666)), unix.Mkdev(1, 5)},
	{"/dev/full", syscall.S_IFCHR | uint32(os.FileMode(0666)), unix.Mkdev(1, 7)},
	{"/dev/tty", syscall.S_IFCHR | uint32(os.FileMode(0666)), unix.Mkdev(5, 0)},
	{"/dev/random", syscall.S_IFCHR | uint32(os.FileMode(0666)), unix.Mkdev(1, 8)},
	{"/dev/urandom", syscall.S_IFCHR | uint32(os.FileMode(0666)), unix.Mkdev(1, 9)},
}

type symlinkType struct {
	from string
	to   string
}

var specialLinks = []symlinkType{
	{"/dev/pts/ptmx", "/dev/ptmx"},
	{"/proc/self/fd", "/dev/fd"},
	{"/proc/self/fd/0", "/dev/stdin"},
	{"/proc/self/fd/1", "/dev/stdout"},
	{"/proc/self/fd/2", "/dev/stderr"},
}

// this is run as a separate process
func containerInit(mountpoint string, systemd bool, args []string) {
	// 0. setup logging
	l := log.New(os.Stdout, "init: ", log.Ldate|log.Ltime|log.Lmsgprefix)

	// 1. Make sure mounts are marked as private, necessary for moving mounts
	l.Print("making all mounts private")
	if err := makeMountsPrivate(); err != nil {
		l.Fatalf("failed to make mounts private: %v", err)
	}

	// 2. Do the root moving dance
	l.Print("preparing image")
	if err := moveRoot(mountpoint); err != nil {
		l.Fatalf("could not prepare image: %v", err)
	}

	// 3. Setup special mounts
	if systemd {
		specialMounts = append(specialMounts, systemdMounts...)
	}
	// we want our perms to be absolute (i.e. no umask) for the next steps
	oldUmask := unix.Umask(int(os.FileMode(0000)))
	for _, m := range specialMounts {
		if err := containerMount(l, m.dev, m.path, m.fstype, m.flags); err != nil {
			l.Fatalf("mount failed for %s: %v", m.path, err)
		}
	}

	// 4. Setup special dev files
	for _, d := range specialDevices {
		l.Printf("making device file %s", d.path)
		if err := unix.Mknod(d.path, d.mode, int(d.dev)); err != nil {
			l.Fatalf("failed to create device %s: %v", d.path, err)
		}
	}
	unix.Umask(oldUmask)

	// 5. Setup special symlinks
	for _, s := range specialLinks {
		l.Printf("creating symlink %s -> %s", s.from, s.to)
		if err := os.Symlink(s.from, s.to); err != nil {
			l.Fatalf("failed to create symlink %s: %v", s.to, err)
		}
	}

	// 6. execute init
	l.Print("executing init")
	if err := unix.Exec(args[0], args, []string{}); err != nil {
		l.Fatalf("containerInit: exec failed: %v", err)
	}
}

func (c *ContainersType) watcher(ctx context.Context, ctn *container, f *fork.Function) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "watcher",
		"id":        ctn.ctn.ID,
		"name":      ctn.ctn.Name,
	})
	end := make(chan error)
	go func() {
		e := f.Wait()
		l.Trace("exited")
		end <- e
	}()
	state := models.ContainerStateExited
	select {
	case e := <-end:
		if e != nil {
			l.WithError(e).Debug("process ended in error state")
			ctn.log.Printf("process ended in error state: %v", e)
			state = models.ContainerStateDead
		}
		ctn.log.Printf("process ended, exit code (0)")
	case <-ctx.Done():
		// signal the process to stop
		// TODO: be smarter about the signals we send
		l.Debug("sending kill signal")
		if ctn.ctn.Systemd {
			// SIGRTMIN+3
			f.Process.Signal(syscall.Signal(37))
		} else {
			f.Process.Kill()
		}
		ctn.log.Printf("process killed")
	}
	// don't report that we're done until fs is synced
	fd, err := unix.Open(ctn.mnt, unix.O_RDWR|unix.O_DIRECTORY, 0)
	if err != nil {
		log.Printf("failed to open mount for sync: %v", err)
	} else {
		unix.Syncfs(fd)
		unix.Close(fd)
	}
	// process is over, set the state
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ctn.ctn.State = state
	ctn.log.Printf("container state: %s", state)
}

func containerMount(l *log.Logger, dev, path, fstype string, flags uintptr) error {
	l.Printf("mounting %s", path)
	// make sure path exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("could not create mount path %s: %v", path, err)
	}
	if isMountpoint(path) {
		if err := unix.Unmount(path, unix.MNT_DETACH); err != nil {
			l.Printf("warn: could not unmount %s, will overlay instead: %v", path, err)
		}
	}
	return unix.Mount(dev, path, fstype, flags, "")
}

// utilities for run

func getDev(fd int) (uint64, error) {
	stat := &unix.Stat_t{}
	if err := unix.Fstat(fd, stat); err != nil {
		return 0, err
	}
	return stat.Dev, nil
}

func isMountpointAt(parentDev uint64, fd int) bool {
	dev, err := getDev(fd)
	if err != nil {
		// note this behavior is slightly arbitrary
		return false
	}
	if dev != parentDev {
		return true
	}
	return false
}

func isMountpoint(path string) bool {
	var fd, pfd int
	var parentDev uint64
	var err error
	parent := filepath.Dir(path)
	if pfd, err = unix.Open(parent, unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		// note this behavior is slightly arbitrary
		return false
	}
	defer unix.Close(pfd)
	if parentDev, err = getDev(pfd); err != nil {
		return false
	}

	if fd, err = unix.Open(path, unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		return false
	}
	defer unix.Close(fd)
	return isMountpointAt(parentDev, fd)
}

func validateInit(newRoot, init string) (err error) {
	var stat os.FileInfo
	var realInit string
	var exit func() error

	if exit, err = chroot(newRoot); err != nil {
		return fmt.Errorf("could not chroot into %s: %v", newRoot, err)
	}

	if realInit, err = filepath.EvalSymlinks(init); err != nil {
		return fmt.Errorf("init file could not be found: %v", err)
	}

	if err := exit(); err != nil {
		return fmt.Errorf("could not exit chroot: %v", err)
	}

	if stat, err = os.Stat(filepath.Join(newRoot, realInit)); err != nil {
		return fmt.Errorf("init file could not be opened: %v", err)
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("init does not reference a regular file: %v", err)
	}
	if stat.Mode()&0111 == 0 {
		return fmt.Errorf("init file is not executable: %v", err)
	}
	return
}

func moveMount(newRoot, mount string) (err error) {
	joined := filepath.Join(newRoot, mount)
	if !isMountpoint(mount) {
		return fmt.Errorf("original mountpoint does not exist")
	}
	if isMountpoint(joined) {
		// we *do* want to unmount at least
		unix.Unmount(mount, unix.MNT_DETACH)
		return fmt.Errorf("new mountpoint already mounted, old mount detached")
	}
	if err = unix.Mount(mount, joined, "", unix.MS_MOVE, ""); err != nil {
		// we still force an unmount
		unix.Unmount(mount, unix.MNT_FORCE)
		return fmt.Errorf("mount move failed, old mount force unmounted: %v", err)
	}
	return
}

func chroot(path string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := unix.Chroot(path); err != nil {
		root.Close()
		return nil, err
	}

	if err := os.Chdir("/"); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return unix.Chroot(".")
	}, nil
}

// we store these as a global var so that we could potentially have a way to update at runtime
//var specialMounts = []string{"/dev", "/proc", "/sys", "/run"}
var moveMounts = []string{}

// this is the workhorse for all schemes
// it preforms the root-moving dance
func moveRoot(newRoot string) (err error) {
	// 1. move special mounts
	for _, mount := range moveMounts {
		if err := moveMount(newRoot, mount); err != nil {
			// this isn't fatal, but we should mention it
			log.Printf("warn: couldn't move mount %s: %v", mount, err)
		}
	}
	// 2. chdir to new root
	if err = os.Chdir(newRoot); err != nil {
		return fmt.Errorf("failed to chdir to new root: %v", err)
	}
	// 3. Move newRoot -> /
	if err = unix.Mount(newRoot, "/", "", unix.MS_MOVE, ""); err != nil {
		return fmt.Errorf("failed to move new root to /: %v", err)
	}
	// 4. chroot "."
	if _, err = chroot("."); err != nil {
		return fmt.Errorf("failed to change root: %v", err)
	}

	// the dance is done
	return
}

func makeMountsPrivate() error {
	return unix.Mount("", "/", "", unix.MS_REC|unix.MS_PRIVATE, "")
}

func bindMountSelf(path string) (err error) {
	// if we're already a mount point, just return
	if isMountpoint(path) {
		return
	}
	// we blindly try this without verifying that it's a directory
	if err = unix.Mount(path, path, "", unix.MS_BIND, ""); err != nil {
		return fmt.Errorf("failed to create root bind mount: %v", err)
	}
	return
}

func validateImage(newRoot string) (err error) {
	var fd int
	// Does the directory exist? Or, is it a directory?
	if fd, err = unix.Open(newRoot, unix.O_DIRECTORY, unix.O_RDONLY); err != nil {
		return fmt.Errorf("new root is not a valid directory")
	}
	unix.Close(fd)

	// Is it a mount point?
	if !isMountpoint(newRoot) {
		return fmt.Errorf("new root is not a mountpoint")
	}
	return
}

func init() {
	fork.RegisterFunc("containerInit", containerInit)
	fork.Init()
}
