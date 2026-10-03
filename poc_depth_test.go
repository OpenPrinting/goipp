package goipp

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// decodeCollection (decoder.go) recursed once per nesting level with no
// depth limit. A small, attacker-crafted IPP message with millions of
// nested collections drove the Go call stack to its limit, crashing the
// whole process with a fatal (unrecoverable) error - confirmed with a
// ~22MB message encoding 2,000,000 nested levels. maxCollectionDepth now
// bounds the recursion and turns that crash into an ordinary decode error.

func buildNestedCollectionMessage(depth int) []byte {
	var buf bytes.Buffer

	// Header: version (2 bytes), code (2 bytes), request-id (4 bytes)
	buf.Write([]byte{0x02, 0x00}) // version 2.0
	buf.Write([]byte{0x00, 0x02}) // code
	buf.Write([]byte{0, 0, 0, 1}) // request-id

	// Group delimiter: operation-attributes-tag
	buf.WriteByte(byte(TagOperationGroup))

	// Top-level attribute: TagBeginCollection, name="c", value=""
	writeAttr(&buf, byte(TagBeginCollection), "c", nil)

	// Nested levels: TagMemberName("m") + TagBeginCollection(name="")
	for i := 0; i < depth; i++ {
		writeAttr(&buf, byte(TagMemberName), "", []byte("m"))
		writeAttr(&buf, byte(TagBeginCollection), "", nil)
	}

	// Close every collection level
	for i := 0; i < depth+1; i++ {
		writeAttr(&buf, byte(TagEndCollection), "", nil)
	}

	buf.WriteByte(byte(TagEnd))

	return buf.Bytes()
}

func writeAttr(buf *bytes.Buffer, tag byte, name string, value []byte) {
	buf.WriteByte(tag)
	var l [2]byte
	binary.BigEndian.PutUint16(l[:], uint16(len(name)))
	buf.Write(l[:])
	buf.WriteString(name)
	binary.BigEndian.PutUint16(l[:], uint16(len(value)))
	buf.Write(l[:])
	buf.Write(value)
}

// Just past the limit: must return a clean decode error, not recurse further.
func TestCollectionDepthLimitRejectsExcess(t *testing.T) {
	msg := buildNestedCollectionMessage(maxCollectionDepth + 10)

	var m Message
	err := m.DecodeBytes(msg)
	if err == nil {
		t.Fatal("expected a decode error for collections nested past maxCollectionDepth, got nil")
	}
}

// Within the limit: nesting that any real IPP device would send must
// still decode successfully.
func TestCollectionDepthLimitAllowsReasonableNesting(t *testing.T) {
	msg := buildNestedCollectionMessage(maxCollectionDepth - 1)

	var m Message
	err := m.DecodeBytes(msg)
	if err != nil {
		t.Fatalf("expected reasonable collection nesting to decode cleanly, got: %v", err)
	}
}

// Original confirmed crash reproduction, now safe to run as a regression
// test: with the fix in place this fails fast with a clean error at
// maxCollectionDepth, long before the stack could be exhausted, instead
// of crashing the process.
func TestPocUnboundedCollectionRecursionNowBounded(t *testing.T) {
	msg := buildNestedCollectionMessage(2000000)

	var m Message
	err := m.DecodeBytes(msg)
	if err == nil {
		t.Fatal("expected a decode error for 2,000,000 nested collections, got nil")
	}
	t.Logf("decode result (expected error): %v", err)
}
