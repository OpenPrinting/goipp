/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Message attributes
 */

package goipp

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Attribute represents a single attribute
type Attribute struct {
	Name   string // Attribute name
	Values Values // Slice of values
}

// Attributes represents a slice of attributes
type Attributes []Attribute

// Add Attribute to Attributes
func (attrs *Attributes) Add(attr Attribute) {
	*attrs = append(*attrs, attr)
}

// AddValue adds value to attribute's values
func (a *Attribute) AddValue(tag Tag, val Value) {
	a.Values.Add(tag, val)
}

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
		return a.unpackBinary(tag, value)

	case TagText, TagName, TagReservedString, TagKeyword, TagURI, TagURIScheme,
		TagCharset, TagLanguage, TagMimeType, TagMemberName:
		return a.unpackString(tag, value)

	case TagDate:
		return a.unpackDate(tag, value)

	case TagResolution:
		return a.unpackResolution(tag, value)

	case TagRange:
		return a.unpackRange(tag, value)

	default:
		return a.unpackBinary(tag, value)
	}

}

// Unpack Integer value
func (a *Attribute) unpackInteger(tag Tag, value []byte) error {
	if len(value) != 4 {
		return fmt.Errorf("Value of %s tag must be 4 bytes", tag)
	}

	a.AddValue(tag, Integer(binary.BigEndian.Uint32(value)))
	return nil
}

// Unpack Boolean value
func (a *Attribute) unpackBoolean(tag Tag, value []byte) error {
	if len(value) != 1 {
		return fmt.Errorf("Value of %s tag must be 1 byte", tag)
	}

	a.AddValue(tag, Boolean(value[0] != 0))
	return nil
}

// Unpack String value
func (a *Attribute) unpackString(tag Tag, value []byte) error {
	a.AddValue(tag, String(value))
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
		name := fmt.Sprintf("UTC%c%d", value[8], value[9])
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

	a.AddValue(tag, Time{t})
	return nil
}

// Unpack Resolution value
func (a *Attribute) unpackResolution(tag Tag, value []byte) error {
	if len(value) != 9 {
		return fmt.Errorf("Value of %s tag must be 9 bytes", tag)
	}

	val := Resolution{
		Xres:  int(binary.BigEndian.Uint32(value[0:4])),
		Yres:  int(binary.BigEndian.Uint32(value[4:8])),
		Units: Units(value[9]),
	}

	a.AddValue(tag, val)
	return nil
}

// Unpack Range value
func (a *Attribute) unpackRange(tag Tag, value []byte) error {
	if len(value) != 8 {
		return fmt.Errorf("Value of %s tag must be 8 bytes", tag)
	}

	val := Range{
		Lower: int(binary.BigEndian.Uint32(value[0:4])),
		Upper: int(binary.BigEndian.Uint32(value[4:8])),
	}

	a.AddValue(tag, val)
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

	// Add a value
	a.AddValue(tag, StringWithLang{Lang: lang, Text: text})
	return nil

ERROR:
	return fmt.Errorf("Value of %s tag has invalid format", tag)
}

// Unpack Binary value
func (a *Attribute) unpackBinary(tag Tag, value []byte) error {
	a.AddValue(tag, Binary(value))
	return nil
}
