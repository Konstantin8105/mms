package mms

import (
	"fmt"
	"sort"
	"sync"
)

// Cache of slices
type Cache struct {
	sync.RWMutex
	ps []pool
}

type pool struct {
	p    *sync.Pool
	size int
}

// Get return slice
func (c *Cache) Get(size int) []float64 {
	// lock
	c.Lock()
	defer func() {
		c.Unlock()
	}()

	// finding index
	index := -1
	for i := range c.ps {
		if c.ps[i].size == size {
			index = i
			break
		}
		if c.ps[i].size > size {
			// typically for creating a new pool
			break
		}
	}

	// creating a new pool
	if index < 0 {
		c.ps = append(c.ps, pool{
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
		index = len(c.ps) - 1
		return c.ps[index].p.New().([]float64)
	}

	// pool is found
	arr := c.ps[index].p.Get().([]float64)

	if len(arr) < size {
		panic(fmt.Errorf("not same sizes: %d != %d", len(arr), size))
	}

	for i := range arr {
		// initialization of slice
		arr[i] = 0.0
	}
	return arr
}

// Put slice into pool
func (c *Cache) Put(arr []float64) {
	c.RLock() // lock

	size := cap(arr)

	// finding index
	index := -1
	for i := range c.ps {
		if c.ps[i].size == size {
			index = i
			break
		}
		if c.ps[i].size > size {
			// typically for creating a new pool
			break
		}
	}

	c.RUnlock() // unlock

	if index < 0 {
		// pool is not exist
		return
	}

	// lock and add
	c.Lock()
	c.ps[index].p.Put(arr)
	c.Unlock()
}
