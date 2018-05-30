// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package dbc

import (
	"math/rand"
	"testing"

	"github.com/irifrance/bb"
)

func TestBr(t *testing.T) {
	const N = 1029
	bio := bb.NewBuffer(N / 8)
	r := &br{i: 8, r: bio}
	var d [N]bool
	for i := range d {
		if rand.Intn(3) == 1 {
			d[i] = true
		}
		bio.WriteBool(d[i])
	}
	bio.SeekBit(0)
	for i := range d {
		v, _ := r.ReadBool()
		if v != d[i] {
			t.Errorf("%d: got %t not %t\n", i, v, d[i])
		}
	}
}

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
	N := 8199
	enc := NewEncoder(bio, uint64(N))
	d := make([]bool, N)
	ps := make([]int, N)
	for i := range d {
		p := rand.Intn(255) + 1
		bit := rand.Intn(255) < p
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
	if err := enc.End(); err != nil {
		t.Fatal(err)
	}
	bio.SeekBit(0)
	dec := NewDecoder(bio, uint64(len(d)))
	dec.SetP(216)
	for i := range d {
		b, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		if b != d[i] {
			t.Errorf("%d: got %t not %t\n", i, b, d[i])
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

func BenchmarkEncode(b *testing.B) {
	b.StopTimer()
	const N = 8192
	var d [N]bool
	for i := range d {
		if rand.Intn(2) == 1 {
			d[i] = true
		}
	}
	bio := bb.NewBuffer(N / 8)
	enc := NewEncoder(bio, N)
	enc.SetP(38)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bio.SeekBit(0)
		for j := 0; j < N; j++ {
			enc.Encode(d[j])
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	b.StopTimer()
	const N = 8192
	var d [N]bool
	for i := range d {
		if rand.Intn(2) == 1 {
			d[i] = true
		}
	}
	bio := bb.NewBuffer(N / 8)
	enc := NewEncoder(bio, N)
	dec := NewDecoder(bio, N)
	enc.SetP(38)
	dec.SetP(38)
	for i := 0; i < N; i++ {
		enc.Encode(d[i])
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bio.SeekBit(0)
		for j := 0; j < N; j++ {
			dec.Decode()
		}
	}
}
