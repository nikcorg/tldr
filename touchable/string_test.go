package touchable

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	var s *String

	s = NewString("hello")
	assert.Equal("hello", s.Val(), "untouched String returns initial value")

	s = NewString("beep")
	s.Set("boop")
	assert.Equal("boop", s.Val(), "touched String returns value")

	s = NewString("beep")
	s.Set("boop")
	s.Set("brrt")
	assert.Equal("brrt", s.Val(), "second Set reassigns value")

	s = NewString("beep")
	s.SetUnlessTouched("brrt")
	assert.Equal("brrt", s.Val(), "SetUnlessTouched assigns untouched value")

	s = NewString("beep")
	s.Set("boop")
	s.SetUnlessTouched("brrt")
	assert.Equal("boop", s.Val(), "SetUnlessTouched after Set leaves value unchanged")

	s = NewString("beep")
	assert.Equal("boop", s.ValOrDefault("boop"), "ValOrDefault returns default for untouched String")

	s = NewString("beep")
	s.Set("brrt")
	assert.Equal("brrt", s.ValOrDefault("boop"), "ValOrDefault returns value for touched String")
}
