/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP protocol messages
 */

package goipp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Type Code represents Op(operation) or Status codes
type Code uint16

// Type Version represents a protocol version. It consist
// of Major and Minor version codes, packed into a single
// 16-bit word
type Version uint16

// Make version
func MakeVersion(major, minor uint8) Version {
	return Version(major)<<16 | Version(minor)
}

// Get Major part of version
func (v Version) Major() uint8 {
	return uint8(v >> 8)
}

// Get Minor part of version
func (v Version) Minor() uint8 {
	return uint8(v)
}

// Convert version to string (i.e., 2.0)
func (v Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

// Type Message represents a single IPP message, which may be either
// client request or server response
type Message struct {
	// Common header
	Version   Version // Protocol version
	Code      Code    // Operation for request, status for response
	RequestId uint32  // Set in request, returned in response

	// Attributes, by group
	Operation         []Attribute // Operation attributes
	Job               []Attribute // Job attributes
	Printer           []Attribute // Printer attributes
	Unsupported       []Attribute // Unsupported attributes
	Subscription      []Attribute // Subscription attributes
	EventNotification []Attribute // Event Notification attributes
	Resource          []Attribute // Resource attributes
	Document          []Attribute // Document attributes
	System            []Attribute // System attributes
	Future11          []Attribute // \
	Future12          []Attribute //  \
	Future13          []Attribute //   | Reserved for future extensions
	Future14          []Attribute //  /
	Future15          []Attribute // /
}

// Pretty-print the message. Request parameter affects interpretation
// of Message.Code: it is interpreted either as Op or as Status
func (m *Message) Print(out io.Writer, request bool) {
	out.Write([]byte("{\n"))

	fmt.Fprintf(out, "\tVERSION %s\n", m.Version)

	if request {
		fmt.Fprintf(out, "\tOPERATION %s\n", Op(m.Code))
	} else {
		fmt.Fprintf(out, "\tSTATUS %s\n", Status(m.Code))
	}

	groups := []struct {
		tag   Tag
		attrs []Attribute
	}{
		{TagOperationGroup, m.Operation},
		{TagJobGroup, m.Job},
		{TagPrinterGroup, m.Printer},
		{TagUnsupportedGroup, m.Unsupported},
		{TagSubscriptionGroup, m.Subscription},
		{TagEventNotificationGroup, m.EventNotification},
		{TagResourceGroup, m.Resource},
		{TagDocumentGroup, m.Document},
		{TagSystemGroup, m.System},
		{TagFuture11Group, m.Future11},
		{TagFuture12Group, m.Future12},
		{TagFuture13Group, m.Future13},
		{TagFuture14Group, m.Future14},
		{TagFuture15Group, m.Future15},
	}

	for _, grp := range groups {
		if grp.attrs != nil {
			fmt.Fprintf(out, "\tGROUP %s\n", grp.tag)
			for _, attr := range grp.attrs {
				tag := attr.Values[0].T
				fmt.Fprintf(out, "\tATTR %s %s", tag, attr.Name)
				for _, val := range attr.Values {
					if val.T != tag {
						fmt.Fprintf(out, " %s", tag)
						tag = val.T
					}
					fmt.Fprintf(out, " %s", val.V)
				}
				out.Write([]byte("\n"))
			}
		}
	}

	out.Write([]byte("}\n"))
}

// Decode the message
func (m *Message) Decode(in io.Reader) error {
	md := messageDecoder{
		in: in,
	}

	return md.decode(m)
}

// Type messageDecoder represents Message decoder
type messageDecoder struct {
	in  io.Reader // Input stream
	off int       // Offset of last tag
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
		m.RequestId, err = md.decodeU32()
	}

	// Now parse attributes
	done := false
	var group *[]Attribute
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
			err = md.error("Invalid tag 0")
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
			attr, err = md.decodeAttribute(tag)

			switch {
			case err != nil:
			case attr.Name == "":
				if prev != nil {
					prev.AddValue(attr.Values[0].T, attr.Values[0].V)
				} else {
					err = md.error("Additional value without preceding attribute")
				}
			case group != nil:
				*group = append(*group, attr)
				prev = &(*group)[len(*group)-1]
			default:
				err = md.error("Attribute without a group")
			}
		}
	}

	return err
}

// Decode a tag
func (md *messageDecoder) decodeTag() (Tag, error) {
	md.off = md.cnt
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
	if err == nil {
		value, err = md.decodeBytes()
	}

	// Unpack value
	if err == nil {
		err = attr.unpack(tag, value)
	}

	if err != nil {
		return Attribute{}, err
	}

	return attr, nil
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
	} else {
		return string(data), nil
	}
}

// Read a piece of raw data from input stream
func (md *messageDecoder) read(data []byte) error {
	for len(data) > 0 {
		n, err := md.in.Read(data)
		if err != nil {
			return err
		} else {
			md.cnt += n
			data = data[n:]
		}
	}

	return nil
}

// Create an error
func (md *messageDecoder) error(format string, args ...interface{}) error {
	s := fmt.Sprintf(format, args...)
	s += fmt.Sprintf(" at 0x%x", md.off)
	return errors.New(s)
}
