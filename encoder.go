// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package dbc

import (
	"fmt"
	"io"
)

type bw struct {
	d byte
	i uint
	w io.ByteWriter
}

func (b *bw) WriteBool(bit bool) error {
	var err error
	if bit {
		b.d |= 1 << b.i
	}
	b.i++
	if b.i == 8 {
		b.i = 0
		err = b.w.WriteByte(b.d)
		b.d = 0
	}
	return err
}

type Encoder struct {
	n         uint64
	w         bw
	p         uint64
	low, high uint64
	writes    uint64
}

func NewEncoder(w io.ByteWriter, n uint64) *Encoder {
	return &Encoder{n: n, w: bw{w: w}, p: 128, low: 0, high: top - 1}
}

func (t *Encoder) Writes() uint64 {
	return t.writes
}

func (t *Encoder) SetP(p int) {
	const mask = (1 << ProbBits) - 1
	t.p = uint64(p & mask)
}

func (t *Encoder) Encode(bit bool) error {
	if t.n == 0 {
		return io.EOF
	}
	t.n--
	if debug {
		fmt.Printf("encode %t l,h= %b %b\n", bit, t.low, t.high)
	}
	var span uint64
	if bit {
		span = t.p
	} else {
		span = 256 - t.p
	}
	scale := span * (1 + t.high - t.low)
	scale >>= 8
	if bit {
		t.high = t.low + scale - 1
	} else {
		t.low = t.high - scale + 1
	}
	l, h := t.low, t.high
	w := &t.w
	var err error
	for {
		if l >= half {
			bit = true
		} else if h < half {
			bit = false
		} else {
			break
		}
		l = (l << 1) & mask
		h = (h << 1) & mask
		h |= 1
		err = w.WriteBool(bit)
		t.writes++
	}
	t.low, t.high = l, h
	if debug {
		fmt.Printf("\tlow %08b...%08b high %08b...%08b\n", t.low>>oneBits, t.low&0xff,
			t.high>>oneBits, t.high&0xff)
	}
	return err
}

func (t *Encoder) End() error {
	if t.n != 0 {
		return io.EOF
	}
	nFlush := t.writes + uint64(oneBits+ProbBits)
	trg := t.low + (t.high-t.low)/2
	l := uint64(0)
	h := uint64(top) - 1
	m := h / 2
	var err error
	w := &t.w
	for err == nil && (l < t.low || h > t.high) {
		t.writes++
		if m <= trg {
			err = w.WriteBool(true)
			l, m = m, m+(h-m)/2
		} else {
			err = w.WriteBool(false)
			m, h = l+(m-l)/2, m
		}
	}
	for t.writes < nFlush {
		w.WriteBool(false)
		t.writes++
	}
	for t.writes%8 != 0 {
		w.WriteBool(false)
		t.writes++
	}
	return err
}
