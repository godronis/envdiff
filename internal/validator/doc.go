// Package validator provides security and correctness validation
// for .env file entries.
//
// It complements the linter package by focusing on runtime-safety
// concerns such as:
//
//   - Detecting potentially exposed secrets (passwords, tokens, API keys)
//     that are not properly redacted in non-production env files.
//   - Flagging insecure URL schemes (http:// instead of https://).
//   - Identifying unresolved template expressions left in values.
//
// Each finding is reported as an Issue with a severity of either
// "warning" (advisory) or "error" (should block usage).
//
// Usage:
//
//	result := validator.Validate(envs)
//	if result.HasErrors() {
//		// handle blocking issues
//	}
package validator
