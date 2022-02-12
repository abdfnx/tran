package tools

import (
	"fmt"
	"regexp"
	"math/rand"
	"encoding/hex"
	"crypto/sha256"

	"github.com/abdfnx/tran/data"
	"github.com/abdfnx/tran/models"
)

const passwordLength = 4

// GeneratePassword generates a random password prefixed with the supplied id.
func GeneratePassword(id int) models.Password {
	var words []string
	hitlistSize := len(data.PasswordList)

	// generate three unique words
	for len(words) != passwordLength {
		candidateWord := data.PasswordList[rand.Intn(hitlistSize)]
		if !Contains(words, candidateWord) {
			words = append(words, candidateWord)
		}
	}

	password := formatPassword(id, words)
	return models.Password(password)
}

func ParsePassword(passStr string) (models.Password, error) {
	re := regexp.MustCompile(`^\d+-[a-z]+-[a-z]+-[a-z]+$`)
	ok := re.MatchString(passStr)

	if !ok {
		return models.Password(""), fmt.Errorf("password: %q is on wrong format", passStr)
	}

	return models.Password(passStr), nil
}

func formatPassword(prefixIndex int, words []string) string {
	return fmt.Sprintf("%d-%s-%s-%s", prefixIndex, words[0], words[1], words[2])
}

func HashPassword(password models.Password) string {
	h := sha256.New()
	h.Write([]byte(password))

	return hex.EncodeToString(h.Sum(nil))
}
