package types

// Secret contains metadata for a secret.
type Secret struct {
	Data map[string]interface{} `json:"data,omitempty" validate:"required"` // Secret data.
	Path string                 `json:"path,omitempty" validate:"required"` // Path.
}
