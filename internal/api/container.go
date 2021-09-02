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

// Unlike other subsystems, we wrap models.Container with extra data
type Container struct {
	Container *models.Container
	log       *log.Logger
	cancel    context.CancelFunc
}

// Make sure Container is an EndpointObject
var _ EndpointObject = (*Container)(nil)

func (c *Container) GetID() models.ID                       { return c.Container.ID }
func (c *Container) SetID(id models.ID)                     { c.Container.ID = id }
func (c *Container) GetRefs() int64                         { return c.Container.Refs }
func (c *Container) RefAdd(i int64)                         { c.Container.Refs += i }
func (c *Container) EndpointObjectType() EndpointObjectType { return EndpointObjectContainer }

type Containers struct {
	names map[models.Name]models.ID
	mutex *sync.Mutex
	log   *logrus.Entry
}

func (c *Containers) Init(log *logrus.Entry) {
	c.names = make(map[models.Name]models.ID)
	c.mutex = &sync.Mutex{}
	c.log = log
	c.log.Info("container subsystem initialized")
}

func (c *Containers) List() (ret []*Container) {
	ret = []*Container{}
	for _, o := range API.Store.ListType(EndpointObjectContainer) {
		ret = append(ret, o.(*Container))
	}
	return
}

func (c *Containers) Get(id models.ID) *Container {
	if eo := API.Store.Get(id); eo != nil {
		if ret, ok := eo.(*Container); ok {
			return ret
		}
	}
	return nil
}

func (c *Containers) Create(n *Container) (ret *Container, err error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "create",
	})
	// This creates a container in our list, and activates its initial state
	// find the mount
	if n.Container == nil {
		l.Debug("container create called with no container definition")
		return nil, ErrInvalDat
	}
	l = l.WithField("name", n.Container.Name)
	ctn := n.Container
	m, err := API.Mounts.GetOrMount((*Mount)(ctn.Mount))
	if err != nil {
		return nil, err
	}
	ctn.Mount = (*models.Mount)(m)
	defer func() {
		if err != nil {
			API.Store.RefAdd(ctn.Mount.ID, -1)
		}
	}()
	l = l.WithField("mount_point", ctn.Mount.Mountpoint)

	// ok, we've got a valid mountpoint
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// fail early on non-unique name
	if _, ok := c.names[ctn.Name]; ctn.Name != "" && ok {
		l.Debugf("container with name %s already exists", ctn.Name)
		return nil, ErrBusy
	}

	// we enter the container into the store so we get an ID
	n = API.Store.Register(n).(*Container)
	ctn = n.Container
	defer func() {
		if err != nil {
			API.Store.Unregister(n)
		}
	}()

	// set up logger
	if err = os.MkdirAll(API.LogDir, 0700); err != nil {
		l.WithError(err).Error("could not make log directory")
		return nil, ErrSrv
	}
	ctn.Logfile = path.Join(API.LogDir, fmt.Sprintf("%d-%d.log", ctn.ID, time.Now().Unix()))
	f, err := os.Create(ctn.Logfile)
	if err != nil {
		l.WithError(err).Error("failed to creat elog file")
		return nil, ErrSrv
	}
	n.log = log.New(f, fmt.Sprintf("container(%d): ", ctn.ID), log.Ldate|log.Ltime|log.Lmsgprefix)
	n.log.Printf("container created")

	n.log.Printf("running script hook: create")
	var hook *models.ContainerScriptHook
	if ctn.Hooks != nil {
		hook = ctn.Hooks.Create
	}
	if h, err := NewHook(hook, "", ctn.Mount.Mountpoint, n.log, map[string]string{
		"logfile":    ctn.Logfile,
		"command":    *ctn.Command,
		"mountpoint": ctn.Mount.Mountpoint,
		"mountkind":  ctn.Mount.Kind,
		"name":       string(ctn.Name),
		"id":         fmt.Sprintf("%d", ctn.ID),
		"systemd":    fmt.Sprintf("%t", ctn.Systemd),
	}); err != nil {
		l.WithError(err).Debug("fatal error building create scripts")
		return nil, ErrFail
	} else {
		if err = h.Run(); err != nil {
			l.WithError(err).Debug("fatal error running create scripts")
			return nil, ErrFail
		}
	}

	// handle initial state
	switch ctn.State {
	case models.ContainerStateRunning:
		if err := c.run(n); err != nil {
			l.WithError(err).Error("failed to start container")
			ctn.State = models.ContainerStateDead
		} else {
			l.Info("container started")
			ctn.State = models.ContainerStateRunning
		}
	case models.ContainerStateStopping,
		models.ContainerStateExited,
		models.ContainerStateDead:
		l.Debug("invalid initial state")
		return nil, ErrInvalDat
	case models.ContainerStateCreated:
		fallthrough
	default: // wasn't specified
		ctn.State = models.ContainerStateCreated
	}

	// update our object entry
	API.Store.Update(n)
	c.names[ctn.Name] = ctn.ID

	l.Info("successfully created")
	return n, nil
}

func (c *Containers) SetState(id models.ID, state models.ContainerState) (ret *Container, err error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "setstate",
		"id":        id,
	})
	ctn := c.Get(id)
	if ctn == nil {
		l.Debug("requested setstate on non-existant container")
		return nil, ErrNotFound
	}
	defer func() {
		API.Store.Update(ctn)
		API.Store.RefAdd(id, -1) // clear our hold
	}()

	// handle state request
	switch state {
	case models.ContainerStateRunning:
		if ctn.Container.State == state {
			return ctn, nil
		}
		ctn.Container.State = models.ContainerStateRunning
		l.Info("starting container")
		if err = c.run(ctn); err != nil {
			l.WithError(err).Error("failed to start")
			ctn.Container.State = models.ContainerStateDead
			return nil, ErrFail
		}
	case models.ContainerStateExited:
		if ctn.Container.State == state {
			return ctn, nil
		}
		l.Info("stopping container")
		c.stop(ctn)
	default: // something not valid
		err := fmt.Errorf("can't set state to: %s.  valid states to request: [ %s, %s ]", state,
			models.ContainerStateRunning,
			models.ContainerStateExited)
		l.WithError(err).Error("failed")
		return nil, ErrInvalDat
	}
	return ctn, nil
}

func (c *Containers) Delete(id models.ID) (ret *Container, err error) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "setstate",
		"id":        id,
	})
	ctn := c.Get(id)
	if ctn == nil {
		l.Debug("delete called on non-existent container")
		return nil, ErrNotFound
	}

	switch ctn.Container.State {
	//case models.ContainerStatePaused:
	//case models.ContainerStateRestarting:
	case models.ContainerStateRunning:
		l.Trace("attempt to delete running container")
		return nil, ErrBusy
	case models.ContainerStateStopping:
		l.Trace("attempt to delete stopping container")
		return nil, ErrBusy
	}
	// run delete hook
	ctn.log.Printf("running script hook: exit")
	var hook *models.ContainerScriptHook
	if ctn.Container.Hooks != nil {
		hook = ctn.Container.Hooks.Exit
	}
	if h, err := NewHook(hook, "", ctn.Container.Mount.Mountpoint, ctn.log, map[string]string{
		"logfile":    ctn.Container.Logfile,
		"command":    *ctn.Container.Command,
		"mountpoint": ctn.Container.Mount.Mountpoint,
		"mountkind":  ctn.Container.Mount.Kind,
		"name":       string(ctn.Container.Name),
		"id":         fmt.Sprintf("%d", ctn.Container.ID),
		"systemd":    fmt.Sprintf("%t", ctn.Container.Systemd),
		"state":      string(ctn.Container.State),
	}); err != nil {
		l.WithError(err).Debug("fatal error building exit scripts")
		// but we don't actually do anything here
	} else {
		if err = h.Run(); err != nil {
			l.WithError(err).Debug("fatal error running exit scripts")
			// but we don't actually do anything here
		}
	}

	ctn.log.Printf("container deleted")
	ctn.log.Writer().(io.WriteCloser).Close()
	if ctn.Container.Name != "" {
		c.mutex.Lock()
		delete(c.names, ctn.Container.Name)
		c.mutex.Unlock()
	}
	API.Store.Unregister(ctn)
	API.Store.RefAdd(ctn.Container.Mount.ID, -1)
	l.Info("container deleted")
	// garbage collection should take care of our mount if it's now unused
	return ctn, nil
}

// NameGetID will return the ID for a given name
// This is used to implement the `byname` calls
// If the name is not found,  it will return -1
func (c *Containers) NameGetID(name models.Name) models.ID {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if id, ok := c.names[name]; ok {
		return id
	}
	return -1
}

func (c *Containers) stop(ctn *Container) error {
	// this will need to be updated if we implement pause and/or restart
	if ctn.Container.State != models.ContainerStateRunning {
		// it's not an error to stop an already stopped container
		return nil
	}
	// we trust the watcher to take care of everything for us
	ctn.cancel()
	ctn.Container.State = models.ContainerStateStopping
	API.Store.Update(ctn)
	return nil
}

// this is the workhorse
func (c *Containers) run(ctn *Container) (err error) {
	log := ctn.log
	log.Print("starting container")

	// 0. For a container, we want to be able to launch bare directory containers
	//    We implement this by bind mounting newRoot on itself.
	// NOTE: this should always be true, but it's a good check anyway
	if !isMountpoint(ctn.Container.Mount.Mountpoint) {
		if err = bindMountSelf(ctn.Container.Mount.Mountpoint); err != nil {
			return fmt.Errorf("clone: could not self-bind mount bare directory: %v", err)
		}
	}

	// 1. Is our image valid?
	log.Print("validating image")
	if err = validateImage(ctn.Container.Mount.Mountpoint); err != nil {
		return fmt.Errorf("clone: image validation failed: %v", err)
	}

	// 2 parse command into args
	args := uinit.SplitCommandLine(*ctn.Container.Command)
	if len(args) < 1 {
		return fmt.Errorf("clone: command appears to be invalid: %s", *ctn.Container.Command)
	}

	// 3. Is our init valid?
	log.Print("validating init")
	if err = validateInit(ctn.Container.Mount.Mountpoint, args[0]); err != nil {
		return fmt.Errorf("clone: init validationfailed: %v", err)
	}

	// 4. Pre-build start hook
	log.Print("building script hook: Init")
	var hook *models.ContainerScriptHook
	if ctn.Container.Hooks != nil {
		hook = ctn.Container.Hooks.Init
	}
	h, err := NewHook(hook, "", ctn.Container.Mount.Mountpoint, ctn.log, map[string]string{
		"logfile":    ctn.Container.Logfile,
		"command":    *ctn.Container.Command,
		"mountpoint": ctn.Container.Mount.Mountpoint,
		"mountkind":  ctn.Container.Mount.Kind,
		"name":       string(ctn.Container.Name),
		"id":         fmt.Sprintf("%d", ctn.Container.ID),
		"systemd":    fmt.Sprintf("%t", ctn.Container.Systemd),
	})
	if err != nil {
		return fmt.Errorf("clone: fatal error loading init script hook: %v", err.Error())
	}

	// 5. Launch new process
	f := fork.NewFork("containerInit", containerInit)
	f.Stdout = log.Writer().(*os.File)
	f.Stderr = log.Writer().(*os.File)
	f.Stdin = nil
	f.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS,
	}
	if err := f.Fork(ctn.Container.Mount.Mountpoint, ctn.Container.Systemd, h, args); err != nil {
		return fmt.Errorf("clone: failed to start pid_init: %v", err)
	}

	// 4. Launch the process watcher
	ctx := context.Background()
	ctx, ctn.cancel = context.WithCancel(ctx)
	API.Store.Update(ctn)
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
func containerInit(mountpoint string, systemd bool, hook *Hook, args []string) {
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

	// 6. Run init script hooks
	if hook != nil {
		l.Printf("executing init script hooks")
		if err := hook.Run(); err != nil {
			l.Fatalf("fatal init script hook error: %v", err)
		}
	}

	// 7. execute init
	l.Print("executing init")
	if err := unix.Exec(args[0], args, []string{}); err != nil {
		l.Fatalf("containerInit: exec failed: %v", err)
	}
}

func (c *Containers) watcher(ctx context.Context, ctn *Container, f *fork.Function) {
	l := c.log.WithFields(logrus.Fields{
		"operation": "watcher",
		"id":        ctn.Container.ID,
		"name":      ctn.Container.Name,
	})
	end := make(chan error)
	go func() {
		e := f.Wait()
		l.Trace("exited")
		end <- e
	}()
	state := models.ContainerStateExited
	var e error
	select {
	case e = <-end:
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
		if ctn.Container.Systemd {
			// SIGRTMIN+3
			f.Process.Signal(syscall.Signal(37))
		} else {
			f.Process.Kill()
		}
		ctn.log.Printf("process killed")
	}
	// don't report that we're done until fs is synced
	fd, err := unix.Open(ctn.Container.Mount.Mountpoint, unix.O_RDWR|unix.O_DIRECTORY, 0)
	if err != nil {
		log.Printf("failed to open mount for sync: %v", err)
	} else {
		unix.Syncfs(fd)
		unix.Close(fd)
	}

	// run the Exit hook
	ctn.log.Printf("running script hook: exit")
	var hook *models.ContainerScriptHook
	if ctn.Container.Hooks != nil {
		hook = ctn.Container.Hooks.Exit
	}
	errstr := ""
	if e != nil {
		errstr = e.Error()
	}
	if h, err := NewHook(hook, "", ctn.Container.Mount.Mountpoint, ctn.log, map[string]string{
		"logfile":     ctn.Container.Logfile,
		"command":     *ctn.Container.Command,
		"mountpoint":  ctn.Container.Mount.Mountpoint,
		"mountkind":   ctn.Container.Mount.Kind,
		"name":        string(ctn.Container.Name),
		"id":          fmt.Sprintf("%d", ctn.Container.ID),
		"systemd":     fmt.Sprintf("%t", ctn.Container.Systemd),
		"error":       fmt.Sprintf("%t", e != nil),
		"errorstring": errstr,
	}); err != nil {
		l.WithError(err).Debug("fatal error building exit scripts")
		state = models.ContainerStateDead
	} else {
		if err = h.Run(); err != nil {
			l.WithError(err).Debug("fatal error running exit scripts")
			state = models.ContainerStateDead
		}
	}

	// process is over, set the state
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ctn.Container.State = state
	API.Store.Update(ctn)
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

func ForkInit() {
	fork.RegisterFunc("containerInit", containerInit)
	fork.Init()
}
