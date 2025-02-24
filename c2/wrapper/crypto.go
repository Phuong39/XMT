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

package wrapper

import (
	"crypto/cipher"
	"io"

	"github.com/iDigitalFlame/xmt/c2"
	"github.com/iDigitalFlame/xmt/data/crypto"
	"github.com/iDigitalFlame/xmt/util/xerr"
)

type cbk [5]byte
type block struct {
	cipher.Block
	v []byte
}
type stream struct {
	_ [0]func()
	w crypto.Writer
	r crypto.Reader
}

// NewXOR is a function that is an alias for 'Stream(crypto.XOR(k), crypto.XOR(k))'
//
// This wil return a Stream-backed XOR Wrapper.
func NewXOR(k []byte) c2.Wrapper {
	return Stream(crypto.XOR(k), crypto.XOR(k))
}

// NewCBK creates a special type of Wrapper for CBK-based encryptors.
//
// NOTE: This function will prevent CBK from using its index based block
// functions, not sure if there's a way to work around this.
func NewCBK(a, b, c, d, size byte) c2.Wrapper {
	var e cbk
	e[0], e[1], e[2], e[3], e[4] = a, b, c, d, size
	return e
}
func (c cbk) Unwrap(r io.Reader) (io.Reader, error) {
	e, err := crypto.NewCBKSource(c[0], c[1], c[2], c[3], c[4])
	if err != nil {
		return nil, err
	}
	return crypto.NewReader(e, r), nil
}
func (b *block) Unwrap(r io.Reader) (io.Reader, error) {
	return crypto.Decrypter(b.Block, b.v, r)
}
func (s *stream) Unwrap(r io.Reader) (io.Reader, error) {
	return crypto.NewReader(s.r, r), nil
}

// Block returns a Wrapper based on a Block Cipher, such as AES.
func Block(b cipher.Block, v []byte) (c2.Wrapper, error) {
	if b == nil || len(v) == 0 {
		return nil, xerr.Sub("arguments cannot be nil or empty", 0x6E)
	}
	return &block{v: v, Block: b}, nil
}

// Stream returns a Wrapper based on the crypto.Writer and crypto.Reader interfaces,
// such as XOR.
func Stream(r crypto.Reader, w crypto.Writer) c2.Wrapper {
	return &stream{r: r, w: w}
}
func (c cbk) Wrap(w io.WriteCloser) (io.WriteCloser, error) {
	e, err := crypto.NewCBKSource(c[0], c[1], c[2], c[3], c[4])
	if err != nil {
		return nil, err
	}
	return crypto.NewWriter(e, w), nil
}
func (b *block) Wrap(w io.WriteCloser) (io.WriteCloser, error) {
	return crypto.Encrypter(b.Block, b.v, w)
}
func (s *stream) Wrap(w io.WriteCloser) (io.WriteCloser, error) {
	return crypto.NewWriter(s.w, w), nil
}
