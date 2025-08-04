// Package contextkeys provides type-safe keys for storing values in request context.
package contextkeys

type contextKey struct{}

// Package-level context keys for storing common request values.
var (
	// UserID is the context key for storing authenticated user ID.
	// Populated by auth middleware after JWT verification.
	UserID = contextKey{}
)
