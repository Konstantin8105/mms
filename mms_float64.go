// Code generated by mmsgen; DO NOT EDIT.

package mms

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"
	"unsafe"
)

// Float64sCache of slices
type Float64sCache struct {
	mutex  sync.RWMutex
	ps     []poolFloat64sCache
	putarr []debugFloat64sCache
}

type poolFloat64sCache struct {
	p    *sync.Pool
	size int
}

type debugFloat64sCache struct {
	ptr  unsafe.Pointer
	size int
	line string
	hash [sha256.Size]byte
}

// Get return slice
func (c *Float64sCache) Get(size int) []float64 {

	// lock
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()

	// finding index
	index := c.index(size)

	// creating a new pool
	if index < 0 {
		c.ps = append(c.ps, poolFloat64sCache{
			p: &sync.Pool{
				New: func() interface{} {
					return make([]float64, size)
				},
			},
			size: size,
		})
		// sort
		sort.SliceStable(c.ps, func(i, j int) bool {
			return c.ps[i].size < c.ps[j].size
		})

		// return
		index = c.index(size)
		return c.ps[index].p.New().([]float64)
	}

	// pool is found
	arr := c.ps[index].p.Get().([]float64)

	if len(arr) == 0 {
		arr = arr[:size]
	}

	if Debug {
		if len(arr) < size {
			panic(fmt.Errorf("not same sizes: %d != %d", len(arr), size))
		}
		if len(arr) != cap(arr) {
			panic(fmt.Errorf("not valid capacity: %d != %d", len(arr), cap(arr)))
		}
	}

	for i := range arr {
		// initialization of slice
		arr[i] = 0
	}
	return arr
}

// Put slice into pool
func (c *Float64sCache) Put(arr *[]float64) {
	c.mutex.RLock() // lock
	var (
		size  = cap(*arr)
		index = c.index(size) // finding index
	)
	c.mutex.RUnlock() // unlock

	if index < 0 {
		// pool is not exist
		return
	}
	if size == 0 {
		// empty size
		return
	}
	if len(*arr) == 0 {
		// propably it is a dublicate putting
		return
	}

	// lock and add
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()
	if index < len(c.ps) && c.ps[index].size == size {
		if Debug {
			// check if putting same arr
			ptr := (unsafe.Pointer(arr))
			hsh := sha256.Sum256([]byte(fmt.Sprintf("%v", *arr)))
			for i := range c.putarr {
				if c.putarr[i].size == size &&
					c.putarr[i].ptr == ptr &&
					c.putarr[i].hash == hsh {
					length := 12
					if cap(*arr) < length {
						length = cap(*arr)
					}
					panic(fmt.Errorf(
						"Dublicate of Put.\n"+
							"Last is called in :\n%v\n"+
							"Present call in   :\n%v\n"+
							"Array   : %v\n"+
							"Pointer : %v\n"+
							"Hash256 : %v\n",
						c.putarr[i].line,
						called(),
						(*arr)[:length],
						ptr,
						hsh,
					))
				}
			}
			c.putarr = append(c.putarr, debugFloat64sCache{
				ptr:  ptr,
				size: size,
				line: called(),
				hash: hsh,
			})
			*arr = (*arr)[:0]
			ptr = (unsafe.Pointer(arr))
			c.putarr = append(c.putarr, debugFloat64sCache{
				ptr:  ptr,
				size: size,
				line: called(),
				hash: hsh,
			})
			return
		}
		*arr = (*arr)[:0]
		c.ps[index].p.Put(*arr)
	}
}

// return index with excepted size
func (c *Float64sCache) index(size int) int {
	index := -1
	for i := range c.ps {
		if c.ps[i].size < size {
			continue
		}
		if c.ps[i].size == size {
			index = i
		}
		break
	}
	return index
}

// Reset internal structure.
// In debug case - better for founding double putting.
// In normal case - for memory management with different memory profile.
//
//	Example of code:
//	w := cache.Get(10)
//	defer func() {
//		if mms.Debug {
//			cache.Reset()
//		}
//	}
//	... // Put memory in cache in next lines of code
//
func (c *Float64sCache) Reset() {
	// lock
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()

	// remove
	c.ps = make([]poolFloat64sCache, 0)
	c.putarr = make([]debugFloat64sCache, 0)
}
