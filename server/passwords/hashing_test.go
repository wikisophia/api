package hash_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/hash"
)

func TestHasher(t *testing.T) {
	hasher := hash.NewHasher(config.Hash{
		Time:        1,
		Memory:      64 * 1024,
		Parallelism: 1,
		SaltLength:  32,
		KeyLength:   32,
	})

	assertMatches(t, hasher, "password")
	assertMatches(t, hasher, "bjkncASDKIXCNH)*(*(!@#412-=_+~`,l./\\z]][p]{}682XDT&^T62t<>1?/.,m 2@#!wy8qasbki nyu")
}

func assertMatches(t *testing.T, hasher *hash.Hasher, value string) {
	t.Helper()
	hash, err := hasher.Hash(value)
	require.NoError(t, err)

	matches, err := hasher.Matches(value, hash)
	require.NoError(t, err)
	assert.True(t, matches)

	matches, err = hasher.Matches("some other value", hash)
	require.NoError(t, err)
	assert.False(t, matches)
}
