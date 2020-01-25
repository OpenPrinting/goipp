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
	"encoding/binary"
	"errors"
	"fmt"
	"math"
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

// Equal checks that two Values are equal
func (values Values) Equal(values2 Values) bool {
	if len(values) != len(values2) {
		return false
	}

	for i, v := range values {
		v2 := values2[i]
		if v.T != v2.T || !ValueEqual(v.V, v2.V) {
			return false
		}
	}

	return true
}

// Value represents an attribute value
type Value interface {
	String() string
	Type() Type
	encode() ([]byte, error)
	decode([]byte) (Value, error)
}

// ValueEqual checks if two values are equal
func ValueEqual(v1, v2 Value) bool {
	if v1.Type() != v2.Type() {
		return false
	}

	switch v1.Type() {
	case TypeDateTime:
		return v1.(Time).Equal(v2.(Time).Time)
	case TypeBinary:
		return bytes.Equal(v1.(Binary), v2.(Binary))
	case TypeCollection:
		c1 := Attributes(v1.(Collection))
		c2 := Attributes(v2.(Collection))
		return c1.Equal(c2)
	}

	return v1 == v2
}

// Void represents "no value"
//
// Use with: TagUnsupportedValue, TagDefault, TagUnknown,
// TagNotSettable, TagDeleteAttr, TagAdminDefine
type Void struct{}

// String() converts Void Value to string
func (Void) String() string { return "" }

// Type returns type of Value
func (Void) Type() Type { return TypeVoid }

// Encode Void Value into wire format
func (v Void) encode() ([]byte, error) {
	return []byte{}, nil
}

// Decode Void Value from wire format
func (Void) decode([]byte) (Value, error) {
	return Void{}, nil
}

// Integer represents an Integer Value
//
// Use with: TagInteger, TagEnum
type Integer int32

// String() converts Integer value to string
func (v Integer) String() string { return fmt.Sprintf("%d", int32(v)) }

// Type returns type of Value
func (Integer) Type() Type { return TypeInteger }

// Encode Integer Value into wire format
func (v Integer) encode() ([]byte, error) {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
}

// Decode Integer Value from wire format
func (Integer) decode(data []byte) (Value, error) {
	if len(data) != 4 {
		return nil, errors.New("value must be 4 bytes")
	}

	return Integer(binary.BigEndian.Uint32(data)), nil
}

// Boolean represents a boolean Value
//
// Use with: TagBoolean
type Boolean bool

// String() converts Boolean value to string
func (v Boolean) String() string { return fmt.Sprintf("%t", bool(v)) }

// Type returns type of Value
func (Boolean) Type() Type { return TypeBoolean }

// Encode Boolean Value into wire format
func (v Boolean) encode() ([]byte, error) {
	if v {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// Decode Boolean Value from wire format
func (Boolean) decode(data []byte) (Value, error) {
	if len(data) != 1 {
		return nil, errors.New("value must be 1 byte")
	}

	return Boolean(data[0] != 0), nil
}

// String represents a string Value
//
// Use with: TagText, TagName, TagReservedString, TagKeyword, TagURI,
// TagURIScheme, TagCharset, TagLanguage, TagMimeType, TagMemberName
type String string

// String() converts String value to string
func (v String) String() string { return string(v) }

// Type returns type of Value
func (String) Type() Type { return TypeString }

// Encode String Value into wire format
func (v String) encode() ([]byte, error) {
	return []byte(v), nil
}

// Decode String Value from wire format
func (String) decode(data []byte) (Value, error) {
	return String(data), nil
}

// Time represents a DateTime Value
//
// Use with: TagTime
type Time struct{ time.Time }

// String() converts Time value to string
func (v Time) String() string { return v.Time.Format(time.RFC3339) }

// Type returns type of Value
func (Time) Type() Type { return TypeDateTime }

// Encode Time Value into wire format
func (v Time) encode() ([]byte, error) {
	// From RFC2579:
	//
	//     field  octets  contents                  range
	//     -----  ------  --------                  -----
	//       1      1-2   year*                     0..65536
	//       2       3    month                     1..12
	//       3       4    day                       1..31
	//       4       5    hour                      0..23
	//       5       6    minutes                   0..59
	//       6       7    seconds                   0..60
	//                    (use 60 for leap-second)
	//       7       8    deci-seconds              0..9
	//       8       9    direction from UTC        '+' / '-'
	//       9      10    hours from UTC*           0..13
	//      10      11    minutes from UTC          0..59
	//
	//     * Notes:
	//     - the value of year is in network-byte order
	//     - daylight saving time in New Zealand is +13

	year := v.Year()
	_, zone := v.Zone()
	dir := byte('+')
	if zone < 0 {
		zone = -zone
		dir = '-'
	}

	return []byte{
		byte(year >> 8), byte(year),
		byte(v.Month()),
		byte(v.Day()),
		byte(v.Hour()),
		byte(v.Minute()),
		byte(v.Second()),
		byte(v.Nanosecond() / 100000000),
		dir,
		byte(zone / 3600),
		byte((zone / 60) % 60),
	}, nil
}

// Decode Time Value from wire format
func (Time) decode(data []byte) (Value, error) {
	// Check size
	if len(data) != 9 && len(data) != 11 {
		return nil, errors.New("value must be 9 or 11 bytes")
	}

	// Decode time zone
	var l *time.Location
	switch {
	case len(data) == 9:
		l = time.UTC
	case data[8] == '+', data[8] == '-':
		name := fmt.Sprintf("UTC%c%d", data[8], data[9])
		if data[10] != 0 {
			name += fmt.Sprintf(":%d", data[10])
		}

		off := 3600*int(data[9]) + 60*int(data[10])
		if data[8] == '-' {
			off = -off
		}

		l = time.FixedZone(name, off)

	default:
		return nil, errors.New("invalid data format")
	}

	// Decode time
	t := time.Date(
		int(binary.BigEndian.Uint16(data[0:2])), // year
		time.Month(data[2]),                     // month
		int(data[3]),                            // day
		int(data[4]),                            // hour
		int(data[5]),                            // min
		int(data[6]),                            // sec
		int(data[7])*100000000,                  // nsec
		l,                                       // time zone
	)

	return Time{t}, nil
}

// Resolution represents a resolution Value
//
// Use with: TagResolution
type Resolution struct {
	Xres, Yres int   // X/Y resolutions
	Units      Units // Resolution units
}

// String() converts Resolution value to string
func (v Resolution) String() string {
	return fmt.Sprintf("%dx%d%s", v.Xres, v.Yres, v.Units)
}

// Type returns type of Value
func (Resolution) Type() Type { return TypeResolution }

// Encode Resolution Value into wire format
func (v Resolution) encode() ([]byte, error) {
	// Wire format
	//    4 bytes: Xres
	//    4 bytes: Yres
	//    1 byte:  Units

	x, y := v.Xres, v.Yres

	return []byte{
		byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x),
		byte(y >> 24), byte(y >> 16), byte(y >> 8), byte(y),
		byte(v.Units),
	}, nil
}

// Decode Resolution Value from wire format
func (Resolution) decode(data []byte) (Value, error) {
	if len(data) != 9 {
		return nil, errors.New("value must be 9 bytes")
	}

	return Resolution{
		Xres:  int(binary.BigEndian.Uint32(data[0:4])),
		Yres:  int(binary.BigEndian.Uint32(data[4:8])),
		Units: Units(data[8]),
	}, nil

}

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
//
// Use with: TagRange
type Range struct {
	Lower, Upper int // Lower/upper bounds
}

// String() converts Range value to string
func (v Range) String() string {
	return fmt.Sprintf("%d-%d", v.Lower, v.Upper)
}

// Type returns type of Value
func (Range) Type() Type { return TypeRange }

// Encode Range Value into wire format
func (v Range) encode() ([]byte, error) {
	// Wire format
	//    4 bytes: Lower
	//    4 bytes: Upper

	l, u := v.Lower, v.Upper

	return []byte{
		byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l),
		byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u),
	}, nil
}

// Decode Range Value from wire format
func (Range) decode(data []byte) (Value, error) {
	if len(data) != 8 {
		return nil, errors.New("value must be 9 bytes")
	}

	return Range{
		Lower: int(binary.BigEndian.Uint32(data[0:4])),
		Upper: int(binary.BigEndian.Uint32(data[4:8])),
	}, nil
}

// TextWithLang represents a combination of two strings:
// one is a name of natural language and second is a text
// on this language
//
// Use with: TagTextLang, TagNameLang
type TextWithLang struct {
	Lang, Text string // Language and text
}

// String() converts TextWithLang value to string
func (v TextWithLang) String() string { return v.Text + " [" + v.Lang + "]" }

// Type returns type of Value
func (TextWithLang) Type() Type { return TypeTextWithLang }

// Encode TextWithLang Value into wire format
func (v TextWithLang) encode() ([]byte, error) {
	// Wire format
	//    2 bytes:  len(Lang)
	//    variable: Lang
	//    2 bytes:  len(Text)
	//    variable: Text

	lang := []byte(v.Lang)
	text := []byte(v.Text)

	if len(lang) > math.MaxUint16 {
		return nil, fmt.Errorf("Lang exceeds %d bytes", math.MaxUint16)
	}

	if len(text) > math.MaxUint16 {
		return nil, fmt.Errorf("Text exceeds %d bytes", math.MaxUint16)
	}

	data := make([]byte, 2+2+len(lang)+len(text))
	binary.BigEndian.PutUint16(data, uint16(len(lang)))
	copy(data[2:], []byte(lang))

	data2 := data[2+len(lang):]
	binary.BigEndian.PutUint16(data2, uint16(len(text)))
	copy(data2[2:], []byte(text))

	return data, nil
}

// Decode TextWithLang Value from wire format
func (TextWithLang) decode(data []byte) (Value, error) {
	var langLen, textLen int
	var lang, text string

	// Unpack language length
	if len(data) < 2 {
		goto ERROR
	}

	langLen = int(binary.BigEndian.Uint16(data[0:2]))
	data = data[2:]

	// Unpack language value
	if len(data) < langLen {
		goto ERROR
	}

	lang = string(data[:langLen])
	data = data[langLen:]

	// Unpack text length
	if len(data) < 2 {
		goto ERROR
	}

	textLen = int(binary.BigEndian.Uint16(data[0:2]))
	data = data[2:]

	// Unpack text value
	if len(data) < textLen {
		goto ERROR
	}

	text = string(data[:textLen])
	data = data[textLen:]

	// We must have consumed all bytes at this point
	if len(data) != 0 {
		goto ERROR
	}

	// Return a value
	return TextWithLang{Lang: lang, Text: text}, nil

ERROR:
	return nil, errors.New("invalid data format")
}

// Binary represents a raw binary Value
type Binary []byte

// String() converts Range value to string
func (v Binary) String() string {
	return fmt.Sprintf("%x", []byte(v))
}

// Type returns type of Value
func (Binary) Type() Type { return TypeBinary }

// Encode TextWithLang Value into wire format
func (v Binary) encode() ([]byte, error) {
	return []byte(v), nil
}

// Decode Binary Value from wire format
func (Binary) decode(data []byte) (Value, error) {
	return Binary(data), nil
}

// Collection represents a collection of attributes
//
// Use with: TagBeginCollection
type Collection Attributes

// Add Attribute to Attributes
func (collection *Collection) Add(attr Attribute) {
	*collection = append(*collection, attr)
}

// Equal checks that two collections are equal
func (c1 Collection) Equal(c2 Attributes) bool {
	return Attributes(c1).Equal(Attributes(c2))
}

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

// Type returns type of Value
func (Collection) Type() Type { return TypeCollection }

// Encode Collection Value into wire format
func (v Collection) encode() ([]byte, error) {
	// Note, TagBeginCollection attribute contains
	// no data, collection itself handled the different way
	return []byte{}, nil
}

// Decode Collection Value from wire format
func (Collection) decode(data []byte) (Value, error) {
	panic("internal error")
}
