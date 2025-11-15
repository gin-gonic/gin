/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package alg

// Algorithm 3-way Radix Quicksort, d means the radix.
// Reference: https://algs4.cs.princeton.edu/51radix/Quick3string.java.html
func radixQsort(kvs []_MapPair, d, maxDepth int) {
    for len(kvs) > 11 {
        // To avoid the worst case of quickSort (time: O(n^2)), use introsort here.
        // Reference: https://en.wikipedia.org/wiki/Introsort and
        // https://github.com/golang/go/issues/467
        if maxDepth == 0 {
            heapSort(kvs, 0, len(kvs))
            return
        }
        maxDepth--

        p := pivot(kvs, d)
        lt, i, gt := 0, 0, len(kvs)
        for i < gt {
            c := byteAt(kvs[i].k, d)
            if c < p {
                swap(kvs, lt, i)
                i++
                lt++
            } else if c > p {
                gt--
                swap(kvs, i, gt)
            } else {
                i++
            }
        }

        // kvs[0:lt] < v = kvs[lt:gt] < kvs[gt:len(kvs)]
        // Native implementation:
        //     radixQsort(kvs[:lt], d, maxDepth)
        //     if p > -1 {
        //         radixQsort(kvs[lt:gt], d+1, maxDepth)
        //     }
        //     radixQsort(kvs[gt:], d, maxDepth)
        // Optimize as follows: make recursive calls only for the smaller parts.
        // Reference: https://www.geeksforgeeks.org/quicksort-tail-call-optimization-reducing-worst-case-space-log-n/
        if p == -1 {
            if lt > len(kvs) - gt {
                radixQsort(kvs[gt:], d, maxDepth)
                kvs = kvs[:lt]
            } else {
                radixQsort(kvs[:lt], d, maxDepth)
                kvs = kvs[gt:]
            }
        } else {
            ml := maxThree(lt, gt-lt, len(kvs)-gt)
            if ml == lt {
                radixQsort(kvs[lt:gt], d+1, maxDepth)
                radixQsort(kvs[gt:], d, maxDepth)
                kvs = kvs[:lt]
            } else if ml == gt-lt {
                radixQsort(kvs[:lt], d, maxDepth)
                radixQsort(kvs[gt:], d, maxDepth)
                kvs = kvs[lt:gt]
                d += 1
            } else {
                radixQsort(kvs[:lt], d, maxDepth)
                radixQsort(kvs[lt:gt], d+1, maxDepth)
                kvs = kvs[gt:] 
            }
        }
    }
    insertRadixSort(kvs, d)
}

func insertRadixSort(kvs []_MapPair, d int) {
    for i := 1; i < len(kvs); i++ {
        for j := i; j > 0 && lessFrom(kvs[j].k, kvs[j-1].k, d); j-- {
            swap(kvs, j, j-1)
        }
    }
}

func pivot(kvs []_MapPair, d int) int {
    m := len(kvs) >> 1
    if len(kvs) > 40 {
        // Tukey's ``Ninther,'' median of three mediankvs of three.
        t := len(kvs) / 8
        return medianThree(
            medianThree(byteAt(kvs[0].k, d), byteAt(kvs[t].k, d), byteAt(kvs[2*t].k, d)),
            medianThree(byteAt(kvs[m].k, d), byteAt(kvs[m-t].k, d), byteAt(kvs[m+t].k, d)),
            medianThree(byteAt(kvs[len(kvs)-1].k, d),
                byteAt(kvs[len(kvs)-1-t].k, d),
                byteAt(kvs[len(kvs)-1-2*t].k, d)))
    }
    return medianThree(byteAt(kvs[0].k, d), byteAt(kvs[m].k, d), byteAt(kvs[len(kvs)-1].k, d))
}

func medianThree(i, j, k int) int {
    if i > j {
        i, j = j, i
    } // i < j
    if k < i {
        return i
    }
    if k > j {
        return j
    }
    return k
}

func maxThree(i, j, k int) int {
    max := i
    if max < j {
        max = j
    }
    if max < k {
        max = k
    }
    return max
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returnkvs 2*ceil(lg(n+1)).
func maxDepth(n int) int {
    var depth int
    for i := n; i > 0; i >>= 1 {
        depth++
    }
    return depth * 2
}

// siftDown implements the heap property on kvs[lo:hi].
// first is an offset into the array where the root of the heap lies.
func siftDown(kvs []_MapPair, lo, hi, first int) {
    root := lo
    for {
        child := 2*root + 1
        if child >= hi {
            break
        }
        if child+1 < hi && kvs[first+child].k < kvs[first+child+1].k {
            child++
        }
        if kvs[first+root].k >= kvs[first+child].k {
            return
        }
        swap(kvs, first+root, first+child)
        root = child
    }
}

func heapSort(kvs []_MapPair, a, b int) {
    first := a
    lo := 0
    hi := b - a

    // Build heap with the greatest element at top.
    for i := (hi - 1) / 2; i >= 0; i-- {
        siftDown(kvs, i, hi, first)
    }

    // Pop elements, the largest first, into end of kvs.
    for i := hi - 1; i >= 0; i-- {
        swap(kvs, first, first+i)
        siftDown(kvs, lo, i, first)
    }
}

// Note that _MapPair.k is NOT pointed to _MapPair.m when map key is integer after swap
func swap(kvs []_MapPair, a, b int) {
    kvs[a].k, kvs[b].k = kvs[b].k, kvs[a].k
    kvs[a].v, kvs[b].v = kvs[b].v, kvs[a].v
}

// Compare two strings from the pos d.
func lessFrom(a, b string, d int) bool {
    l := len(a)
    if l > len(b) {
        l = len(b)
    }
    for i := d; i < l; i++ {
        if a[i] == b[i] {
            continue
        }
        return a[i] < b[i]
    }
    return len(a) < len(b)
}

func byteAt(b string, p int) int {
    if p < len(b) {
        return int(b[p])
    }
    return -1
}
