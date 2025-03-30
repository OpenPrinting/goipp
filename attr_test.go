/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Message attributes tests
 */

package goipp

import (
	"errors"
	"testing"
	"time"
	"unsafe"
)

// TestAttributesEqualSimilar tests Attributes.Equal and
// Attributes.Similar methods
func TestAttributesEqualSimilar(t *testing.T) {
	type testData struct {
		a1, a2  Attributes // A pair of Attributes slice
		equal   bool       // Expected a1.Equal(a2) output
		similar bool       // Expected a2.Similar(a2) output

	}

	tests := []testData{
		{
			// nil Attributes are equal and similar
			a1:      nil,
			a2:      nil,
			equal:   true,
			similar: true,
		},

		{
			// Empty Attributes are equal and similar
			a1:      Attributes{},
			a2:      Attributes{},
			equal:   true,
			similar: true,
		},

		{
			// Attributes(nil) vs Attributes{} are similar but not equal
			a1:      Attributes{},
			a2:      nil,
			equal:   false,
			similar: true,
		},

		{
			// Attributes of different length are neither equal nor similar
			a1: Attributes{
				MakeAttr("attr1", TagInteger, Integer(0)),
			},
			a2:      Attributes{},
			equal:   false,
			similar: false,
		},
		{
			// Same Attributes are equal and similar
			a1: Attributes{
				MakeAttr("attr1", TagInteger, Integer(0)),
			},
			a2: Attributes{
				MakeAttr("attr1", TagInteger, Integer(0)),
			},
			equal:   true,
			similar: true,
		},
		{
			// Same tag, different value: not equal or similar
			a1: Attributes{
				MakeAttr("attr1", TagInteger, Integer(0)),
			},
			a2: Attributes{
				MakeAttr("attr1", TagInteger, Integer(1)),
			},
			equal:   false,
			similar: false,
		},
		{
			// Same value, tag value: not equal or similar
			a1: Attributes{
				MakeAttr("attr1", TagInteger, Integer(0)),
			},
			a2: Attributes{
				MakeAttr("attr1", TagEnum, Integer(0)),
			},
			equal:   false,
			similar: false,
		},

		{
			// Different but similar Value types:
			// Attributes are not equal but similar
			a1: Attributes{
				MakeAttr("attr1", TagString, Binary("hello")),
				MakeAttr("attr2", TagString, String("world")),
			},
			a2: Attributes{
				MakeAttr("attr1", TagString, String("hello")),
				MakeAttr("attr2", TagString, Binary("world")),
			},
			equal:   false,
			similar: true,
		},

		{
			// Different order: not equal but similar
			a1: Attributes{
				MakeAttr("attr1", TagString, String("hello")),
				MakeAttr("attr2", TagString, String("world")),
			},
			a2: Attributes{
				MakeAttr("attr2", TagString, String("world")),
				MakeAttr("attr1", TagString, String("hello")),
			},
			equal:   false,
			similar: true,
		},
	}

	for _, test := range tests {
		equal := test.a1.Equal(test.a2)
		similar := test.a1.Similar(test.a2)

		if equal != test.equal {
			t.Errorf("testing Attributes.Equal:\n"+
				"attrs 1:   %s\n"+
				"attrs 2:   %s\n"+
				"expected:  %v\n"+
				"present:   %v\n",
				test.a1, test.a2,
				test.equal, equal,
			)
		}

		if similar != test.similar {
			t.Errorf("testing Attributes.Similar:\n"+
				"attrs 1:  %s\n"+
				"attrs 2:  %s\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.a1, test.a2,
				test.similar, similar,
			)
		}
	}
}

// TestAttributesConstructors tests Attributes.Add and MakeAttr
func TestAttributesConstructors(t *testing.T) {
	attrs1 := Attributes{
		Attribute{
			Name: "attr1",
			Values: Values{
				{TagString, String("hello")},
				{TagString, String("world")},
			},
		},
		Attribute{
			Name: "attr2",
			Values: Values{
				{TagInteger, Integer(1)},
				{TagInteger, Integer(2)},
				{TagInteger, Integer(3)},
			},
		},
	}

	attrs2 := Attributes{}
	attrs2.Add(MakeAttr("attr1", TagString, String("hello"), String("world")))
	attrs2.Add(MakeAttr("attr2", TagInteger, Integer(1), Integer(2), Integer(3)))

	if !attrs1.Equal(attrs2) {
		t.Errorf("Attributes constructors test failed")
	}
}

// TestMakeAttribute tests MakeAttribute function
func TestMakeAttribute(t *testing.T) {
	a1 := Attribute{
		Name:   "attr",
		Values: Values{{TagInteger, Integer(1)}},
	}

	a2 := MakeAttribute("attr", TagInteger, Integer(1))

	if !a1.Equal(a2) {
		t.Errorf("MakeAttribute test failed")
	}
}

// TestAttributesCopy tests Attributes.Clone and Attributes.DeepCopy
func TestAttributesCopy(t *testing.T) {
	type testData struct {
		attrs Attributes
	}

	tests := []testData{
		{nil},
		{Attributes{}},
		{
			Attributes{
				MakeAttr("attr1", TagString, String("hello"), String("world")),
				MakeAttr("attr2", TagInteger, Integer(1), Integer(2), Integer(3)),
				MakeAttr("attr2", TagBoolean, Boolean(true), Boolean(false)),
			},
		},
	}

	for _, test := range tests {
		clone := test.attrs.Clone()

		if !test.attrs.Equal(clone) {
			t.Errorf("testing Attributes.Clone\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				test.attrs,
				clone,
			)
		}

		copy := test.attrs.DeepCopy()
		if !test.attrs.Equal(copy) {
			t.Errorf("testing Attributes.DeepCopy\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				test.attrs,
				copy,
			)
		}
	}
}

// TestAttributeUnpack tests Attribute.unpack method for all kinds
// of Value types
func TestAttributeUnpack(t *testing.T) {
	loc := time.FixedZone("UTC+3:30", 3*3600+1800)
	tm, _ := time.ParseInLocation(
		time.RFC3339, "2025-03-29T16:48:53+03:30", loc)

	values := Values{
		{TagBoolean, Boolean(true)},
		{TagExtension, Binary{}},
		{TagString, Binary{1, 2, 3}},
		{TagInteger, Integer(123)},
		{TagEnum, Integer(-321)},
		{TagRange, Range{-100, 200}},
		{TagRange, Range{-100, -50}},
		{TagResolution, Resolution{150, 300, UnitsDpi}},
		{TagResolution, Resolution{100, 200, UnitsDpcm}},
		{TagResolution, Resolution{75, 150, 10}},
		{TagName, String("hello")},
		{TagTextLang, TextWithLang{"en-US", "hello"}},
		{TagDateTime, Time{tm}},
		{TagNoValue, Void{}},
	}

	for _, v := range values {
		expected := Attribute{Name: "attr", Values: Values{v}}
		present := Attribute{Name: "attr"}
		data, _ := v.V.encode()
		present.unpack(v.T, data)

		if !expected.Equal(present) {
			t.Errorf("%x %d", data, unsafe.Sizeof(int(5)))
			t.Errorf("testing Attribute.unpack:\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				expected, present,
			)
		}
	}
}

// TestAttributeUnpackError tests that Attribute.unpack properly
// handles errors from the Value.decode
func TestAttributeUnpackError(t *testing.T) {
	noError := errors.New("")

	type testData struct {
		t    Tag    // Input value tag
		data []byte // Input data
		err  string // Expected error
	}

	tests := []testData{
		{
			t:    TagInteger,
			data: []byte{},
			err:  "integer: value must be 4 bytes",
		},

		{
			t:    TagBoolean,
			data: []byte{},
			err:  "boolean: value must be 1 byte",
		},
	}

	for _, test := range tests {
		attr := Attribute{Name: "attr"}
		err := attr.unpack(test.t, test.data)
		if err == nil {
			err = noError
		}

		if err.Error() != test.err {
			t.Errorf("testing Attribute.unpack:\n"+
				"input tag:      %s\n"+
				"input data:     %x\n"+
				"error expected: %s\n"+
				"error present:  %s\n",
				test.t, test.data,
				test.err, err,
			)

		}
	}
}

// TestAttributeUnpackPanic tests that Attribute.unpack panics
// on invalid Tag
func TestAttributeUnpackPanic(t *testing.T) {
	defer func() {
		recover()
	}()

	attr := Attribute{Name: "attr"}
	attr.unpack(TagOperationGroup, []byte{})

	t.Errorf("Attribute.unpack must panic on the invalid Tag")
}
