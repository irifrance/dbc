// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package dbc

import (
	"fmt"
	"io"

	"github.com/irifrance/bb"
)

type Decoder struct {
	r         bb.Reader
	n         uint64
	err       error
	p         uint64
	low, high uint64
	bio       uint64
	reads     uint64
}

func NewDecoder(r bb.Reader, n uint64) *Decoder {
	return &Decoder{n: n, r: r, p: 128, low: top - 1, high: top - 1}
}

func (d *Decoder) SetP(p int) {
	d.p = uint64(p & 0xff)
}

func (d *Decoder) fmt(msg string) {
	fmt.Printf("%s\n\tlow %08b...%08b\n\thigh %08b...%08b\n\tbio %08b...%08b\n",
		msg, d.low>>oneBits, d.low&0xff, d.high>>oneBits, d.high&0xff, d.bio>>oneBits,
		d.bio&0xff)
}

func (d *Decoder) Reads() uint64 {
	return d.reads
}

func (d *Decoder) slurp() error {
	r := d.r
	var bit bool
	var err error
	for {
		if d.low >= half {
		} else if d.high < half {
		} else {
			break
		}
		d.low = (d.low << 1) & mask
		d.high = (d.high << 1) & mask
		d.high |= 1
		d.bio = (d.bio << 1) & mask
		bit, err = r.ReadBool()
		if err != nil {
			return err
		}
		d.reads++
		if bit {
			d.bio |= 1
		}
	}
	return nil
}

func (d *Decoder) Decode() (bool, error) {
	if d.n == 0 {
		return false, io.EOF
	}
	d.n--
	var outBit bool
	var nDrain uint64
	if d.n == 0 {
		nDrain = d.reads
		_ = nDrain
	}
	// here
	if err := d.slurp(); err != nil {
		return false, err
	}
	span := 1 + d.high - d.low
	bDiff := (d.bio - d.low) << 8
	val := bDiff / span

	if debug {
		fmt.Printf("=> val %d, p %d\n", val, d.p)
	}
	if val < d.p {
		outBit = true
		d.high = d.low + (span*d.p)>>8 - 1
	} else {
		outBit = false
		d.low = d.high - (span*(256-d.p))>>8 + 1
	}
	if d.n == 0 {
		if err := d.slurp(); err != nil {
			return false, err
		}
	}
	return outBit, nil
}
