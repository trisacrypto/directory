package trtl

const (
	NamespacePeers    = "peers"
	NamespaceIndex    = "index"
	NamespaceDefault  = "default"
	NamespaceSequence = "sequence"
)

// Reserved namespaces that cannot be used by the caller since they are in use by trtl.
// If necessary this can be moved to configuration in the future.
var reservedNamespaces = map[string]struct{}{
	NamespacePeers:    {},
	NamespaceSequence: {},
	NamespaceIndex:    {},
	NamespaceDefault:  {}, // if the user does not specify a namespace
}
