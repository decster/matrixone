// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hash

import (
	"math"
	"math/rand"
	"unsafe"
)

const (
	ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const
	c0      = uintptr((8-ptrSize)/4*2860486313 + (ptrSize-4)/4*33054211828000289)
	c1      = uintptr((8-ptrSize)/4*3267000013 + (ptrSize-4)/4*23344194077549503)
	// Constants for multiplication: four random odd 64-bit numbers.
	m1 = 16877499708836156737
	m2 = 2820277070424839065
	m3 = 9497967016996688599
	m4 = 15839092249703872147
)

// hashKey is used to seed the hash function.
var hashKey [4]uintptr

func init() {
	for i := range hashKey {
		hashKey[i] = 1
	}
}

func readUnaligned32(p unsafe.Pointer) uint32 {
	return *(*uint32)(p)
}

func readUnaligned64(p unsafe.Pointer) uint64 {
	return *(*uint64)(p)
}

// Should be a built-in for unsafe.Pointer?
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// This function is copied from the Go runtime.
// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	//lint:ignore SA4016 x ^ 0 is a no-op that fools escape analysis.
	return unsafe.Pointer(x ^ 0)
}

func Memhash(p unsafe.Pointer, seed, s uintptr) uintptr {
	h := uint64(seed + s*hashKey[0])
tail:
	switch {
	case s == 0:
	case s < 4:
		h ^= uint64(*(*byte)(p))
		h ^= uint64(*(*byte)(add(p, s>>1))) << 8
		h ^= uint64(*(*byte)(add(p, s-1))) << 16
		h = rotl31(h*m1) * m2
	case s <= 8:
		h ^= uint64(readUnaligned32(p))
		h ^= uint64(readUnaligned32(add(p, s-4))) << 32
		h = rotl31(h*m1) * m2
	case s <= 16:
		h ^= readUnaligned64(p)
		h = rotl31(h*m1) * m2
		h ^= readUnaligned64(add(p, s-8))
		h = rotl31(h*m1) * m2
	case s <= 32:
		h ^= readUnaligned64(p)
		h = rotl31(h*m1) * m2
		h ^= readUnaligned64(add(p, 8))
		h = rotl31(h*m1) * m2
		h ^= readUnaligned64(add(p, s-16))
		h = rotl31(h*m1) * m2
		h ^= readUnaligned64(add(p, s-8))
		h = rotl31(h*m1) * m2
	default:
		v1 := h
		v2 := uint64(seed * hashKey[1])
		v3 := uint64(seed * hashKey[2])
		v4 := uint64(seed * hashKey[3])
		for s >= 32 {
			v1 ^= readUnaligned64(p)
			v1 = rotl31(v1*m1) * m2
			p = add(p, 8)
			v2 ^= readUnaligned64(p)
			v2 = rotl31(v2*m2) * m3
			p = add(p, 8)
			v3 ^= readUnaligned64(p)
			v3 = rotl31(v3*m3) * m4
			p = add(p, 8)
			v4 ^= readUnaligned64(p)
			v4 = rotl31(v4*m4) * m1
			p = add(p, 8)
			s -= 32
		}
		h = v1 ^ v2 ^ v3 ^ v4
		goto tail
	}

	h ^= h >> 29
	h *= m3
	h ^= h >> 32
	return uintptr(h)
}

func Memhash16(p unsafe.Pointer, h uintptr) uintptr {
	return Memhash(p, h, 2)
}

func Memhash32(p unsafe.Pointer, seed uintptr) uintptr {
	h := uint64(seed + 4*hashKey[0])
	v := uint64(readUnaligned32(p))
	h ^= v
	h ^= v << 32
	h = rotl31(h*m1) * m2
	h ^= h >> 29
	h *= m3
	h ^= h >> 32
	return uintptr(h)
}

func Memhash64(p unsafe.Pointer, seed uintptr) uintptr {
	h := uint64(seed + 8*hashKey[0])
	h ^= uint64(readUnaligned32(p)) | uint64(readUnaligned32(add(p, 4)))<<32
	h = rotl31(h*m1) * m2
	h ^= h >> 29
	h *= m3
	h ^= h >> 32
	return uintptr(h)
}

// Note: in order to get the compiler to issue rotl instructions, we
// need to constant fold the shift amount by hand.
// TODO: convince the compiler to issue rotl instructions after inlining.
func rotl31(x uint64) uint64 {
	return (x << 31) | (x >> (64 - 31))
}

// NOTE: Because NaN != NaN, a map can contain any
// number of (mostly useless) entries keyed with NaNs.
// To avoid long hash chains, we assign a random number
// as the hash value for a NaN.
func F64hash(p unsafe.Pointer, h uintptr) uintptr {
	f := *(*float64)(p)
	if math.IsNaN(f) {
		f = 0
	}
	switch {
	case f == 0:
		return c1 * (c0 ^ h) // +0, -0
	case f != f:
		// TODO(asubiotto): fastrand relies on some stack internals.
		//return c1 * (c0 ^ h ^ uintptr(fastrand())) // any kind of NaN
		return c1 * (c0 ^ h ^ uintptr(rand.Uint32())) // any kind of NaN
	default:
		return Memhash(p, h, 8)
	}
}

func F32hash(p unsafe.Pointer, h uintptr) uintptr {
	f := float64(*(*float32)(p))
	if math.IsNaN(f) {
		f = 0
	}
	switch {
	case f == 0:
		return c1 * (c0 ^ h) // +0, -0
	case f != f:
		// TODO(asubiotto): fastrand relies on some stack internals.
		//return c1 * (c0 ^ h ^ uintptr(fastrand())) // any kind of NaN
		return c1 * (c0 ^ h ^ uintptr(rand.Uint32())) // any kind of NaN
	default:
		return Memhash(p, h, 4)
	}
}
