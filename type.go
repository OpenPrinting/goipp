/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Enumeration of value types
 */

package goipp

import (
	"fmt"
)

// Type enumerates all possible value types
type Type int

const (
	TypeInvalid Type = -1
	TypeVoid    Type = iota
	TypeInteger
	TypeBoolean
	TypeString
	TypeDateTime
	TypeResolution
	TypeRange
	TypeTextWithLang
	TypeBinary
	TypeCollection
)

// String converts Type to string, for debugging
func (t Type) String() string {
	switch t {
	case TypeInvalid:
		return "Invalid"
	case TypeVoid:
		return "Void"
	case TypeInteger:
		return "Integer"
	case TypeBoolean:
		return "Boolean"
	case TypeString:
		return "String"
	case TypeDateTime:
		return "DateTime"
	case TypeResolution:
		return "Resolution"
	case TypeRange:
		return "Range"
	case TypeTextWithLang:
		return "TextWithLang"
	case TypeBinary:
		return "Binary"
	case TypeCollection:
		return "Collection"
	}

	return fmt.Sprintf("Unknown type %d", int(t))
}
