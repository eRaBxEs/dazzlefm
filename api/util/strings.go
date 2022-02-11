//
// strings.go
// a collection of strings manipulation functions
//

package util

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789" +
	"ABCDEGHIJKLMNOPQRSTUVWXYZ"

const numset = "0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset ...
func StringWithCharset(length int, charset string) string {

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]

	}
	return string(b)
}

// StringNumbersCharset ...
func StringNumbersCharset(length int, numset string) string {

	b := make([]byte, length)

	for i := range b {
		b[i] = numset[seededRand.Intn(len(numset))]
	}
	return string(b)
}

// String ...
func String(length int) string {

	return StringWithCharset(length, charset)
}

// Numbers ...
func Numbers(length int) string {

	return StringNumbersCharset(length, numset)
}

// NigeriaCode ...
func NigeriaCode(mobileNo string) string {

	data := string(mobileNo[0])
	var mobile string

	if data == "0" {

		mobile = strings.Replace(mobileNo, "0", "234", 1)
	} else {
		mobile = mobileNo
	}

	return mobile
}
