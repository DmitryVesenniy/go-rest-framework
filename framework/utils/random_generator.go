package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@$&"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func RandomGenerator(length int, prefix string) string {
	rand.Seed(time.Now().UnixNano())
	var key strings.Builder

	for i := 0; i < length; i++ {
		random := rand.Intn(len(allCharSet))
		key.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(key.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, string(inRune))
	}

	return string(inRune)
}
