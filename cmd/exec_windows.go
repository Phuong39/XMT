//go:build windows

// Copyright (C) 2020 - 2022 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/iDigitalFlame/xmt/cmd/filter"
	"github.com/iDigitalFlame/xmt/device/winapi"
	"github.com/iDigitalFlame/xmt/util/bugtrack"
	"github.com/iDigitalFlame/xmt/util/xerr"
)

// modernGo is a backwards compatibility flag used to disable post-WinVista features.
const modernGo = true

// NOTE(dij): This needs to be a var as if it's a const 'UpdateProcThreadAttribute'
//            will throw an access violation.
//
// 0x100100000000 - PROCESS_CREATION_MITIGATION_POLICY_EXTENSION_POINT_DISABLE_ALWAYS_ON |
//                  PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_ON
var secProtect uint64 = 0x100100000000
var envOnce struct {
	sync.Once
	r string
	e []string
}
var secOnce struct {
	_ [0]func()
	sync.Once
	e bool
}

type closer uintptr
type file interface {
	File() (*os.File, error)
}
type fileFd interface {
	Fd() uintptr
}
type executable struct {
	r                  *os.File
	filter             *filter.Filter
	title              string
	user, domain, pass string
	closers            []io.Closer
	i                  winapi.ProcessInformation
	token, parent, m   uintptr
	sf, x, y, w, h     uint32
	mode               uint16
}
type closeFunc func() error

func envOnceFunc() {
	envOnce.e = syscall.Environ()[4:] // Removes all '=' prefixed vars
	if envOnce.r, _ = syscall.Getenv(sysRoot); len(envOnce.e) == 0 {
		envOnce.r = sysRootVar
	}
}
func secOnceFunc() {
	if m, x, _ := winapi.GetVersionNumbers(); m > 6 {
		secOnce.e = true
	} else {
		secOnce.e = m >= 6 && x >= 3
	}
}
func checkVersionSec() bool {
	secOnce.Do(secOnceFunc)
	return secOnce.e
}
func (e *executable) close() {
	if atomic.LoadUintptr(&e.i.Process) == 0 {
		return
	}
	for i := range e.closers {
		e.closers[i].Close()
	}
	e.parent, e.closers = 0, nil
	if atomic.StoreUintptr(&e.i.Process, 0); e.m > 0 {
		winapi.SetEvent(e.m)
	}
}
func (c closer) Close() error {
	return winapi.CloseHandle(uintptr(c))
}
func wait(h, m uintptr) error {
	if m == 0 {
		if r, err := winapi.WaitForSingleObject(h, -1); r != 0 {
			return err
		}
		return nil
	}
	if r, err := winapi.WaitForMultipleObjects([]uintptr{h, m}, false, -1); r != 0 {
		return err
	}
	return nil
}
func (c closeFunc) Close() error {
	return c()
}
func (e *executable) Pid() uint32 {
	return e.i.ProcessID
}

// ResumeProcess will attempt to resume the process via its PID. This will
// attempt to resume the process using an OS-dependent syscall.
//
// This will not affect already running processes.
func ResumeProcess(p uint32) error {
	// 0x800 - PROCESS_SUSPEND_RESUME
	h, err := winapi.OpenProcess(0x800, false, p)
	if err != nil {
		return err
	}
	err = winapi.ResumeProcess(h)
	winapi.CloseHandle(h)
	return err
}

// SuspendProcess will attempt to suspend the process via its PID. This will
// attempt to suspend the process using an OS-dependent syscall.
//
// This will not affect already suspended processes.
func SuspendProcess(p uint32) error {
	// 0x800 - PROCESS_SUSPEND_RESUME
	h, err := winapi.OpenProcess(0x800, false, p)
	if err != nil {
		return err
	}
	err = winapi.SuspendProcess(h)
	winapi.CloseHandle(h)
	return err
}
func (e *executable) Resume() error {
	return winapi.ResumeProcess(e.i.Process)
}
func (e *executable) Suspend() error {
	return winapi.SuspendProcess(e.i.Process)
}
func (e *executable) isStarted() bool {
	// BUG(dij): I think this causes some 'ErrAlreadyStarted' issues.
	//           keep watch on this.
	//
	// return e.i.ProcessID > 0 // && e.i.Process > 0
	return atomic.LoadUint32(&e.i.ProcessID) > 0 || e.i.Process > 0
}
func (e *executable) isRunning() bool {
	return e.i.Process > 0
}
func (e *executable) Handle() uintptr {
	return e.i.Process
}
func (e *executable) SetToken(t uintptr) {
	e.token = t
}
func (e *executable) SetFullscreen(f bool) {
	// 0x20 - STARTF_RUNFULLSCREEN
	if f {
		e.sf |= 0x20
	} else {
		e.sf = e.sf &^ 0x20
	}
}
func (e *executable) SetWindowDisplay(m int) {
	// 0x1 - STARTF_USESHOWWINDOW
	if m < 0 {
		e.sf = e.sf &^ 0x1
	} else {
		e.sf |= 0x1
		e.mode = uint16(m)
	}
}
func (e *executable) SetWindowTitle(s string) {
	// 0x1000 - STARTF_TITLEISAPPID
	if len(s) > 0 {
		e.sf |= 0x1000
		e.title = s
	} else {
		e.sf, e.title = e.sf&^0x1000, ""
	}
}
func (e *executable) SetLogin(u, d, p string) {
	if e.user, e.domain, e.pass = u, d, p; len(d) == 0 {
		d = "."
	}
}
func (executable) SetUID(_ int32, _ *Process) {}
func (executable) SetGID(_ int32, _ *Process) {}
func (e *executable) SetWindowSize(w, h uint32) {
	// 0x2 - STARTF_USESIZE
	e.sf |= 0x2
	e.w, e.h = w, h
}
func (executable) SetNoWindow(h bool, p *Process) {
	// 0x8000000 - CREATE_NO_WINDOW
	if h {
		p.flags |= 0x8000000
	} else {
		p.flags = p.flags &^ 0x8000000
	}
}
func (executable) SetDetached(d bool, p *Process) {
	// 0x8  - DETACHED_PROCESS
	// 0x10 - CREATE_NEW_CONSOLE
	if d {
		p.flags = (p.flags | 0x8) &^ 0x10
	} else {
		p.flags = p.flags &^ 0x8
	}
}
func (executable) SetChroot(_ string, _ *Process) {}
func (executable) SetSuspended(s bool, p *Process) {
	// 0x4 - CREATE_SUSPENDED
	if s {
		p.flags |= 0x4
	} else {
		p.flags = p.flags &^ 0x4
	}
}
func (e *executable) SetWindowPosition(x, y uint32) {
	// 0x4 - STARTF_USEPOSITION
	e.sf |= 0x4
	e.x, e.y = x, y
}
func (*executable) SetNewConsole(c bool, p *Process) {
	// 0x10 - CREATE_NEW_CONSOLE
	if c {
		p.flags |= 0x10
	} else {
		p.flags = p.flags &^ 0x10
	}
}
func (e *executable) kill(x uint32, p *Process) error {
	p.exit = x
	return winapi.TerminateProcess(e.i.Process, x)
}
func createEnvBlock(env []string, split bool) []string {
	if envOnce.Do(envOnceFunc); len(env) == 0 && !split {
		return envOnce.e
	}
	r := make([]string, len(env), len(env)+len(envOnce.e))
	if copy(r, env); !split {
		// NOTE(dij): If split == true, do NOT add any env vars, but DO
		//            check and add %SYSTEMROOT% if it doesn't exist in the
		//            supplied block.
		return append(r, envOnce.e...)
	}
	for i := range r {
		if len(r) > 11 {
			if x := strings.IndexByte(r[i], '='); x > 9 {
				if strings.EqualFold(r[i][:x], sysRoot) {
					return r
				}
			}
		}
	}
	return append(r, sysRoot+"="+envOnce.r)
}
func (e *executable) wait(x context.Context, p *Process) {
	if bugtrack.Enabled {
		defer bugtrack.Recover("cmd.executable.wait()")
	}
	var (
		w   = make(chan error)
		err error
	)
	if e.m, err = winapi.CreateEvent(nil, false, false, ""); err != nil {
		if bugtrack.Enabled {
			bugtrack.Track("cmd.executable.wait(): Creating Event failed, falling back to single wait: %s", err.Error())
		}
		err = nil
	}
	go func() {
		if atomic.LoadUintptr(&e.i.Process) > 0 {
			if bugtrack.Enabled {
				defer bugtrack.Recover("cmd.executable.wait.func1()")
			}
			w <- wait(e.i.Process, e.m)
		}
		close(w)
	}()
	select {
	case err = <-w:
	case <-x.Done():
	}
	if e.m > 0 {
		winapi.CloseHandle(e.m)
		e.m = 0
	}
	if err != nil {
		p.stopWith(exitStopped, err)
		return
	}
	if err2 := x.Err(); err2 != nil {
		p.stopWith(exitStopped, err2)
		return
	}
	if atomic.SwapUint32(&p.cookie, 2) == 2 || atomic.LoadUintptr(&e.i.Process) == 0 {
		p.stopWith(0, nil)
		return
	}
	if err = winapi.GetExitCodeProcess(e.i.Process, &p.exit); err != nil {
		p.stopWith(exitStopped, err)
		return
	}
	if p.exit != 0 {
		p.stopWith(p.exit, &ExitError{Exit: p.exit})
		return
	}
	p.stopWith(p.exit, nil)
}
func (e *executable) writer(w io.Writer) (uintptr, error) {
	h, c, err := e.writerToHandle(w)
	if err != nil {
		return 0, err
	}
	return e.addRetHandle(c, h)
}
func (e *executable) reader(r io.Reader) (uintptr, error) {
	h, c, err := e.readerToHandle(r)
	if err != nil {
		return 0, err
	}
	return e.addRetHandle(c, h)
}
func (e *executable) SetParent(f *filter.Filter, p *Process) {
	if e.filter = f; f != nil {
		e.SetNewConsole(true, p)
	}
}
func (e *executable) addRetHandle(c bool, h uintptr) (uintptr, error) {
	if e.parent == 0 {
		if c {
			e.closers = append(e.closers, closer(h))
		}
		return h, nil
	}
	var (
		n   uintptr
		err = winapi.DuplicateHandle(winapi.CurrentProcess, h, e.parent, &n, 0, true, 0x2)
		// 0x2 - DUPLICATE_SAME_ACCESS
	)
	if c {
		winapi.CloseHandle(h)
	}
	if err != nil {
		return 0, xerr.Wrap("DuplicateHandle", err)
	}
	return n, nil
}
func (e *executable) readerToHandle(r io.Reader) (uintptr, bool, error) {
	if r == nil {
		// 0 - READONLY
		f, err := os.OpenFile(os.DevNull, 0, 0)
		if err != nil {
			return 0, false, xerr.Wrap("cannot open NULL device", err)
		}
		e.closers = append(e.closers, f)
		return f.Fd(), false, nil
	}
	switch i := r.(type) {
	case file:
		f, err := i.File()
		if err != nil {
			return 0, false, xerr.Wrap("cannot obtain file handle", err)
		}
		// Closeable "c" is true, since this /should/ be a separate
		// handle from the initial "File" type.
		//
		// NOTE(dij): Technically on Windows this will always fail
		// See: https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/net/fd_windows.go;l=175
		//
		// BUT if we're going off the *nix implementation, it would be a
		// duplicate handle, and safe to close.
		return f.Fd(), true, nil
	case fileFd:
		return i.Fd(), false, nil
	}
	x, y, err := os.Pipe()
	if err != nil {
		return 0, false, xerr.Wrap("cannot create Pipe", err)
	}
	e.closers = append(e.closers, x)
	// e.closers = append(e.closers, y)
	go func() {
		if bugtrack.Enabled {
			defer bugtrack.Recover("cmd.options.readerToHandle.func1()")
		}
		io.Copy(y, r)
		y.Close()
	}()
	return x.Fd(), false, nil
}
func (e *executable) writerToHandle(w io.Writer) (uintptr, bool, error) {
	if w == nil {
		// 1 - WRITEONLY
		f, err := os.OpenFile(os.DevNull, 1, 0)
		if err != nil {
			return 0, false, xerr.Wrap("cannot open NULL device", err)
		}
		e.closers = append(e.closers, f)
		return f.Fd(), false, nil
	}
	switch i := w.(type) {
	case file:
		f, err := i.File()
		if err != nil {
			return 0, false, xerr.Wrap("cannot obtain file handle", err)
		}
		// Closeable "c" is true, since this /should/ be a separate handle from
		// the initial "File" type.
		//
		// NOTE(dij): Technically on Windows this will always fail
		// See: https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/net/fd_windows.go;l=175
		//
		// BUT if we're going off the *nix implementation, it would be a duplicate
		// handle, and safe to close.
		return f.Fd(), true, nil
	case fileFd:
		return i.Fd(), false, nil
	}
	x, y, err := os.Pipe()
	if err != nil {
		return 0, false, xerr.Wrap("cannot create Pipe", err)
	}
	e.closers = append(e.closers, y)
	e.closers = append(e.closers, x)
	go func() {
		if bugtrack.Enabled {
			defer bugtrack.Recover("cmd.options.writerToHandle.func1()")
		}
		io.Copy(w, x)
	}()
	return y.Fd(), false, nil
}
func (e *executable) start(x context.Context, p *Process, sus bool) error {
	r, err := exec.LookPath(p.Args[0])
	if err != nil {
		return err
	}
	v, y, err := e.startInfo()
	if err != nil {
		return err
	}
	if v != nil && v.AttributeList != nil {
		e.closers = append(e.closers, closeFunc(func() error {
			return winapi.DeleteProcThreadAttributeList(v.AttributeList)
		}))
	}
	if p.Stderr != nil || p.Stdout != nil || p.Stdin != nil {
		var si, so, se uintptr
		if si, err = e.reader(p.Stdin); err != nil {
			if e.parent > 0 {
				winapi.CloseHandle(e.parent)
				e.parent = 0
			}
			return err
		}
		if so, err = e.writer(p.Stdout); err != nil {
			if e.parent > 0 {
				winapi.CloseHandle(e.parent)
				e.parent = 0
			}
			return err
		}
		if p.Stderr == p.Stdout {
			se = so
		} else if se, err = e.writer(p.Stderr); err != nil {
			if e.parent > 0 {
				winapi.CloseHandle(e.parent)
				e.parent = 0
			}
			return err
		}
		if v != nil {
			v.StartupInfo.StdErr = se
			v.StartupInfo.StdInput = si
			v.StartupInfo.StdOutput = so
			v.StartupInfo.Flags |= 0x100
			// 0x100 - STARTF_USESTDHANDLES
		} else if y != nil {
			y.StdErr, y.StdInput, y.StdOutput = se, si, so
			y.Flags |= 0x100
			// 0x100 - STARTF_USESTDHANDLES
		}
	}
	u := e.token
	if runtime.LockOSThread(); u == 0 && e.parent == 0 {
		// NOTE(dij): Handle threads that currently have an impersonated Token
		//            set. This will trigger this function call to use
		//            'CreateProcessWithToken' instead of 'CreateProcess'.
		//            This is only called IF there is no parent Process set, as
		//            Windows permissions cause some fucky stuff to happen.
		//            Failing silently is fine.
		//
		// NOTE(dij): Added a 'IsUserLoginToken' token to check the Token origin
		//            to see if it's an impersonated user token or a stolen elevated
		//            process token, as impersonated user tokens do NOT like being
		//            ran with 'CreateProcessWithToken'.
		//
		// BUG(dij):  Watch this function call, as it can cause problems when launching
		//            non-parent processes under impersonation.
		//
		// (old was 0xF01FF - TOKEN_ALL_ACCESS)
		// 0x200EF - TOKEN_ASSIGN_PRIMARY | TOKEN_DUPLICATE | TOKEN_IMPERSONATE |
		//            TOKEN_QUERY | TOKEN_WRITE (STANDARD_RIGHTS_WRITE | TOKEN_ADJUST_PRIVILEGES |
		//            TOKEN_ADJUST_GROUPS | TOKEN_ADJUST_DEFAULT)
		if winapi.OpenThreadToken(winapi.CurrentThread, 0xF01FF, true, &u); u > 0 && winapi.IsUserLoginToken(u) {
			if u = 0; bugtrack.Enabled {
				bugtrack.Track("cmd.executable.start(): Removing user login token.")
			}
		}
	}
	if sus {
		// 0x4 - CREATE_SUSPENDED
		p.flags |= 0x4
	}
	if e.r != nil {
		e.r.Close()
		e.r = nil
	}
	switch z := createEnvBlock(p.Env, p.split); {
	case len(e.user) > 0:
		if bugtrack.Enabled {
			bugtrack.Track("cmd.executable.start(): Using API call CreateProcessWithLogin for execution.")
		}
		// NOTE(dij): Network Only (0x2) logins seem to fail here.
		//            Apparently a 0 value works!
		// BUG(dij): Tracking "A logon request contained an invalid logon type value."
		//           when launching with "0x2".
		// 0x1 - LOGON_WITH_PROFILE
		// 0x0 - *shrug*
		err = winapi.CreateProcessWithLogin(e.user, e.domain, e.pass, 0x0, r, strings.Join(p.Args, " "), p.flags, z, p.Dir, y, v, &e.i)
	case u > 0:
		if bugtrack.Enabled {
			bugtrack.Track("cmd.executable.start(): Using API call CreateProcessWithToken for execution.")
		}
		// 0x2 - LOGON_NETCREDENTIALS_ONLY
		err = winapi.CreateProcessWithToken(u, 0x2, r, strings.Join(p.Args, " "), p.flags, z, p.Dir, y, v, &e.i)
	default:
		if bugtrack.Enabled {
			bugtrack.Track("cmd.executable.start(): Using API call CreateProcess for execution.")
		}
		err = winapi.CreateProcess(r, strings.Join(p.Args, " "), nil, nil, true, p.flags, z, p.Dir, y, v, &e.i)
	}
	if runtime.UnlockOSThread(); e.parent > 0 {
		winapi.CloseHandle(e.parent)
		e.parent = 0
	}
	if err != nil {
		for i := range e.closers {
			e.closers[i].Close()
		}
		return err
	}
	winapi.CloseHandle(e.i.Thread)
	if e.closers = append(e.closers, closer(e.i.Process)); sus {
		return nil
	}
	go e.wait(x, p)
	return nil
}
func (e *executable) startInfo() (*winapi.StartupInfoEx, *winapi.StartupInfo, error) {
	var (
		x   winapi.StartupInfoEx
		err error
	)
	e.close()
	x.StartupInfo.XSize, x.StartupInfo.YSize = e.w, e.h
	x.StartupInfo.Flags, x.StartupInfo.ShowWindow = e.sf, e.mode
	if x.StartupInfo.X, x.StartupInfo.Y = e.x, e.y; len(e.title) > 0 {
		if x.StartupInfo.Title, err = winapi.UTF16PtrFromString(e.title); err != nil {
			return nil, nil, xerr.Wrap("cannot convert title", err)
		}
	}
	if x.StartupInfo.Cb = uint32(unsafe.Sizeof(x)); !modernGo {
		return nil, &x.StartupInfo, nil
	}
	if e.filter != nil && !e.filter.Empty() {
		// (old 0x100CC1 - SYNCHRONIZE | PROCESS_DUP_HANDLE | PROCESS_CREATE_PROCESS |
		//                  PROCESS_QUERY_INFORMATION | PROCESS_SUSPEND_RESUME | PROCESS_TERMINATE)
		// (old 0x4C0 - PROCESS_QUERY_INFORMATION | PROCESS_DUP_HANDLE | PROCESS_CREATE_PROCESS)
		//
		// 0x10C0 - PROCESS_CREATE_PROCESS | PROCESS_DUP_HANDLE | PROCESS_QUERY_LIMITED_INFORMATION
		if e.parent, err = e.filter.HandleFunc(0x10C0, nil); err != nil {
			return nil, nil, err
		}
		// FIXED(dij): Apparently sometimes this isn't closed? It seems to /only/
		//             happen during spawn? Look into this later.
		//
		//             FIX: Close handle immediately after spawning process!
		// e.closers = append(e.closers, closer(e.parent))
	}
	var c uint32
	// NOTE(dij): SecProtect isn't allowed until Windows 8.1 and Windows Server 2012R2
	//            Thanks for the super small blurb of text on it Microsoft >:[
	switch v := checkVersionSec(); {
	case !v && e.parent == 0: // No sec and no parent
		return nil, &x.StartupInfo, nil
	case !v && e.parent > 0: // No sec, has parent (1 slot)
		fallthrough
	case v && e.parent == 0: // Has sec, no parent (1 slot)
		c = 1
	case v && e.parent > 0: // Has sec, has parent (2 slots)
		c = 2
	}
	x.AttributeList = &winapi.StartupAttributes{Count: c}
	if x.StartupInfo.Cb = uint32(unsafe.Sizeof(x)); e.parent > 0 {
		// 0x20000 - PROC_THREAD_ATTRIBUTE_PARENT_PROCESS
		if err = winapi.UpdateProcThreadAttribute(x.AttributeList, 0x20000, unsafe.Pointer(&e.parent), uint64(unsafe.Sizeof(e.parent)), nil, nil); err != nil {
			winapi.DeleteProcThreadAttributeList(x.AttributeList)
			winapi.CloseHandle(e.parent)
			e.parent = 0
			return nil, nil, xerr.Wrap("UpdateProcThreadAttribute", err)
		}
		c--
	}
	if c == 1 {
		// 0x20007 - PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY
		if err = winapi.UpdateProcThreadAttribute(x.AttributeList, 0x20007, unsafe.Pointer(&secProtect), uint64(unsafe.Sizeof(secProtect)), nil, nil); err != nil {
			if winapi.DeleteProcThreadAttributeList(x.AttributeList); e.parent > 0 {
				winapi.CloseHandle(e.parent)
				e.parent = 0
			}
			return nil, nil, xerr.Wrap("UpdateProcThreadAttribute", err)
		}
	}
	return &x, nil, nil
}
