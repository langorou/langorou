package utils

import (
	"log"
	"math/rand"
	"os"
	"unicode"
)

// IsASCII verifies is the string is composed of only ASCII characters
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// CreateDirIfNotExist does exactly what it said
// thx https://siongui.github.io/2017/03/28/go-create-directory-if-not-exist/
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRandIntInRange returns a random int from specific range
func GetRandIntInRange(min, max int) int {
	return min + rand.Intn(max-min)
}

// FailIf logs a fatal error and exit
func FailIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}
