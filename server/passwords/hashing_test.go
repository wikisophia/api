package passwords_test

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/passwords"
)

func TestHasher(t *testing.T) {
	hasher := passwords.NewHasher(config.Hash{
		Time:        1,
		Memory:      64 * 1024,
		Parallelism: 1,
		SaltLength:  32,
		KeyLength:   32,
	})

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		go doHashTests(t, hasher, &wg)
	}
	wg.Add(2)
	assertMatches(t, hasher, "password", &wg)
	assertMatches(t, hasher, "bjkncASDKIXCNH)*(*(!@#412-=_+~`,l./\\z]][p]{}682XDT&^T62t<>1?/.,m 2@#!wy8qasbki nyu", &wg)
	wg.Wait()
}

func doHashTests(t *testing.T, hasher *passwords.Hasher, wg *sync.WaitGroup) {
	thisPassword := make([]byte, rand.Intn(50))
	wg.Add(10)
	for i := 0; i < 10; i++ {
		_, err := rand.Read(thisPassword)
		require.NoError(t, err)
		assertMatches(t, hasher, string(thisPassword), wg)
	}
}

func assertMatches(t *testing.T, hasher *passwords.Hasher, value string, wg *sync.WaitGroup) {
	t.Helper()
	hash, err := hasher.Hash(value)
	require.NoError(t, err)

	matches, err := hasher.Matches(value, hash)
	require.NoError(t, err)
	assert.True(t, matches)

	matches, err = hasher.Matches("some other value", hash)
	require.NoError(t, err)
	assert.False(t, matches)
	wg.Done()
}
