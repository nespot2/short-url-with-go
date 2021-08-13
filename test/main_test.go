package tests

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var base62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func TestEncode(t *testing.T) {
	encode1 := encode(100)
	assert.Equal(t, "mB", encode1)
}

func TestDecode(t *testing.T) {
	decode1 := decode("mB")

	fmt.Println(decode1)
	assert.Equal(t, uint64(100), decode1)
}

func encode(id uint64) string {

	var b bytes.Buffer

	i := id % 62
	id /= 62

	b.WriteString(string(base62[i]))

	for id > 0 {
		i := id % 62
		id /= 62
		b.WriteString(string(base62[i]))
	}
	return b.String()
}

func decode(url string) uint64 {
	result := uint64(0)
	power := uint64(1)

	for _, char := range url {
		digit := uint64(strings.Index(base62, string(char)))
		result += digit * power
		power *= 62
	}
	return result
}
