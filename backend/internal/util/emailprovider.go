package util

import (
	"fmt"
	"strconv"

	"github.com/cespare/xxhash"
)

func GenerateEmailProviderHash(email string, provider string) string {
	return strconv.FormatUint(xxhash.Sum64String(fmt.Sprintf("%s:%s", email, provider)), 10)
}
