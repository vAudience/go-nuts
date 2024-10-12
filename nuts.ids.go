package gonuts

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// idAlphabet is the set of characters used for generating IDs.
const idAlphabet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const NID_Prefix_Separator = "_"

var (
	UUIDRegEx            = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	NotLegalIdCharacters = regexp.MustCompile("[^A-Za-z0-9-_]")

	ErrBadUUID     = errors.New("uuid format error")
	ErrBadId       = errors.New("bad id format")
	ErrIllegalId   = errors.New("illegal id")
	ErrUnknownId   = errors.New("unknown id")
	ErrMalformedId = errors.New("malformed id")
)

// NanoID generates a unique ID with a given prefix.
//
// It uses a cryptographically secure random number generator to create
// a unique identifier, then prepends the given prefix to it.
//
// Parameters:
//   - prefix: A string to be prepended to the generated ID.
//
// Returns:
//   - A string containing the prefix followed by a unique identifier.
//
// Example usage:
//
//	id := gonuts.NanoID("user")
//	fmt.Println(id) // Output: user_6ByTSYmGzT2c
func NanoID(prefix string) string {
	nid, err := gonanoid.Generate(idAlphabet, 12)
	if err != nil {
		L.Error(err)
		// Fallback to using timestamp if nanoid generation fails
		nid = strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return prefix + NID_Prefix_Separator + nid
}

// NID generates a unique ID with a specified length and optional prefix.
//
// It creates a unique identifier of the given length using the idAlphabet.
// If a prefix is provided, it's prepended to the generated ID.
//
// Parameters:
//   - prefix: An optional string to be prepended to the generated ID. Use "" for no prefix.
//   - length: The desired length of the generated part of the ID (excluding prefix).
//
// Returns:
//   - A string containing the optional prefix followed by a unique identifier of the specified length.
//
// Example usage:
//
//	id1 := gonuts.NID("doc", 8)
//	fmt.Println(id1) // Output: doc_r3tM9wK1
//
//	id2 := gonuts.NID("", 10)
//	fmt.Println(id2) // Output: 7mHxL2pQ4R
func NID(prefix string, length int) string {
	nid, err := gonanoid.Generate(idAlphabet, length)
	if err != nil {
		L.Error(err)
		// Fallback to using timestamp and random bytes if nanoid generation fails
		nid = fallbackIDGeneration(length)
	}
	if prefix != "" {
		return prefix + NID_Prefix_Separator + nid
	}
	return nid
}

// IsNID checks if the given ID is a valid NID with the specified prefix and length.
func IsNID(id string, prefix string, length int) bool {
	if len(id) != length {
		return false
	}
	if prefix != "" {
		return id[:len(prefix)+1] == prefix+NID_Prefix_Separator
	}
	return true
}

// fallbackIDGeneration creates an ID using timestamp and random bytes.
// This is used as a fallback method if the primary ID generation fails.
func fallbackIDGeneration(length int) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 36)
	remainingLength := length - len(timestamp)
	if remainingLength <= 0 {
		return timestamp[:length]
	}

	randomPart := make([]byte, remainingLength)
	_, err := rand.Read(randomPart)
	if err != nil {
		// If random generation fails, pad with '0'
		return timestamp + string(make([]byte, remainingLength))
	}

	for i := 0; i < remainingLength; i++ {
		randomPart[i] = idAlphabet[int(randomPart[i])%len(idAlphabet)]
	}
	return timestamp + string(randomPart)
}

// GenerateRandomString creates a random string of a given length using the provided character set.
//
// Parameters:
//   - chars: A slice of runes representing the character set to use.
//   - length: The desired length of the random string.
//
// Returns:
//   - A randomly generated string of the specified length.
//
// Example usage:
//
//	chars := []rune("abcdefghijklmnopqrstuvwxyz")
//	randomStr := gonuts.GenerateRandomString(chars, 10)
//	fmt.Println(randomStr) // Output: (a random 10-character string using the given alphabet)
func GenerateRandomString(chars []rune, length int) string {
	b := make([]rune, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			L.Error(err)
			// Fallback to less secure random if crypto/rand fails
			b[i] = chars[time.Now().Nanosecond()%len(chars)]
		} else {
			b[i] = chars[n.Int64()]
		}
	}
	return string(b)
}
