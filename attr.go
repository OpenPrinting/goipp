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
		var v Void
		val, err = v, v.decode(value)

	case TypeInteger:
		var v Integer
		val, err = v, v.decode(value)

	case TypeBoolean:
		var v Boolean
		val, err = v, v.decode(value)

	case TypeString:
		var v String
		val, err = v, v.decode(value)

	case TypeDateTime:
		var v Time
		val, err = v, v.decode(value)

	case TypeResolution:
		var v Resolution
		val, err = v, v.decode(value)

	case TypeRange:
		var v Range
		val, err = v, v.decode(value)

	case TypeTextWithLang:
		var v TextWithLang
		val, err = v, v.decode(value)

	case TypeBinary:
		var v Binary
		val, err = v, v.decode(value)

	default:
		panic(fmt.Sprintf("(Attribute) uppack(): tag=%s type=%s", tag, tag.Type()))
	}

	if err == nil {
		a.Values.Add(tag, val)
	} else {
		err = fmt.Errorf("%s: %s", tag, err)
	}

	return err
}
