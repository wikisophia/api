package endpoints

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/wikisophia/api/server/accounts"
)

func postSessionHandler(key *ecdsa.PrivateKey, authenticator accounts.Authenticator) http.HandlerFunc {
	type request struct {
		Email    string
		Password string
	}
	type response struct {
		Token string `json:"token"`
	}
	responsePrefix := []byte(`{"token":"`)
	responseSuffix := []byte(`"}`)
	responseOverhead := len(responsePrefix) + len(responseSuffix)

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var req request
		if err = json.Unmarshal(body, &req); err != nil {
			http.Error(w, "invald request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Email == "" {
			http.Error(w, "missing required property: \"email\"", http.StatusBadRequest)
			return
		}
		if req.Password == "" {
			http.Error(w, "missing required property: \"password\"", http.StatusBadRequest)
			return
		}
		_, err = authenticator.Authenticate(context.Background(), req.Email, req.Password)
		if err != nil {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		jwt, err := newJwt(key, 1)
		if err != nil {
			http.Error(w, "error signing token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(jwt)+responseOverhead))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Set-Cookie", "auth="+jwt+"; SameSite=Strict; Secure; HttpOnly")
		w.WriteHeader(http.StatusOK)
		w.Write(responsePrefix)
		w.Write([]byte(jwt))
		w.Write(responseSuffix)
	}
}

const jwtHeader = `{"alg":"ES384","typ":"JWT"}`

const hashFunction = crypto.SHA384
const expectedCurveBitSize = 384
const keySize = expectedCurveBitSize / 8
const expectedSignatureSize = 2 * keySize

// newJwt makes a JWT for the given user, signing it with key.
func newJwt(key *ecdsa.PrivateKey, userID int64) (string, error) {
	header := base64.StdEncoding.EncodeToString([]byte(jwtHeader))
	payload := base64.StdEncoding.EncodeToString([]byte(`{"userId":` + strconv.FormatInt(userID, 10) + "}"))

	hash := hashFunction.New()
	hash.Write([]byte(header + "." + payload))
	rSig, sSig, err := ecdsa.Sign(rand.Reader, key, hash.Sum(nil))
	if err != nil {
		return "", err
	}

	return header + "." + payload + "." + encodeSignature(rSig, sSig), nil
}

// JWT is the contract class for our auth objects
type JWT struct {
	UserID int64 `json:"userId"`
}

// parseUserID returns the UserID claim from this jwt
func parseJwt(key *ecdsa.PublicKey, jwt string) (JWT, error) {
	jwtParts := strings.Split(jwt, ".")
	if len(jwtParts) != 3 {
		return JWT{}, errors.New("jwt parse failed: a jwt should have three parts separated by decimals")
	}
	header := jwtParts[0]
	claims := jwtParts[1]

	signature, err := base64.StdEncoding.DecodeString(string(jwtParts[2]))
	if err != nil {
		return JWT{}, errors.New("jwt parse failed: signature was not base64 encoded")
	}
	if len(signature) != expectedSignatureSize {
		return JWT{}, errors.New("jwt parse failed: invalid signature length")
	}
	sigR := big.NewInt(0).SetBytes(signature[:keySize])
	sigS := big.NewInt(0).SetBytes(signature[keySize:])

	hash := hashFunction.New()
	hash.Write([]byte(header))
	hash.Write([]byte("."))
	hash.Write([]byte(claims))
	if !ecdsa.Verify(key, hash.Sum(nil), sigR, sigS) {
		return JWT{}, errors.New("jwt rejected: signatures did not match")
	}
	decodedClaims, err := base64.StdEncoding.DecodeString(string(claims))
	if err != nil {
		return JWT{}, errors.New("error decoding jwt: claims were not base64 encoded")
	}

	var parsed JWT
	if err := json.Unmarshal([]byte(decodedClaims), &parsed); err != nil {
		return JWT{}, errors.New("malformed jwt: this shouldn't happen")
	}
	return parsed, nil
}

func encodeSignature(r, s *big.Int) string {
	return base64.StdEncoding.EncodeToString(append(bigIntToBytes(r), bigIntToBytes(s)...))
}

func bigIntToBytes(n *big.Int) []byte {
	b := n.Bytes()
	bPadded := make([]byte, keySize)
	copy(bPadded[keySize-len(b):], b)
	return bPadded
}
