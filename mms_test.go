package mms

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

// mm - interface memory management
type mm interface {
	Get(size int) []float64
	Put(arr []float64)
}

// direct allocation of memory
type Direct struct{}

func (c *Direct) Get(size int) []float64 {
	return make([]float64, size)
}

func (c *Direct) Put(arr []float64) {
	_ = arr
}

// test

func getChan() chan int {
	ch := make(chan int, 100) // with buffer

	// generate sizes
	var wg sync.WaitGroup
	wg.Add(gos)
	for j := 0; j < gos; j++ {
		go func() {
			for _, s := range sizesExpect {
				for k := 0; k < s.amount; k++ {
					ch <- s.size
				}
			}
			wg.Done()
		}()
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
	{43, 5},
}

const gos = 3

func Test(t *testing.T) {
	memory := []mm{
		&Direct{},
		&Cache{},
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
					atomic.AddInt64(&profile[i][amount], 1)
					memory[i].Put(arr)
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
			if profile[0][sizesExpect[i].size] != int64(sizesExpect[i].amount*gos) {
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
		{"Cache", &Cache{}},
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
							memory[im].m.Put(arr)
						}
						wg.Done()
					}()
				}
				wg.Wait()
			}
		})
	}
}
