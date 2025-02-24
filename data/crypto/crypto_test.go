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
	"bytes"
	"crypto/aes"
	"testing"

	"github.com/iDigitalFlame/xmt/util"
)

func TestXOR(t *testing.T) {
	var (
		x = make(XOR, 32)
		i = make([]byte, 32)
		w bytes.Buffer
	)
	util.Rand.Read(x)
	util.Rand.Read(i)
	z, err := Encrypter(x, i, &w)
	if err != nil {
		t.Fatalf("Encrypter failed with error: %s!", err.Error())
	}
	if _, err = z.Write([]byte("hello there")); err != nil {
		t.Fatalf("Write failed with error: %s!", err.Error())
	}
	if err = z.Close(); err != nil {
		t.Fatalf("Close failed with error: %s!", err.Error())
	}
	r, err := Decrypter(x, i, bytes.NewReader(w.Bytes()))
	if err != nil {
		t.Fatalf("Encrypter failed with error: %s!", err.Error())
	}
	o := make([]byte, 11)
	if _, err := r.Read(o); err != nil {
		t.Fatalf("Read failed with error: %s!", err.Error())
	}
	if string(o) != "hello there" {
		t.Fatalf(`Result output "%s" did not match "hello there"!`, o)
	}
}
func TestCBK(t *testing.T) {
	var (
		c, _ = NewCBKEx(0x90, 128, nil)
		v, _ = NewCBKSource(c.A, c.B, c.C, c.D, 128)
		b    bytes.Buffer
		w    = NewWriter(c, &b)
	)
	if _, err := w.Write([]byte("hello there")); err != nil {
		t.Fatalf("Write failed with error: %s!", err.Error())
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close failed with error: %s!", err.Error())
	}
	var (
		r = NewReader(v, bytes.NewReader(b.Bytes()))
		o = make([]byte, 11)
	)
	if _, err := r.Read(o); err != nil {
		t.Fatalf("Read failed with error: %s!", err.Error())
	}
	if string(o) != "hello there" {
		t.Fatalf(`Result output "%s" did not match "hello there"!`, o)
	}
}
func TestAES(t *testing.T) {
	k := make([]byte, 32)
	util.Rand.Read(k)
	b, err := aes.NewCipher(k)
	if err != nil {
		t.Fatalf("aes.NewCipher failed with error: %s!", err.Error())
	}
	var (
		i = make([]byte, 16)
		w bytes.Buffer
	)
	util.Rand.Read(i)
	z, err := Encrypter(b, i, &w)
	if err != nil {
		t.Fatalf("Encrypter failed with error: %s!", err.Error())
	}
	if _, err = z.Write([]byte("hello there")); err != nil {
		t.Fatalf("Write failed with error: %s!", err.Error())
	}
	if err = z.Close(); err != nil {
		t.Fatalf("Close failed with error: %s!", err.Error())
	}
	r, err := Decrypter(b, i, bytes.NewReader(w.Bytes()))
	if err != nil {
		t.Fatalf("Encrypter failed with error: %s!", err.Error())
	}
	o := make([]byte, 11)
	if _, err := r.Read(o); err != nil {
		t.Fatalf("Read failed with error: %s!", err.Error())
	}
	if string(o) != "hello there" {
		t.Fatalf(`Result output "%s" did not match "hello there"!`, o)
	}
}
