package main

import (
	"fmt"
)

// Type Tag represents a tag used in a binary representation
// of the IPP message
type Tag uint8

const (
	// Delimiter tags
	TagZero                   Tag = 0x00 // Zero tag - used for separators
	TagOperationGroup         Tag = 0x01 // Operation group
	TagJobGroup               Tag = 0x02 // Job group
	TagEnd                    Tag = 0x03 // End-of-attributes
	TagPrinterGroup           Tag = 0x04 // Printer group
	TagUnsupportedGroup       Tag = 0x05 // Unsupported attributes group
	TagSubscriptionGroup      Tag = 0x06 // Subscription group
	TagEventNotificationGroup Tag = 0x07 // Event group
	TagResourceGroup          Tag = 0x08 // Resource group
	TagDocumentGroup          Tag = 0x09 // Document group
	TagSystemGroup            Tag = 0x0a // System group
	TagFuture11Group          Tag = 0x0b // Future group 11
	TagFuture12Group          Tag = 0x0c // Future group 11
	TagFuture13Group          Tag = 0x0d // Future group 11
	TagFuture14Group          Tag = 0x0e // Future group 11
	TagFuture15Group          Tag = 0x0f // Future group 11

	// Value tags
	TagUnsupportedValue Tag = 0x10 // Unsupported value
	TagDefault          Tag = 0x11 // Default value
	TagUnknown          Tag = 0x12 // Unknown value
	TagNoValue          Tag = 0x13 // No-value value
	TagNotSettable      Tag = 0x15 // Not-settable value
	TagDeleteAttr       Tag = 0x16 // Delete-attribute value
	TagAdminDefine      Tag = 0x17 // Admin-defined value
	TagInteger          Tag = 0x21 // Integer value
	TagBoolean          Tag = 0x22 // Boolean value
	TagEnum             Tag = 0x23 // Enumeration value
	TagString           Tag = 0x30 // Octet string value
	TagDate             Tag = 0x31 // Date/time value
	TagResolution       Tag = 0x32 // Resolution value
	TagRange            Tag = 0x33 // Range value
	TagBeginCollection  Tag = 0x34 // Beginning of collection value @exclude all@
	TagTextLang         Tag = 0x35 // Text-with-language value
	TagNameLang         Tag = 0x36 // Name-with-language value
	TagEndCollection    Tag = 0x37 // End of collection value
	TagText             Tag = 0x41 // Text value
	TagName             Tag = 0x42 // Name value
	TagReservedString   Tag = 0x43 // Reserved for future string value
	TagKeyword          Tag = 0x44 // Keyword value
	TagUri              Tag = 0x45 // URI value
	TagUriScheme        Tag = 0x46 // URI scheme value
	TagCharset          Tag = 0x47 // Character set value
	TagLanguage         Tag = 0x48 // Language value
	TagMimeType         Tag = 0x49 // MIME media type value
	TagMemberName       Tag = 0x4a // Collection member name value
	TagExtension        Tag = 0x7f // Extension point for 32-bit tags
)

// Check if Tag is delimiter
func (tag Tag) IsDelimiter() bool {
	return tag < 0x10
}

// String() returns a tag name in a form, as defined by RFC 8010
func (tag Tag) String() string {
	switch tag {
	case TagZero:
		return "zero"
	case TagOperationGroup:
		return "operation-attributes-tag"
	case TagJobGroup:
		return "job-attributes-tag"
	case TagEnd:
		return "end-of-attributes-tag"
	case TagPrinterGroup:
		return "printer-attributes-tag"
	case TagUnsupportedGroup:
		return "unsupported-attributes-tag"
	case TagSubscriptionGroup:
		return "subscription-attributes-tag"
	case TagEventNotificationGroup:
		return "event-notification-attributes-tag"
	case TagResourceGroup:
		return "resource-attributes-tag"
	case TagDocumentGroup:
		return "document-attributes-tag"
	case TagSystemGroup:
		return "system-attributes-tag"

	// Value tags
	case TagUnsupportedValue:
		return "unsupported"
	case TagDefault:
		return "default"
	case TagUnknown:
		return "unknown"
	case TagNoValue:
		return "no-value"
	case TagNotSettable:
		return "not-settable"
	case TagDeleteAttr:
		return "delete-attribute"
	case TagAdminDefine:
		return "admin-define"
	case TagInteger:
		return "integer"
	case TagBoolean:
		return "boolean"
	case TagEnum:
		return "enum"
	case TagString:
		return "octetString"
	case TagDate:
		return "dateTime"
	case TagResolution:
		return "resolution"
	case TagRange:
		return "rangeOfInteger"
	case TagBeginCollection:
		return "collection"
	case TagTextLang:
		return "textWithLanguage"
	case TagNameLang:
		return "nameWithLanguage"
	case TagEndCollection:
		return "endCollection"
	case TagText:
		return "textWithoutLanguage"
	case TagName:
		return "nameWithoutLanguage"
	case TagKeyword:
		return "keyword"
	case TagUri:
		return "uri"
	case TagUriScheme:
		return "uriScheme"
	case TagCharset:
		return "charset"
	case TagLanguage:
		return "naturalLanguage"
	case TagMimeType:
		return "mimeMediaType"
	case TagMemberName:
		return "memberAttrName"
	}

	return fmt.Sprintf("0x%2.2x", int(tag))
}
