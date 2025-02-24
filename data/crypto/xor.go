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

package crypto

import (
	"io"
	"sync"

	"github.com/iDigitalFlame/xmt/data/crypto/subtle"
	"github.com/iDigitalFlame/xmt/util/bugtrack"
)

const bufMax = 2 << 14 // Should cover the default chunk.Write buffer size

var bufs = sync.Pool{
	New: func() any {
		b := make([]byte, 512, bufMax)
		return &b
	},
}

// XOR is an alias for a byte array that acts as the XOR
// key data buffer.
type XOR []byte

// BlockSize returns the cipher's block size.
func (x XOR) BlockSize() int {
	return len(x)
}

// Operate preforms the XOR operation on the specified byte
// array using the cipher as the key.
func (x XOR) Operate(b []byte) {
	if len(x) == 0 {
		return
	}
	subtle.XorOp(b, x)
}

// Flush satisfies the crypto.Writer interface.
func (XOR) Flush(w io.Writer) error {
	if f, ok := w.(flusher); ok {
		return f.Flush()
	}
	return nil
}

// Decrypt preforms the XOR operation on the specified byte array using the cipher
// as the key.
func (x XOR) Decrypt(dst, src []byte) {
	subtle.XorBytes(dst, x, src)
}

// Encrypt preforms the XOR operation on the specified byte array using the cipher
// as the key.
func (x XOR) Encrypt(dst, src []byte) {
	subtle.XorBytes(dst, x, src)
}

// Read satisfies the crypto.Reader interface.
func (x XOR) Read(r io.Reader, b []byte) (int, error) {
	n, err := io.ReadFull(r, b)
	//        NOTE(dij): ErrUnexpectedEOF happens here on short (< buf size)
	//                   Reads, though is completely normal.
	x.Operate(b[:n])
	return n, err
}

// Write satisfies the crypto.Writer interface.
func (x XOR) Write(w io.Writer, b []byte) (int, error) {
	n := len(b)
	if n > bufMax {
		if bugtrack.Enabled {
			bugtrack.Track("crypto.XOR.Write(): Creating non-heap buffer, n=%d, bufMax=%d", n, bufMax)
		}
		o := make([]byte, n)
		copy(o, b)
		x.Operate(o)
		z, err := w.Write(o)
		// NOTE(dij): Make the GCs job easy
		o = nil
		return z, err
	}
	o := bufs.Get().(*[]byte)
	if len(*o) < n {
		if bugtrack.Enabled {
			bugtrack.Track("crypto.XOR.Write(): Increasing heap buffer size len(*o)=%d, n=%d", len(*o), n)
		}
		*o = append(*o, make([]byte, n-len(*o))...)
	}
	copy(*o, b)
	x.Operate((*o)[:n])
	z, err := w.Write((*o)[:n])
	bufs.Put(o)
	return z, err
}
