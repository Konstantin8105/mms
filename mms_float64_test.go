package mms

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mm - interface memory management
type mm interface {
	Get(size int) []float64
	Put(arr *[]float64)
}

// direct allocation of memory
type Direct struct{}

func (c *Direct) Get(size int) []float64 {
	return make([]float64, size)
}

func (c *Direct) Put(arr *[]float64) {
	_ = arr
}

// test

func getChan() chan int {
	ch := make(chan int, 10) // with buffer

	// generate sizes
	var wg sync.WaitGroup
	wg.Add(len(sizesExpect))
	for j := 0; j < len(sizesExpect); j++ {
		go func(j int) {
			for k := 0; k < sizesExpect[j].amount; k++ {
				ch <- sizesExpect[j].size
			}
			wg.Done()
		}(j)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

var sizesExpect = [...]struct {
	size   int
	amount int
}{
	{12, 50},
	{99, 14},
	{15, 32},
	{4, 22},
	{43, 2},
}

const gos = 10

func Test(t *testing.T) {
	memory := []mm{
		&Direct{},
		&Float64sCache{},
	}

	var profile [2][100]int64

	for i := range memory {
		ch := getChan()

		var wg sync.WaitGroup
		wg.Add(gos)
		for j := 0; j < gos; j++ {
			go func() {
				for amount := range ch {
					arr := memory[i].Get(amount)
					// check input data
					for o := range arr {
						if arr[o] != 0.0 {
							t.Errorf("Not zero initialization for %d flow: %e", i, arr[o])
						}
					}
					// change input data
					for o := range arr {
						arr[o] = rand.Float64()
					}
					size := cap(arr)
					atomic.AddInt64(&profile[i][size], 1)
					time.Sleep(time.Millisecond)
					memory[i].Put(&arr)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	for i := range profile[0] {
		if profile[0][i] != profile[1][i] {
			t.Errorf("Size = %d. Not same %d != %d", i, profile[0][i], profile[1][i])
		}
	}
	for i := range sizesExpect {
		for j := 0; j < 2; j++ {
			if profile[0][sizesExpect[i].size] != int64(sizesExpect[i].amount) {
				t.Errorf("not expectd size for pos = %d. %d != %d ",
					sizesExpect[i].size,
					profile[0][sizesExpect[i].size],
					sizesExpect[i].amount*gos,
				)
			}
		}
	}
}

func Benchmark(b *testing.B) {
	memory := []struct {
		name string
		m    mm
	}{
		{"Direct", &Direct{}},
		{"Float64sCache", &Float64sCache{}},
	}

	for im := range memory {
		b.Run(memory[im].name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ch := getChan()

				var wg sync.WaitGroup
				wg.Add(gos)
				for j := 0; j < gos; j++ {
					go func() {
						for amount := range ch {
							arr := memory[im].m.Get(amount)
							// change input data
							for o := range arr {
								arr[o] = rand.Float64()
							}
							memory[im].m.Put(&arr)
						}
						wg.Done()
					}()
				}
				wg.Wait()
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	var c Float64sCache
	arr := c.Get(3)
	if len(arr) != 3 {
		t.Errorf("not valid len")
	}
	arr[0] = 42
	c.Put(&arr)
	for _, size := range []int{2, 5, 5, 3, 100, 2, 5, 3, 100} {
		arr2 := c.Get(size)
		if len(arr2) != size {
			t.Errorf("not valid len: %d:%d with size = %d",
				len(arr2), cap(arr2), size)
		}
		if arr2[0] != 0 {
			t.Errorf("not same arrays")
		}
		c.Put(&arr2)
	}
}

func TestDouble(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Cannot found double putting")
		}
	}()
	oldDebug := Debug
	defer func() {
		Debug = oldDebug
	}()

	Debug = true

	size := 5

	var c Float64sCache
	arr := c.Get(size)
	c.Put(&arr)
	c.Put(&arr)
}

func TestMemoryAccessAfterPut(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Cannot found memory access after putting")
		}
	}()
	oldDebug := Debug
	defer func() {
		Debug = oldDebug
	}()

	Debug = true

	size := 5

	var c Float64sCache
	arr := c.Get(size)
	c.Put(&arr)
	arr[3] = 42
}

func TestReset(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Cannot reset cache")
		}
	}()
	oldDebug := Debug
	defer func() {
		Debug = oldDebug
	}()

	Debug = true

	size := 5

	var c Float64sCache
	arr := c.Get(size)
	c.Put(&arr)
	c.Reset() // reset
	c.Put(&arr)
}
