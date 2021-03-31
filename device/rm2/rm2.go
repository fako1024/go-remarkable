package rm2

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/fako1024/go-remarkable/device/common"
	"github.com/fako1024/go-remarkable/internal/procs"
	"github.com/fako1024/go-remarkable/multiwriter"
)

const (

	// ProcName denotes the name of the framebuffer binary
	ProcName = "xochitl"

	width  = 1872
	height = 1404

	idleInputWait  = 3 * time.Second
	forceFrameWait = 10 * time.Second

	penEventDevPath   = "/dev/input/event1"
	touchEventDevPath = "/dev/input/event2"

	pidFile  = "/tmp/" + ProcName + ".pid"
	dataPath = "/home/root/.local/share/remarkable/xochitl"
)

// RM2 denotes a Remarkable 2 device
type RM2 struct {
	offset int64

	pid           string
	fbPtr         uintptr
	fbDevice      *os.File
	penEventDev   *os.File
	touchEventDev *os.File

	lastInput time.Time
	writers   *multiwriter.MultiWriter

	*common.Device
	sync.Mutex
}

// New instantiates a new Remarkable 2 device
func New() (*RM2, error) {

	r := &RM2{
		writers: multiwriter.New(),
		Device:  common.NewDevice(dataPath),
	}

	if err := r.ensureRunning(); err != nil {
		return nil, err
	}

	return r, nil
}

// Close closes the device
func (r *RM2) Close() error {

	var errs []error
	if err := r.penEventDev.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := r.touchEventDev.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := r.fbDevice.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

func (r *RM2) openInputDevices() (err error) {

	r.penEventDev, err = os.OpenFile(penEventDevPath, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return fmt.Errorf("error opening pen input device: %s", err)
	}
	r.touchEventDev, err = os.OpenFile(touchEventDevPath, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return fmt.Errorf("error opening touch input device: %s", err)
	}

	return
}

func (r *RM2) closeInputDevices() error {

	var errs []error
	if err := r.penEventDev.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := r.touchEventDev.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

// Frame retrieves a single frame
func (r *RM2) Frame() ([]byte, error) {

	if err := r.ensureRunning(); err != nil {
		return nil, err
	}

	buf := make([]byte, width*height)
	if err := r.readFrame(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

// Stream continuously streams frames from the device
func (r *RM2) stream() error {

	if err := r.ensureRunning(); err != nil {
		return err
	}

	frameLen := width * height
	buf := make([]byte, frameLen)
	c := make(chan error)

	// Track device input in the background
	ctx, cancel := context.WithCancel(context.Background())
	go r.waitEvent(ctx, c, r.penEventDev)
	go r.waitEvent(ctx, c, r.touchEventDev)

	// Ensure termination of input detection routines
	defer func() {
		cancel()
		close(c)
	}()

	for {

		// If there are no more receiving writers left, terminate
		if r.writers.N() == 0 {
			return nil
		}

		// If no input has been detected recently, wait for it
		if time.Since(r.lastInput) > idleInputWait {
			if err := r.readInput(c); err != nil {
				return err
			}
		}

		// Read a single frame from the framebuffer
		if err := r.readFrame(buf); err != nil {
			return err
		}

		// Write / broadcast the frame to all recipients
		if _, err := r.writers.Write(buf); err != nil {
			return err
		}
	}
}

// NewStream adds a new stream recipient on the provided writer
func (r *RM2) NewStream(w io.Writer) error {

	// If there are no existing recipients, start streaming
	if r.writers.N() == 0 {
		if err := r.openInputDevices(); err != nil {
			return err
		}

		go r.stream()
	}

	// Add recipient and wait until connection is terminated
	err := <-r.writers.Append(w)

	// If there are no more recipients left, stop input detection
	// and by extension, streaming
	if r.writers.N() == 0 {
		if err := r.closeInputDevices(); err != nil {
			return err
		}
	}

	return err
}

// Dimensions returns the width + height of the underlying frame(buffer)
func (r *RM2) Dimensions() (int, int) {
	return width, height
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// init (re-)initializes the device
func (r *RM2) init(pid string) (err error) {

	r.pid = pid
	r.offset, err = procs.MemoryOffset(r.pid)
	if err != nil {
		return fmt.Errorf("error ascertaining memory offset of framebuffer device: %s", err)
	}

	r.fbDevice, err = os.OpenFile("/proc/"+r.pid+"/mem", os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return fmt.Errorf("error opening framebuffer device: %s", err)
	}
	r.fbPtr = r.fbDevice.Fd()

	return
}

// ensureRunning makes sure that the PID of the underlying xochitl process hasn't changed and
// (re-)initializes the device if necessary
func (r *RM2) ensureRunning() error {

	pid, err := procs.PIDOf(pidFile)
	if err != nil {
		return fmt.Errorf("error ascertaining PID of %s process: %s", ProcName, err)
	}
	if pid != r.pid {
		return r.init(pid)
	}

	return nil
}

func (r *RM2) readInput(c chan error) error {
	select {
	case <-time.After(forceFrameWait):
		// Force broadcast of a frame to ensure there is no deadlock
		return nil
	case err := <-c:
		// Input was detected
		return err
	}
}

func (r *RM2) waitEvent(ctx context.Context, c chan error, dev *os.File) {
	b := make([]byte, 16)
	for {

		// As soon as any data is read successfully, update timestamp
		_, err := dev.Read(b)
		if err == nil {
			r.lastInput = time.Now()
		}

		select {
		case <-ctx.Done():
			// Stop waiting and return once the context is done
			return
		case c <- err:
			// Notify a potential waiting routine about the detected input
		default:
			// Discard data if no one is waiting for it
		}
	}
}

func (r *RM2) readFrame(buf []byte) error {

	// Perform a Pread() on the framebuffer using a single syscall as it is more
	// efficient than lseek() + read()
	n, err := syscall.Pread(int(r.fbPtr), buf, r.offset)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("unexpected number of bytes read from framebuffer, want %d, have %d", len(buf), n)
	}

	return nil
}
