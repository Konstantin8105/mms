package main

const code = `
// Code generated by mmsgen; DO NOT EDIT.

package {{ .Package }}

import (
	"sort"
	"sync"
	"fmt"
	"unsafe"
	{{ range .Imports }}
	"{{ . }}"
	{{ end }}
)

// {{ .CacheName }} of slices
type {{ .CacheName }} struct {
	mutex  sync.RWMutex
	ps     []pool{{ .CacheName }}
	putarr []uintptr
}

type pool{{ .CacheName }} struct {
	p    *sync.Pool
	size int
}

// Get return slice
func (c *{{ .CacheName }}) Get(size int) {{ .Type }} {
	// lock
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()

	// finding index
	index := c.index(size)

	// creating a new pool
	if index < 0 {
		c.ps = append(c.ps, pool{{ .CacheName }}{
			p: &sync.Pool{
				New: func() interface{} {
					return {{ .CodeNew }}
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
		return c.ps[index].p.New().({{ .Type }})
	}

	// pool is found
	arr := c.ps[index].p.Get().({{ .Type }})

	if len(arr) == 0{
		arr = arr[:size]
	}

	if Debug {
		if len(arr) < size {
			panic(fmt.Errorf("not same sizes: %d != %d", len(arr), size))
		}
	}

	for i := range arr {
		// initialization of slice
		arr[i] = 0
	}
	return arr
}

// Put slice into pool
func (c *{{ .CacheName }}) Put(arr *{{ .Type }}) {
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

	// lock and add
	c.mutex.Lock()
	if index < len(c.ps) && c.ps[index].size == size {
		if Debug {
			// check if putting same arr
			ptr := uintptr(unsafe.Pointer(arr))
			for i := range c.putarr {
				if c.putarr[i] == ptr {
					panic(fmt.Errorf("dublicate of putting"))
				}
			}
			c.putarr = append(c.putarr, ptr)
		}
		*arr = (*arr)[:0]
		c.ps[index].p.Put(*arr)
	}
	c.mutex.Unlock()
}

// return index with excepted size
func (c *{{ .CacheName }}) index(size int) int {
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

`
