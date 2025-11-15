// Copyright (c) 2012-2020 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

//go:build safe || codec.safe || !gc

package codec

// growCap will return a new capacity for a slice, given the following:
//   - oldCap: current capacity
//   - unit: in-memory size of an element
//   - num: number of elements to add
func growCap(oldCap, unit, num uint) (newCap uint) {
	// appendslice logic (if cap < 1024, *2, else *1.25):
	//   leads to many copy calls, especially when copying bytes.
	//   bytes.Buffer model (2*cap + n): much better for bytes.
	// smarter way is to take the byte-size of the appended element(type) into account

	// maintain 1 thresholds:
	// t1: if cap <= t1, newcap = 2x
	//     else          newcap = 1.5x
	//
	// t1 is always >= 1024.
	// This means that, if unit size >= 16, then always do 2x or 1.5x (ie t1, t2, t3 are all same)
	//
	// With this, appending for bytes increase by:
	//    100% up to 4K
	//     50% beyond that

	// unit can be 0 e.g. for struct{}{}; handle that appropriately
	maxCap := num + (oldCap * 3 / 2)
	if unit == 0 || maxCap > maxArrayLen || maxCap < oldCap { // handle wraparound, etc
		return maxArrayLen
	}

	var t1 uint = 1024 // default thresholds for large values
	if unit <= 4 {
		t1 = 8 * 1024
	} else if unit <= 16 {
		t1 = 2 * 1024
	}

	newCap = 2 + num
	if oldCap > 0 {
		if oldCap <= t1 { // [0,t1]
			newCap = num + (oldCap * 2)
		} else { // (t1,infinity]
			newCap = maxCap
		}
	}

	// ensure newCap takes multiples of a cache line (size is a multiple of 64)
	t1 = newCap * unit
	if t2 := t1 % 64; t2 != 0 {
		t1 += 64 - t2
		newCap = t1 / unit
	}

	return
}
