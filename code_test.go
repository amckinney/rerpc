package rerpc

import (
	"strings"
	"testing"
)

func TestCodeMarshaling(t *testing.T) {
	valid := make([]Code, 0)
	for c := minCode; c <= maxCode; c++ {
		valid = append(valid, c)
	}

	unmarshal := func(t testing.TB, c *Code, text []byte) {
		if err := c.UnmarshalText(text); err != nil {
			t.Errorf("unexpected error unmarshaling Code from %q", text)
		}
	}

	t.Run("round-trip", func(t *testing.T) {
		for _, c := range valid {
			out, err := c.MarshalText()
			if err != nil {
				t.Errorf("failed to marshal code %v as text: %v", c, err)
			}
			in := new(Code)
			unmarshal(t, in, out)
			if *in != c {
				t.Errorf("failed to round-trip code %v", c)
			}
		}
	})

	t.Run("out of bounds", func(t *testing.T) {
		if _, err := Code(42).MarshalText(); err == nil {
			t.Log("expected error marshaling invalid code")
			t.Fail()
		}

		Code(42).String() // shouldn't panic

		c := new(Code)
		if err := c.UnmarshalText([]byte("42")); err == nil {
			t.Log("expected error unmarshaling invalid code")
			t.Fail()
		}
	})

	t.Run("from string", func(t *testing.T) {
		c := new(Code)
		text := []byte(`"UNIMPLEMENTED"`)
		unmarshal(t, c, text)
		if *c != CodeUnimplemented {
			t.Errorf("unmarshaled %q as %v", text, *c)
		}
	})

	t.Run("to string", func(t *testing.T) {
		for _, c := range valid {
			if strings.Contains(c.String(), "(") {
				t.Errorf("regenerate stringer method for code: %v is out of date", c)
			}
		}
	})
}
