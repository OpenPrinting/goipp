package main

import (
	"time"
)

// Type Values represents a slice of Attribute values with tags
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

// Type Value represents a attribute value
type Value interface {
	isValue()
}

// Type Integer represents an Integer value
type Integer uint32

func (Integer) isValue() {}

// Type Boolean represents a Boolean value
type Boolean bool

func (Boolean) isValue() {}

// Type Strings represents a string value
type String string

func (String) isValue() {}

// Type Time represents a DateTime value
type Time struct{ time.Time }

func (Time) isValue() {}

// Type Resolution represents a resolution value
type Resolution struct {
	Xres, Yres int   // X/Y resolutions
	Units      uint8 // Resolution units
}

func (Resolution) isValue() {}

// Type Range represents a range of integers
type Range struct {
	Lower, Upper int // Lower/upper bounds
}

func (Range) isValue() {}

// Type StringWithLang represents a combination of
// Two strings: one is a name of natural language and
// second is a text on this language
type StringWithLang struct {
	Lang, Text string // Language and text
}

func (StringWithLang) isValue() {}
