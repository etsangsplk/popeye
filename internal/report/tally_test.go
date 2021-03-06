package report

import (
	"bytes"
	"testing"

	"github.com/derailed/popeye/internal/linter"
	"github.com/stretchr/testify/assert"
)

func TestTallyWrite(t *testing.T) {
	uu := []struct {
		jurassic bool
		e        string
	}{
		{false, "💥 0 😱 0 🔊 0 ✅ 0 \x1b[38;5;196;m0\x1b[0m٪"},
		{true, "E:0 W:0 I:0 OK:0 0%%"},
	}

	for _, u := range uu {
		ta := NewTally()
		b := bytes.NewBuffer([]byte(""))
		s := NewSanitizer(b, 0, u.jurassic)
		ta.write(b, s)

		assert.Equal(t, u.e, b.String())
	}
}

func TestTallyRollup(t *testing.T) {
	uu := []struct {
		issues linter.Issues
		e      *Tally
	}{
		{
			linter.Issues{},
			&Tally{counts: []int{0, 0, 0, 0}, score: 0, valid: false},
		},
		{
			linter.Issues{
				"a": {
					linter.NewError(linter.InfoLevel, ""),
					linter.NewError(linter.WarnLevel, ""),
				},
				"b": {
					linter.NewError(linter.ErrorLevel, ""),
				},
				"c": {},
			},
			&Tally{counts: []int{1, 1, 1, 1}, score: 50, valid: true},
		},
	}

	for _, u := range uu {
		ta := NewTally()
		ta.Rollup(u.issues)

		assert.Equal(t, u.e, ta)
	}
}

func TestTallyScore(t *testing.T) {
	uu := []struct {
		issues linter.Issues
		e      int
	}{
		{
			linter.Issues{
				"a": {
					linter.NewError(linter.InfoLevel, ""),
					linter.NewError(linter.WarnLevel, ""),
				},
				"b": {
					linter.NewError(linter.ErrorLevel, ""),
				},
				"c": {},
			},
			50,
		},
	}

	for _, u := range uu {
		ta := NewTally()
		ta.Rollup(u.issues)

		assert.Equal(t, u.e, ta.Score())
	}
}

func TestTallyWidth(t *testing.T) {
	uu := []struct {
		issues linter.Issues
		e      string
	}{
		{
			linter.Issues{
				"a": {
					linter.NewError(linter.InfoLevel, ""),
					linter.NewError(linter.WarnLevel, ""),
				},
				"b": {
					linter.NewError(linter.ErrorLevel, ""),
				},
				"c": {},
			},
			"💥 1 😱 1 🔊 1 ✅ 1 \x1b[38;5;196;m50\x1b[0m٪",
		},
	}

	s := new(Sanitizer)
	for _, u := range uu {
		ta := NewTally()
		ta.Rollup(u.issues)

		assert.Equal(t, u.e, ta.Dump(s))
	}
}

func TestToPerc(t *testing.T) {
	uu := []struct {
		v1, v2 float64
		e      float64
	}{
		{0, 0, 0},
		{100, 50, 200},
		{50, 100, 50},
	}

	for _, u := range uu {
		assert.Equal(t, u.e, toPerc(u.v1, u.v2))
	}
}

func TestMarshalJSON(t *testing.T) {
	uu := []struct {
		t *Tally
		e string
	}{
		{NewTally(), `{"ok":0,"info":0,"warning":0,"error":0,"score":0}`},
	}

	for _, u := range uu {
		s, err := u.t.MarshalJSON()
		assert.Nil(t, err)
		assert.Equal(t, u.e, string(s))
	}
}

func TestMarshalYAML(t *testing.T) {
	uu := []struct {
		t *Tally
		e interface{}
	}{
		{NewTally(), struct {
			OK    int `yaml:"ok"`
			Info  int `yaml:"info"`
			Warn  int `yaml:"warning"`
			Error int `yaml:"error"`
			Score int `yaml:"score"`
		}{
			OK:    0,
			Info:  0,
			Warn:  0,
			Error: 0,
			Score: 0,
		}},
	}

	for _, u := range uu {
		s, err := u.t.MarshalYAML()
		assert.Nil(t, err)
		assert.Equal(t, u.e, s)
	}
}
