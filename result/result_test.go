package result_test

import (
	"errors"
	"testing"

	"github.com/andrew-werdna/fpgo/result"
	"github.com/stretchr/testify/assert"
)

func divide(a, b int) result.Result[int] {
	if b == 0 {
		return result.Err[int](errors.New("division by zero"))
	}
	return result.Ok(a / b)
}

func TestOkAndErr(t *testing.T) {
	ok := result.Ok(10)
	err := result.Err[int](errors.New("fail"))

	assert.True(t, ok.IsOk())
	assert.False(t, ok.IsErr())
	assert.Equal(t, 10, ok.Unwrap())

	assert.True(t, err.IsErr())
	assert.Equal(t, "fail", err.UnwrapErr().Error())
}

func TestUnwrapOr(t *testing.T) {
	r := divide(10, 2)
	assert.Equal(t, 5, r.UnwrapOr(-1))

	r = divide(10, 0)
	assert.Equal(t, -1, r.UnwrapOr(-1))
}

func TestMapAndMapErr(t *testing.T) {
	r := result.Ok(4)
	mapped := result.Map(r, func(v int) string {
		t.Logf("current value is: %d\n", rune(v))
		return "x" + string(rune(v+48))
	})
	assert.True(t, mapped.IsOk())
	assert.Equal(t, "x4", mapped.Unwrap())

	e := result.Err[int](errors.New("fail"))
	errMapped := result.MapErr(e, func(err error) error {
		return errors.New("mapped: " + err.Error())
	})
	assert.True(t, errMapped.IsErr())
	assert.Contains(t, errMapped.UnwrapErr().Error(), "mapped:")
}

func TestMatch(t *testing.T) {
	r := result.Ok("hello")
	e := result.Err[string](errors.New("nope"))

	var val string
	r.Match(func(s string) {
		val = s
	}, func(_ error) {
		val = "err"
	})
	assert.Equal(t, "hello", val)

	e.Match(func(s string) {
		val = s
	}, func(_ error) {
		val = "err"
	})
	assert.Equal(t, "err", val)
}

func TestAndThenFlatMap(t *testing.T) {
	addOne := func(x int) result.Result[int] {
		return result.Ok(x + 1)
	}
	fail := func(_ int) result.Result[int] {
		return result.Err[int](errors.New("fail"))
	}

	r := result.Ok(5)
	r2 := r.AndThen(addOne).AndThen(addOne)
	assert.True(t, r2.IsOk())
	assert.Equal(t, 7, r2.Unwrap())

	r3 := r.AndThen(addOne).AndThen(fail)
	assert.True(t, r3.IsErr())
}

func TestFlatten(t *testing.T) {
	nested := result.Ok(result.Ok("hello"))
	flat := result.Flatten(nested)
	assert.True(t, flat.IsOk())
	assert.Equal(t, "hello", flat.Unwrap())

	nestedErr := result.Err[result.Result[string]](errors.New("fail"))
	flatErr := result.Flatten(nestedErr)
	assert.True(t, flatErr.IsErr())
}
