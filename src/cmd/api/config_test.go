package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_stringToTime(t *testing.T) {
	type testCase struct {
		input    string
		expected time.Duration
	}

	testCases := []testCase{
		{input: "592ms", expected: time.Millisecond * 592},
		{input: "595 ms", expected: time.Millisecond * 595},
		{input: "12s", expected: time.Second * 12},
		{input: "15 s", expected: time.Second * 15},
		{input: "12m", expected: time.Minute * 12},
		{input: "15 m", expected: time.Minute * 15},
		{input: "12h", expected: time.Hour * 12},
		{input: "15 h", expected: time.Hour * 15},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("input:'%s'", tc.input)
		t.Run(name, func(t *testing.T) {
			duration, err := stringToDuration(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, duration)
		})
	}
}

func Test_stringToTimeErrors(t *testing.T) {
	type testCase struct {
		input string
	}

	testCases := []testCase{
		{input: ""},
		{input: "xd s"},
		{input: "ms 254"},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("input:'%s'", tc.input)
		t.Run(name, func(t *testing.T) {
			duration, err := stringToDuration(tc.input)
			assert.NotNil(t, err)
			assert.Equal(t, time.Duration(0), duration)
		})
	}
}

func Test_stringToBool(t *testing.T) {
	type testCase struct {
		input    string
		expected bool
	}
	testCases := []testCase{
		{input: "1", expected: true},
		{input: "true", expected: true},
		{input: "TrUe", expected: true},
		{input: "0", expected: false},
		{input: "false", expected: false},
		{input: "FaLsE", expected: false},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("input:'%s'", tc.input)
		t.Run(name, func(t *testing.T) {
			val, err := stringToBool(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, val)
		})
	}

}
func Test_stringToBoolErrors(t *testing.T) {
	type testCase struct {
		input string
	}
	testCases := []testCase{
		{input: ""},
		{input: "-1"},
		{input: "-0"},
		{input: "2"},
		{input: "999"},
		{input: "asdf"},
		{input: "QWERTY"},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("input:'%s'", tc.input)
		t.Run(name, func(t *testing.T) {
			val, err := stringToBool(tc.input)
			assert.NotNil(t, err)
			assert.Equal(t, false, val)
		})
	}
}
