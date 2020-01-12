/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Message decoder
 */

package goipp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Type messageDecoder represents Message decoder
type messageDecoder struct {
	in  io.Reader // Input stream
	off int       // Offset of last read
	cnt int       // Count of read bytes
}

// Decode the message
func (md *messageDecoder) decode(m *Message) error {
	/*
	   1 byte:   VersionMajor
	   1 byte:   VersionMinor
	   2 bytes:  operation-id or status-code
	   variable: attributes
	   1 byte:   end-of-attributes-tag
	*/

	// Parse message header
	var err error
	m.Version, err = md.decodeVersion()
	if err == nil {
		m.Code, err = md.decodeCode()
	}
	if err == nil {
		m.RequestID, err = md.decodeU32()
	}

	// Now parse attributes
	done := false
	var group *Attributes
	var attr Attribute
	var prev *Attribute

	for err == nil && !done {
		var tag Tag
		tag, err = md.decodeTag()

		if tag.IsDelimiter() {
			prev = nil
		}

		switch tag {
		case TagZero:
			err = errors.New("Invalid tag 0")
		case TagEnd:
			done = true

		case TagOperationGroup:
			group = &m.Operation
		case TagJobGroup:
			group = &m.Job
		case TagPrinterGroup:
			group = &m.Printer
		case TagUnsupportedGroup:
			group = &m.Unsupported
		case TagSubscriptionGroup:
			group = &m.Subscription
		case TagEventNotificationGroup:
			group = &m.EventNotification
		case TagResourceGroup:
			group = &m.Resource
		case TagDocumentGroup:
			group = &m.Document
		case TagSystemGroup:
			group = &m.System
		case TagFuture11Group:
			group = &m.Future11
		case TagFuture12Group:
			group = &m.Future12
		case TagFuture13Group:
			group = &m.Future13
		case TagFuture14Group:
			group = &m.Future14
		case TagFuture15Group:
			group = &m.Future15

		default:
			// Decode attribute
			if tag == TagMemberName || tag == TagEndCollection {
				err = fmt.Errorf("Unexpected tag %s", tag)
			} else {
				attr, err = md.decodeAttribute(tag)
			}

			if err == nil && tag == TagBeginCollection {
				attr.Values[0].V, err = md.decodeCollection()
			}

			// If everything is OK, save attribute
			switch {
			case err != nil:
			case attr.Name == "":
				if prev != nil {
					prev.Values.Add(attr.Values[0].T, attr.Values[0].V)
				} else {
					err = errors.New("Additional value without preceding attribute")
				}
			case group != nil:
				group.Add(attr)
				prev = &(*group)[len(*group)-1]
			default:
				err = errors.New("Attribute without a group")
			}
		}
	}

	if err != nil {
		err = fmt.Errorf("%s at 0x%x", err, md.off)
	}

	return err
}

// Decode a Collection
func (md *messageDecoder) decodeCollection() (Collection, error) {
	collection := make(Collection, 0)

	for {
		// Decode next TagEndCollection or next TagMemberName
		tag, err := md.decodeTag()
		if err != nil {
			return nil, err
		}

		if tag != TagEndCollection && tag != TagMemberName {
			err = fmt.Errorf(
				"Collection: expected %s or %s, got %s",
				TagMemberName, TagEndCollection, tag)
			return nil, err
		}

		attrName, err := md.decodeAttribute(tag)
		if err != nil {
			return nil, err
		}

		if tag == TagEndCollection {
			return collection, nil
		}

		// Decode member value
		tag, err = md.decodeTag()
		if err != nil {
			return nil, err
		}

		if tag.IsDelimiter() ||
			tag == TagEndCollection || tag == TagMemberName {
			err = fmt.Errorf("Collection: unexpected %s", tag)
			return nil, err
		}

		attr, err := md.decodeAttribute(tag)
		if err != nil {
			return nil, err
		}

		attr.Name = string(attrName.Values[0].V.(String))
		if err == nil && tag == TagBeginCollection {
			attr.Values[0].V, err = md.decodeCollection()
		}

		if err != nil {
			return nil, err
		}

		collection = append(collection, attr)
	}
}

// Decode a tag
func (md *messageDecoder) decodeTag() (Tag, error) {
	t, err := md.decodeU8()
	return Tag(t), err
}

// Decode a Version
func (md *messageDecoder) decodeVersion() (Version, error) {
	code, err := md.decodeU16()
	return Version(code), err
}

// Decode a Code
func (md *messageDecoder) decodeCode() (Code, error) {
	code, err := md.decodeU16()
	return Code(code), err
}

// Decode a single attribute
func (md *messageDecoder) decodeAttribute(tag Tag) (Attribute, error) {
	var attr Attribute
	var value []byte
	var err error

	// Obtain attribute name and raw value
	attr.Name, err = md.decodeString()
	if err != nil {
		goto ERROR
	}

	value, err = md.decodeBytes()
	if err != nil {
		goto ERROR
	}

	// Handle TagExtension
	if tag == TagExtension {
		if len(value) < 4 {
			err = errors.New("Extension tag truncated")
			goto ERROR
		}

		t := binary.BigEndian.Uint32(value[:4])
		value = value[4:]

		if t > 0x7fffffff {
			err = errors.New("Extension tag out of range")
			goto ERROR
		}

		tag = Tag(t)
	}

	// Unpack value
	err = attr.unpack(tag, value)
	if err != nil {
		goto ERROR
	}

	return attr, nil

	// Return a error
ERROR:
	return Attribute{}, err
}

// Decode a 8-bit integer
func (md *messageDecoder) decodeU8() (uint8, error) {
	buf := make([]byte, 1)
	err := md.read(buf)
	return buf[0], err
}

// Decode a 16-bit integer
func (md *messageDecoder) decodeU16() (uint16, error) {
	buf := make([]byte, 2)
	err := md.read(buf)
	return binary.BigEndian.Uint16(buf[:]), err
}

// Decode a 32-bit integer
func (md *messageDecoder) decodeU32() (uint32, error) {
	buf := make([]byte, 4)
	err := md.read(buf)
	return binary.BigEndian.Uint32(buf[:]), err
}

// Decode sequence of bytes
func (md *messageDecoder) decodeBytes() ([]byte, error) {
	length, err := md.decodeU16()
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	err = md.read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Decode string
func (md *messageDecoder) decodeString() (string, error) {
	data, err := md.decodeBytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Read a piece of raw data from input stream
func (md *messageDecoder) read(data []byte) error {
	md.off = md.cnt

	for len(data) > 0 {
		n, err := md.in.Read(data)
		if err != nil {
			md.off = md.cnt
			return err
		}

		md.cnt += n
		data = data[n:]
	}

	return nil
}
