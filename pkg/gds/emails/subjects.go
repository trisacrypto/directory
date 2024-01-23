package emails

import "github.com/trisacrypto/directory/pkg/gds/admin/v2"

// Subjects for creating emails on demand.
const (
	VerifyContactRE               = "TRISA: Please verify your email address"
	ReviewRequestRE               = "New TRISA Global Directory Registration Request"
	RejectRegistrationRE          = "TRISA Global Directory Registration Update"
	DeliverCertsRE                = "Welcome to the TRISA network!"
	ExpiresAdminNotificationRE    = "A TRISA Identity Certificate is Expiring Soon"
	ReissuanceReminderRE          = "TRISA Identity Certificate Expiration"
	ReissuanceStartedRE           = "TRISA PKCS12 Password for Certificate Reissuance"
	ReissuanceAdminNotificationRE = "A TRISA Identity Certificate Reissuance has been Completed"
)

// These variables are currently used for tests only since the ResendAction constants
// are only defined on the API and are not generally used in the code base except for
// tests.
var (
	Reason2Subject = map[string]string{
		string(admin.ResendVerifyContact): VerifyContactRE,
		string(admin.ResendReview):        ReviewRequestRE,
		string(admin.ResendDeliverCerts):  DeliverCertsRE,
		string(admin.ResendRejection):     RejectRegistrationRE,
		string(admin.ReissuanceReminder):  ReissuanceReminderRE,
		string(admin.ReissuanceStarted):   ReissuanceStartedRE,
	}
	Subject2Reason = map[string]string{
		VerifyContactRE:      string(admin.ResendVerifyContact),
		ReviewRequestRE:      string(admin.ResendReview),
		RejectRegistrationRE: string(admin.ResendRejection),
		DeliverCertsRE:       string(admin.ResendDeliverCerts),
		ReissuanceReminderRE: string(admin.ReissuanceReminder),
		ReissuanceStartedRE:  string(admin.ReissuanceStarted),
	}
)
