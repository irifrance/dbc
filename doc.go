// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

// Package dbc provides a compressor of bit sequences with probabilities
// dynamically specifiable for each read/write.
//
// dbc implements a dynamic arithmetic entropy encoder for a 1 bit alphabet.
// The encoder encodes bits under an assumption that the probability of the bit
// being 1 has been set by SetP(), which accepts only values in the range
//
//  [1...2**ProbBits)
//
// A dbc encoding/decoding round trip requires that the number of bits to be
// written is known in advance by both the encoder and decoder, and it requires
// that the encoder and decoder both use exactly the same probability values
// for each bit encoded or decoded.
//
// The degree of compression is a function of how accurate the probabilities
// sent to the encoder are and how far they are from the middle of the
// probability space. If the probabilities are accurate, then in the worst
// case, only a few extra bits will be output, and in the best case you'll get
// about 256x compression.
//
// While these restrictions require the caller to figure out how to model their
// data and how to encode/decode the model, these requirements also give a lot
// of freedom in modelling without worrying about the effects on the
// compression implementation.
//
// As dbc works on the bit level, it is most well suited to compressing
// relatively small blocks of data in applications such as low-latency streams
// or data stores with seek support.
//
package dbc
