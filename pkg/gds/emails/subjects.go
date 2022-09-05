package emails

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
