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

	"github.com/jlowellwofford/go-fork"
	"github.com/jlowellwofford/imageapi/models"
	"github.com/jlowellwofford/uinit"
	"golang.org/x/sys/unix"
)

type containerStateChange struct {
	id    int64
	state models.ContainerState
}

type container struct {
	ctn    *models.Container
	log    *log.Logger
	mnt    string
	cancel context.CancelFunc
}

type ContainersType struct {
	next  int64
	ctns  map[int64]*container
	mutex *sync.Mutex
}

func (c *ContainersType) Init() {
	c.next = 0
	c.ctns = make(map[int64]*container)
	c.mutex = &sync.Mutex{}
}

func (c *ContainersType) List() (r []*models.Container) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, ctn := range c.ctns {
		r = append(r, ctn.ctn)
	}
	return
}

func (c *ContainersType) Get(id int64) (*models.Container, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if r, ok := c.ctns[id]; ok {
		return r.ctn, nil
	}
	return nil, fmt.Errorf("no container by id %d", id)
}

func (c *ContainersType) Create(ctn *models.Container) (r *models.Container, err error) {
	// This creates a container in our list, and activates its initial state
	// find the mount
	n := &container{}
	switch *ctn.Mount.Kind {
	case models.MountKindOverlay:
		mnt, e := MountsOverlay.Get(*ctn.Mount.ID)
		if e != nil {
			return nil, fmt.Errorf("failed to get mount for container: %v", e)
		}
		n.mnt = mnt.Mountpoint
		MountsOverlay.RefAdd(mnt.ID, 1)
	case models.MountKindRbd:
		mnt, e := MountsRbd.Get(*ctn.Mount.ID)
		if e != nil {
			return nil, fmt.Errorf("failed to get mount for container: %v", e)
		}
		n.mnt = mnt.Mountpoint
		MountsRbd.RefAdd(*mnt.ID, 1)
	}
	// ok, we've got a valid mountpoint
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ctn.ID = c.next

	// set up logger
	ctn.Logfile = path.Join(logDir, fmt.Sprintf("%d-%d.log", ctn.ID, time.Now().Unix()))
	f, e := os.Create(ctn.Logfile)
	if e != nil {
		return nil, fmt.Errorf("could not create log file: %v", e)
	}
	n.log = log.New(f, fmt.Sprintf("container(%d): ", ctn.ID), log.Ldate|log.Ltime|log.Lmsgprefix)
	n.log.Printf("container created")

	// handle initial state
	switch ctn.State {
	case models.ContainerStateRunning:
		// run it
	case models.ContainerStateRestarting,
		models.ContainerStatePaused,
		models.ContainerStateExited,
		models.ContainerStateDead:
		return nil, fmt.Errorf("requested invalid initial container state: %s.  valid initial states: [ %s, %s ]", ctn.State, models.ContainerStateCreated, models.ContainerStateRunning)
	case models.ContainerStateCreated:
		fallthrough
	default: // wasn't specified
		ctn.State = models.ContainerStateCreated
	}

	// container is ready to be entered
	c.ctns[ctn.ID] = n
	c.next++

	return ctn, nil
}

func (c *ContainersType) SetState(id int64, state models.ContainerState) (err error) {
	var ctn *container
	var ok bool
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if ctn, ok = c.ctns[id]; !ok {
		return fmt.Errorf("invalid container id: %d", id)
	}
	// handle state request
	switch state {
	case models.ContainerStateRunning:
		if ctn.ctn.State == state {
			return
		}
		// run it
	case models.ContainerStateExited:
		if ctn.ctn.State == state {
			return
		}
		// stop it
	case models.ContainerStatePaused,
		models.ContainerStateRestarting:
		return fmt.Errorf("container state is not yet implemented: %s", state)
	default: // something not valid
		return fmt.Errorf("can't set state to: %s.  valid initial states: [ %s, %s, %s, %s ]", state,
			models.ContainerStateRunning,
			models.ContainerStateExited,
			models.ContainerStateRestarting,
			models.ContainerStatePaused)
	}
	return
}

func (c *ContainersType) Delete(id int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var ctn *container
	var ok bool
	if ctn, ok = c.ctns[id]; !ok {
		return fmt.Errorf("invalid container id: %d", id)
	}
	switch ctn.ctn.State {
	//case models.ContainerStatePaused:
	//case models.ContainerStateRestarting:
	case models.ContainerStateRunning:
		// stop it
	}
	ctn.log.Printf("container deleted")
	ctn.log.Writer().(io.WriteCloser).Close()
	delete(c.ctns, id)
	switch *ctn.ctn.Mount.Kind {
	case models.MountKindOverlay:
		MountsOverlay.RefAdd(*ctn.ctn.Mount.ID, -1)
	case models.MountKindRbd:
		MountsRbd.RefAdd(*ctn.ctn.Mount.ID, -1)
	}
	return
}

func (c *ContainersType) stop(ctn *container) error {
	// this will need to be updated if we implement pause and/or restart
	if ctn.ctn.State != models.ContainerStateRunning {
		// it's not an error to stop an already stopped container
		return nil
	}
	// we trust the watcher to take care of everything for us
	ctn.cancel()
	return nil
}

var systemdTmpfs = []string{"/run", "/run/lock", "/tmp", "/sys/fs/cgrooup/systemd", "/var/lib/journal"}

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
	f.SysProcAttr.Cloneflags = syscall.CLONE_NEWNS | syscall.CLONE_NEWPID
	if err := f.Fork(ctn.mnt, args); err != nil {
		return fmt.Errorf("clone: failed to start pid_init: %v", err)
	}

	// 4. Launch the process watcher
	ctx := context.Background()
	ctx, ctn.cancel = context.WithCancel(ctx)
	go c.watcher(ctx, ctn, f)
	return
}

// this is run as a separate process
func containerInit(mountpoint string, args []string) {
	// 0. setup logging
	l := log.New(os.Stdout, fmt.Sprintf("init: "), log.Ldate|log.Ltime|log.Lmsgprefix)

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

	// 3. Prep pid stuff
	l.Print("mounting /proc")
	if isMountpoint("/proc") {
		if err := unix.Unmount("/proc", unix.MNT_DETACH); err != nil {
			l.Printf("warn: could not unmount /proc, will overlay instead: %v", err)
		}
	}
	if err := unix.Mount("proc", "/proc", "proc", unix.MS_NODEV|unix.MS_NOSUID|unix.MS_NOEXEC, ""); err != nil {
		l.Fatalf("/proc mount failed: %v", err)
	}
	l.Print("executing init")
	if err := unix.Exec(args[0], args, []string{}); err != nil {
		l.Fatalf("containerInit: exec failed: %v", err)
	}
	return
}

func (c *ContainersType) watcher(ctx context.Context, ctn *container, f *fork.Function) {
	end := make(chan error)
	go func() {
		e := f.Wait()
		end <- e
	}()
	state := models.ContainerStateExited
	select {
	case e := <-end:
		if e != nil {
			ctn.log.Printf("process ended in error state: %v", e)
			state = models.ContainerStateDead
		}
		ctn.log.Printf("process ended, exit code (0)")
	case <-ctx.Done():
		// signal the process to stop
		// TODO: be smarter about the signals we send
		f.Process.Kill()
		ctn.log.Printf("process killed")
	}
	// process is over, set the state
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ctn.ctn.State = state
	ctn.log.Printf("container state: %s", state)
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
var specialMounts = []string{"/dev", "/proc", "/sys", "/run"}

// this is the workhorse for all schemes
// it preforms the root-moving dance
func moveRoot(newRoot string) (err error) {
	// 1. move special mounts
	for _, mount := range specialMounts {
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
