/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Status Codes
 */

package goipp

import (
	"fmt"
)

// Status represents an IPP Status Code
type Status Code

const (
	StatusOk                              Status = 0x0000 // successful-ok
	StatusOkIgnoredOrSubstituted          Status = 0x0001 // successful-ok-ignored-or-substituted-attributes
	StatusOkConflicting                   Status = 0x0002 // successful-ok-conflicting-attributes
	StatusOkIgnoredSubscriptions          Status = 0x0003 // successful-ok-ignored-subscriptions
	StatusOkIgnoredNotifications          Status = 0x0004 // successful-ok-ignored-notifications
	StatusOkTooManyEvents                 Status = 0x0005 // successful-ok-too-many-events
	StatusOkButCancelSubscription         Status = 0x0006 // successful-ok-but-cancel-subscription
	StatusOkEventsComplete                Status = 0x0007 // successful-ok-events-complete
	StatusRedirectionOtherSite            Status = 0x0200 // redirection-other-site
	StatusCupsSeeOther                    Status = 0x0280 // cups-see-other
	StatusErrorBadRequest                 Status = 0x0400 // client-error-bad-request
	StatusErrorForbidden                  Status = 0x0401 // client-error-forbidden
	StatusErrorNotAuthenticated           Status = 0x0402 // client-error-not-authenticated
	StatusErrorNotAuthorized              Status = 0x0403 // client-error-not-authorized
	StatusErrorNotPossible                Status = 0x0404 // client-error-not-possible
	StatusErrorTimeout                    Status = 0x0405 // client-error-timeout
	StatusErrorNotFound                   Status = 0x0406 // client-error-not-found
	StatusErrorGone                       Status = 0x0407 // client-error-gone
	StatusErrorRequestEntity              Status = 0x0408 // client-error-request-entity-too-large
	StatusErrorRequestValue               Status = 0x0409 // client-error-request-value-too-long
	StatusErrorDocumentFormatNotSupported Status = 0x040a // client-error-document-format-not-supported
	StatusErrorAttributesOrValues         Status = 0x040b // client-error-attributes-or-values-not-supported
	StatusErrorURIScheme                  Status = 0x040c // client-error-uri-scheme-not-supported
	StatusErrorCharset                    Status = 0x040d // client-error-charset-not-supported
	StatusErrorConflicting                Status = 0x040e // client-error-conflicting-attributes
	StatusErrorCompressionNotSupported    Status = 0x040f // client-error-compression-not-supported
	StatusErrorCompressionError           Status = 0x0410 // client-error-compression-error
	StatusErrorDocumentFormatError        Status = 0x0411 // client-error-document-format-error
	StatusErrorDocumentAccess             Status = 0x0412 // client-error-document-access-error
	StatusErrorAttributesNotSettable      Status = 0x0413 // client-error-attributes-not-settable
	StatusErrorIgnoredAllSubscriptions    Status = 0x0414 // client-error-ignored-all-subscriptions
	StatusErrorTooManySubscriptions       Status = 0x0415 // client-error-too-many-subscriptions
	StatusErrorIgnoredAllNotifications    Status = 0x0416 // client-error-ignored-all-notifications
	StatusErrorPrintSupportFileNotFound   Status = 0x0417 // client-error-print-support-file-not-found
	StatusErrorDocumentPassword           Status = 0x0418 // client-error-document-password-error
	StatusErrorDocumentPermission         Status = 0x0419 // client-error-document-permission-error
	StatusErrorDocumentSecurity           Status = 0x041a // client-error-document-security-error
	StatusErrorDocumentUnprintable        Status = 0x041b // client-error-document-unprintable-error
	StatusErrorAccountInfoNeeded          Status = 0x041c // client-error-account-info-needed
	StatusErrorAccountClosed              Status = 0x041d // client-error-account-closed
	StatusErrorAccountLimitReached        Status = 0x041e // client-error-account-limit-reached
	StatusErrorAccountAuthorizationFailed Status = 0x041f // client-error-account-authorization-failed
	StatusErrorNotFetchable               Status = 0x0420 // client-error-not-fetchable
	StatusErrorInternal                   Status = 0x0500 // server-error-internal-error
	StatusErrorOperationNotSupported      Status = 0x0501 // server-error-operation-not-supported
	StatusErrorServiceUnavailable         Status = 0x0502 // server-error-service-unavailable
	StatusErrorVersionNotSupported        Status = 0x0503 // server-error-version-not-supported
	StatusErrorDevice                     Status = 0x0504 // server-error-device-error
	StatusErrorTemporary                  Status = 0x0505 // server-error-temporary-error
	StatusErrorNotAcceptingJobs           Status = 0x0506 // server-error-not-accepting-jobs
	StatusErrorBusy                       Status = 0x0507 // server-error-busy
	StatusErrorJobCanceled                Status = 0x0508 // server-error-job-canceled
	StatusErrorMultipleJobsNotSupported   Status = 0x0509 // server-error-multiple-document-jobs-not-supported
	StatusErrorPrinterIsDeactivated       Status = 0x050a // server-error-printer-is-deactivated
	StatusErrorTooManyJobs                Status = 0x050b // server-error-too-many-jobs
	StatusErrorTooManyDocuments           Status = 0x050c // server-error-too-many-documents
)

// String() returns a Status name, as defined by RFC 8010
func (s Status) String() string {
	switch s {
	case StatusOk:
		return "successful-ok"
	case StatusOkIgnoredOrSubstituted:
		return "successful-ok-ignored-or-substituted-attributes"
	case StatusOkConflicting:
		return "successful-ok-conflicting-attributes"
	case StatusOkIgnoredSubscriptions:
		return "successful-ok-ignored-subscriptions"
	case StatusOkIgnoredNotifications:
		return "successful-ok-ignored-notifications"
	case StatusOkTooManyEvents:
		return "successful-ok-too-many-events"
	case StatusOkButCancelSubscription:
		return "successful-ok-but-cancel-subscription"
	case StatusOkEventsComplete:
		return "successful-ok-events-complete"
	case StatusRedirectionOtherSite:
		return "redirection-other-site"
	case StatusCupsSeeOther:
		return "cups-see-other"
	case StatusErrorBadRequest:
		return "client-error-bad-request"
	case StatusErrorForbidden:
		return "client-error-forbidden"
	case StatusErrorNotAuthenticated:
		return "client-error-not-authenticated"
	case StatusErrorNotAuthorized:
		return "client-error-not-authorized"
	case StatusErrorNotPossible:
		return "client-error-not-possible"
	case StatusErrorTimeout:
		return "client-error-timeout"
	case StatusErrorNotFound:
		return "client-error-not-found"
	case StatusErrorGone:
		return "client-error-gone"
	case StatusErrorRequestEntity:
		return "client-error-request-entity-too-large"
	case StatusErrorRequestValue:
		return "client-error-request-value-too-long"
	case StatusErrorDocumentFormatNotSupported:
		return "client-error-document-format-not-supported"
	case StatusErrorAttributesOrValues:
		return "client-error-attributes-or-values-not-supported"
	case StatusErrorURIScheme:
		return "client-error-uri-scheme-not-supported"
	case StatusErrorCharset:
		return "client-error-charset-not-supported"
	case StatusErrorConflicting:
		return "client-error-conflicting-attributes"
	case StatusErrorCompressionNotSupported:
		return "client-error-compression-not-supported"
	case StatusErrorCompressionError:
		return "client-error-compression-error"
	case StatusErrorDocumentFormatError:
		return "client-error-document-format-error"
	case StatusErrorDocumentAccess:
		return "client-error-document-access-error"
	case StatusErrorAttributesNotSettable:
		return "client-error-attributes-not-settable"
	case StatusErrorIgnoredAllSubscriptions:
		return "client-error-ignored-all-subscriptions"
	case StatusErrorTooManySubscriptions:
		return "client-error-too-many-subscriptions"
	case StatusErrorIgnoredAllNotifications:
		return "client-error-ignored-all-notifications"
	case StatusErrorPrintSupportFileNotFound:
		return "client-error-print-support-file-not-found"
	case StatusErrorDocumentPassword:
		return "client-error-document-password-error"
	case StatusErrorDocumentPermission:
		return "client-error-document-permission-error"
	case StatusErrorDocumentSecurity:
		return "client-error-document-security-error"
	case StatusErrorDocumentUnprintable:
		return "client-error-document-unprintable-error"
	case StatusErrorAccountInfoNeeded:
		return "client-error-account-info-needed"
	case StatusErrorAccountClosed:
		return "client-error-account-closed"
	case StatusErrorAccountLimitReached:
		return "client-error-account-limit-reached"
	case StatusErrorAccountAuthorizationFailed:
		return "client-error-account-authorization-failed"
	case StatusErrorNotFetchable:
		return "client-error-not-fetchable"
	case StatusErrorInternal:
		return "server-error-internal-error"
	case StatusErrorOperationNotSupported:
		return "server-error-operation-not-supported"
	case StatusErrorServiceUnavailable:
		return "server-error-service-unavailable"
	case StatusErrorVersionNotSupported:
		return "server-error-version-not-supported"
	case StatusErrorDevice:
		return "server-error-device-error"
	case StatusErrorTemporary:
		return "server-error-temporary-error"
	case StatusErrorNotAcceptingJobs:
		return "server-error-not-accepting-jobs"
	case StatusErrorBusy:
		return "server-error-busy"
	case StatusErrorJobCanceled:
		return "server-error-job-canceled"
	case StatusErrorMultipleJobsNotSupported:
		return "server-error-multiple-document-jobs-not-supported"
	case StatusErrorPrinterIsDeactivated:
		return "server-error-printer-is-deactivated"
	case StatusErrorTooManyJobs:
		return "server-error-too-many-jobs"
	case StatusErrorTooManyDocuments:
		return "server-error-too-many-documents"
	}

	return fmt.Sprintf("0x%4.4x", int(s))
}
