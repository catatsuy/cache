package cache_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/catatsuy/cache"
)

func TestDo(t *testing.T) {
	sf := cache.NewSingleflightGroup[string]()
	v, err, _ := sf.Do("key", func() (string, error) {
		return "bar", nil
	})
	if got, want := v, "bar"; got != want {
		t.Errorf("Do = %s; want %s", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestDoErr(t *testing.T) {
	sf := cache.NewSingleflightGroup[string]()
	someErr := errors.New("Some error")
	v, err, _ := sf.Do("key", func() (string, error) {
		return "", someErr
	})
	if err != someErr {
		t.Errorf("Do error = %v; want someErr %v", err, someErr)
	}
	if v != "" {
		t.Errorf("unexpected empty value %#v", v)
	}
}

func TestDoDupSuppress(t *testing.T) {
	sf := cache.NewSingleflightGroup[string]()
	var wg1, wg2 sync.WaitGroup
	c := make(chan string, 1)
	var calls int32
	fn := func() (string, error) {
		if atomic.AddInt32(&calls, 1) == 1 {
			// First invocation.
			wg1.Done()
		}
		v := <-c
		c <- v // pump; make available for any future calls

		time.Sleep(10 * time.Millisecond) // let more goroutines enter Do

		return v, nil
	}

	const n = 100

	wg1.Add(1)
	for range n {
		wg1.Add(1)
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			wg1.Done()
			v, err, _ := sf.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
				return
			}
			if v != "bar" {
				t.Errorf("Do = %s; want %s", v, "bar")
			}
		}()
	}
	wg1.Wait()
	// At least one goroutine is in fn now and all of them have at
	// least reached the line before the Do.
	c <- "bar"
	wg2.Wait()
	if got := atomic.LoadInt32(&calls); got <= 0 || got >= n {
		t.Errorf("number of calls = %d; want over 0 and less than %d", got, n)
	}
}

func TestDoTimeout(t *testing.T) {
	sf := cache.NewSingleflightGroup[string]()
	start := time.Now()
	v, err, _ := sf.Do("key", func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "bar", nil
	})
	if err != nil {
		t.Errorf("Do error: %v", err)
	}
	if v != "bar" {
		t.Errorf("Do = %s; want %s", v, "bar")
	}
	if time.Since(start) < 100*time.Millisecond {
		t.Errorf("Do executed too quickly; expected delay")
	}
}

func TestDoMultipleErrors(t *testing.T) {
	sf := cache.NewSingleflightGroup[string]()
	var calls int32
	someErr := errors.New("Some error")

	const n = 10
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err, _ := sf.Do("key", func() (string, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(10 * time.Millisecond)
				return "", someErr
			})
			if err != someErr {
				t.Errorf("Do error = %v; want %v", err, someErr)
			}
			if v != "" {
				t.Errorf("Do = %v; want empty string", v)
			}
		}()
	}
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
}
