package rotation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRotation(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Unset, NewFromString("unset"))
	assert.Equal(Unset, NewFromString("u"))
	assert.Equal(None, NewFromString("none"))
	assert.Equal(None, NewFromString("n"))
	assert.Equal(Daily, NewFromString("daily"))
	assert.Equal(Daily, NewFromString("d"))
	assert.Equal(Weekly, NewFromString("weekly"))
	assert.Equal(Weekly, NewFromString("w"))
	assert.Equal(Monthly, NewFromString("monthly"))
	assert.Equal(Monthly, NewFromString("m"))
	assert.Equal(Yearly, NewFromString("yearly"))
	assert.Equal(Yearly, NewFromString("y"))

	assert.Equal("unset", Unset.String())
	assert.Equal("none", None.String())
	assert.Equal("daily", Daily.String())
	assert.Equal("weekly", Weekly.String())
	assert.Equal("monthly", Monthly.String())
	assert.Equal("yearly", Yearly.String())
}

func TestUnmarshalYAML(t *testing.T) {
	assert := assert.New(t)

	tmp := struct {
		Beep Period `yaml:"beep"`
	}{}

	if err := yaml.Unmarshal([]byte(`beep: monthly`), &tmp); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(Monthly, tmp.Beep)

	if err := yaml.Unmarshal([]byte(`beep: m`), &tmp); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(Monthly, tmp.Beep)
}

func TestMarshalYAML(t *testing.T) {
	assert := assert.New(t)

	tmp := struct {
		Beep Period `yaml:"beep"`
	}{
		Beep: Weekly,
	}

	yamlSrc, err := yaml.Marshal(tmp)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("yaml:\n%v\n", string(yamlSrc))

	assert.Equal("beep: weekly\n", string(yamlSrc))
}
