//go:build !windows && crypt

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
	"os"
	"sort"
	"strconv"

	"github.com/iDigitalFlame/xmt/util/crypt"
)

func Processes() ([]ProcessInfo, error) {
	l, err := os.ReadDir(crypt.Get(12)) // /proc/
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, nil
	}
	var (
		n string
		b []byte
		v uint64
		p uint32
		r = make(processList, 0, len(l)/2)
	)
	for i := range l {
		if !l[i].IsDir() {
			continue
		}
		if n = l[i].Name(); len(n) < 1 || n[0] < 48 || n[0] > 57 {
			continue
		}
		if v, err = strconv.ParseUint(n, 10, 32); err != nil {
			continue
		}
		b, err = os.ReadFile(
			crypt.Get(12) + // /proc/
				n +
				crypt.Get(13), // /status
		)
		if err != nil {
			continue
		}
		u := getProcUser(crypt.Get(12) + n) // /proc/
		if n, p = readProcStats(b); len(n) == 0 {
			continue
		}
		r = append(r, ProcessInfo{Name: n, User: u, PID: uint32(v), PPID: p})
	}
	sort.Sort(r)
	return r, nil
}
