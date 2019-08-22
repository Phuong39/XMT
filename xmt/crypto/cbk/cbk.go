package cbk

import (
	"crypto/rand"
	"io"
	"math"
	"sync"

	"golang.org/x/xerrors"
)

const (
	// BlockSize is the default block buffer size of this Cipher.
	BlockSize = 16

	// BlockSizeMax is the maximum block buffer size of this Cipher
	BlockSizeMax = 128
)

var (
	// ErrSize is returned when an array is read that does not contain enough slots for keys
	// which is three.
	ErrSize = xerrors.New("byte array size must be greather than or equal to three (3)")

	// ErrBlockSize is an error returned when an invalid value for the block size is given
	// when creating the Cipher.
	ErrBlockSize = xerrors.New("block size must be between 16 and 128 and a power of two")

	bufs = &sync.Pool{
		New: func() interface{} {
			return make([]byte, BlockSize+1)
		},
	}
	tables = &sync.Pool{
		New: func() interface{} {
			b := make([][]byte, BlockSize+1)
			for i := 0; i < len(b); i++ {
				b[i] = make([]byte, 256)
			}
			return b
		},
	}
)

// Cipher is the repersentation of the CBK Cipher.
// CBK is a block based cipher that allows for a variable size index in encoding.
type Cipher struct {
	A       byte
	B       byte
	C       byte
	D       byte
	Entropy source

	buf   []byte
	pos   int
	index uint8
	total int

	compat bool
}
type source interface {
	Reset() error
	Next(uint16) uint16
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	bufs.Put(b)
}

// NewCipher returns a new CBK Cipher with the D value specified. The other A, B and C values
// are randomally generated at runtime.
func NewCipher(d int) *Cipher {
	c := &Cipher{
		D:   byte(d),
		buf: make([]byte, BlockSize+1),
	}
	c.Reset()
	return c
}

// Reset resets the encryption keys and sets them to new random bytes.
func (e *Cipher) Reset() error {
	b := bufs.Get().([]byte)
	defer clear(b)
	_, err := rand.Read(b[0:3])
	if err != nil {
		return xerrors.Errorf("unabled to get random values: %w", err)
	}
	e.A = b[0]
	e.B = b[1]
	e.C = b[2]
	return nil
}

// BlockSize returns the cipher's block BlockSize.
func (e *Cipher) BlockSize() int {
	if e.buf == nil {
		return BlockSize
	}
	return len(e.buf) - 1
}

// Shuffle will switch around the bytes in the array based on the Cipher bytes.
func (e *Cipher) Shuffle(b []byte) {
	if e == nil {
		return
	}
	if len(b) > 1 {
		b[0] += e.A
	}
	for i := byte(0); i < byte(len(b)); i++ {
		switch {
		case i%e.A == 0:
			b[i] += (e.A - i)
		case e.C%i == 0:
			b[i] += (e.B - e.D)
		case i == e.D:
			b[i] -= (e.A + i)
		default:
			if i%2 == 0 {
				b[i] += e.B / 3
			} else {
				b[i] += e.C / 5
			}
		}
	}
}

// Deshuffle will reverse the switch around the bytes in the array based on the Cipher bytes.
func (e *Cipher) Deshuffle(b []byte) {
	if e == nil {
		return
	}
	if len(b) > 1 {
		b[0] -= e.A
	}
	for i := byte(0); i < byte(len(b)); i++ {
		switch {
		case i%e.A == 0:
			b[i] -= (e.A - i)
		case e.C%i == 0:
			b[i] -= (e.B - e.D)
		case i == e.D:
			b[i] += (e.A + i)
		default:
			if i%2 == 0 {
				b[i] -= e.B / 3
			} else {
				b[i] -= e.C / 5
			}
		}
	}
}
func (e *Cipher) cipherTable(b []byte) {
	b[0] = byte(uint16(e.index+1)*uint16(e.D+1) + e.adjust(uint16(e.D)))
	for i := byte(1); i < byte(len(b))-1; i++ {
		switch {
		case i <= 6:
			if i%2 == 0 {
				b[i] = byte(uint16(e.index) - uint16(e.A) + uint16(e.B-(i-e.C)) + uint16(i) - e.adjust(uint16(e.A)))
			} else {
				b[i] = byte(uint16(e.index) - uint16(e.A) + uint16(e.B-(i-3)) + uint16(i) - e.adjust(uint16(e.A)))
			}
		case i > 6 && i <= 11:
			b[i] = byte(uint16(e.C) - uint16(e.B) + uint16((e.index+1)*i) + e.adjust(uint16(e.C)))
		case i > 11:
			b[i] = byte(e.adjust(uint16(e.B+e.C)) + uint16(e.D) - uint16(len(b)-1) - uint16(e.D) + uint16(e.A-e.C))
		}
	}
	b[len(b)-1] = byte(e.adjust(uint16(e.B+e.C)) + uint16(e.index) - uint16(len(b)-1) - uint16(e.D) + uint16(e.A-e.C))
}
func switchHalf(r bool, a, b byte) byte {
	if r {
		return byte((((b & 0xF) & 0xF) << 4) | ((a >> 4) & 0xF))
	}
	return byte(((((b >> 4) & 0xF) & 0xF) << 4) | (a & 0xF))
}
func (e *Cipher) adjust(i uint16) uint16 {
	if e.Entropy != nil {
		return e.Entropy.Next(i)
	}
	return uint16(math.Max(float64((e.A^e.B)-e.C), 1))
}

// Encrypt encrypts the first block in src into dst.
// Dst and src must overlap entirely or not at all.
func (e *Cipher) Encrypt(dst, src []byte) {
	copy(dst, src)
	e.Shuffle(dst)
}

// Decrypt decrypts the first block in src into dst.
// Dst and src must overlap entirely or not at all.
func (e *Cipher) Decrypt(dst, src []byte) {
	copy(dst, src)
	e.Deshuffle(dst)
}
func (e *Cipher) readIndex(b []byte) error {
	if b == nil || len(b) < 3 {
		return ErrSize
	}
	a := byte(uint16(b[0]) - uint16(e.D) - (uint16(e.D) / 2) - (uint16(e.D) - (2 + e.adjust(uint16(e.D)))))
	c := byte(uint16(b[1]) - e.adjust(uint16(e.D)) - (uint16(e.D) / 3) - ((1 + e.adjust(uint16(e.D))) * (uint16(e.D) + 1) * 2) - (uint16(e.D) + 5))
	if e.D%2 == 0 {
		e.A = switchHalf(true, c, a)
		e.B = switchHalf(true, a, c)
	} else {
		e.A = switchHalf(false, a, c)
		e.B = switchHalf(false, c, a)
	}
	e.C = byte(uint16(b[2]) - (uint16(e.D) / 5) - 7 - ((uint16(e.D)+1)*3 + e.adjust(uint16(e.D)) + 8 + uint16(e.D)))
	return nil
}
func (e *Cipher) writeIndex(b []byte) error {
	if b == nil || len(b) < 3 {
		return ErrSize
	}
	b[0] = byte(uint16(switchHalf(e.D%2 == 0, e.A, e.B)) + uint16(e.D) + (uint16(e.D) / 2) + (uint16(e.D) - (2 + e.adjust(uint16(e.D)))))
	b[1] = byte(uint16(switchHalf(e.D%2 == 0, e.B, e.A)) + e.adjust(uint16(e.D)) + (uint16(e.D) / 3) + ((1 + e.adjust(uint16(e.D))) * (uint16(e.D) + 1) * 2) + (uint16(e.D) + 5))
	b[2] = byte(uint16(e.C) + (uint16(e.D) / 5) + 7 + ((uint16(e.D)+1)*3 + e.adjust(uint16(e.D)) + 8 + uint16(e.D)))
	return nil
}
func (e *Cipher) scramble(b []byte, t, i byte) {
	o := bufs.Get().([]byte)
	defer clear(o)
	for x := 0; x < 6; x++ {
		a := byte(math.Abs(float64(e.blockIndex(true, uint16(t+byte(x)), uint16(t+byte(x))) % 8)))
		c := byte(math.Abs(float64(e.blockIndex(false, uint16(t+byte(x)), uint16(t+byte(x))) % 8)))
		if a != c {
			copy(o[0:2], b[a*2:a*2+2])
			copy(b[a*2:], b[c*2:c*2+2])
			copy(b[c*2:], o[0:2])
		}
	}
}
func (e *Cipher) descramble(b []byte, t, i byte) {
	o := bufs.Get().([]byte)
	defer clear(o)
	for x := 5; x >= 0; x-- {
		a := byte(math.Abs(float64(e.blockIndex(true, uint16(t+byte(x)), uint16(t+byte(x))) % 8)))
		c := byte(math.Abs(float64(e.blockIndex(false, uint16(t+byte(x)), uint16(t+byte(x))) % 8)))
		if a != c {
			copy(o[0:2], b[a*2:a*2+2])
			copy(b[a*2:], b[c*2:c*2+2])
			copy(b[c*2:], o[0:2])
		}
	}
}
func (e *Cipher) syncInput(r io.Reader) (int, error) {
	n, err := r.Read(e.buf)
	if err != nil {
		e.total = 0
		return 0, err
	}
	if n <= 0 {
		e.total = 0
		return 0, io.EOF
	}
	e.index++
	if e.index > 30 {
		e.index = 0
	}
	t := bufs.Get().([]byte)
	c := tables.Get().([][]byte)
	defer clear(t)
	defer tables.Put(c)
	e.cipherTable(t)
	e.descramble(e.buf, e.D, e.index)
	for x := 0; x < len(c); x++ {
		for z := 0; z < len(c[x]); z++ {
			c[x][t[x]&0xFF] = byte(z)
			t[x]++
		}
	}
	for i := 0; i < len(e.buf); i++ {
		e.buf[i] = c[i][e.buf[i]&0xFF]
	}
	if !e.compat {
		e.Deshuffle(e.buf)
	}
	e.total = int(e.buf[len(e.buf)-1])
	e.pos = 0
	if e.total == 0 {
		return 0, io.EOF
	}
	return n, nil
}
func (e *Cipher) syncOutput(w io.Writer) (int, error) {
	e.index++
	if e.index > 30 {
		e.index = 0
	}
	e.total += e.pos
	e.buf[len(e.buf)-1] = byte(e.pos)
	t := bufs.Get().([]byte)
	c := tables.Get().([][]byte)
	defer clear(t)
	defer tables.Put(c)
	e.cipherTable(t)
	for x := 0; x < len(c); x++ {
		for z := 0; z < len(c[x]); z++ {
			c[x][z] = t[x]
			t[x]++
		}
	}
	if !e.compat {
		e.Shuffle(e.buf)
	}
	for i := 0; i < len(e.buf); i++ {
		e.buf[i] = c[i][e.buf[i]&0xFF]
	}
	e.scramble(e.buf, e.D, e.index)
	n, err := w.Write(e.buf)
	if err != nil {
		return 0, err
	}
	e.pos = 0
	return n, nil
}
func (e *Cipher) blockIndex(a bool, t, i uint16) byte {
	switch v := t % 8; {
	case v == 0 && a:
		return byte((((t+1)*(1+i+e.adjust(t)) + t + 5) / 3) + 4 + (5 * t) + (i / 5))
	case v == 1 && a:
		return byte((t / 5) + i + ((i + 1) * 7) + ((1 + t) * 3) + (i / 2) + t)
	case v == 2 && a:
		return byte((((3+t+e.adjust(i))/4+1)+i)/2 + (3 * t) + (t / 5) + i + 3)
	case v == 3 && a:
		return byte(((t / 2) * 3) + 7 + ((t + i) * 3) - 2 + ((t * (i + 5 + e.adjust(i))) * 3))
	case v == 4 && a:
		return byte((((i*6)+2)/5)*3 + ((4 * i) / 5) + 3 + (t / 4))
	case v == 5 && a:
		return byte((((t*3)/5)+(5+i))*3 + (t * (2 - e.adjust(i))) + (i / (t + 1)) + (6 + t))
	case v == 6 && a:
		return byte((((((i + 5) / 3) * 7) + 3 + e.adjust(t)) / (t + 1)) + 3 + (t/(i+1))*3)
	case v == 7 && a:
		return byte(((((t / (i + 1) * 2) + 5) / 4) + 10) + (3 * t) + ((i / 2) + (t * 3)) + 4)
	case v == 0 && !a:
		return byte((((3/(2+i) + 3) / (t + 1)) * 9) + 6 - e.adjust(i))
	case v == 1 && !a:
		return byte(((((4*i)/3 + (t * 2)) / 3) + 8) / 3)
	case v == 2 && !a:
		return byte(((((9 + i + e.adjust(i)) / 4) + (t / 2) + (2*i + 1 + e.adjust(t))) / (((i + 3) / (5 + t)) + 6)))
	case v == 3 && !a:
		return byte(((((4+(t-5)/2)/6)+3)*2)*((5+i)/3) + 4)
	case v == 4 && !a:
		return byte((((((t/3)/(3+i) + e.adjust(t)) / 9) * 2) + 8) + (5+i)/(3+t))
	case v == 5 && !a:
		return byte(((i * 4) + (t / 3) - e.adjust(t) + (6 / (1 + i))) + (6 / (3 + t)) + (i * 3))
	case v == 6 && !a:
		return byte((((((t*9)/6)+(i*3)/9)*5 + i) - e.adjust((t + i))) + (t+2)/4)
	case v == 7 && !a:
		return byte((((((i/3)*7)+3-e.adjust(t))*5 + t) * (t + 3) / 7) + e.adjust(i))
	}
	return 0
}
func (e *Cipher) read(r io.Reader, b []byte) (int, error) {
	if e.buf == nil {
		e.buf = make([]byte, BlockSize+1)
	}
	if e.total-e.pos > len(b) {
		n := copy(b, e.buf[e.pos:e.pos+len(b)])
		e.pos += len(b)
		return n, nil
	}
	if e.pos >= e.total {
		if o, err := e.syncInput(r); err != nil && err != io.EOF || o == 0 {
			return o, err
		}
	}
	n := 0
	for i := 0; n < len(b); n += i {
		if e.total <= 0 {
			return n, io.EOF
		}
		i = copy(b[n:], e.buf[e.pos:e.total])
		e.pos += i
		if e.pos >= e.total {
			if _, err := e.syncInput(r); err != nil && err != io.EOF {
				return n, err
			}
		}
	}
	return n, nil
}
func (e *Cipher) write(w io.Writer, b []byte) (int, error) {
	if e.buf == nil {
		e.buf = make([]byte, BlockSize+1)
	}
	n := 0
	for i := len(e.buf) - 1; n < len(b) && i == len(e.buf)-1; n += i {
		if e.pos >= len(e.buf)-1 {
			if o, err := e.syncOutput(w); err != nil {
				return o, err
			}
		}
		i = copy(e.buf[e.pos:len(e.buf)-1], b[n:])
		e.pos += i
	}
	if e.pos >= len(e.buf)-1 {
		if o, err := e.syncOutput(w); err != nil {
			return o, err
		}
	}
	return n, nil
}

// NewCipherEx returns a new CBK Cipher with the D value, BlockSize and Entropy source specified. The other A, B and C values
// are randomally generated at runtime.
func NewCipherEx(d int, size int, ent source) (*Cipher, error) {
	if size < BlockSize || size > BlockSizeMax || math.Floor(math.Log2(float64(size))) != math.Ceil(math.Log2(float64(size))) {
		return nil, ErrBlockSize
	}
	c := &Cipher{
		D:       byte(d),
		buf:     make([]byte, size+1),
		Entropy: ent,
	}
	c.Reset()
	return c, nil
}
