// Package stream provides a translation of the python standard library module stream.
// Many of the functions have been brought over, althought not all.
// In this implementation, chan interface{} has been used as all Streamators; if more specific types are necessary,
// feel free to copy the code to your project to be implemented with more specific types.
package stream

import (
	"sync"
)

type Stream chan interface{}
type MultiMapper func(...interface{}) interface{}
type Reducer func(memo interface{}, element interface{}) interface{}

type MapCallback func(current interface{}, index int, all Stream) interface{}
type FilterCallback func(current interface{}, index int, all Stream) bool
type ReduceCallback func(memo interface{}, element interface{}) interface{}

func New(els ...interface{}) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func Int64(els ...int64) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func Int32(els ...int32) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func Float64(els ...float64) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func Float32(els ...float32) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func Uint(els ...uint) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}
func Uint64(els ...uint64) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func Uint32(els ...uint32) Stream {
	c := make(Stream)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

func List(it Stream) []interface{} {
	arr := make([]interface{}, 0, 1)
	for el := range it {
		arr = append(arr, el)
	}
	return arr
}

// Count from i to infinity
func Count(i int) Stream {
	c := make(Stream)
	go func() {
		for ; true; i++ {
			c <- i
		}
	}()
	return c
}

// Cycle through an Streamator infinitely (requires memory)
func Cycle(it Stream) Stream {
	c, a := make(Stream), make([]interface{}, 0, 1)
	go func() {
		for el := range it {
			a = append(a, el)
			c <- el
		}
		for {
			for _, el := range a {
				c <- el
			}
		}
	}()
	return c
}

// Repeat an element n times or infinitely
func Repeat(el interface{}, n ...int) Stream {
	c := make(Stream)
	go func() {
		for i := 0; len(n) == 0 || i < n[0]; i++ {
			c <- el
		}
		close(c)
	}()
	return c
}

// Chain together multiple Streamators
func Chain(its ...Stream) Stream {
	c := make(Stream)
	go func() {
		for _, it := range its {
			for el := range it {
				c <- el
			}
		}
		close(c)
	}()
	return c
}

// Elements after pred(el) == true
func (it Stream) DropWhile(fn FilterCallback) Stream {
	c := make(Stream)
	go func() {
		i := 0
		for el := range it {
			if !fn(el, i, it) {
				c <- el
				break
			}
			i++
		}
		for el := range it {
			c <- el
		}
		close(c)
	}()
	return c
}

// Elements before pred(el) == false
func (it Stream) TakeWhile(fn FilterCallback) Stream {
	c := make(Stream)
	go func() {
		i := 0
		for el := range it {
			if fn(el, i, it) {
				c <- el
			} else {
				break
			}
			i++
		}
		close(c)
	}()
	return c
}

// Filter out any elements where pred(el) == false
func (it Stream) Filter(fn FilterCallback) Stream {
	c := make(Stream)
	go func() {
		i := 0
		for el := range it {
			if fn(el, i, it) {
				c <- el
			}
			i++
		}
		close(c)
	}()
	return c
}

// Every test if all element verify callback condition
func (it Stream) Every(fn FilterCallback) chan bool {
	c := make(chan bool)
	go func() {
		i := 0
		for el := range it {
			if !fn(el, i, it) {
				break
			}
			i++
		}
		c <- i != len(it)
		close(c)
	}()
	return c
}

func (it Stream) Some(fn FilterCallback) chan bool {
	c := make(chan bool)
	go func() {
		i := 0
		found := false
		for el := range it {
			if fn(el, i, it) {
				found = true
				break
			}
			i++
		}
		c <- found
		close(c)
	}()
	return c
}

// Sub-Streamator from start (inclusive) to [stop (exclusive) every [step (default 1)]]
func (it Stream) Slice(startstopstep ...int) Stream {
	start, stop, step := 0, 0, 1
	if len(startstopstep) == 1 {
		start = startstopstep[0]
	} else if len(startstopstep) == 2 {
		start, stop = startstopstep[0], startstopstep[1]
	} else if len(startstopstep) >= 3 {
		start, stop, step = startstopstep[0], startstopstep[1], startstopstep[2]
	}

	c := make(Stream)
	go func() {
		i := 0
		// Start
		for el := range it {
			if i >= start {
				c <- el // inclusive
				break
			}
			i += 1
		}

		// Stop
		i, j := i+1, 1
		for el := range it {
			if stop > 0 && i >= stop {
				break
			} else if j%step == 0 {
				c <- el
			}

			i, j = i+1, j+1
		}

		close(c)
	}()
	return c
}

// Map an Streamator to fn(el) for el in it
func (it Stream) Map(fn MapCallback) Stream {
	c := make(Stream)
	go func() {
		i := 0
		for el := range it {
			c <- fn(el, i, it)
			i++
		}
		close(c)
	}()
	return c
}

// Map p, q, ... to fn(pEl, qEl, ...)
// Breaks on first closed channel
func MultiMap(fn MultiMapper, its ...Stream) Stream {
	c := make(Stream)
	go func() {
	Outer:
		for {
			els := make([]interface{}, len(its))
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					break Outer
				}
			}
			c <- fn(els...)
		}
		close(c)
	}()
	return c
}

// Map p, q, ... to fn(pEl, qEl, ...)
// Breaks on last closed channel
func MultiMapLongest(fn MultiMapper, its ...Stream) Stream {
	c := make(Stream)
	go func() {
		for {
			els := make([]interface{}, len(its))
			n := 0
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					n += 1
				}
			}
			if n < len(its) {
				c <- fn(els...)
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

// Map an Streamator if arrays to a fn(els...)
// Stream must be an Streamator of []interface{} (possibly created by Zip)
// If not, Starmap will act like MultiMap with a single Streamator
func (it Stream) Starmap(fn MultiMapper) Stream {
	c := make(Stream)
	go func() {
		for els := range it {
			if elements, ok := els.([]interface{}); ok {
				c <- fn(elements...)
			} else {
				c <- fn(els)
			}
		}
		close(c)
	}()
	return c
}

// Zip up multiple interators into one
// Close on shortest Streamator
func Zip(its ...Stream) Stream {
	c := make(Stream)
	go func() {
		defer close(c)
		for {
			els := make([]interface{}, len(its))
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					return
				}
			}
			c <- els
		}
	}()
	return c
}

// Zip up multiple Streamators into one
// Close on longest Streamator
func ZipLongest(its ...Stream) Stream {
	c := make(Stream)
	go func() {
		for {
			els := make([]interface{}, len(its))
			n := 0
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					n += 1
				}
			}
			if n < len(its) {
				c <- els
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

// Reduce the Streamator (aka fold) from the left
func (it Stream) Reduce(fn ReduceCallback, memo interface{}) interface{} {
	for el := range it {
		memo = fn(memo, el)
	}
	return memo
}

// Split an Streamator into n multiple Streamators
// Requires memory to keep values for n Streamators
func (it Stream) Tee(n int) []Stream {
	deques := make([][]interface{}, n)
	Streams := make([]Stream, n)
	for i := 0; i < n; i++ {
		Streams[i] = make(Stream)
	}

	mutex := new(sync.Mutex)

	gen := func(myStream Stream, i int) {
		for {
			if len(deques[i]) == 0 {
				mutex.Lock()
				if len(deques[i]) == 0 {
					if newval, ok := <-it; ok {
						for i, d := range deques {
							deques[i] = append(d, newval)
						}
					} else {
						mutex.Unlock()
						close(myStream)
						break
					}
				}
				mutex.Unlock()
			}
			var popped interface{}
			popped, deques[i] = deques[i][0], deques[i][1:]
			myStream <- popped
		}
	}
	for i, Stream := range Streams {
		go gen(Stream, i)
	}
	return Streams
}

// Helper to tee just into two Streamators
func (it Stream) Tee2() (Stream, Stream) {
	Streams := it.Tee(2)
	return Streams[0], Streams[1]
}
