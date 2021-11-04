package trtl

// Reserved namespaces that cannot be used by the caller since they are in use by trtl.
// If necessary this can be moved to configuration in the future.
var reservedNamespaces = map[string]struct{}{
	"peers":    struct{}{},
	"sequence": struct{}{},
	"index":    struct{}{},
	"default":  struct{}{}, // if the user does not specify a namespace
}
