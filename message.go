/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP protocol messages
 */

package goipp

import (
	"bytes"
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
	Operation         Attributes // Operation attributes
	Job               Attributes // Job attributes
	Printer           Attributes // Printer attributes
	Unsupported       Attributes // Unsupported attributes
	Subscription      Attributes // Subscription attributes
	EventNotification Attributes // Event Notification attributes
	Resource          Attributes // Resource attributes
	Document          Attributes // Document attributes
	System            Attributes // System attributes
	Future11          Attributes // \
	Future12          Attributes //  \
	Future13          Attributes //   | Reserved for future extensions
	Future14          Attributes //  /
	Future15          Attributes // /
}

// NewRequest creates a new request message
//
// Use DefaultVersion as a first argument, if you don't
// have any specific needs
func NewRequest(v Version, op Op, id uint32) *Message {
	return &Message{
		Version:   v,
		Code:      Code(op),
		RequestID: id,
	}
}

// NewResponse creates a new response message
//
// Use DefaultVersion as a first argument, if you don't
func NewResponse(v Version, status Status, id uint32) *Message {
	return &Message{
		Version:   v,
		Code:      Code(status),
		RequestID: id,
	}
}

// Equal checks that two messages are equal
func (m Message) Equal(m2 Message) bool {
	if m.Version != m2.Version ||
		m.Code != m2.Code ||
		m.RequestID != m2.RequestID {
		return false
	}

	groups1 := m.attrGroups()
	groups2 := m2.attrGroups()

	if len(groups1) != len(groups2) {
		return false
	}

	for i, grp1 := range groups1 {
		grp2 := groups2[i]

		if grp1.tag != grp2.tag || !grp1.attrs.Equal(grp2.attrs) {
			return false
		}
	}

	return true
}

// Reset the message into initial state
func (m *Message) Reset() {
	*m = Message{}
}

// Encode message
func (m *Message) Encode(out io.Writer) error {
	me := messageEncoder{
		out: out,
	}

	return me.encode(m)
}

// EncodeBytes encodes message to byte slice
func (m *Message) EncodeBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := m.Encode(&buf)
	return buf.Bytes(), err
}

// Decode message
func (m *Message) Decode(in io.Reader) error {
	md := messageDecoder{
		in: in,
	}

	m.Reset()
	return md.decode(m)
}

// DecodeBytes decodes message from byte slice
func (m *Message) DecodeBytes(data []byte) error {
	return m.Decode(bytes.NewBuffer(data))
}

// Print pretty-prints the message. The 'request' parameter affects
// interpretation of Message.Code: it is interpreted either
// as Op or as Status
func (m *Message) Print(out io.Writer, request bool) {
	out.Write([]byte("{\n"))

	fmt.Fprintf(out, msgPrintIndent+"VERSION %s\n", m.Version)

	if request {
		fmt.Fprintf(out, msgPrintIndent+"OPERATION %s\n", Op(m.Code))
	} else {
		fmt.Fprintf(out, msgPrintIndent+"STATUS %s\n", Status(m.Code))
	}

	for _, grp := range m.attrGroups() {
		fmt.Fprintf(out, "\n"+msgPrintIndent+"GROUP %s\n", grp.tag)
		for _, attr := range grp.attrs {
			m.printAttribute(out, attr, 1)
			out.Write([]byte("\n"))
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
		out.Write([]byte(msgPrintIndent))
	}
}

// Get attributes by group. Groups with nil Attributes are skipped,
// but groups with non-nil are not, even if len(Attributes) == 0
//
// This is a helper function for message encoder and pretty-printer
func (m *Message) attrGroups() []struct {
	tag   Tag
	attrs Attributes
} {
	// Initialize slice of groups
	groups := []struct {
		tag   Tag
		attrs Attributes
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

	// Skip all empty groups
	out := 0
	for in := 0; in < len(groups); in++ {
		if groups[in].attrs != nil {
			groups[out] = groups[in]
			out++
		}
	}

	return groups[:out]
}
