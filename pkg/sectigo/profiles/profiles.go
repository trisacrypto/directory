package profiles

// Valid Sectigo Certificate Profile Names and IDs
// TODO: do not hardcode this, but get programatically from Sectigo API
const (
	ProfileCipherTraceEE                     = "CipherTrace EE"
	ProfileIDCipherTraceEE                   = "17"
	ProfileCipherTraceEndEntityCertificate   = "CipherTrace End Entity Certificate"
	ProfileIDCipherTraceEndEntityCertificate = "85"
)

var AllProfiles = [4]string{
	ProfileCipherTraceEE, ProfileIDCipherTraceEE,
	ProfileCipherTraceEndEntityCertificate, ProfileIDCipherTraceEndEntityCertificate,
}
