package trtl

import "github.com/trisacrypto/directory/pkg/utils/wire"

const (
	NamespacePeers    = wire.NamespaceReplicas
	NamespaceIndex    = wire.NamespaceIndices
	NamespaceDefault  = "default"
	NamespaceSequence = wire.NamespaceSequence
	NamespaceVASPs    = wire.NamespaceVASPs
	NamespaceCertReqs = wire.NamespaceCertReqs
)

// Reserved namespaces that cannot be used by the caller since they are in use by trtl.
// If necessary this can be moved to configuration in the future.
var reservedNamespaces = map[string]struct{}{
	NamespacePeers:    {},
	NamespaceSequence: {},
	NamespaceIndex:    {},
	NamespaceDefault:  {}, // if the user does not specify a namespace
}

// Replicated namespaces are the namespaces that are used in anti-entropy by default.
var replicatedNamespaces = []string{
	NamespaceVASPs, NamespaceCertReqs, NamespacePeers, NamespaceDefault,
}
