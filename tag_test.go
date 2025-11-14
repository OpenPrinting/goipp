/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Tags tests
 */

package goipp

import (
	"fmt"
	"testing"
)

// TestTagIsDelimiter tests Tag.IsDelimiter function
func TestTagIsDelimiter(t *testing.T) {
	type testData struct {
		t      Tag
		answer bool
	}

	tests := []testData{
		{TagZero, true},
		{TagOperationGroup, true},
		{TagJobGroup, true},
		{TagEnd, true},
		{TagFuture15Group, true},
		{TagUnsupportedValue, false},
		{TagUnknown, false},
		{TagInteger, false},
		{TagBeginCollection, false},
		{TagEndCollection, false},
		{TagExtension, false},
	}

	for _, test := range tests {
		answer := test.t.IsDelimiter()
		if answer != test.answer {
			t.Errorf("testing Tag.IsDelimiter:\n"+
				"tag:      %s (0x%.2x)\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.t, uint32(test.t), test.answer, answer,
			)
		}
	}
}

// TestTagIsGroup tests Tag.IsGroup function
func TestTagIsGroup(t *testing.T) {
	type testData struct {
		t      Tag
		answer bool
	}

	tests := []testData{
		{TagZero, false},
		{TagOperationGroup, true},
		{TagJobGroup, true},
		{TagEnd, false},
		{TagPrinterGroup, true},
		{TagUnsupportedGroup, true},
		{TagSubscriptionGroup, true},
		{TagEventNotificationGroup, true},
		{TagResourceGroup, true},
		{TagDocumentGroup, true},
		{TagSystemGroup, true},
		{TagFuture11Group, true},
		{TagFuture12Group, true},
		{TagFuture13Group, true},
		{TagFuture14Group, true},
		{TagFuture15Group, true},
		{TagInteger, false},
	}

	for _, test := range tests {
		answer := test.t.IsGroup()
		if answer != test.answer {
			t.Errorf("testing Tag.IsGroup:\n"+
				"tag:      %s (0x%.2x)\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.t, uint32(test.t), test.answer, answer,
			)
		}
	}
}

// TestTagType tests Tag.Type function
func TestTagType(t *testing.T) {
	type testData struct {
		t      Tag
		answer Type
	}

	tests := []testData{
		{TagZero, TypeInvalid},
		{TagInteger, TypeInteger},
		{TagEnum, TypeInteger},
		{TagBoolean, TypeBoolean},
		{TagUnsupportedValue, TypeVoid},
		{TagDefault, TypeVoid},
		{TagUnknown, TypeVoid},
		{TagNotSettable, TypeVoid},
		{TagNoValue, TypeVoid},
		{TagDeleteAttr, TypeVoid},
		{TagAdminDefine, TypeVoid},
		{TagText, TypeString},
		{TagName, TypeString},
		{TagReservedString, TypeString},
		{TagKeyword, TypeString},
		{TagURI, TypeString},
		{TagURIScheme, TypeString},
		{TagCharset, TypeString},
		{TagLanguage, TypeString},
		{TagMimeType, TypeString},
		{TagMemberName, TypeString},
		{TagDateTime, TypeDateTime},
		{TagResolution, TypeResolution},
		{TagRange, TypeRange},
		{TagTextLang, TypeTextWithLang},
		{TagNameLang, TypeTextWithLang},
		{TagBeginCollection, TypeCollection},
		{TagEndCollection, TypeVoid},
		{TagExtension, TypeBinary},
		{0x1234, TypeBinary},
	}

	for _, test := range tests {
		answer := test.t.Type()
		if answer != test.answer {
			t.Errorf("testing Tag.Type:\n"+
				"tag:      %s (0x%.2x)\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.t, uint32(test.t), test.answer, answer,
			)
		}
	}
}

// TestTagString tests Tag.String function
func TestTagString(t *testing.T) {
	type testData struct {
		t      Tag
		answer string
	}

	tests := []testData{
		{TagZero, "zero"},
		{TagUnsupportedValue, "unsupported"},
		{-1, "0xffffffff"},
		{0xff, "0xff"},
		{0x1234, "0x00001234"},
	}

	for _, test := range tests {
		answer := test.t.String()
		if answer != test.answer {
			t.Errorf("testing Tag.String:\n"+
				"tag:      %s (0x%.2x)\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.t, uint32(test.t), test.answer, answer,
			)
		}
	}
}

// TestTagGoString tests Tag.GoString function
func TestTagGoString(t *testing.T) {
	type testData struct {
		t      Tag
		answer string
	}

	tests := []testData{
		{TagZero, "goipp.TagZero"},
		{TagUnsupportedValue, "goipp.TagUnsupportedValue"},
		{-1, "goipp.Tag(0xffffffff)"},
		{0xff, "goipp.Tag(0xff)"},
		{0x1234, "goipp.Tag(0x00001234)"},
	}

	for _, test := range tests {
		answer := fmt.Sprintf("%#v", test.t)
		if answer != test.answer {
			t.Errorf("testing Tag.GoString:\n"+
				"tag:      %s (0x%.2x)\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.t, uint32(test.t), test.answer, answer,
			)
		}
	}
}
