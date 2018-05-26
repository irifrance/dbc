// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package dbc

import (
	"fmt"
	"io"

	"github.com/irifrance/bb"
)

type Encoder struct {
	n         uint64
	w         bb.Writer
	p         uint64
	low, high uint64
	writes    uint64
}

func NewEncoder(w bb.Writer, n uint64) *Encoder {
	return &Encoder{n: n, w: w, p: 128, low: 0, high: top - 1}
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
	w := t.w
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
		if err := w.WriteBool(bit); err != nil {
			return err
		}
		if debug {
			fmt.Printf("\t\twrote %t\n", bit)
		}
		t.writes++
	}
	t.low, t.high = l, h
	if debug {
		fmt.Printf("\tlow %08b...%08b high %08b...%08b\n", t.low>>oneBits, t.low&0xff,
			t.high>>oneBits, t.high&0xff)
	}
	return nil
}

func (t *Encoder) End() error {
	if t.n != 0 {
		return io.EOF
	}
	trg := t.low + (t.high-t.low)/2
	l := uint64(0)
	h := uint64(top) - 1
	m := h / 2
	w := t.w
	var err error
	if debug {
		fmt.Printf("end:\n")
	}
	for err == nil && (l < t.low || h > t.high) {
		if m <= trg {
			if debug {
				fmt.Printf("\ttrue\n")
			}
			err = w.WriteBool(true)
			l, m = m, m+(h-m)/2
		} else {
			if debug {
				fmt.Printf("\tfalse\n")
			}
			err = w.WriteBool(false)
			m, h = l+(m-l)/2, m
		}
	}
	return err
}
