package token

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func Test_generateRandomCharacters(t *testing.T) {
	strLen := 12
	for strLen != 6 {
		randomString := generateRandomCharacters(strLen)
		assert.Equal(t, len(randomString), strLen)
		strLen--
	}
}
