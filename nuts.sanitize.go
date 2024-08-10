package gonuts

import (
	"strings"
	"unicode"
)

// SANITIZE_SQLSAFER is a list of strings that are commonly used in SQL injection attacks.
// NOTE: This is not a comprehensive list and should not be relied upon for complete SQL injection protection.
var SANITIZE_SQLSAFER = []string{"`", "´", "'", " OR ", " or ", "=", ";", ":", "(", ")", "--", "/*", "*/", "@@", "@"}

// SanitizeString removes potentially dangerous strings from the input.
//
// WARNING: This function provides basic sanitization and should not be considered
// a complete solution for preventing SQL injection or other security vulnerabilities.
// Always use parameterized queries and proper input validation in addition to this function.
//
// Parameters:
//   - badStringsList: A slice of strings to be removed from the input.
//   - stringToClean: The input string to be sanitized.
//
// Returns:
//   - A sanitized version of the input string.
//
// Example usage:
//
//	input := "SELECT * FROM users WHERE name = 'John' OR 1=1; --"
//	sanitized := gonuts.SanitizeString(gonuts.SANITIZE_SQLSAFER, input)
//	fmt.Println(sanitized) // Output: SELECT * FROM users WHERE name  John
func SanitizeString(badStringsList []string, stringToClean string) string {
	cleanString := stringToClean

	// Remove all strings from the badStringsList
	for _, badString := range badStringsList {
		cleanString = strings.ReplaceAll(cleanString, badString, "")
	}

	// Remove control characters
	cleanString = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, cleanString)

	// Trim spaces
	cleanString = strings.TrimSpace(cleanString)

	return cleanString
}

// SafeSQLString prepares a string for safe use in SQL queries by escaping single quotes.
//
// WARNING: This function should be used in conjunction with parameterized queries,
// not as a replacement for them. It does not provide complete protection against SQL injection.
//
// Parameters:
//   - s: The input string to be escaped.
//
// Returns:
//   - An escaped version of the input string, safe for use in SQL queries.
//
// Example usage:
//
//	userInput := "O'Reilly"
//	safeName := gonuts.SafeSQLString(userInput)
//	query := fmt.Sprintf("SELECT * FROM authors WHERE name = '%s'", safeName)
//	// Use parameterized queries instead of string formatting in production code
func SafeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
