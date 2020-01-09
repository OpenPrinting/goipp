package main

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Type Attribute represents a single attribute
type Attribute struct {
	Name  string
	Value Value
}

// Type Value represents an attribute value
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

// Unpack attribute value
func (a *Attribute) unpack(tag Tag, value []byte) error {
	switch tag {
	case TagInteger, TagEnum:
		return a.unpackInteger(tag, value)

	case TagBoolean:
		return a.unpackBoolean(tag, value)

	case TagUnsupportedValue, TagDefault, TagUnknown, TagNotSettable,
		TagDeleteAttr, TagAdminDefine:
		// These tags not expected to have value
		return nil

	case TagText, TagName, TagReservedString, TagKeyword, TagUri, TagUriScheme,
		TagCharset, TagLanguage, TagMimeType:
		return a.unpackString(tag, value)

	case TagDate:
		return a.unpackDate(tag, value)

	case TagResolution:
		return a.unpackResolution(tag, value)

	case TagRange:
		return a.unpackRange(tag, value)
	}

	return nil
}

// Unpack Integer value
func (a *Attribute) unpackInteger(tag Tag, value []byte) error {
	if len(value) != 4 {
		return fmt.Errorf("Value of %s tag must be 4 bytes", tag)
	}

	a.Value = Integer(binary.BigEndian.Uint32(value))
	return nil
}

// Unpack Boolean value
func (a *Attribute) unpackBoolean(tag Tag, value []byte) error {
	if len(value) != 1 {
		return fmt.Errorf("Value of %s tag must be 1 byte", tag)
	}

	a.Value = Boolean(value[0] != 0)
	return nil
}

// Unpack String value
func (a *Attribute) unpackString(tag Tag, value []byte) error {
	a.Value = String(value)
	return nil
}

// Unpack Time value
func (a *Attribute) unpackDate(tag Tag, value []byte) error {
	if len(value) != 9 && len(value) != 11 {
		return fmt.Errorf("Value of %s tag must be 9 or 11 bytes", tag)
	}

	/*
		From RFC2579:

		    field  octets  contents                  range
		    -----  ------  --------                  -----
		      1      1-2   year*                     0..65536
		      2       3    month                     1..12
		      3       4    day                       1..31
		      4       5    hour                      0..23
		      5       6    minutes                   0..59
		      6       7    seconds                   0..60
				   (use 60 for leap-second)
		      7       8    deci-seconds              0..9
		      8       9    direction from UTC        '+' / '-'
		      9      10    hours from UTC*           0..13
		     10      11    minutes from UTC          0..59

		    * Notes:
		    - the value of year is in network-byte order
		    - daylight saving time in New Zealand is +13
	*/

	var l *time.Location
	switch {
	case len(value) == 9:
		l = time.UTC
	case value[8] == '+', value[8] == '-':
		name := fmt.Sprintf("UTC%c%d", value[9])
		if value[10] != 0 {
			name += fmt.Sprintf(":%d", value[10])
		}

		off := 3600*int(value[9]) + 60*int(value[10])
		if value[8] == '-' {
			off = -off
		}

		l = time.FixedZone(name, off)

	default:
		return fmt.Errorf("Invalid format of %s value", tag)
	}

	t := time.Date(
		int(binary.BigEndian.Uint16(value[0:2])), // year
		time.Month(value[2]),                     // month
		int(value[3]),                            // day
		int(value[4]),                            // hour
		int(value[5]),                            // min
		int(value[6]),                            // sec
		int(value[6])*100000000,                  // nsec
		l,                                        // FIXME
	)

	a.Value = Time{t}
	return nil
}

// Unpack Resolution value
func (a *Attribute) unpackResolution(tag Tag, value []byte) error {
	if len(value) != 9 {
		return fmt.Errorf("Value of %s tag must be 9 bytes", tag)
	}

	a.Value = Resolution{
		Xres:  int(binary.BigEndian.Uint32(value[0:4])),
		Yres:  int(binary.BigEndian.Uint32(value[4:8])),
		Units: value[9],
	}

	return nil
}

// Unpack Range value
func (a *Attribute) unpackRange(tag Tag, value []byte) error {
	if len(value) != 8 {
		return fmt.Errorf("Value of %s tag must be 8 bytes", tag)
	}

	a.Value = Range{
		Lower: int(binary.BigEndian.Uint32(value[0:4])),
		Upper: int(binary.BigEndian.Uint32(value[4:8])),
	}

	return nil
}

// Unpack StringWithLang value
func (a *Attribute) unpackStringWithLang(tag Tag, value []byte) error {
	var langLen, textLen int
	var lang, text string

	// Unpack language length
	if len(value) < 2 {
		goto ERROR
	}

	langLen = int(binary.BigEndian.Uint16(value[0:2]))
	value = value[2:]

	// Unpack language value
	if len(value) < langLen {
		goto ERROR
	}

	lang = string(value[:langLen])
	value = value[langLen:]

	// Unpack text length
	if len(value) < 2 {
		goto ERROR
	}

	textLen = int(binary.BigEndian.Uint16(value[0:2]))
	value = value[2:]

	// Unpack text value
	if len(value) < textLen {
		goto ERROR
	}

	text = string(value[:textLen])
	value = value[textLen:]

	// We must have consumed all bytes at this point
	if len(value) != 0 {
		goto ERROR
	}

	// Construct a value
	a.Value = StringWithLang{Lang: lang, Text: text}
	return nil

ERROR:
	return fmt.Errorf("Value of %s tag has invalid format", tag)
}
