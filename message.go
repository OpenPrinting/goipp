package main

import (
	"encoding/binary"
	"errors"
	"io"
)

// Type Version represents a protocol version
type Version struct {
	Major, Minor uint8
}

// Type Message represents a single IPP message, which may be either
// client request or server response
type Message struct {
	// Common header
	Version   Version // Protocol version
	Code      uint16  // Operation for request, status for response
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

// Decode the message
func (m *Message) Decode(in io.Reader) error {
	/*
	   1 byte:   VersionMajor
	   1 byte:   VersionMinor
	   2 bytes:  operation-id or status-code
	   variable: attributes
	   1 byte:   end-of-attributes-tag
	*/

	// Parse message header
	var err error
	m.Version.Major, err = m.decodeU8(in)
	if err == nil {
		m.Version.Minor, err = m.decodeU8(in)
	}
	if err == nil {
		m.Code, err = m.decodeU16(in)
	}
	if err == nil {
		m.RequestId, err = m.decodeU32(in)
	}

	// Now parse attributes
	done := false
	var group *[]Attribute
	var attr Attribute

	for err == nil && !done {
		var tag Tag
		tag, err = m.decodeTag(in)

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
			attr, err = m.decodeAttribute(in, tag)
			if err == nil {
				if group != nil {
					*group = append(*group, attr)
				} else {
					err = errors.New("Attribute without a group")
				}
			}
		}
	}

	return err
}

// Decode a tag
func (m *Message) decodeTag(in io.Reader) (Tag, error) {
	t, err := m.decodeU8(in)
	return Tag(t), err
}

// Decode a single attribute
func (m *Message) decodeAttribute(in io.Reader, tag Tag) (Attribute, error) {
	var attr Attribute
	var value []byte
	var err error

	// Obtain attribute name and raw value
	attr.Name, err = m.decodeString(in)
	if err == nil {
		value, err = m.decodeBytes(in)
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
func (m *Message) decodeU8(in io.Reader) (uint8, error) {
	buf := make([]byte, 1)
	err := m.decodeRaw(in, buf)
	return buf[0], err
}

// Decode a 16-bit integer
func (m *Message) decodeU16(in io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	err := m.decodeRaw(in, buf)
	return binary.BigEndian.Uint16(buf[:]), err
}

// Decode a 32-bit integer
func (m *Message) decodeU32(in io.Reader) (uint32, error) {
	buf := make([]byte, 4)
	err := m.decodeRaw(in, buf)
	return binary.BigEndian.Uint32(buf[:]), err
}

// Decode sequence of bytes
func (m *Message) decodeBytes(in io.Reader) ([]byte, error) {
	length, err := m.decodeU16(in)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	err = m.decodeRaw(in, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Decode string
func (m *Message) decodeString(in io.Reader) (string, error) {
	data, err := m.decodeBytes(in)
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// Read a piece of raw data from input stream
func (m *Message) decodeRaw(in io.Reader, data []byte) error {
	for len(data) > 0 {
		n, err := in.Read(data)
		if err != nil {
			return err
		} else {
			data = data[n:]
		}
	}

	return nil
}
