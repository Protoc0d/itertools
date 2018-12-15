package stream

import (
	"reflect"
	"testing"
)

// Test Streamators for element equality. Allow it1 to be longer than it2
func testStream(t *testing.T, it1, it2 Stream) {
	t.Log("Start")
	for el1 := range it1 {
		if el2, ok := <-it2; !ok {
			t.Error("it2 shorter than it1!", el1)
			return
		} else if !reflect.DeepEqual(el1, el2) {
			t.Error("Elements are not equal", el1, el2)
		} else {
			t.Log(el1, el2)
		}
	}
	t.Log("Stop")
}

func TestMap(t *testing.T) {
	mapper := func(current interface{}, index int, all Stream) interface{} {
		return len(current.(string))
	}
	s := New("a", "ab", "abc", "abcd")
	testStreamEq(t, New(1, 2, 3, 4), s.Map(mapper))
}

// Test Streamators for element equality. Don't allow it1 to be longer than it2
func testStreamEq(t *testing.T, it1, it2 Stream) {
	t.Log("Start")
	for el1 := range it1 {
		if el2, ok := <-it2; !ok {
			t.Error("it2 shorter than it1!", el1)
			return
		} else if !reflect.DeepEqual(el1, el2) {
			t.Error("Elements are not equal", el1, el2)
		} else {
			t.Log(el1, el2)
		}
	}
	if el2, ok := <-it2; ok {
		t.Error("it1 shorter than it2!", el2)
	}
	t.Log("Stop")
}

func TestReduce(t *testing.T) {
	summer := func(memo interface{}, el interface{}) interface{} {
		return memo.(float64) + el.(float64)
	}
	s := Float64(.1, .2, .3, .22)
	if float64(.82)-s.Reduce(summer, float64(0)).(float64) > .000001 {
		t.Error("Sum Reduce failed")
	}
}

func TestList(t *testing.T) {
	list := List(New(1, 2, 3))
	if !reflect.DeepEqual(list, []interface{}{1, 2, 3}) {
		t.Error("List didn't make a list", list)
	}
}

func TestCount(t *testing.T) {
	testStream(t, New(1, 2, 3, 4, 5, 6, 7, 8, 9), Count(1))
}

func TestCycle(t *testing.T) {
	testStream(t, New("a", "b", "ccc", "a", "b", "ccc", "a"), Cycle(New("a", "b", "ccc")))
}

func TestRepeat(t *testing.T) {
	testStreamEq(t, Uint64(100, 100, 100, 100), Repeat(uint64(100), 4))
	testStream(t, Uint64(100, 100, 100, 100), Repeat(uint64(100)))
}

func TestChain(t *testing.T) {
	testStreamEq(t, Int32(1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 100), Chain(Int32(1, 2, 3, 4, 5), Int32(5, 4, 3, 2, 1), Int32(100)))
}

func TestDropWhile(t *testing.T) {
	pred := func(current interface{}, index int, all Stream) bool {
		return current.(int) < 10
	}
	s := Count(0)
	testStream(t, New(10, 11, 12, 13, 14, 15), s.DropWhile(pred))
}

func TestTakeWhile(t *testing.T) {
	pred := func(current interface{}, index int, all Stream) bool {
		return current.(string)[:3] == "abc"
	}
	s := Cycle(New("abcdef", "abcdaj", "ajcde"))
	testStreamEq(t, New("abcdef", "abcdaj"), s.TakeWhile(pred))
}

func TestFilter(t *testing.T) {
	pred := func(current interface{}, index int, all Stream) bool {
		return current.(uint64)%2 == 1
	}
	s := Uint64(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	testStreamEq(t, Uint64(1, 3, 5, 7, 9), s.Filter(pred))
}

func TestTee2(t *testing.T) {
	s := New(5, 4, 3, 2, 1)
	it1, it2 := s.Tee2()
	for i := range it1 {
		j := <-it2
		if i != j {
			t.Error("Tees are not coming off equal")
		}
	}
	s2 := New(1, 2, 3, 4, 5, 6)
	it1, it2 = s2.Tee2()
	testStreamEq(t, New(1, 2, 3, 4, 5, 6), it1)
	testStreamEq(t, New(1, 2, 3, 4, 5, 6), it2)
}

func TestTee(t *testing.T) {
	s := New(3, 4, 5)
	its := s.Tee(3)
	if len(its) != 3 {
		t.Error("its length wrong")
	}
	for _, it := range its {
		testStream(t, New(3, 4, 5), it)
	}
}

func TestZip(t *testing.T) {
	a, b, c := []interface{}{1, "a"}, []interface{}{2, nil}, []interface{}{3, nil}
	test1, test2 := New(a), New(a, b, c)

	testStreamEq(t, test1, Zip(Count(1), New("a")))
	s := Count(1)
	testStreamEq(t, test2, ZipLongest(s.Slice(0, 3), New("a")))
}

func TestStarmap(t *testing.T) {
	multiMapper := func(is ...interface{}) interface{} {
		s := 1
		for _, i := range is {
			s *= i.(int)
		}
		return s
	}
	s := Zip(New(1, 2, 3), Repeat(10, 3))
	testStreamEq(t, New(10, 20, 30), s.Starmap(multiMapper))
}

func TestMultiMap(t *testing.T) {
	multiMapper := func(is ...interface{}) interface{} {
		var s float64
		for _, i := range is {
			s += i.(float64)
		}
		return s
	}
	testStreamEq(t, Float64(10.4, 3.2), MultiMap(multiMapper, Float64(5.2, 1.6, 2.2), Float64(5.2, 1.0), Float64(0, 0.6, 0)))
}

func TestSlice(t *testing.T) {
	s := Count(0)
	s2 := Count(0)
	s3 := Count(0)
	testStream(t, New(5, 6, 7, 8, 9, 10), s.Slice(5))
	testStreamEq(t, New(2, 3, 4, 5, 6, 7, 8), s2.Slice(2, 9))
	testStreamEq(t, New(3, 6, 9), s3.Slice(3, 11, 3))
}

func TestEvery(t *testing.T) {
	s := New(2, 4, 6, 8, 10)

	isPair := func(current interface{}, index int, all Stream) bool {
		return current.(int)%2 == 0
	}

	result := <-s.Every(isPair)

	if !result {
		t.Error("true value expected for Some method, false returned")
	}
}

func TestSome(t *testing.T) {
	s := New(1, 2, 3, 4, 5)

	isPair := func(current interface{}, index int, all Stream) bool {
		return current.(int)%2 == 0
	}

	result := <-s.Some(isPair)

	if !result {
		t.Error("true value expected for Some method, false returned")
	}
}
