// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package dbc

import (
	"math/rand"
	"testing"

	"github.com/irifrance/bb"
)

func TestDbcStaticGood(t *testing.T) {
	for i := 1; i < 256; i++ {
		testDbcStatic(i, false, t)
	}
}

func TestDbcStaticBad(t *testing.T) {
	for i := 1; i < 256; i++ {
		testDbcStatic(i, true, t)
	}
}

func TestDbcDynamic(t *testing.T) {
	bio := bb.NewBuffer(1024)
	N := 16389
	enc := NewEncoder(bio, uint64(N))
	d := make([]bool, N)
	ps := make([]int, N)
	for i := range d {
		p := rand.Intn(255) + 1
		bit := rand.Intn(p) < p
		d[i] = bit
		ps[i] = p
		enc.SetP(p)
		if err := enc.Encode(bit); err != nil {
			t.Fatal(err)
		}
	}
	if err := enc.End(); err != nil {
		t.Fatal(err)
	}
	bio.SeekBit(0)
	dec := NewDecoder(bio, uint64(N))
	for i := range d {
		dec.SetP(ps[i])
		bit, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		if bit != d[i] {
			t.Errorf("%d: got %t not %t\n", i, bit, d[i])
		}
	}
	if dec.Reads() != enc.Writes() {
		t.Errorf("enc wrote %d dec read %d\n", enc.Writes(), dec.Reads())
	}
}

func TestDbcIz(t *testing.T) {
	d := []bool{true,
		true, true, true, true, true, true, true, true, true,
		true, true, true, false, true, true, true, true, true,
		false, true, false, true, true, true, true, true, false,
		true, false, true, true}

	bio := bb.NewBuffer(1024)
	enc := NewEncoder(bio, uint64(len(d)))
	enc.SetP(216)
	for _, v := range d {
		if err := enc.Encode(v); err != nil {
			t.Fatal(err)
		}
	}
	enc.End()
	bio.SeekBit(0)
	dec := NewDecoder(bio, uint64(len(d)))
	dec.SetP(216)
	e := make([]bool, len(d))
	for i := range d {
		b, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		e[i] = b
	}
	for i, v := range d {
		if e[i] != v {
			t.Errorf("%d: got %t not %t\n", i, e[i], v)
		}
	}
}

func testDbcStatic(p int, flip bool, t *testing.T) {
	bio := bb.NewBuffer(4)
	N := 1029
	d := make([]bool, N)
	enc := NewEncoder(bio, uint64(N))
	enc.SetP(p)
	for i := 0; i < N; i++ {
		bit := rand.Intn(256) <= p
		if flip {
			bit = !bit
		}
		d[i] = bit
		if err := enc.Encode(bit); err != nil {
			t.Fatal(err)
		}
	}
	if err := enc.End(); err != nil {
		t.Fatal(err)
	}
	//t.Logf("p=%d bad=%t wrote %d bits with %d\n", p, flip, N, bio.BitsWritten())
	bio.SeekBit(0)
	dec := NewDecoder(bio, uint64(N))
	dec.SetP(p)
	for i := 0; i < N; i++ {
		bit, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		if d[i] != bit {
			t.Errorf("p=%d, bad=%t %d: got %t not %t\n", p, flip, i, bit, d[i])
		}
	}
	if enc.Writes() != dec.Reads() {
		t.Errorf("dec reads %d enc writes %d\n", dec.Reads(), enc.Writes())
	}
}
