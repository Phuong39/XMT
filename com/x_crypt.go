//go:build crypt

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

package com

import "github.com/iDigitalFlame/xmt/util/crypt"

// Named Network Constants
var (
	NameIP   = crypt.Get(14) // ip
	NameTCP  = crypt.Get(15) // tcp
	NameUDP  = crypt.Get(16) // udp
	NameUnix = crypt.Get(17) // unix
	NamePipe = crypt.Get(18) // pipe
	NameHTTP = crypt.Get(19) // http
)
