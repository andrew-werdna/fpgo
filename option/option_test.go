package option_test

import (
	"encoding/json"
	"testing"

	"github.com/andrew-werdna/fpgo/option"
	"github.com/stretchr/testify/assert"
)

func TestSomeAndNone(t *testing.T) {
	s := option.Some(42)
	n := option.None[int]()

	assert.True(t, s.IsSome())
	assert.False(t, s.IsNone())

	assert.False(t, n.IsSome())
	assert.True(t, n.IsNone())

	assert.Equal(t, 42, s.Unwrap())
	assert.Equal(t, 0, n.UnwrapOr(0))
}

func TestMap(t *testing.T) {
	opt := option.Some(2)
	result := option.Map(opt, func(x int) int {
		return x * 10
	})
	assert.Equal(t, 20, result.Unwrap())

	none := option.None[int]()
	mapped := option.Map(none, func(x int) int {
		return x * 10
	})
	assert.True(t, mapped.IsNone())
}

func TestAndOr(t *testing.T) {
	a := option.Some(1)
	b := option.Some(2)
	none := option.None[int]()

	assert.Equal(t, b, a.And(b))
	assert.Equal(t, option.None[int](), none.And(b))

	assert.Equal(t, a, a.Or(b))
	assert.Equal(t, b, none.Or(b))
}

func TestFilter(t *testing.T) {
	opt := option.Some(5)
	filtered := opt.Filter(func(x int) bool {
		return x > 3
	})
	assert.True(t, filtered.IsSome())

	filtered = opt.Filter(func(x int) bool {
		return x > 10
	})
	assert.True(t, filtered.IsNone())
}

func TestMatch(t *testing.T) {
	opt := option.Some("hello")
	none := option.None[string]()

	var result string
	opt.Match(func(v string) {
		result = v
	}, func() {
		result = "none"
	})
	assert.Equal(t, "hello", result)

	none.Match(func(v string) {
		result = v
	}, func() {
		result = "none"
	})
	assert.Equal(t, "none", result)
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	s := option.Some("go")
	n := option.None[string]()

	b, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, `"go"`, string(b))

	bn, err := json.Marshal(n)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(bn))

	var deserialized option.Option[string]
	err = json.Unmarshal([]byte(`"hello"`), &deserialized)
	assert.NoError(t, err)
	assert.True(t, deserialized.IsSome())
	assert.Equal(t, "hello", deserialized.Unwrap())

	var nullOpt option.Option[string]
	err = json.Unmarshal([]byte(`null`), &nullOpt)
	assert.NoError(t, err)
	assert.True(t, nullOpt.IsNone())
}

func TestFlatten(t *testing.T) {
	nested := option.Some(option.Some("inner"))
	flat := option.Flatten(nested)
	assert.True(t, flat.IsSome())
	assert.Equal(t, "inner", flat.Unwrap())

	outerNone := option.None[option.Option[string]]()
	assert.True(t, option.Flatten(outerNone).IsNone())

	innerNone := option.Some(option.None[string]())
	assert.True(t, option.Flatten(innerNone).IsNone())
}
