// Code generated by mmsgen; DO NOT EDIT.

package mms

import (
	"fmt"
	"sort"
	"sync"
	"unsafe"
)

// Float64sCache of slices
type Float64sCache struct {
	mutex  sync.RWMutex
	ps     []poolFloat64sCache
	putarr []uintptr
}

type poolFloat64sCache struct {
	p    *sync.Pool
	size int
}

// Get return slice
func (c *Float64sCache) Get(size int) []float64 {

	if size == 0 {
		// empty size
		return
	}

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

	// lock and add
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()
	if index < len(c.ps) && c.ps[index].size == size {
		*arr = (*arr)[:0]
		if Debug {
			// check if putting same arr
			ptr := uintptr(unsafe.Pointer(arr))
			for i := range c.putarr {
				if c.putarr[i] == ptr {
					panic(fmt.Errorf("dublicate of putting"))
				}
			}
			c.putarr = append(c.putarr, ptr)
			return
		}
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
	c.putarr = make([]uintptr, 0)
}
