package touchable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	var b *Bool

	assert := assert.New(t)

	b = NewBool(false)
	assert.Equal(false, b.Val(), "untouched Bool returns initial value")

	b = NewBool(false)
	b.Set(true)
	assert.Equal(true, b.Val(), "Set updates value")

	b = NewBool(false)
	b.Set(true)
	b.Set(false)
	assert.Equal(false, b.Val(), "Set updates value")

	b = NewBool(false)
	b.Set(true)
	b.SetUnlessTouched(false)
	assert.Equal(true, b.Val(), "SetUnlessTouched leaves touched value unchanged")

	b = NewBool(false)
	assert.Equal(true, b.ValOrDefault(true), "untouched Bool defers to default")

	b = NewBool(false)
	b.Set(true)
	assert.Equal(true, b.ValOrDefault(false), "touched Bool returns value")
}
