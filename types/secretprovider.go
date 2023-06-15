// Package types hosts the SecretProvider type, specifying where secrets are stored.
package types

// SecretProvider contains location metadata where secrets are stored.
type SecretProvider struct {
	// IDs.
	ID      string `bson:"_id,omitempty" json:"id,omitempty" validate:"required"` // ID of the secret store.
	Inherit bool   `json:"inherit,omitempty" default:"true"`                      // Whether to inherit settings from ancestors.

	// Instance metadata.
	CollectionType string `json:"-"` // Collection type.
	ParentNodeID   string `json:"-"` // ID of the parent node.

	// Persistent metadata.
	ClientID     string   `env:"SECRETSTORE_CLIENTID" json:"clientID,omitempty"`             // Client ID used by the secret store provider.
	ClientSecret string   `env:"SECRETSTORE_CLIENTSECRET" json:"clientSecret,omitempty"`     // Shared secret used by the secret store provider.
	ClientToken  string   `env:"SECRETSTORE_CLIENTTOKEN" json:"clientToken,omitempty"`       // Optional token used by the secret store provider.
	Region       string   `env:"SECRETSTORE_REGION" json:"region,omitempty"`                 // Region of the secret store provider.
	Type         string   `env:"SECRETSTORE_TYPE" json:"type,omitempty" validate:"required"` // Type of secret storage (e.g., Vault).
	UnsealShards []string `env:"SECRETSTORE_UNSEALSHARDS" json:"unsealShards,omitempty"`     // Shared secrets to unseal the secret store.
	URI          string   `env:"SECRETSTORE_URI" json:"uri,omitempty"`                       // Address of the secret store.
}
