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

// Attributes represents a slice of attributes
type Attributes []Attribute

// Add Attribute to Attributes
func (attrs *Attributes) Add(attr Attribute) {
	*attrs = append(*attrs, attr)
}

// Attribute represents a single attribute
type Attribute struct {
	Name   string // Attribute name
	Values Values // Slice of values
}

// Make Attribute with single value
func MakeAttribute(name string, tag Tag, value Value) Attribute {
	attr := Attribute{Name: name}
	attr.Values.Add(tag, value)
	return attr
}

// Unpack attribute value
func (a *Attribute) unpack(tag Tag, value []byte) error {
	switch tag.Type() {
	case TypeVoid, TypeCollection:
		return a.unpackVoid(tag, value)

	case TypeInteger:
		return a.unpackInteger(tag, value)

	case TypeBoolean:
		return a.unpackBoolean(tag, value)

	case TypeString:
		return a.unpackString(tag, value)

	case TypeDateTime:
		return a.unpackDate(tag, value)

	case TypeResolution:
		return a.unpackResolution(tag, value)

	case TypeRange:
		return a.unpackRange(tag, value)

	case TypeTextWithLang:
		return a.unpackTextWithLang(tag, value)

	case TypeBinary:
		return a.unpackBinary(tag, value)
	}

	panic(fmt.Sprintf("(Attribute) uppack(): tag=%s type=%s", tag, tag.Type()))
}

// Unpack Void value
func (a *Attribute) unpackVoid(tag Tag, value []byte) error {
	a.Values.Add(tag, Void{})
	return nil
}

// Unpack Integer value
func (a *Attribute) unpackInteger(tag Tag, value []byte) error {
	if len(value) != 4 {
		return fmt.Errorf("Value of %s tag must be 4 bytes", tag)
	}

	a.Values.Add(tag, Integer(binary.BigEndian.Uint32(value)))
	return nil
}

// Unpack Boolean value
func (a *Attribute) unpackBoolean(tag Tag, value []byte) error {
	if len(value) != 1 {
		return fmt.Errorf("Value of %s tag must be 1 byte", tag)
	}

	a.Values.Add(tag, Boolean(value[0] != 0))
	return nil
}

// Unpack String value
func (a *Attribute) unpackString(tag Tag, value []byte) error {
	a.Values.Add(tag, String(value))
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

	a.Values.Add(tag, Time{t})
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

	a.Values.Add(tag, val)
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

	a.Values.Add(tag, val)
	return nil
}

// Unpack TextWithLang value
func (a *Attribute) unpackTextWithLang(tag Tag, value []byte) error {
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
	a.Values.Add(tag, TextWithLang{Lang: lang, Text: text})
	return nil

ERROR:
	return fmt.Errorf("Value of %s tag has invalid format", tag)
}

// Unpack Binary value
func (a *Attribute) unpackBinary(tag Tag, value []byte) error {
	a.Values.Add(tag, Binary(value))
	return nil
}
