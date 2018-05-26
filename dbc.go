// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package dbc

const (
	ProbBits = 8
	oneBits  = 48
	one      = 1 << oneBits
	top      = one << ProbBits
	half     = top / 2
	mask     = top - 1
)

const (
	debug = false
)
