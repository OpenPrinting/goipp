/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Values for message attributes
 */

package goipp

import (
	"bytes"
	"fmt"
	"time"
)

// Values represents a slice of Attribute values with tags
type Values []struct {
	T Tag   // The tag
	V Value // The value
}

// Add value to Values
func (values *Values) Add(t Tag, v Value) {
	*values = append(*values, struct {
		T Tag
		V Value
	}{t, v})
}

// String() converts Values to string
func (values Values) String() string {
	if len(values) == 1 {
		return values[0].V.String()
	}

	var buf bytes.Buffer
	buf.Write([]byte("["))
	for i, v := range values {
		if i != 0 {
			buf.Write([]byte(","))
		}
		buf.Write([]byte(v.V.String()))
	}
	buf.Write([]byte("]"))

	return buf.String()
}

// Value represents an attribute value
type Value interface {
	String() string
	Type() Type
	isValue()
}

// Void represents "no value"
type Void struct{}

func (Void) isValue() {}

// String() converts Void Value to string
func (Void) String() string { return "" }

// Type() returns type of Value
func (Void) Type() Type { return TypeVoid }

// Integer represents an Integer Value
type Integer uint32

func (Integer) isValue() {}

// String() converts Integer value to string
func (v Integer) String() string { return fmt.Sprintf("%d", uint32(v)) }

// Type() returns type of Value
func (Integer) Type() Type { return TypeInteger }

// Boolean represents a boolean Value
type Boolean bool

func (Boolean) isValue() {}

// String() converts Boolean value to string
func (v Boolean) String() string { return fmt.Sprintf("%t", bool(v)) }

// Type() returns type of Value
func (Boolean) Type() Type { return TypeBoolean }

// String represents a string Value
type String string

func (String) isValue() {}

// String() converts String value to string
func (v String) String() string { return string(v) }

// Type() returns type of Value
func (String) Type() Type { return TypeString }

// Time represents a DateTime Value
type Time struct{ time.Time }

func (Time) isValue() {}

// String() converts Time value to string
func (v Time) String() string { return v.Time.Format(time.RFC3339) }

// Type() returns type of Value
func (Time) Type() Type { return TypeDateTime }

// Resolution represents a resolution Value
type Resolution struct {
	Xres, Yres int   // X/Y resolutions
	Units      Units // Resolution units
}

func (Resolution) isValue() {}

// String() converts Resolution value to string
func (v Resolution) String() string {
	return fmt.Sprintf("%dx%d%s", v.Xres, v.Yres, v.Units)
}

// Type() returns type of Value
func (Resolution) Type() Type { return TypeResolution }

// Units represents resolution units
type Units uint8

const (
	UnitsDpi  Units = 3 // Dots per inch
	UnitsDpcm Units = 4 // Dots per cm
)

// String() converts Units to string
func (u Units) String() string {
	switch u {
	case UnitsDpi:
		return "dpi"
	case UnitsDpcm:
		return "dpcm"
	default:
		return fmt.Sprintf("0x%2.2x", uint8(u))
	}
}

// Range represents a range of integers Value
type Range struct {
	Lower, Upper int // Lower/upper bounds
}

func (Range) isValue() {}

// String() converts Range value to string
func (v Range) String() string {
	return fmt.Sprintf("%d-%d", v.Lower, v.Upper)
}

// Type() returns type of Value
func (Range) Type() Type { return TypeRange }

// TextWithLang represents a combination of two strings:
// one is a name of natural language and second is a text
// on this language
type TextWithLang struct {
	Lang, Text string // Language and text
}

func (TextWithLang) isValue() {}

// String() converts TextWithLang value to string
func (v TextWithLang) String() string { return v.Text + " [" + v.Lang + "]" }

// Type() returns type of Value
func (TextWithLang) Type() Type { return TypeTextWithLang }

// Binary represents a raw binary Value
type Binary []byte

func (Binary) isValue() {}

// String() converts Range value to string
func (v Binary) String() string {
	return fmt.Sprintf("%x", []byte(v))
}

// Type() returns type of Value
func (Binary) Type() Type { return TypeBinary }

// Collection represents a collection of attributes
type Collection []Attribute

func (Collection) isValue() {}

// String() converts Collection to string
func (v Collection) String() string {
	var buf bytes.Buffer
	buf.Write([]byte("{"))
	for i, attr := range v {
		if i > 0 {
			buf.Write([]byte(" "))
		}
		fmt.Fprintf(&buf, "%s=%s", attr.Name, attr.Values)
	}
	buf.Write([]byte("}"))

	return buf.String()
}

// Type() returns type of Value
func (Collection) Type() Type { return TypeCollection }
