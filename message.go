/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP protocol messages
 */

package goipp

import (
	"fmt"
	"io"
)

// Code represents Op(operation) or Status codes
type Code uint16

// Version represents a protocol version. It consist
// of Major and Minor version codes, packed into a single
// 16-bit word
type Version uint16

// DefaultVersion is the default IPP version
const DefaultVersion Version = 0x0200

// MakeVersion makes version from major and minor parts
func MakeVersion(major, minor uint8) Version {
	return Version(major)<<8 | Version(minor)
}

// Major returns a major part of version
func (v Version) Major() uint8 {
	return uint8(v >> 8)
}

// Minor returns a minor part of version
func (v Version) Minor() uint8 {
	return uint8(v)
}

// String() converts version to string (i.e., 2.0)
func (v Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

// Message represents a single IPP message, which may be either
// client request or server response
type Message struct {
	// Common header
	Version   Version // Protocol version
	Code      Code    // Operation for request, status for response
	RequestID uint32  // Set in request, returned in response

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

// NewRequest creates a new request message
//
// Use DefaultVersion as a first argument, if you don't
// have any specific needs
func NewRequest(v Version, op Op) *Message {
	return &Message{
		Version: v,
		Code:    Code(op),
	}
}

// NewResponse creates a new response message
//
// Use DefaultVersion as a first argument, if you don't
func NewResponse(v Version, status Status) *Message {
	return &Message{
		Version: v,
		Code:    Code(status),
	}
}

// Print pretty-prints the message. The 'request' parameter affects
// interpretation of Message.Code: it is interpreted either
// as Op or as Status
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
			fmt.Fprintf(out, "\n\tGROUP %s\n", grp.tag)
			for _, attr := range grp.attrs {
				m.printAttribute(out, attr, 1)
				out.Write([]byte("\n"))
			}
		}
	}

	out.Write([]byte("}\n"))
}

// Pretty-print an attribute. Handles Collection attributes
// recursively
func (m *Message) printAttribute(out io.Writer, attr Attribute, indent int) {
	m.printIndent(out, indent)
	fmt.Fprintf(out, "ATTR %q", attr.Name)

	tag := TagZero
	for _, val := range attr.Values {
		if val.T != tag {
			fmt.Fprintf(out, " %s:", val.T)
			tag = val.T
		}

		if collection, ok := val.V.(Collection); ok {
			out.Write([]byte(" {\n"))
			for _, attr2 := range collection {
				m.printAttribute(out, attr2, indent+1)
				out.Write([]byte("\n"))
			}
			m.printIndent(out, indent)
			out.Write([]byte("}"))
		} else {
			fmt.Fprintf(out, " %s", val.V)
		}
	}
}

// Print indentation
func (m *Message) printIndent(out io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		out.Write([]byte("\t"))
	}
}

// Decode the message
func (m *Message) Decode(in io.Reader) error {
	md := messageDecoder{
		in: in,
	}

	return md.decode(m)
}
