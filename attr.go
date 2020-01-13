/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Message attributes
 */

package goipp

import (
	"fmt"
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
	//var decoder valueDecoder
	var err error
	var val Value

	switch tag.Type() {
	case TypeVoid, TypeCollection:
		val = Void{}

	case TypeInteger:
		val = Integer(0)

	case TypeBoolean:
		val = Boolean(false)

	case TypeString:
		val = String("")

	case TypeDateTime:
		val = Time{}

	case TypeResolution:
		val = Resolution{}

	case TypeRange:
		val = Range{}

	case TypeTextWithLang:
		val = TextWithLang{}

	case TypeBinary:
		val = Binary(nil)

	default:
		panic(fmt.Sprintf("(Attribute) uppack(): tag=%s type=%s", tag, tag.Type()))
	}

	val, err = val.decode(value)

	if err == nil {
		a.Values.Add(tag, val)
	} else {
		err = fmt.Errorf("%s: %s", tag, err)
	}

	return err
}
