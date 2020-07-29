package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/wikisophia/api-arguments/server/config"
	"golang.org/x/crypto/argon2"
)

// Hasher can hash and verify hashes of incoming strings.
// Create these with the NewHasher() function
type Hasher struct {
	params config.Hash
	salts  *sync.Pool
}

// NewHasher makes a Hasher which runs Argon2 with the given params.
func NewHasher(params config.Hash) *Hasher {
	return &Hasher{
		params: params,
		salts: &sync.Pool{
			New: func() interface{} {
				return make([]byte, params.SaltLength)
			},
		},
	}
}

// Hash hashes a string and returns the hashed value.
func (h *Hasher) Hash(value string) (string, error) {
	salt := h.salts.Get().([]byte)
	defer h.salts.Put(salt)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return h.doHash([]byte(value), salt, h.params.Time, h.params.Memory, h.params.Parallelism, h.params.KeyLength), nil
}

// Matches returns true if the value matches the hash, and false otherwise.
// An error will be thrown if the hash isn't formatted properly. This shouldn't happen with hashes generated
// by the Hash() function
func (h *Hasher) Matches(value string, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false, errors.New("hash does not have the five expected $ symbols")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, errors.New("failed to parse the hash version")
	}
	if version != argon2.Version {
		return false, fmt.Errorf("the golang library implements hash version %d, but the password was hashed with %d", argon2.Version, version)
	}

	var memory uint32
	var time uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism); err != nil {
		return false, errors.New("Could not parse the time, memory, and parallelism params from the hashed value")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("salt was not base64 encoded properly")
	}

	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("hash was not base64 encoded properly")
	}

	thisHashedValue := h.doHash([]byte(value), []byte(salt), time, memory, parallelism, uint32(len(key)))

	// ConstantTimeCompare to help protect against timing attacks
	return subtle.ConstantTimeCompare([]byte(hash), []byte(thisHashedValue)) == 1, nil
}

func (h *Hasher) doHash(value []byte, salt []byte, time uint32, memory uint32, parallelism uint8, keyLength uint32) string {
	key := argon2.IDKey([]byte(value), salt, time, memory, parallelism, keyLength)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedKey := base64.RawStdEncoding.EncodeToString(key)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, parallelism, encodedSalt, encodedKey)
}
