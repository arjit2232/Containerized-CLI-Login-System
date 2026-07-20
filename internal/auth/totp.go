package auth

import (
	"github.com/pquerna/otp/totp"
)

// GenerateSecret creates a new TOTP secret for the given username, along
// with the otpauth:// URI an authenticator app can consume (as a QR code
// or manual entry).
func GenerateSecret(username string) (secret string, qrURI string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OstoLoginCLI",
		AccountName: username,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// ValidateCode checks a 6-digit code against the stored secret, allowing
// for the default +/-1 time-step clock skew.
func ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// BuildQRURI reconstructs the otpauth:// URI for an already-known secret,
// useful if a user wants to re-display the manual-entry string without
// generating a brand new secret.
func BuildQRURI(username, secret string) string {
	return "otpauth://totp/" + "OstoLoginCLI:" + username +
		"?secret=" + secret + "&issuer=OstoLoginCLI"
}