//go:build windows && !crypt && !altload

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

var (
	funcLoadLibraryEx = dllKernelBase.proc("LoadLibraryExW")
	funcFormatMessage = dllKernelBase.proc("FormatMessageW")

	funcNtClose                     = dllNtdll.proc("NtClose")
	funcRtlFreeHeap                 = dllNtdll.proc("RtlFreeHeap")
	funcNtTraceEvent                = dllNtdll.proc("NtTraceEvent")
	funcEtwEventWrite               = dllNtdll.proc("EtwEventWrite")
	funcDbgBreakPoint               = dllNtdll.proc("DbgBreakPoint")
	funcNtOpenProcess               = dllNtdll.proc("NtOpenProcess")
	funcNtResumeThread              = dllNtdll.proc("NtResumeThread")
	funcNtCreateSection             = dllNtdll.proc("NtCreateSection")
	funcNtSuspendThread             = dllNtdll.proc("NtSuspendThread")
	funcNtResumeProcess             = dllNtdll.proc("NtResumeProcess")
	funcRtlAllocateHeap             = dllNtdll.proc("RtlAllocateHeap")
	funcNtDuplicateToken            = dllNtdll.proc("NtDuplicateToken")
	funcEtwEventRegister            = dllNtdll.proc("EtwEventRegister")
	funcNtSuspendProcess            = dllNtdll.proc("NtSuspendProcess")
	funcNtCreateThreadEx            = dllNtdll.proc("NtCreateThreadEx")
	funcNtDuplicateObject           = dllNtdll.proc("NtDuplicateObject")
	funcNtTerminateThread           = dllNtdll.proc("NtTerminateThread")
	funcNtOpenThreadToken           = dllNtdll.proc("NtOpenThreadToken")
	funcEtwEventWriteFull           = dllNtdll.proc("EtwEventWriteFull")
	funcRtlReAllocateHeap           = dllNtdll.proc("RtlReAllocateHeap")
	funcNtMapViewOfSection          = dllNtdll.proc("NtMapViewOfSection")
	funcNtTerminateProcess          = dllNtdll.proc("NtTerminateProcess")
	funcNtOpenProcessToken          = dllNtdll.proc("NtOpenProcessToken")
	funcRtlCopyMappedMemory         = dllNtdll.proc("RtlCopyMappedMemory")
	funcNtFreeVirtualMemory         = dllNtdll.proc("NtFreeVirtualMemory")
	funcNtImpersonateThread         = dllNtdll.proc("NtImpersonateThread")
	funcNtUnmapViewOfSection        = dllNtdll.proc("NtUnmapViewOfSection")
	funcNtWriteVirtualMemory        = dllNtdll.proc("NtWriteVirtualMemory")
	funcNtProtectVirtualMemory      = dllNtdll.proc("NtProtectVirtualMemory")
	funcNtSetInformationThread      = dllNtdll.proc("NtSetInformationThread")
	funcRtlGetNtVersionNumbers      = dllNtdll.proc("RtlGetNtVersionNumbers")
	funcEtwNotificationRegister     = dllNtdll.proc("EtwNotificationRegister")
	funcNtAllocateVirtualMemory     = dllNtdll.proc("NtAllocateVirtualMemory")
	funcRtlSetProcessIsCritical     = dllNtdll.proc("RtlSetProcessIsCritical")
	funcNtFlushInstructionCache     = dllNtdll.proc("NtFlushInstructionCache")
	funcNtAdjustTokenPrivileges     = dllNtdll.proc("NtAdjustPrivilegesToken")
	funcNtQueryInformationThread    = dllNtdll.proc("NtQueryInformationThread")
	funcNtQuerySystemInformation    = dllNtdll.proc("NtQuerySystemInformation")
	funcNtQueryInformationProcess   = dllNtdll.proc("NtQueryInformationProcess")
	funcRtlLengthSecurityDescriptor = dllNtdll.proc("RtlLengthSecurityDescriptor")

	funcSetEvent                      = dllKernelBase.proc("SetEvent")
	funcReadFile                      = dllKernelBase.proc("ReadFile")
	funcWriteFile                     = dllKernelBase.proc("WriteFile")
	funcOpenMutex                     = dllKernelBase.proc("OpenMutexW")
	funcLocalFree                     = dllKernelBase.proc("LocalFree")
	funcOpenEvent                     = dllKernelBase.proc("OpenEventW")
	funcOpenThread                    = dllKernelBase.proc("OpenThread")
	funcHeapCreate                    = dllKernelBase.proc("HeapCreate")
	funcCreateFile                    = dllKernelBase.proc("CreateFileW")
	funcCancelIoEx                    = dllKernelBase.proc("CancelIoEx")
	funcDebugBreak                    = dllKernelBase.proc("DebugBreak")
	funcFreeLibrary                   = dllKernelBase.proc("FreeLibrary")
	funcHeapDestroy                   = dllKernelBase.proc("HeapDestroy")
	funcCreateMutex                   = dllKernelBase.proc("CreateMutexW")
	funcCreateEvent                   = dllKernelBase.proc("CreateEventW")
	funcWaitNamedPipe                 = dllKernelBase.proc("WaitNamedPipeW")
	funcOpenSemaphore                 = dllKernelBase.proc("OpenSemaphoreW")
	funcSetThreadToken                = dllKernelBase.proc("SetThreadToken")
	funcIsWow64Process                = dllKernelBase.proc("IsWow64Process")
	funcIsWellKnownSID                = dllKernelBase.proc("IsWellKnownSid")
	funcCreateNamedPipe               = dllKernelBase.proc("CreateNamedPipeW")
	funcGetLogicalDrives              = dllKernelBase.proc("GetLogicalDrives")
	funcConnectNamedPipe              = dllKernelBase.proc("ConnectNamedPipe")
	funcGetModuleHandleEx             = dllKernelBase.proc("GetModuleHandleExW")
	funcIsDebuggerPresent             = dllKernelBase.proc("IsDebuggerPresent")
	funcCreateWellKnownSid            = dllKernelBase.proc("CreateWellKnownSid")
	funcGetCurrentThreadID            = dllKernelBase.proc("GetCurrentThreadId")
	funcGetOverlappedResult           = dllKernelBase.proc("GetOverlappedResult")
	funcDisconnectNamedPipe           = dllKernelBase.proc("DisconnectNamedPipe")
	funcSetTokenInformation           = dllKernelBase.proc("SetTokenInformation")
	funcGetTokenInformation           = dllKernelBase.proc("GetTokenInformation")
	funcGetCurrentProcessID           = dllKernelBase.proc("GetCurrentProcessId")
	funcWaitForSingleObject           = dllKernelBase.proc("WaitForSingleObjectEx")
	funcWaitForMultipleObjects        = dllKernelBase.proc("WaitForMultipleObjectsEx")
	funcImpersonateLoggedOnUser       = dllKernelBase.proc("ImpersonateLoggedOnUser")
	funcUpdateProcThreadAttribute     = dllKernelBase.proc("UpdateProcThreadAttribute")
	funcImpersonateNamedPipeClient    = dllKernelBase.proc("ImpersonateNamedPipeClient")
	funcDeleteProcThreadAttributeList = dllKernelBase.proc("DeleteProcThreadAttributeList")

	funcLoadLibrary                = dllKernel32.proc("LoadLibraryW")
	funcThread32Next               = dllKernel32.proc("Thread32Next")
	funcCreateProcess              = dllKernel32.proc("CreateProcessW")
	funcProcess32Next              = dllKernel32.proc("Process32NextW")
	funcThread32First              = dllKernel32.proc("Thread32First")
	funcProcess32First             = dllKernel32.proc("Process32FirstW")
	funcCreateMailslot             = dllKernel32.proc("CreateMailslotW")
	funcCreateSemaphore            = dllKernel32.proc("CreateSemaphoreW")
	funcK32EnumDeviceDrivers       = dllKernel32.proc("K32EnumDeviceDrivers")
	funcK32GetModuleInformation    = dllKernel32.proc("K32GetModuleInformation")
	funcCreateToolhelp32Snapshot   = dllKernel32.proc("CreateToolhelp32Snapshot")
	funcK32GetDeviceDriverBaseName = dllKernel32.proc("K32GetDeviceDriverBaseNameW")
	funcSetProcessWorkingSetSizeEx = dllKernel32.proc("SetProcessWorkingSetSizeEx")

	funcLsaClose                                            = dllAdvapi32.proc("LsaClose")
	funcLogonUser                                           = dllAdvapi32.proc("LogonUserW")
	funcRegFlushKey                                         = dllAdvapi32.proc("RegFlushKey")
	funcRegEnumValue                                        = dllAdvapi32.proc("RegEnumValueW")
	funcRegSetValueEx                                       = dllAdvapi32.proc("RegSetValueExW")
	funcLsaOpenPolicy                                       = dllAdvapi32.proc("LsaOpenPolicy")
	funcRegDeleteTree                                       = dllAdvapi32.proc("RegDeleteTreeW")
	funcRegDeleteKeyEx                                      = dllAdvapi32.proc("RegDeleteKeyExW")
	funcRegDeleteValue                                      = dllAdvapi32.proc("RegDeleteValueW")
	funcRegCreateKeyEx                                      = dllAdvapi32.proc("RegCreateKeyExW")
	funcSetServiceStatus                                    = dllAdvapi32.proc("SetServiceStatus")
	funcLookupAccountSid                                    = dllAdvapi32.proc("LookupAccountSidW")
	funcLookupPrivilegeValue                                = dllAdvapi32.proc("LookupPrivilegeValueW")
	funcConvertSIDToStringSID                               = dllAdvapi32.proc("ConvertSidToStringSidW")
	funcCreateProcessWithToken                              = dllAdvapi32.proc("CreateProcessWithTokenW")
	funcCreateProcessWithLogon                              = dllAdvapi32.proc("CreateProcessWithLogonW")
	funcInitiateSystemShutdownEx                            = dllAdvapi32.proc("InitiateSystemShutdownExW")
	funcLsaQueryInformationPolicy                           = dllAdvapi32.proc("LsaQueryInformationPolicy")
	funcStartServiceCtrlDispatcher                          = dllAdvapi32.proc("StartServiceCtrlDispatcherW")
	funcRegisterServiceCtrlHandlerEx                        = dllAdvapi32.proc("RegisterServiceCtrlHandlerExW")
	funcQueryServiceDynamicInformation                      = dllAdvapi32.proc("QueryServiceDynamicInformation")
	funcConvertStringSecurityDescriptorToSecurityDescriptor = dllAdvapi32.proc("ConvertStringSecurityDescriptorToSecurityDescriptorW")

	funcGetDC                      = dllUser32.proc("GetDC")
	funcIsZoomed                   = dllUser32.proc("IsZoomed")
	funcIsIconic                   = dllUser32.proc("IsIconic")
	funcSetFocus                   = dllUser32.proc("SetFocus")
	funcReleaseDC                  = dllUser32.proc("ReleaseDC")
	funcSendInput                  = dllUser32.proc("SendInput")
	funcBlockInput                 = dllUser32.proc("BlockInput")
	funcShowWindow                 = dllUser32.proc("ShowWindow")
	funcMessageBox                 = dllUser32.proc("MessageBoxW")
	funcEnumWindows                = dllUser32.proc("EnumWindows")
	funcEnableWindow               = dllUser32.proc("EnableWindow")
	funcSetWindowPos               = dllUser32.proc("SetWindowPos")
	funcGetWindowText              = dllUser32.proc("GetWindowTextW")
	funcGetWindowInfo              = dllUser32.proc("GetWindowInfo")
	funcGetMonitorInfo             = dllUser32.proc("GetMonitorInfoW")
	funcGetWindowLongW             = dllUser32.proc("GetWindowLongW")
	funcSetWindowLongW             = dllUser32.proc("SetWindowLongW")
	funcIsWindowVisible            = dllUser32.proc("IsWindowVisible")
	funcGetDesktopWindow           = dllUser32.proc("GetDesktopWindow")
	funcSendNotifyMessage          = dllUser32.proc("SendNotifyMessageW")
	funcEnumDisplayMonitors        = dllUser32.proc("EnumDisplayMonitors")
	funcEnumDisplaySettings        = dllUser32.proc("EnumDisplaySettingsW")
	funcGetWindowTextLength        = dllUser32.proc("GetWindowTextLengthW")
	funcSetForegroundWindow        = dllUser32.proc("SetForegroundWindow")
	funcSystemParametersInfo       = dllUser32.proc("SystemParametersInfoW")
	funcSetLayeredWindowAttributes = dllUser32.proc("SetLayeredWindowAttributes")

	funcBitBlt                 = dllGdi32.proc("BitBlt")
	funcDeleteDC               = dllGdi32.proc("DeleteDC")
	funcGetDIBits              = dllGdi32.proc("GetDIBits")
	funcSelectObject           = dllGdi32.proc("SelectObject")
	funcDeleteObject           = dllGdi32.proc("DeleteObject")
	funcCreateCompatibleDC     = dllGdi32.proc("CreateCompatibleDC")
	funcCreateCompatibleBitmap = dllGdi32.proc("CreateCompatibleBitmap")

	funcWTSFreeMemory              = dllWtsapi32.proc("WTSFreeMemory")
	funcWTSOpenServer              = dllWtsapi32.proc("WTSOpenServerW")
	funcWTSCloseServer             = dllWtsapi32.proc("WTSCloseServer")
	funcWTSSendMessage             = dllWtsapi32.proc("WTSSendMessageW")
	funcWTSLogoffSession           = dllWtsapi32.proc("WTSLogoffSession")
	funcWTSEnumerateSessions       = dllWtsapi32.proc("WTSEnumerateSessionsW")
	funcWTSDisconnectSession       = dllWtsapi32.proc("WTSDisconnectSession")
	funcWTSEnumerateProcesses      = dllWtsapi32.proc("WTSEnumerateProcessesW")
	funcWTSQuerySessionInformation = dllWtsapi32.proc("WTSQuerySessionInformationW")

	funcMiniDumpWriteDump = dllDbgHelp.proc("MiniDumpWriteDump")

	funcWinHTTPGetDefaultProxyConfiguration = dllWinhttp.proc("WinHttpGetDefaultProxyConfiguration")

	funcAmsiScanBuffer = dllAmsi.proc("AmsiScanBuffer")
	funcAmsiInitialize = dllAmsi.proc("AmsiInitialize")
	funcAmsiScanString = dllAmsi.proc("AmsiScanString")
)

func doSearchSystem32() bool {
	searchSystem32.Do(func() {
		searchSystem32.v = dllKernel32.proc("AddDllDirectory").find() == nil
	})
	return searchSystem32.v
}

//funcRevertToSelf                                        = dllAdvapi32.proc("RevertToSelf")
