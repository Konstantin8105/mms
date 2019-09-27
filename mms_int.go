// Code generated by mmsgen; DO NOT EDIT.

package mms

import (
	"sort"
	"sync"
)

// IntsCache of slices
type IntsCache struct {
	mutex sync.RWMutex
	ps    []poolIntsCache
}

type poolIntsCache struct {
	p    *sync.Pool
	size int
}

// Get return slice
func (c *IntsCache) Get(size int) []int {
	// lock
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()

	// finding index
	index := c.index(size)

	// creating a new pool
	if index < 0 {
		c.ps = append(c.ps, poolIntsCache{
			p: &sync.Pool{
				New: func() interface{} {
					return make([]int, size)
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
		return c.ps[index].p.New().([]int)
	}

	// pool is found
	arr := c.ps[index].p.Get().([]int)

	// Only for debugging:
	//	if len(arr) < size {
	//		panic(fmt.Errorf("not same sizes: %d != %d", len(arr), size))
	//	}

	for i := range arr {
		// initialization of slice
		arr[i] = 0
	}
	return arr
}

// Put slice into pool
func (c *IntsCache) Put(arr []int) {
	c.mutex.RLock() // lock
	var (
		size  = cap(arr)
		index = c.index(size) // finding index
	)
	c.mutex.RUnlock() // unlock

	if index < 0 {
		// pool is not exist
		return
	}

	// lock and add
	c.mutex.Lock()
	if index < len(c.ps) && c.ps[index].size == size {
		c.ps[index].p.Put(arr)
	}
	c.mutex.Unlock()
}

// return index with excepted size
func (c *IntsCache) index(size int) int {
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