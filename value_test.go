/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Values test
 */

package goipp

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestValueEncode tests Value.encode for all value types
func TestValueEncode(t *testing.T) {
	noError := errors.New("")
	longstr := strings.Repeat("x", 65536)
	loc1 := time.FixedZone("UTC+3:30", 3*3600+1800)
	tm1, _ := time.ParseInLocation(time.RFC3339, "2025-03-29T16:48:53+03:30", loc1)
	loc2 := time.FixedZone("UTC-3", -3*3600)
	tm2, _ := time.ParseInLocation(time.RFC3339, "2025-03-29T16:48:53-03:00", loc2)

	type testData struct {
		v    Value  // Input value
		data []byte // Expected output data
		err  string // Expected error string ("" if no error)
	}

	tests := []testData{
		// Simple values
		{Binary{}, []byte{}, ""},
		{Binary{1, 2, 3}, []byte{1, 2, 3}, ""},
		{Boolean(false), []byte{0}, ""},
		{Boolean(true), []byte{1}, ""},
		{Integer(0), []byte{0, 0, 0, 0}, ""},
		{Integer(0x01020304), []byte{1, 2, 3, 4}, ""},
		{String(""), []byte{}, ""},
		{String("Hello"), []byte("Hello"), ""},
		{Void{}, []byte{}, ""},

		// Range
		{
			v:    Range{0x01020304, 0x05060708},
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			v: Range{-100, 100},
			data: []byte{
				0xff, 0xff, 0xff, 0x9c, 0x00, 0x00, 0x00, 0x64,
			},
		},

		// Resolution
		{
			v:    Resolution{0x01020304, 0x05060708, 0x09},
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			v: Resolution{150, 300, UnitsDpi},
			data: []byte{
				0x00, 0x00, 0x00, 0x96, // 150
				0x00, 0x00, 0x01, 0x2c, // 300
				0x03,
			},
		},

		// TextWithLang
		{
			v: TextWithLang{"en-US", "Hello!"},
			data: []byte{
				0x00, 0x05,
				'e', 'n', '-', 'U', 'S',
				0x00, 0x06,
				'H', 'e', 'l', 'l', 'o', '!',
			},
		},

		{
			v: TextWithLang{"ru-RU", "Привет!"},
			data: []byte{
				0x00, 0x05,
				'r', 'u', '-', 'R', 'U',
				0x00, 0x0d,
				0xd0, 0x9f, 0xd1, 0x80, 0xd0, 0xb8, 0xd0, 0xb2,
				0xd0, 0xb5, 0xd1, 0x82, 0x21,
			},
		},

		{
			v:   TextWithLang{"en-US", longstr},
			err: "Text exceeds 65535 bytes",
		},

		{
			v:   TextWithLang{longstr, "hello"},
			err: "Lang exceeds 65535 bytes",
		},

		// Time
		{
			v: Time{tm1},
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x1e, // Minutes from UTC
			},
		},

		{
			v: Time{tm2},
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'-',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
		},

		// Collection
		//
		// Note, Collection.encode is the stub and encodes
		// as Void. Actual collection encoding handled the
		// different way.
		{
			v: Collection{
				MakeAttribute("test", TagString, String("")),
			},
			data: []byte{},
		},
	}

	for _, test := range tests {
		data, err := test.v.encode()
		if err == nil {
			err = noError
		}

		vstr := test.v.String()
		if len(vstr) > 40 {
			vstr = vstr[:40] + "..."
		}

		if err.Error() != test.err {
			t.Errorf("testing %s.encode:\n"+
				"value:          %s\n"+
				"error expected: %q\n"+
				"error present:  %q\n",
				reflect.TypeOf(test.v).String(),
				vstr,
				test.err,
				err,
			)
			continue
		}

		if test.err == "" && !bytes.Equal(data, test.data) {
			t.Errorf("testing %s.encode:\n"+
				"value:         %s\n"+
				"data expected: %x\n"+
				"data present:  %x\n",
				reflect.TypeOf(test.v).String(),
				vstr,
				test.data,
				data,
			)
		}
	}
}

// TestValueEncode tests Value.decode for all value types
func TestValueDecode(t *testing.T) {
	noError := errors.New("")
	loc1 := time.FixedZone("UTC+3:30", 3*3600+1800)
	tm1, _ := time.ParseInLocation(time.RFC3339, "2025-03-29T16:48:53+03:30", loc1)
	loc2 := time.FixedZone("UTC-3", -3*3600)
	tm2, _ := time.ParseInLocation(time.RFC3339, "2025-03-29T16:48:53-03:00", loc2)

	type testData struct {
		data []byte // Input data
		v    Value  // Expected output value
		err  string // Expected error string ("" if no error)
	}

	tests := []testData{
		// Simple types
		{[]byte{}, Binary{}, ""},
		{[]byte{1, 2, 3, 4}, Binary{1, 2, 3, 4}, ""},
		{[]byte{0}, Boolean(false), ""},
		{[]byte{1}, Boolean(true), ""},
		{[]byte{0, 1}, Boolean(false), "value must be 1 byte"},
		{[]byte{1, 2, 3, 4}, Integer(0x01020304), ""},
		{[]byte{}, Integer(0), "value must be 4 bytes"},
		{[]byte{1, 2, 3, 4, 5}, Integer(0), "value must be 4 bytes"},
		{[]byte{}, Void{}, ""},
		{[]byte("hello"), String("hello"), ""},
		{[]byte{1, 2, 3, 4, 5}, Void{}, ""},

		// Range
		{
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			v:    Range{0x01020304, 0x05060708},
		},

		{
			data: []byte{1, 2, 3, 4, 5, 6, 7},
			v:    Range{},
			err:  "value must be 8 bytes",
		},

		{
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			v:    Range{},
			err:  "value must be 8 bytes",
		},

		// Resolution
		{
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			v:    Resolution{0x01020304, 0x05060708, 0x09},
		},

		{
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			v:    Resolution{},
			err:  "value must be 9 bytes",
		},

		{
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			v:    Resolution{},
			err:  "value must be 9 bytes",
		},

		// Time
		{
			// Good time, UTC+...
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x1e, // Minutes from UTC
			},
			v: Time{tm1},
		},

		{
			// Good time, UTC-...
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'-',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v: Time{tm2},
		},

		{
			// Truncated data
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
			},
			v:   Time{},
			err: "value must be 11 bytes",
		},

		{
			// Extra data
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
				0,
			},
			v:   Time{},
			err: "value must be 11 bytes",
		},

		{
			// Bad month
			data: []byte{
				0x07, 0xe9,
				0,    // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad month 0",
		},

		{
			// Bad day
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				32,   // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad day 32",
		},

		{
			// Bad hours
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				99,   // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad hours 99",
		},

		{
			// Bad minutes
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				88,   // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad minutes 88",
		},

		{
			// Bad seconds
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				77,   // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad seconds 77",
		},

		{
			// Bad deciseconds
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				100,  // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad deciseconds 100",
		},

		{
			// Bad UTC sign
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'?',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad UTC sign",
		},

		{
			// Bad UTC hours
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				66,   // Hours from UTC
				0x00, // Minutes from UTC
			},
			v:   Time{},
			err: "bad UTC hours 66",
		},

		{
			// Bad UTC minutes
			data: []byte{
				0x07, 0xe9,
				0x03, // Month, 1...12
				0x1d, // Day, 1...31
				0x10, // Hour, 0...23
				0x30, // Minutes, 0...59
				0x35, // Seconds, 0...59
				0x00, // Deci-seconds, 0...9
				'+',  // Direction from UTC, +/-
				0x03, // Hours from UTC
				166,  // Minutes from UTC
			},
			v:   Time{},
			err: "bad UTC minutes 166",
		},

		// TextWithLang
		{
			data: []byte{
				0x00, 0x05,
				'e', 'n', '-', 'U', 'S',
				0x00, 0x06,
				'H', 'e', 'l', 'l', 'o', '!',
			},
			v: TextWithLang{"en-US", "Hello!"},
		},

		{
			data: []byte{
				0x00, 0x05,
				'r', 'u', '-', 'R', 'U',
				0x00, 0x0d,
				0xd0, 0x9f, 0xd1, 0x80, 0xd0, 0xb8, 0xd0, 0xb2,
				0xd0, 0xb5, 0xd1, 0x82, 0x21,
			},
			v: TextWithLang{"ru-RU", "Привет!"},
		},

		{
			// truncated language length
			data: []byte{},
			v:    TextWithLang{},
			err:  "truncated language length",
		},

		{
			// truncated language name
			data: []byte{
				0x00, 0x05,
				'e',
			},
			v:   TextWithLang{},
			err: "truncated language name",
		},

		{
			// truncated text length
			data: []byte{
				0x00, 0x05,
				'e', 'n', '-', 'U', 'S',
				0x00,
			},
			v:   TextWithLang{},
			err: "truncated text length",
		},

		{
			// truncated text string
			data: []byte{
				0x00, 0x05,
				'e', 'n', '-', 'U', 'S',
				0x00, 0x06,
				'H', 'e',
			},
			v:   TextWithLang{},
			err: "truncated text string",
		},

		{
			// extra data bytes
			data: []byte{
				0x00, 0x05,
				'e', 'n', '-', 'U', 'S',
				0x00, 0x06,
				'H', 'e', 'l', 'l', 'o', '!',
				0, 2, 3,
			},
			v:   TextWithLang{},
			err: "extra 3 bytes at the end of value",
		},
	}

	for _, test := range tests {
		v, err := test.v.decode(test.data)
		if err == nil {
			err = noError
		}

		if err.Error() != test.err {
			t.Errorf("testing %s.decode:\n"+
				"value:          %s\n"+
				"error expected: %q\n"+
				"error present:  %q\n",
				reflect.TypeOf(test.v).String(),
				v,
				test.err,
				err,
			)
			continue
		}

		if test.err == "" && !reflect.DeepEqual(v, test.v) {
			t.Errorf("testing %s.decode:\n"+
				"data:           %x\n"+
				"value expected: %#v\n"+
				"value present:  %#v\n",
				reflect.TypeOf(test.v).String(),
				test.data,
				test.v,
				v,
			)
		}

	}
}

// TestValueCollectionDecode tests Collection.decode for all value types
func TestValueCollectionDecode(t *testing.T) {
	// Collection.decode is a stub and must panic
	defer func() {
		recover()
	}()

	v := Collection{}
	v.decode([]byte{})

	t.Errorf("Collection.decode() method is a stub and must panic")
}

// TestValueString rests Value.String method for various
// kinds of the Value
func TestValueString(t *testing.T) {
	loc1 := time.FixedZone("UTC+3:30", 3*3600+1800)
	tm1, _ := time.ParseInLocation(time.RFC3339, "2025-03-29T16:48:53+03:30", loc1)

	type testData struct {
		v Value  // Input value
		s string // Expected output string
	}

	tests := []testData{
		// Simple types
		{Binary{}, ""},
		{Binary{1, 2, 3}, "010203"},
		{Integer(123), "123"},
		{Integer(-321), "-321"},
		{Range{-100, 200}, "-100-200"},
		{Range{-100, -50}, "-100--50"},
		{Resolution{150, 300, UnitsDpi}, "150x300dpi"},
		{Resolution{100, 200, UnitsDpcm}, "100x200dpcm"},
		{Resolution{75, 150, 10}, "75x150unknown(0x0a)"},
		{String("hello"), "hello"},
		{TextWithLang{"en-US", "hello"}, "hello [en-US]"},
		{Time{tm1}, "2025-03-29T16:48:53+03:30"},
		{Void{}, ""},

		// Collections
		{Collection{}, "{}"},

		{
			v: Collection{
				MakeAttr("attr1", TagInteger, Integer(1)),
				MakeAttr("attr2", TagString, String("hello")),
			},
			s: "{attr1=1 attr2=hello}",
		},
	}

	for _, test := range tests {
		s := test.v.String()
		if s != test.s {
			t.Errorf("testing %s.String:\n"+
				"value:    %#v\n"+
				"expected: %q\n"+
				"present:  %q\n",
				reflect.TypeOf(test.v).String(),
				test.v, test.s, s,
			)
		}
	}
}

// TestValueType rests Value.Type method for various
// kinds of the Value
func TestValueType(t *testing.T) {
	type testData struct {
		v  Value // Input value
		tp Type  // Expected output type
	}

	tests := []testData{
		{Binary(nil), TypeBinary},
		{Boolean(false), TypeBoolean},
		{Collection(nil), TypeCollection},
		{Integer(0), TypeInteger},
		{Range{}, TypeRange},
		{Resolution{}, TypeResolution},
		{String(""), TypeString},
		{TextWithLang{}, TypeTextWithLang},
		{Time{time.Time{}}, TypeDateTime},
		{Void{}, TypeVoid},
	}

	for _, test := range tests {
		tp := test.v.Type()
		if tp != test.tp {
			t.Errorf("testing %s.Type:\n"+
				"expected: %q\n"+
				"present:  %q\n",
				reflect.TypeOf(test.v).String(),
				test.tp, tp,
			)
		}
	}
}

// TestValueEqualSimilar tests ValueEqual and ValueSimilar
func TestValueEqualSimilar(t *testing.T) {
	tm1 := time.Now()
	tm2 := tm1.Add(time.Hour)

	type testData struct {
		v1, v2  Value // A pair of values
		equal   bool  // Expected ValueEqual(v1,v2) output
		similar bool  // Expected ValueSimilar(v1,v2) output
	}

	tests := []testData{
		// Simple types
		{Integer(0), Integer(0), true, true},
		{Integer(0), Integer(1), false, false},
		{Integer(0), String("hello"), false, false},
		{Time{tm1}, Time{tm1}, true, true},
		{Time{tm1}, Time{tm2}, false, false},
		{Binary{}, Binary{}, true, true},
		{Binary{}, Binary{1, 2, 3}, false, false},
		{Binary{1, 2, 3}, Binary{4, 5, 6}, false, false},
		{Binary("hello"), Binary("hello"), true, true},
		{String("hello"), String("hello"), true, true},
		{Binary("hello"), String("hello"), false, true},
		{String("hello"), Binary("hello"), false, true},

		// Collections
		//
		// Note, ValueEqual for Collection values is a thin wrapper
		// around Attributes.Equal. So the serious testing will be
		// performed there. Here we only test a couple of simple
		// cases.
		//
		// The same is true for ValueSimilar.
		{Collection{}, Collection{}, true, true},

		{
			v1: Collection{
				MakeAttr("attr1", TagInteger, Integer(1)),
				MakeAttr("attr2", TagString, String("hello")),
			},

			v2: Collection{
				MakeAttr("attr2", TagString, String("hello")),
				MakeAttr("attr1", TagInteger, Integer(1)),
			},

			equal:   false,
			similar: true,
		},
	}

	for _, test := range tests {
		equal := ValueEqual(test.v1, test.v2)
		similar := ValueSimilar(test.v1, test.v2)

		if equal != test.equal {
			t.Errorf("testing ValueEqual:\n"+
				"value 1:  %s\n"+
				"value 2:  %s\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.v1, test.v2,
				test.equal, equal,
			)
		}

		if similar != test.similar {
			t.Errorf("testing ValueSimilar:\n"+
				"value 1:  %s\n"+
				"value 2:  %s\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.v1, test.v2,
				test.similar, similar,
			)
		}
	}
}

// TestValuesString tests Values.String function
func TestValuesString(t *testing.T) {
	type testData struct {
		v Values // Input Values
		s string // Expected output
	}

	tests := []testData{
		{
			v: nil,
			s: "[]",
		},

		{
			v: Values{},
			s: "[]",
		},

		{
			v: Values{{TagInteger, Integer(5)}},
			s: "5",
		},

		{
			v: Values{{TagInteger, Integer(5)}, {TagEnum, Integer(6)}},
			s: "[5,6]",
		},
	}

	for _, test := range tests {
		s := test.v.String()
		if s != test.s {
			t.Errorf("testing Values.String:\n"+
				"value:    %#v\n"+
				"expected: %q\n"+
				"present:  %q\n",
				test.v, test.s, s,
			)
		}
	}
}

// TestValuesEqualSimilar tests Values.Equal and Values.Similar
func TestValuesEqualSimilar(t *testing.T) {
	type testData struct {
		v1, v2  Values // A pair of values
		equal   bool   // Expected v1.Equal(v2) output
		similar bool   // Expected v2.Similar(v1,v2) output
	}

	tests := []testData{
		{
			v1:      nil,
			v2:      nil,
			equal:   true,
			similar: true,
		},

		{
			v1:      Values{},
			v2:      Values{},
			equal:   true,
			similar: true,
		},

		{
			v1:      Values{},
			v2:      nil,
			equal:   false,
			similar: true,
		},

		{
			v1:      Values{},
			v2:      Values{{TagInteger, Integer(5)}},
			equal:   false,
			similar: false,
		},

		{
			v1: Values{
				{TagInteger, Integer(5)},
				{TagEnum, Integer(6)},
			},
			v2: Values{
				{TagInteger, Integer(5)},
				{TagEnum, Integer(6)},
			},
			equal:   true,
			similar: true,
		},

		{
			v1: Values{
				{TagInteger, Integer(5)},
				{TagEnum, Integer(6)}},
			v2: Values{
				{TagInteger, Integer(5)},
				{TagInteger, Integer(6)}},
			equal:   false,
			similar: false,
		},

		{
			v1:      Values{{TagInteger, Integer(6)}, {TagEnum, Integer(5)}},
			v2:      Values{{TagInteger, Integer(5)}, {TagEnum, Integer(6)}},
			equal:   false,
			similar: false,
		},

		{
			v1: Values{
				{TagString, String("hello")},
				{TagString, Binary("world")},
			},
			v2: Values{
				{TagString, String("hello")},
				{TagString, Binary("world")},
			},
			equal:   true,
			similar: true,
		},

		{
			v1: Values{
				{TagString, Binary("hello")},
				{TagString, String("world")},
			},
			v2: Values{
				{TagString, String("hello")},
				{TagString, Binary("world")},
			},
			equal:   false,
			similar: true,
		},
	}

	for _, test := range tests {
		equal := test.v1.Equal(test.v2)
		similar := test.v1.Similar(test.v2)

		if equal != test.equal {
			t.Errorf("testing Values.Equal:\n"+
				"values 1:  %s\n"+
				"values 2:  %s\n"+
				"expected:  %v\n"+
				"present:   %v\n",
				test.v1, test.v2,
				test.equal, equal,
			)
		}

		if similar != test.similar {
			t.Errorf("testing Values.Similar:\n"+
				"values 1: %s\n"+
				"values 2: %s\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.v1, test.v2,
				test.similar, similar,
			)
		}
	}
}

// TestValuesCopy tests Values.Clone and Values.DeepCopy methods
func TestValuesCopy(t *testing.T) {
	values := Values{}
	values.Add(TagBoolean, Boolean(true))
	values.Add(TagExtension, Binary{})
	values.Add(TagString, Binary{1, 2, 3})
	values.Add(TagInteger, Integer(123))
	values.Add(TagEnum, Integer(-321))
	values.Add(TagRange, Range{-100, 200})
	values.Add(TagRange, Range{-100, -50})
	values.Add(TagResolution, Resolution{150, 300, UnitsDpi})
	values.Add(TagResolution, Resolution{100, 200, UnitsDpcm})
	values.Add(TagResolution, Resolution{75, 150, 10})
	values.Add(TagString, String("hello"))
	values.Add(TagTextLang, TextWithLang{"en-US", "hello"})
	values.Add(TagDateTime, Time{time.Now()})
	values.Add(TagNoValue, Void{})
	values.Add(TagBeginCollection,
		Collection{MakeAttribute("test", TagString, String(""))})

	clone := values.Clone()
	copy := values.DeepCopy()

	if !values.Equal(clone) {
		t.Errorf("Values.Clone test failed")
	}

	if !values.Equal(copy) {
		t.Errorf("Values.DeepCopy test failed")
	}
}

// TestIntegerOrRange tests IntegerOrRange interface implementation
// by Integer and Range types
func TestIntegerOrRange(t *testing.T) {
	type testData struct {
		v      IntegerOrRange // Value being tested
		x      int            // IntegerOrRange.Within input
		within bool           // IntegerOrRange.Within expected output
	}

	tests := []testData{
		{
			v:      Integer(5),
			x:      5,
			within: true,
		},

		{
			v:      Integer(5),
			x:      6,
			within: false,
		},

		{
			v:      Range{-5, 5},
			x:      0,
			within: true,
		},

		{
			v:      Range{-5, 5},
			x:      -5,
			within: true,
		},

		{
			v:      Range{-5, 5},
			x:      5,
			within: true,
		},

		{
			v:      Range{-5, 5},
			x:      -6,
			within: false,
		},

		{
			v:      Range{-5, 5},
			x:      6,
			within: false,
		},
	}

	for _, test := range tests {
		within := test.v.Within(test.x)
		if within != test.within {
			t.Errorf("testing %s.Within:\n"+
				"value:    %#v\n"+
				"param:    %v\n"+
				"expected: %v\n"+
				"present:  %v\n",
				reflect.TypeOf(test.v).String(),
				test.v, test.x,
				test.within, within,
			)
		}
	}
}

// TestCollection tests Collection.Add method
func TestCollectionAdd(t *testing.T) {
	col1 := Collection{
		MakeAttr("attr1", TagInteger, Integer(1)),
		MakeAttr("attr2", TagString, String("hello")),
	}

	col2 := Collection{}
	col2.Add(MakeAttr("attr1", TagInteger, Integer(1)))
	col2.Add(MakeAttr("attr2", TagString, String("hello")))

	if !ValueEqual(col1, col2) {
		t.Errorf("Collection.Add test failed")
	}
}
