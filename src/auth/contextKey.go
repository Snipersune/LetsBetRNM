package auth

// Context Key
type contextKey string

// Context keys used for checking session (cookie) variables
const UserIDKey contextKey = "user_id"
const UserNameKey contextKey = "username"
