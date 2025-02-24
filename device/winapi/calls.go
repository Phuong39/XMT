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

package winapi

import (
	"syscall"
	"unsafe"

	"github.com/iDigitalFlame/xmt/util/bugtrack"
)

// RevertToSelf Windows API Call
//   The RevertToSelf function terminates the impersonation of a client application.
//
// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-revqerttoself
//
// Calls 'NtSetInformationThread' under the hood.
func RevertToSelf() error {
	//r, _, err := syscall.SyscallN(funcRevertToSelf.address())
	// 0x5 - ThreadImpersonationToken
	if r, _, _ := syscall.SyscallN(funcNtSetInformationThread.address(), CurrentThread, 0x5, 0, 0); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// IsDebuggerPresent Windows API Call
//   Determines whether the calling process is being debugged by a user-mode
//   debugger.
//
// https://docs.microsoft.com/en-us/windows/win32/api/debugapi/nf-debugapi-isdebuggerpresent
func IsDebuggerPresent() bool {
	r, _, _ := syscall.SyscallN(funcIsDebuggerPresent.address())
	return r > 0
}

// BlockInput Windows API Call
//   Blocks keyboard and mouse input events from reaching applications.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-blockinput
func BlockInput(e bool) error {
	var v uint32
	if e {
		v = 1
	}
	r, _, err := syscall.SyscallN(funcBlockInput.address(), uintptr(v))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// SetEvent Windows API Call
//   Sets the specified event object to the signaled state.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-setevent
func SetEvent(h uintptr) error {
	if r, _, err := syscall.SyscallN(funcSetEvent.address(), h); r == 0 {
		return unboxError(err)
	}
	return nil
}

// CloseHandle Windows API Call
//   Closes an open object handle.
//
// https://docs.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
//
// Re-targeted to use 'NtClose' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntclose
func CloseHandle(h uintptr) error {
	r, _, _ := syscall.SyscallN(funcNtClose.address(), h)
	if bugtrack.Enabled { // Trace Bad Handles
		bugtrack.Track("winapi.CloseHandle() h=0x%X, r=0x%X", h, r)
	}
	if r > 0 {
		return formatNtError(r)
	}
	return nil
}

// GetCurrentProcessID Windows API Call
//   Retrieves the process identifier of the calling process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentprocessid
func GetCurrentProcessID() uint32 {
	r, _, _ := syscall.SyscallN(funcGetCurrentProcessID.address())
	return uint32(r)
}

// RegFlushKey Windows API Call
//   Writes all the attributes of the specified open registry key into the registry.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regflushkey
func RegFlushKey(h uintptr) error {
	r, _, err := syscall.SyscallN(funcRegFlushKey.address(), h)
	if r > 0 {
		return unboxError(err)
	}
	return nil
}

// IsWow64Process Windows API Call
//   Determines whether the specified process is running under WOW64 or an
//   Intel64 of x64 processor.
//
// https://docs.microsoft.com/en-us/windows/win32/api/wow64apiset/nf-wow64apiset-iswow64process
func IsWow64Process() (bool, error) {
	if funcIsWow64Process.find() != nil {
		// Running on "true" x86.
		return false, nil
	}
	var (
		o         bool
		r, _, err = syscall.SyscallN(funcIsWow64Process.address(), CurrentProcess, uintptr(unsafe.Pointer(&o)))
	)
	if r == 0 {
		return false, unboxError(err)
	}
	return o, nil
}

// ResumeProcess Windows API Call
//   Resumes a process and all it's threads.
//
// http://www.pinvoke.net/default.aspx/ntdll/NtResumeProcess.html
func ResumeProcess(h uintptr) error {
	if r, _, _ := syscall.SyscallN(funcNtResumeProcess.address(), h); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// SuspendProcess Windows API Call
//   Suspends a process and all it's threads.
//
// http://www.pinvoke.net/default.aspx/ntdll/NtSuspendProcess.html
func SuspendProcess(h uintptr) error {
	if r, _, _ := syscall.SyscallN(funcNtSuspendProcess.address(), h); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// GetLogicalDrives Windows API Call
//   Retrieves a bitmask representing the currently available disk drives.
//
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getlogicaldrives
func GetLogicalDrives() (uint32, error) {
	r, _, err := syscall.SyscallN(funcGetLogicalDrives.address())
	if r == 0 {
		return 0, unboxError(err)
	}
	return uint32(r), nil
}

// DisconnectNamedPipe Windows API Call
//   Disconnects the server end of a named pipe instance from a client process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/namedpipeapi/nf-namedpipeapi-disconnectnamedpipe
func DisconnectNamedPipe(h uintptr) error {
	r, _, err := syscall.SyscallN(funcDisconnectNamedPipe.address(), h)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// ResumeThread Windows API Call
//   Decrements a thread's suspend count. When the suspend count is decremented
//   to zero, the execution of the thread is resumed.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-resumethread
//
// Re-targeted to use 'NtResumeThread' instead.
// https://docs.rs/ntapi/0.3.1/ntapi/ntpsapi/type.NtResumeThread.html
func ResumeThread(h uintptr) (uint32, error) {
	var c uint32
	if r, _, _ := syscall.SyscallN(funcNtResumeThread.address(), h, uintptr(unsafe.Pointer(&c))); r > 0 {
		return 0, formatNtError(r)
	}
	return c, nil
}

// GetProcessID Windows API Call
//   Retrieves the process identifier of the specified process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getprocessid
//
// Calls 'NtQueryInformationProcess' instead under the hood.
func GetProcessID(h uintptr) (uint32, error) {
	var (
		p       processBasicInfo
		r, _, _ = syscall.SyscallN(funcNtQueryInformationProcess.address(), h, 0, uintptr(unsafe.Pointer(&p)), unsafe.Sizeof(p), 0)
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return uint32(p.UniqueProcessID), nil
}

// ImpersonateLoggedOnUser Windows API Call
//   The ImpersonateLoggedOnUser function lets the calling thread impersonate the
//   security context of a logged-on user. The user is represented by a token handle.
//
// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-impersonateloggedonuser
func ImpersonateLoggedOnUser(h uintptr) error {
	r, _, err := syscall.SyscallN(funcImpersonateLoggedOnUser.address(), h)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// SuspendThread Windows API Call
//   Suspends the specified thread.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-suspendthread
//
// Re-targeted to use 'NtSuspendThread' instead.
// https://docs.rs/ntapi/0.3.1/ntapi/ntpsapi/type.NtSuspendThread.html
func SuspendThread(h uintptr) (uint32, error) {
	var c uint32
	if r, _, _ := syscall.SyscallN(funcNtSuspendThread.address(), h, uintptr(unsafe.Pointer(&c))); r > 0 {
		return 0, formatNtError(r)
	}
	return c, nil
}

// CancelIoEx Windows API Call
//   Marks any outstanding I/O operations for the specified file handle. The
//   function only cancels I/O operations in the current process, regardless of
//   which thread created the I/O operation.
//
// https://docs.microsoft.com/en-us/windows/win32/fileio/cancelioex-func
func CancelIoEx(h uintptr, o *Overlapped) error {
	r, _, err := syscall.SyscallN(funcCancelIoEx.address(), h, uintptr(unsafe.Pointer(o)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// TerminateThread Windows API Call
//   Terminates a thread.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-terminatethread
//
// Re-targeted to use 'NtTerminateThread' instead.
// http://pinvoke.net/default.aspx/ntdll/NtTerminateThread.html
func TerminateThread(h uintptr, e uint32) error {
	if r, _, _ := syscall.SyscallN(funcNtTerminateThread.address(), h, uintptr(e)); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// RegDeleteKey Windows API Call
//   Deletes a subkey and its values. Note that key names are not case sensitive.
//   ONLY DELETES EMPTY SUBKEYS. (invalid argument if non-empty)
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regdeletekeyw
func RegDeleteKey(h uintptr, path string) error {
	return RegDeleteKeyEx(h, path, 0)
}

// RegDeleteTree Windows API Call
//   Deletes the subkeys and values of the specified key recursively.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regdeletetreew
func RegDeleteTree(h uintptr, path string) error {
	p, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(funcRegDeleteTree.address(), h, uintptr(unsafe.Pointer(p)))
	if r > 0 {
		return unboxError(err1)
	}
	return nil
}

// TerminateProcess Windows API Call
//   Terminates the specified process and all of its threads.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-terminateprocess
//
// Re-targeted to use 'NtTerminateProcess' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntddk/nf-ntddk-zwterminateprocess
func TerminateProcess(h uintptr, e uint32) error {
	if r, _, _ := syscall.SyscallN(funcNtTerminateProcess.address(), h, uintptr(e)); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// ImpersonateNamedPipeClient Windows API Call
//   The ImpersonateNamedPipeClient function impersonates a named-pipe client
//   application.
//
// https://docs.microsoft.com/en-us/windows/win32/api/namedpipeapi/nf-namedpipeapi-impersonatenamedpipeclient
func ImpersonateNamedPipeClient(h uintptr) error {
	r, _, err := syscall.SyscallN(funcImpersonateNamedPipeClient.address(), h)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// SetThreadToken Windows API Call
//   The SetThreadToken function assigns an impersonation token to a thread. The
//   function can also cause a thread to stop using an impersonation token.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-setthreadtoken
func SetThreadToken(h *uintptr, t uintptr) error {
	r, _, err := syscall.SyscallN(funcSetThreadToken.address(), uintptr(unsafe.Pointer(h)), t)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// RegDeleteValue Windows API Call
//   Removes a named value from the specified registry key. Note that value names
//   are not case sensitive.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regdeletevaluew
func RegDeleteValue(h uintptr, path string) error {
	p, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(funcRegDeleteValue.address(), h, uintptr(unsafe.Pointer(p)))
	if r > 0 {
		return unboxError(err1)
	}
	return nil
}

// GetExitCodeThread Windows API Call
//   Retrieves the termination status of the specified thread.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getexitcodethread
//
// Calls 'NtQueryInformationThread' under the hood.
func GetExitCodeThread(h uintptr, e *uint32) error {
	var (
		t       threadBasicInfo
		r, _, _ = syscall.SyscallN(funcNtQueryInformationThread.address(), h, 0, uintptr(unsafe.Pointer(&t)), unsafe.Sizeof(t), 0)
	)
	if r > 0 {
		return formatNtError(r)
	}
	*e = t.ExitStatus
	return nil
}

// RtlSetProcessIsCritical Windows API Call
//   Set process system critical status.
//   Returns the last Critical status.
//
// https://www.codeproject.com/articles/43405/protecting-your-process-with-rtlsetprocessiscriti
func RtlSetProcessIsCritical(c bool) (bool, error) {
	var s, o byte
	if c {
		s = 1
	}
	r, _, err := syscall.SyscallN(funcRtlSetProcessIsCritical.address(), uintptr(s), uintptr(unsafe.Pointer(&o)), 0)
	if r > 0 {
		return false, unboxError(err)
	}
	return o == 1, nil
}

// NtFreeVirtualMemory Windows API Call
//   The NtFreeVirtualMemory routine releases, decommits, or both releases and
//   decommits, a region of pages within the virtual address space of a specified
//   process.
//
// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntfreevirtualmemory
func NtFreeVirtualMemory(h, address uintptr) error {
	var (
		s       uint32
		r, _, _ = syscall.SyscallN(
			funcNtFreeVirtualMemory.address(), h, uintptr(unsafe.Pointer(&address)), uintptr(unsafe.Pointer(&s)),
			0x8000,
		)
		// 0x8000 - MEM_RELEASE
	)
	if r > 0 {
		return formatNtError(r)
	}
	return nil
}

// GetExitCodeProcess Windows API Call
//   Retrieves the termination status of the specified process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getexitcodeprocess
//
// Calls 'NtQueryInformationProcess' under the hood.
func GetExitCodeProcess(h uintptr, e *uint32) error {
	var (
		p       processBasicInfo
		r, _, _ = syscall.SyscallN(funcNtQueryInformationProcess.address(), h, 0, uintptr(unsafe.Pointer(&p)), unsafe.Sizeof(p), 0)
	)
	if r > 0 {
		return formatNtError(r)
	}
	*e = p.ExitStatus
	return nil
}

// Thread32Next Windows API Call
//   Retrieves information about the next thread of any process encountered in
//   the system memory snapshot.
//
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-thread32next
func Thread32Next(h uintptr, e *ThreadEntry32) error {
	r, _, err := syscall.SyscallN(funcThread32Next.address(), h, uintptr(unsafe.Pointer(e)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// Thread32First Windows API Call
//   Retrieves information about the first thread of any process encountered in
//   a system snapshot.
//
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-thread32first
func Thread32First(h uintptr, e *ThreadEntry32) error {
	r, _, err := syscall.SyscallN(funcThread32First.address(), h, uintptr(unsafe.Pointer(e)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// WaitNamedPipe Windows API Call
//   Waits until either a time-out interval elapses or an instance of the
//   specified named pipe is available for connection (that is, the pipe's server
//   process has a pending ConnectNamedPipe operation on the pipe).
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-waitnamedpipea
func WaitNamedPipe(name string, timeout uint32) error {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(funcWaitNamedPipe.address(), uintptr(unsafe.Pointer(n)), uintptr(timeout))
	if r == 0 {
		return unboxError(err1)
	}
	return nil
}

// ConnectNamedPipe Windows API Call
//   Enables a named pipe server process to wait for a client process to connect
//   to an instance of a named pipe. A client process connects by calling either
//   the CreateFile or CallNamedPipe function.
//
// https://docs.microsoft.com/en-us/windows/win32/api/namedpipeapi/nf-namedpipeapi-connectnamedpipe
func ConnectNamedPipe(h uintptr, o *Overlapped) error {
	r, _, err := syscall.SyscallN(funcConnectNamedPipe.address(), h, uintptr(unsafe.Pointer(o)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// NtUnmapViewOfSection Windows API Call
//   The NtUnmapViewOfSection routine unmaps a view of a section from the virtual
//   address space of a subject process.
//
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/wdm/nf-wdm-zwunmapviewofsection
func NtUnmapViewOfSection(proc, section uintptr) error {
	if r, _, _ := syscall.SyscallN(funcNtUnmapViewOfSection.address(), proc, section); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// Process32Next Windows API Call
//   Retrieves information about the next process recorded in a system snapshot.
//
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-process32nextw
func Process32Next(h uintptr, e *ProcessEntry32) error {
	r, _, err := syscall.SyscallN(funcProcess32Next.address(), h, uintptr(unsafe.Pointer(e)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// Process32First Windows API Call
//   Retrieves information about the next process recorded in a system snapshot.
//
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-process32next
func Process32First(h uintptr, e *ProcessEntry32) error {
	r, _, err := syscall.SyscallN(funcProcess32First.address(), h, uintptr(unsafe.Pointer(e)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// SetServiceStatus Windows API Call
//   Contains status information for a service. The ControlService, EnumDependentServices,
//   EnumServicesStatus, and QueryServiceStatus functions use this structure. A
//   service uses this structure in the SetServiceStatus function to report its
//   current status to the service control manager.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winsvc/ns-winsvc-service_status
func SetServiceStatus(h uintptr, s *ServiceStatus) error {
	r, _, err := syscall.SyscallN(funcSetServiceStatus.address(), h, uintptr(unsafe.Pointer(s)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// CheckRemoteDebuggerPresent Windows API Call
//   Determines whether the specified process is being debugged.
//
// https://docs.microsoft.com/en-us/windows/win32/api/debugapi/nf-debugapi-checkremotedebuggerpresent
//
// Calls 'NtQueryInformationProcess' under the hood.
func CheckRemoteDebuggerPresent(h uintptr, b *bool) error {
	var (
		p       uintptr
		r, _, _ = syscall.SyscallN(funcNtQueryInformationProcess.address(), h, 0x7, uintptr(unsafe.Pointer(&p)), ptrSize, 0)
		// 0x7 - ProcessDebugPort
	)
	if r > 0 {
		return formatNtError(r)
	}
	*b = p > 0
	return nil
}

// RegDeleteKeyEx Windows API Call
//   Deletes a subkey and its values. Note that key names are not case sensitive.
//   ONLY DELETES EMPTY SUBKEYS. (invalid argument if non-empty)
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regdeletekeyexw
func RegDeleteKeyEx(h uintptr, path string, f uint32) error {
	p, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(funcRegDeleteKeyEx.address(), h, uintptr(unsafe.Pointer(p)), uintptr(f), 0)
	if r > 0 {
		return unboxError(err1)
	}
	return nil
}

// StartServiceCtrlDispatcher Windows API Call
//   Connects the main thread of a service process to the service control manager,
//   which causes the thread to be the service control dispatcher thread for the
//   calling process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winsvc/nf-winsvc-startservicectrldispatcherw
func StartServiceCtrlDispatcher(t *ServiceTableEntry) error {
	r, _, err := syscall.SyscallN(funcStartServiceCtrlDispatcher.address(), uintptr(unsafe.Pointer(t)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// LoadLibraryEx Windows API Call
//   Loads the specified module into the address space of the calling process.
//   The specified module may cause other modules to be loaded.
//
// https://docs.microsoft.com/en-us/windows/win32/api/libloaderapi/nf-libloaderapi-loadlibraryexw
func LoadLibraryEx(s string, flags uintptr) (uintptr, error) {
	n, err := UTF16PtrFromString(s)
	if err != nil {
		return 0, err
	}
	r, _, e := syscall.SyscallN(funcLoadLibraryEx.address(), uintptr(unsafe.Pointer(n)), 0, flags)
	if r == 0 {
		return 0, unboxError(e)
	}
	return r, nil
}

// WinHTTPGetDefaultProxyConfiguration Windows API Call
//   The WinHttpGetDefaultProxyConfiguration function retrieves the default WinHTTP
//   proxy configuration from the registry.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winhttp/nf-winhttp-winhttpgetdefaultproxyconfiguration
func WinHTTPGetDefaultProxyConfiguration(i *ProxyInfo) error {
	r, _, err := syscall.SyscallN(funcWinHTTPGetDefaultProxyConfiguration.address(), uintptr(unsafe.Pointer(&i)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// LookupPrivilegeValue Windows API Call
//   The LookupPrivilegeValue function retrieves the locally unique identifier
//   (LUID) used on a specified system to locally represent the specified privilege
//   name.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-lookupprivilegevaluew
func LookupPrivilegeValue(system, name string, l *LUID) error {
	var (
		s, n *uint16
		err  error
	)
	if len(system) > 0 {
		if s, err = UTF16PtrFromString(system); err != nil {
			return err
		}
	}
	if n, err = UTF16PtrFromString(name); err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(
		funcLookupPrivilegeValue.address(), uintptr(unsafe.Pointer(s)), uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(l)),
	)
	if r == 0 {
		return unboxError(err1)
	}
	return nil
}

// DeleteProcThreadAttributeList Windows API Call
//   Deletes the specified list of attributes for process and thread creation.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-deleteprocthreadattributelist
//
// TODO(dij): Are we sure this function does anything? Ghidra shows that it just
//            returns!
func DeleteProcThreadAttributeList(a *StartupAttributes) error {
	r, _, err := syscall.SyscallN(funcDeleteProcThreadAttributeList.address(), uintptr(unsafe.Pointer(a)))
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// CreateToolhelp32Snapshot Windows API Call
//   Takes a snapshot of the specified processes, as well as the heaps, modules,
//   and threads used by these processes.
//
// https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
func CreateToolhelp32Snapshot(flags, pid uint32) (uintptr, error) {
	r, _, err := syscall.SyscallN(funcCreateToolhelp32Snapshot.address(), uintptr(flags), uintptr(pid))
	if r == invalid {
		return 0, unboxError(err)
	}
	return r, nil
}

// WaitForSingleObject Windows API Call
//   Waits until the specified object is in the signaled state or the time-out
//   interval elapses.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-waitforsingleobject
func WaitForSingleObject(h uintptr, timeout int32) (uint32, error) {
	r, _, err := syscall.SyscallN(funcWaitForSingleObject.address(), h, uintptr(uint32(timeout)))
	// 0xFFFFFFFF - WAIT_FAILED
	if r == 0xFFFFFFFF {
		return 0, unboxError(err)
	}
	return uint32(r), nil
}

// ReadFile Windows API Call
//   Reads data from the specified file or input/output (I/O) device. Reads
//   occur at the position specified by the file pointer if supported by the device.
//
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-readfile
func ReadFile(h uintptr, b []byte, n *uint32, o *Overlapped) error {
	var v *byte
	if len(b) > 0 {
		v = &b[0]
	}
	r, _, err := syscall.SyscallN(
		funcReadFile.address(), h, uintptr(unsafe.Pointer(v)), uintptr(len(b)), uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(o)),
	)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// WriteFile Windows API Call
//   Writes data to the specified file or input/output (I/O) device.
//
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-writefile
func WriteFile(h uintptr, b []byte, n *uint32, o *Overlapped) error {
	var v *byte
	if len(b) > 0 {
		v = &b[0]
	}
	r, _, err := syscall.SyscallN(
		funcWriteFile.address(), h, uintptr(unsafe.Pointer(v)), uintptr(len(b)), uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(o)),
	)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// OpenProcessToken Windows API Call
//   The OpenProcessToken function opens the access token associated with a process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocesstoken
//
// Re-targeted to use 'NtOpenProcessToken' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntopenprocesstoken
func OpenProcessToken(h uintptr, access uint32, res *uintptr) error {
	if r, _, _ := syscall.SyscallN(funcNtOpenProcessToken.address(), h, uintptr(access), uintptr(unsafe.Pointer(res))); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// NtWriteVirtualMemory Windows API Call
//   This function copies the specified address range from the current process
//   into the specified address range of the specified process.
//
// http://www.codewarrior.cn/ntdoc/winnt/mm/NtWriteVirtualMemory.htm
func NtWriteVirtualMemory(h, address uintptr, b []byte) (uint32, error) {
	var (
		s       uint32
		r, _, _ = syscall.SyscallN(
			funcNtWriteVirtualMemory.address(), h, address, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)),
			uintptr(unsafe.Pointer(&s)),
		)
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return s, nil
}

// SecurityDescriptorFromString converts an SDDL string describing a security
// descriptor into a self-relative security descriptor object allocated on the
// Go heap.
func SecurityDescriptorFromString(s string) (*SecurityDescriptor, error) {
	var (
		h   *SecurityDescriptor
		err = securityDescriptorFromString(s, 1, &h, nil)
	)
	if err != nil {
		return nil, err
	}
	c := h.copyRelative()
	localFree(uintptr(unsafe.Pointer(h)))
	return c, nil
}

// QueryServiceDynamicInformation Windows API Call
//   Retrieves dynamic information related to the current service start.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winsvc/nf-winsvc-queryservicedynamicinformation
func QueryServiceDynamicInformation(h uintptr, l uint32) (uint32, error) {
	if funcQueryServiceDynamicInformation.find() != nil {
		return 0, syscall.EINVAL
	}
	var (
		a         *uint32
		r, _, err = syscall.SyscallN(
			funcQueryServiceDynamicInformation.address(), h, uintptr(l), uintptr(unsafe.Pointer(&a)),
		)
	)
	if r == 0 {
		return 0, unboxError(err)
	}
	v := *a
	localFree(uintptr(unsafe.Pointer(&a)))
	return v, nil
}

// OpenThread Windows API Call
//   Opens an existing thread object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openthread
func OpenThread(access uint32, inherit bool, tid uint32) (uintptr, error) {
	var i uint32
	if inherit {
		i = 1
	}
	r, _, err := syscall.SyscallN(funcOpenThread.address(), uintptr(access), uintptr(i), uintptr(tid))
	if r == 0 {
		return 0, unboxError(err)
	}
	return r, nil
}

// OpenMutex Windows API Call
//   Opens an existing named mutex object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-OpenMutexw
func OpenMutex(access uint32, inherit bool, name string) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	var i uint32
	if inherit {
		i = 1
	}
	r, _, err1 := syscall.SyscallN(funcOpenMutex.address(), uintptr(access), uintptr(i), uintptr(unsafe.Pointer(n)))
	if r == 0 {
		return 0, unboxError(err1)
	}
	return r, nil
}

// OpenEvent Windows API Call
//   Opens an existing named event object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-openeventw
func OpenEvent(access uint32, inherit bool, name string) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	var i uint32
	if inherit {
		i = 1
	}
	r, _, err1 := syscall.SyscallN(funcOpenEvent.address(), uintptr(access), uintptr(i), uintptr(unsafe.Pointer(n)))
	if r == 0 {
		return 0, unboxError(err1)
	}
	return r, nil
}

// OpenProcess Windows API Call
//   Opens an existing local process object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
//
// Re-targeted to use 'NtOpenProcess' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntddk/nf-ntddk-ntopenprocess
func OpenProcess(access uint32, inherit bool, pid uint32) (uintptr, error) {
	var (
		o objAttrs
		h uintptr
		i clientID
	)
	if i.Process = uintptr(pid); inherit {
		// 0x2 - OBJ_INHERIT
		o.Attributes = 0x2
	}
	o.Length = uint32(unsafe.Sizeof(o))
	r, _, _ := syscall.SyscallN(
		funcNtOpenProcess.address(), uintptr(unsafe.Pointer(&h)), uintptr(access), uintptr(unsafe.Pointer(&o)),
		uintptr(unsafe.Pointer(&i)),
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return h, nil
}

// GetOverlappedResult Windows API Call
//   Retrieves the results of an overlapped operation on the specified file,
//   named pipe, or communications device. To specify a timeout interval or wait
//   on an alertable thread, use GetOverlappedResultEx.
//
// https://docs.microsoft.com/en-us/windows/win32/api/ioapiset/nf-ioapiset-getoverlappedresult
func GetOverlappedResult(h uintptr, o *Overlapped, n *uint32, w bool) error {
	var z uint32
	if w {
		z = 1
	}
	r, _, err := syscall.SyscallN(
		funcGetOverlappedResult.address(), h, uintptr(unsafe.Pointer(o)), uintptr(unsafe.Pointer(n)), uintptr(z),
	)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// OpenThreadToken Windows API Call
//   The OpenThreadToken function opens the access token associated with a thread.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openthreadtoken
//
// Re-targeted to use 'NtOpenThreadToken' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntopenthreadtoken
func OpenThreadToken(h uintptr, access uint32, self bool, t *uintptr) error {
	var s uint32
	if self {
		s = 1
	}
	r, _, _ := syscall.SyscallN(funcNtOpenThreadToken.address(), h, uintptr(access), uintptr(s), uintptr(unsafe.Pointer(t)))
	if r > 0 {
		return formatNtError(r)
	}
	return nil
}

// OpenSemaphore Windows API Call
//   Opens an existing named semaphore object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-OpenSemaphorew
func OpenSemaphore(access uint32, inherit bool, name string) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	var i uint32
	if inherit {
		i = 1
	}
	r, _, err1 := syscall.SyscallN(funcOpenSemaphore.address(), uintptr(access), uintptr(i), uintptr(unsafe.Pointer(n)))
	if r == 0 {
		return 0, unboxError(err1)
	}
	return r, nil
}

// NtAllocateVirtualMemory Windows API Call
//   The NtAllocateVirtualMemory routine reserves, commits, or both, a region of
//   pages within the user-mode virtual address space of a specified process.
//
// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntallocatevirtualmemory
func NtAllocateVirtualMemory(h uintptr, size, access uint32) (uintptr, error) {
	var (
		a       uintptr
		x       = uintptr(size)
		r, _, _ = syscall.SyscallN(
			funcNtAllocateVirtualMemory.address(), h, uintptr(unsafe.Pointer(&a)), 0, uintptr(unsafe.Pointer(&x)),
			0x3000, uintptr(access),
		)
		// 0x300 - MEM_COMMIT | MEM_RESERVE
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return a, nil
}

// NtImpersonateThread Windows API Call
//   This routine is used to cause the server thread to impersonate the client
//   thread.  The impersonation is done according to the specified quality
//   of service parameters.
//
// http://web.archive.org/web/20190822133735/https://www.codewarrior.cn/ntdoc/winnt/ps/NtImpersonateThread.htm
//
// Thanks to: https://www.tiraniddo.dev/2017/08/the-art-of-becoming-trustedinstaller.html
func NtImpersonateThread(h, client uintptr, s *SecurityQualityOfService) error {
	if r, _, _ := syscall.SyscallN(funcNtImpersonateThread.address(), h, client, uintptr(unsafe.Pointer(s))); r > 0 {
		return formatNtError(r)
	}
	return nil
}

// NtCreateThreadEx Windows API Call
//   Creates a thread that runs in the virtual address space of another process
//   and optionally specifies extended attributes such as processor group affinity.
//
// http://pinvoke.net/default.aspx/ntdll/NtCreateThreadEx.html
func NtCreateThreadEx(h, address, args uintptr, suspended bool) (uintptr, error) {
	// TODO(dij): Add additional injection types?
	//            - NtQueueApcThread
	//            - Kernel Table Callback
	f := uint32(0x0004)
	// 0x0004 - THREAD_CREATE_FLAGS_HIDE_FROM_DEBUGGER
	if suspended {
		// 0x0001 - CREATE_SUSPENDED
		f |= 0x0001
	}
	var (
		t       uintptr
		r, _, _ = syscall.SyscallN(
			funcNtCreateThreadEx.address(), uintptr(unsafe.Pointer(&t)), 0x10000000, 0, h, address, args, uintptr(f),
			0, 0, 0, 0,
		)
		// 0x10000000 - THREAD_ALL_ACCESS
		// NOTE(dij): Should we move to this?
		//   ZwCreateThreadEx(
		//      ref IntPtr threadHandle,
		//      AccessMask desiredAccess,
		//      IntPtr objectAttributes,
		//      IntPtr processHandle,
		//      IntPtr startAddress,
		//      IntPtr parameter,
		//      bool inCreateSuspended,
		//      Int32 stackZeroBits,
		//      Int32 sizeOfStack,
		//      Int32 maximumStackSize,
		//      IntPtr attributeList
		//   );
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return t, nil
}

// WaitForMultipleObjects Windows API Call
//   Waits until one or all of the specified objects are in the signaled state or
//   the time-out interval elapses. To enter an alertable wait state, use the
//   WaitForMultipleObjectsEx function.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-waitformultipleobjects
//
// Calls 'WaitForMultipleObjectsEx' under the hood.
func WaitForMultipleObjects(h []uintptr, all bool, timeout int32) (uint32, error) {
	n, a := uint32(len(h)), 0
	if n <= 0 || n > 64 {
		return 0, syscall.EINVAL
	}
	if all {
		a = 1
	}
	r, _, err := syscall.SyscallN(
		funcWaitForMultipleObjects.address(), uintptr(n), uintptr(unsafe.Pointer(&h[0])), uintptr(a), uintptr(uint32(timeout)), 0,
	)
	// 0xFFFFFFFF - WAIT_FAILED
	if r == 0xFFFFFFFF {
		return 0, unboxError(err)
	}
	return uint32(r), nil
}

// CreateMutex Windows API Call
//   Creates or opens a named or unnamed mutex object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-CreateMutexw
func CreateMutex(sa *SecurityAttributes, initial bool, name string) (uintptr, error) {
	var (
		n   *uint16
		err error
	)
	if len(name) > 0 {
		if n, err = UTF16PtrFromString(name); err != nil {
			return 0, err
		}
	}
	var i uint32
	if initial {
		i = 1
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateMutex.address(), uintptr(unsafe.Pointer(sa)), uintptr(i), uintptr(unsafe.Pointer(n)),
	)
	if r == 0 {
		return 0, unboxError(err1)
	}
	if err1 == syscall.ERROR_ALREADY_EXISTS {
		return r, unboxError(err1)
	}
	return r, nil
}

// NtProtectVirtualMemory Windows API Call
//   Changes the protection on a region of committed pages in the virtual address
//   space of a specified process.
//
// http://pinvoke.net/default.aspx/ntdll/NtProtectVirtualMemory.html
func NtProtectVirtualMemory(h, address uintptr, size, access uint32) (uint32, error) {
	var (
		x, v    uint32 = size, 0
		r, _, _        = syscall.SyscallN(
			funcNtProtectVirtualMemory.address(), h, uintptr(unsafe.Pointer(&address)), uintptr(unsafe.Pointer(&x)),
			uintptr(access), uintptr(unsafe.Pointer(&v)),
		)
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return v, nil
}

// LoginUser Windows API Call
//   The LogonUser function attempts to log a user on to the local computer. The
//   local computer is the computer from which LogonUser was called. You cannot
//   use LogonUser to log on to a remote computer. You specify the user with a
//   user name and domain and authenticate the user with a plaintext password.
//   If the function succeeds, you receive a handle to a token that represents
//   the logged-on user. You can then use this token handle to impersonate the
//   specified user or, in most cases, to create a process that runs in the
//   context of the specified user.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-logonuserw
func LoginUser(user, domain, pass string, logintype, provider uint32) (uintptr, error) {
	if len(domain) == 0 {
		domain = "."
	}
	u, err := UTF16PtrFromString(user)
	if err != nil {
		return 0, err
	}
	var p, d *uint16
	if d, err = UTF16PtrFromString(domain); err != nil {
		return 0, err
	}
	if len(pass) > 0 {
		if p, err = UTF16PtrFromString(pass); err != nil {
			return 0, err
		}
	}
	var (
		t          uintptr
		r, _, err1 = syscall.SyscallN(
			funcLogonUser.address(), uintptr(unsafe.Pointer(u)), uintptr(unsafe.Pointer(d)),
			uintptr(unsafe.Pointer(p)), uintptr(logintype), uintptr(provider), uintptr(unsafe.Pointer(&t)),
		)
	)
	if r == 0 {
		return 0, unboxError(err1)
	}
	return t, nil
}

// RegSetValueEx Windows API Call
//   Sets the data and type of a specified value under a registry key.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-RegSetValueExw
func RegSetValueEx(h uintptr, path string, t uint32, data *byte, dataLen uint32) error {
	p, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(
		funcRegSetValueEx.address(), h, uintptr(unsafe.Pointer(p)), 0, uintptr(t), uintptr(unsafe.Pointer(data)),
		uintptr(dataLen),
	)
	if r > 0 {
		return unboxError(err1)
	}
	return nil
}

// CreateEvent Windows API Call
//   Creates or opens a named or unnamed event object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-CreateEventw
func CreateEvent(sa *SecurityAttributes, manual, initial bool, name string) (uintptr, error) {
	var (
		n   *uint16
		err error
	)
	if len(name) > 0 {
		if n, err = UTF16PtrFromString(name); err != nil {
			return 0, err
		}
	}
	var i, m uint32
	if initial {
		i = 1
	}
	if manual {
		m = 1
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateEvent.address(), uintptr(unsafe.Pointer(sa)), uintptr(m), uintptr(i), uintptr(unsafe.Pointer(n)),
	)
	if r == 0 {
		return 0, unboxError(err1)
	}
	if err1 == syscall.ERROR_ALREADY_EXISTS {
		return r, unboxError(err1)
	}
	return r, nil
}
func securityDescriptorFromString(s string, v uint32, i **SecurityDescriptor, n *uint32) error {
	p, err := UTF16PtrFromString(s)
	if err != nil {
		return err
	}
	r, _, err2 := syscall.SyscallN(
		funcConvertStringSecurityDescriptorToSecurityDescriptor.address(), uintptr(unsafe.Pointer(p)), uintptr(v),
		uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(n)),
	)
	if r == 0 {
		return unboxError(err2)
	}
	return nil
}

// RegisterServiceCtrlHandlerEx Windows API Call
//   Registers a function to handle extended service control requests.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winsvc/nf-winsvc-registerservicectrlhandlerexw
func RegisterServiceCtrlHandlerEx(name string, handler uintptr, args uintptr) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, err1 := syscall.SyscallN(funcRegisterServiceCtrlHandlerEx.address(), uintptr(unsafe.Pointer(n)), handler, args)
	if r == 0 {
		return 0, unboxError(err1)
	}
	return r, nil
}

// CreateSemaphore Windows API Call
//   Creates or opens a named or unnamed semaphore object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-CreateSemaphorew
func CreateSemaphore(sa *SecurityAttributes, initial, max uint32, name string) (uintptr, error) {
	var (
		n   *uint16
		err error
	)
	if len(name) > 0 {
		if n, err = UTF16PtrFromString(name); err != nil {
			return 0, err
		}
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateSemaphore.address(), uintptr(unsafe.Pointer(sa)), uintptr(initial), uintptr(max),
		uintptr(unsafe.Pointer(n)),
	)
	if r == 0 {
		return 0, unboxError(err1)
	}
	if err1 == syscall.ERROR_ALREADY_EXISTS {
		return r, unboxError(err1)
	}
	return r, nil
}

// GetTokenInformation Windows API Call
//   The GetTokenInformation function retrieves a specified type of information
//   about an access token. The calling process must have appropriate access
//   rights to obtain the information.
//
// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-gettokeninformation
func GetTokenInformation(t uintptr, class uint32, info *byte, length uint32, ret *uint32) error {
	r, _, err := syscall.SyscallN(
		funcGetTokenInformation.address(), t, uintptr(class), uintptr(unsafe.Pointer(info)), uintptr(length),
		uintptr(unsafe.Pointer(ret)),
	)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// InitiateSystemShutdownEx Windows API Call
//   Initiates a shutdown and optional restart of the specified computer, and
//   optionally records the reason for the shutdown.
//
// NOTE: The caller must have the "SeShutdownPrivilege" privilege enabled. This
//       function does NOT automatically request it.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-initiatesystemshutdownexa
func InitiateSystemShutdownEx(t, msg string, secs uint32, force, reboot bool, reason uint32) error {
	var (
		c, m *uint16
		err  error
	)
	if len(t) > 0 {
		if c, err = UTF16PtrFromString(t); err != nil {
			return err
		}
	}
	if len(msg) > 0 {
		if m, err = UTF16PtrFromString(msg); err != nil {
			return err
		}
	}
	var f, x uint32
	if force {
		f = 1
	}
	if reboot {
		x = 1
	}
	r, _, err1 := syscall.SyscallN(
		funcInitiateSystemShutdownEx.address(), uintptr(unsafe.Pointer(c)), uintptr(unsafe.Pointer(m)),
		uintptr(secs), uintptr(f), uintptr(x), uintptr(reason),
	)
	if r == 0 {
		return unboxError(err1)
	}
	return nil
}

// NtCreateSection Windows API Call
//   The NtCreateSection routine creates a section object.
//
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/wdm/nf-wdm-zwcreatesection
func NtCreateSection(access uint32, size uint64, protect, attrs uint32, file uintptr) (uintptr, error) {
	var (
		x = size
		h uintptr
	)
	r, _, _ := syscall.SyscallN(
		funcNtCreateSection.address(), uintptr(unsafe.Pointer(&h)), uintptr(access), 0, uintptr(unsafe.Pointer(&x)),
		uintptr(protect), uintptr(attrs), file,
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return h, nil
}

// CreateMailslot Windows API Call
//   Creates a mailslot with the specified name and returns a handle that a
//   mailslot server can use to perform operations on the mailslot. The mailslot
//   is local to the computer that creates it. An error occurs if a mailslot
//   with the specified name already exists.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createmailslotw
func CreateMailslot(name string, maxSize uint32, timeout int32, sa *SecurityAttributes) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateMailslot.address(), uintptr(unsafe.Pointer(n)), uintptr(maxSize), uintptr(uint32(timeout)),
		uintptr(unsafe.Pointer(sa)),
	)
	if r == invalid {
		return 0, unboxError(err1)
	}
	if err1 == syscall.ERROR_ALREADY_EXISTS {
		return r, unboxError(err1)
	}
	return r, nil
}

// DuplicateTokenEx Windows API Call
//   The DuplicateTokenEx function creates a new access token that duplicates an
//   existing token. This function can create either a primary token or an
//   impersonation token.
//
// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-duplicatetokenex
//
// Re-targeted to use 'NtDuplicateToken' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntduplicatetoken
func DuplicateTokenEx(h uintptr, access uint32, sa *SecurityAttributes, level, p uint32, new *uintptr) error {
	var (
		o objAttrs
		q SecurityQualityOfService
	)
	if q.ImpersonationLevel = level; sa != nil {
		if o.SecurityDescriptor = sa.SecurityDescriptor; sa.InheritHandle == 1 {
			o.Attributes |= 0x2
		}
	}
	o.Length = uint32(unsafe.Sizeof(o))
	q.Length = uint32(unsafe.Sizeof(q))
	o.SecurityQualityOfService = &q
	r, _, _ := syscall.SyscallN(
		funcNtDuplicateToken.address(), h, uintptr(access), uintptr(unsafe.Pointer(&o)), 0, uintptr(p), uintptr(unsafe.Pointer(new)),
	)
	if r > 0 {
		return formatNtError(r)
	}
	return nil
}

// DuplicateHandle Windows API Call
//   Duplicates an object handle.
//
// https://docs.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-duplicatehandle
//
// Re-targeted to use 'NtDuplicateObject' instead.
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-zwduplicateobject
func DuplicateHandle(srcProc, src, dstProc uintptr, dst *uintptr, access uint32, inherit bool, options uint32) error {
	// NOTE(dij): This is to catch and emulate the actions that normally the
	//            'DuplicateHandle' function will do. This will catch any pusudo
	//            handles for StdErr/StdOut/Stdin and grab the real ones from the
	//            process PEB.
	//
	//            I also think there's a bug in this function as 'DuplicateHandle'
	//            does NOT check the 'srcProc' argument when attempting to resolve
	//            pusudo handles and only resolves them for the current process.
	//
	// Handles constant source https://learn.microsoft.com/en-us/windows/console/getstdhandle
	if v := int64(src); v >= -12 && v <= -10 {
		if p, err := getProcessPeb(); err == nil {
			switch v {
			case -10: // STD_INPUT_HANDLE
				src = p.ProcessParameters.StandardInput
			case -11: // STD_OUTPUT_HANDLE
				src = p.ProcessParameters.StandardOutput
			case -12: // STD_ERROR_HANDLE
				src = p.ProcessParameters.StandardError
			}
		}
	}
	var i uint32
	if inherit {
		// 0x2 - OBJ_INHERIT
		i = 0x2
	}
	r, _, _ := syscall.SyscallN(
		funcNtDuplicateObject.address(), srcProc, src, dstProc, uintptr(unsafe.Pointer(dst)), uintptr(access), uintptr(i),
		uintptr(options),
	)
	if r > 0 {
		return formatNtError(r)
	}
	return nil
}

// NtMapViewOfSection Windows API Call
//   The NtMapViewOfSection routine maps a view of a section into the virtual
//   address space of a subject process.
//
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/wdm/nf-wdm-zwmapviewofsection
func NtMapViewOfSection(section, proc uintptr, offset, size uint64, dis, allocType, protect uint32) (uintptr, error) {
	var (
		a       uintptr
		o, x    = offset, size
		r, _, _ = syscall.SyscallN(
			funcNtMapViewOfSection.address(), section, proc, uintptr(unsafe.Pointer(&a)), 0, 0, uintptr(unsafe.Pointer(&o)),
			uintptr(unsafe.Pointer(&x)), uintptr(dis), uintptr(allocType), uintptr(protect),
		)
	)
	if r > 0 {
		return 0, formatNtError(r)
	}
	return a, nil
}

// RegEnumValue Windows API Call
//   Enumerates the values for the specified open registry key. The function
//   copies one indexed value name and data block for the key each time it is
//   called.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regenumvaluew
func RegEnumValue(h uintptr, index uint32, path *uint16, pathLen, valType *uint32, data *byte, dataLen *uint32) error {
	r, _, err := syscall.SyscallN(
		funcRegEnumValue.address(), h, uintptr(index), uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(pathLen)),
		0, uintptr(unsafe.Pointer(valType)), uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(dataLen)),
	)
	if r > 0 {
		return unboxError(err)
	}
	return nil
}

// CreateNamedPipe Windows API Call
//   Creates an instance of a named pipe and returns a handle for subsequent pipe
//   operations. A named pipe server process uses this function either to create
//   the first instance of a specific named pipe and establish its basic attributes
//   or to create a new instance of an existing named pipe.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createnamedpipea
func CreateNamedPipe(name string, flags, mode, max, out, in, timeout uint32, sa *SecurityAttributes) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateNamedPipe.address(), uintptr(unsafe.Pointer(n)), uintptr(flags), uintptr(mode), uintptr(max),
		uintptr(out), uintptr(in), uintptr(timeout), uintptr(unsafe.Pointer(sa)),
	)
	if r == invalid {
		return 0, unboxError(err1)
	}
	return r, nil
}

// AdjustTokenPrivileges Windows API Call
//   The AdjustTokenPrivileges function enables or disables privileges in the
//   specified access token. Enabling or disabling privileges in an access token
//   requires TOKEN_ADJUST_PRIVILEGES access.
//
// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-adjusttokenprivileges
//
// Re-targeted to use 'NtAdjustPrivilegesToken' instead.
// https://docs.rs/ntapi/0.3.6/aarch64-pc-windows-msvc/ntapi/ntseapi/fn.NtAdjustPrivilegesToken.html
func AdjustTokenPrivileges(h uintptr, disableAll bool, new unsafe.Pointer, newLen uint32, old unsafe.Pointer, oldLen *uint32) error {
	var d uint32
	if disableAll {
		d = 1
	}
	r, _, _ := syscall.SyscallN(
		funcNtAdjustTokenPrivileges.address(), h, uintptr(d), uintptr(new), uintptr(newLen), uintptr(old),
		uintptr(unsafe.Pointer(oldLen)),
	)
	if r > 0 {
		return formatNtError(r)
	}
	return nil
}

// RegCreateKeyEx Windows API Call
//   Creates the specified registry key. If the key already exists, the function
//   opens it. Note that key names are not case sensitive.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regcreatekeyexw
func RegCreateKeyEx(h uintptr, path, class string, options, access uint32, sa *SecurityAttributes, out *uintptr, result *uint32) error {
	var (
		p, c *uint16
		err  error
	)
	if len(class) > 0 {
		if c, err = UTF16PtrFromString(class); err != nil {
			return err
		}
	}
	if p, err = UTF16PtrFromString(path); err != nil {
		return err
	}
	r, _, err1 := syscall.SyscallN(
		funcRegCreateKeyEx.address(), h, uintptr(unsafe.Pointer(p)), 0, uintptr(unsafe.Pointer(c)), uintptr(options),
		uintptr(access), uintptr(unsafe.Pointer(sa)), uintptr(unsafe.Pointer(out)), uintptr(unsafe.Pointer(result)),
	)
	if r > 0 {
		return unboxError(err1)
	}
	return nil
}

// CreateFile Windows API Call
//   Creates or opens a file or I/O device. The most commonly used I/O devices
//   are as follows: file, file stream, directory, physical disk, volume, console
//   buffer, tape drive, communications resource, mailslot, and pipe. The function
//   returns a handle that can be used to access the file or device for various
//   types of I/O depending on the file or device and the flags and attributes
//   specified.
//
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew
func CreateFile(name string, access, mode uint32, sa *SecurityAttributes, disposition, attrs uint32, template uintptr) (uintptr, error) {
	n, err := UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateFile.address(), uintptr(unsafe.Pointer(n)), uintptr(access), uintptr(mode), uintptr(unsafe.Pointer(sa)),
		uintptr(disposition), uintptr(attrs), template,
	)
	if r == invalid {
		return 0, unboxError(err1)
	}
	return r, nil
}

// UpdateProcThreadAttribute Windows API Call
//   Updates the specified attribute in a list of attributes for process and
//   thread creation.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribute
func UpdateProcThreadAttribute(a *StartupAttributes, attr uintptr, val unsafe.Pointer, valLen uint64, old *StartupAttributes, oldLen *uint64) error {
	r, _, err := syscall.SyscallN(
		funcUpdateProcThreadAttribute.address(), uintptr(unsafe.Pointer(a)), 0, attr, uintptr(val), uintptr(valLen),
		uintptr(unsafe.Pointer(old)), uintptr(unsafe.Pointer(oldLen)),
	)
	if r == 0 {
		return unboxError(err)
	}
	return nil
}

// CreateProcessWithToken Windows API Call
//   Creates a new process and its primary thread. The new process runs in the
//   security context of the specified token. It can optionally load the user
//   profile for the specified user.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createprocesswithtokenw
func CreateProcessWithToken(t uintptr, loginFlags uint32, name, cmd string, flags uint32, env []string, dir string, y *StartupInfo, x *StartupInfoEx, i *ProcessInformation) error {
	var (
		n, c, d, e *uint16
		err        error
	)
	if len(name) > 0 {
		if n, err = UTF16PtrFromString(name); err != nil {
			return err
		}
	}
	if len(cmd) > 0 {
		if c, err = UTF16PtrFromString(cmd); err != nil {
			return err
		}
	}
	if len(dir) > 0 {
		if d, err = UTF16PtrFromString(dir); err != nil {
			return err
		}
	}
	if len(env) > 0 {
		if e, err = StringListToUTF16Block(env); err != nil {
			return err
		}
		// 0x400 - CREATE_UNICODE_ENVIRONMENT
		flags |= 0x400
	}
	var j unsafe.Pointer
	if y == nil && x != nil {
		// BUG(dij): For some reason adding this flag causes the function
		//           to return "invalid parameter", even this this IS THE ACCEPTED
		//           thing to do???!
		//
		// flags |= 0x80000
		j = unsafe.Pointer(x)
	} else {
		j = unsafe.Pointer(y)
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateProcessWithToken.address(), t, uintptr(loginFlags), uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(c)),
		uintptr(flags), uintptr(unsafe.Pointer(e)), uintptr(unsafe.Pointer(d)), uintptr(j), uintptr(unsafe.Pointer(i)),
	)
	if r == 0 {
		return unboxError(err1)
	}
	return nil
}

// CreateProcess Windows API Call
//   Creates a new process and its primary thread. The new process runs in the
//   security context of the calling process.
//
// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createprocessw
func CreateProcess(name, cmd string, procSa, threadSa *SecurityAttributes, inherit bool, flags uint32, env []string, dir string, y *StartupInfo, x *StartupInfoEx, i *ProcessInformation) error {
	var (
		z          uint32
		n, c, d, e *uint16
		err        error
	)
	if inherit {
		z = 1
	}
	if len(name) > 0 {
		if n, err = UTF16PtrFromString(name); err != nil {
			return err
		}
	}
	if len(cmd) > 0 {
		if c, err = UTF16PtrFromString(cmd); err != nil {
			return err
		}
	}
	if len(dir) > 0 {
		if d, err = UTF16PtrFromString(dir); err != nil {
			return err
		}
	}
	if len(env) > 0 {
		if e, err = StringListToUTF16Block(env); err != nil {
			return err
		}
		// 0x400 - CREATE_UNICODE_ENVIRONMENT
		flags |= 0x400
	}
	var j unsafe.Pointer
	if y == nil && x != nil {
		flags |= 0x80000
		j = unsafe.Pointer(x)
	} else {
		j = unsafe.Pointer(y)
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateProcess.address(), uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(c)), uintptr(unsafe.Pointer(procSa)),
		uintptr(unsafe.Pointer(threadSa)), uintptr(z), uintptr(flags), uintptr(unsafe.Pointer(e)), uintptr(unsafe.Pointer(d)),
		uintptr(j), uintptr(unsafe.Pointer(i)),
	)
	if r == 0 {
		return unboxError(err1)
	}
	return nil
}

// CreateProcessWithLogin Windows API Call
//   Creates a new process and its primary thread. Then the new process runs the
//   specified executable file in the security context of the specified credentials
//   (user, domain, and password). It can optionally load the user profile for a
//   specified user.
//
// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createprocesswithlogonw
func CreateProcessWithLogin(user, domain, pass string, loginFlags uint32, name, cmd string, flags uint32, env []string, dir string, y *StartupInfo, x *StartupInfoEx, i *ProcessInformation) error {
	var n, c, d, e, p, do *uint16
	u, err := UTF16PtrFromString(user)
	if err != nil {
		return err
	}
	if len(domain) > 0 {
		if do, err = UTF16PtrFromString(domain); err != nil {
			return err
		}
	}
	if len(pass) > 0 {
		if p, err = UTF16PtrFromString(pass); err != nil {
			return err
		}
	}
	if len(name) > 0 {
		if n, err = UTF16PtrFromString(name); err != nil {
			return err
		}
	}
	if len(cmd) > 0 {
		if c, err = UTF16PtrFromString(cmd); err != nil {
			return err
		}
	}
	if len(dir) > 0 {
		if d, err = UTF16PtrFromString(dir); err != nil {
			return err
		}
	}
	if len(env) > 0 {
		if e, err = StringListToUTF16Block(env); err != nil {
			return err
		}
		// 0x400 - CREATE_UNICODE_ENVIRONMENT
		flags |= 0x400
	}
	var j unsafe.Pointer
	if y == nil && x != nil {
		// NOTE(dij): For some reason adding this flag causes the function
		//            to return "invalid parameter", even this IS THE ACCEPTED
		//            thing to do???!
		//
		// flags |= 0x80000
		j = unsafe.Pointer(x)
	} else {
		j = unsafe.Pointer(y)
	}
	r, _, err1 := syscall.SyscallN(
		funcCreateProcessWithLogon.address(), uintptr(unsafe.Pointer(u)), uintptr(unsafe.Pointer(do)), uintptr(unsafe.Pointer(p)),
		uintptr(loginFlags), uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(c)), uintptr(flags), uintptr(unsafe.Pointer(e)),
		uintptr(unsafe.Pointer(d)), uintptr(j), uintptr(unsafe.Pointer(i)),
	)
	if r == 0 {
		return unboxError(err1)
	}
	return nil
}
