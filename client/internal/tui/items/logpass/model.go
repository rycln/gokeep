// Logpass package manages login/password content storage and display.
package logpass

// LogPass represents login credentials for a service or application.
// Contains sensitive authentication data that should always be encrypted.
type LogPass struct {
	Login    string `json:"login"`    // Username or email for authentication
	Password string `json:"password"` // Secret password
}
