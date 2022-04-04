package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/johnjones4/hal-9000/hal9000/core"
)

var (
	ErrorTokenExpired = errors.New("token is expired")
)

type TokenManager struct {
	c cipher.Block
}

func NewTokenManager(key []byte) (*TokenManager, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &TokenManager{c}, nil
}

type tokenInternal struct {
	Username   string    `json:"username"`
	Expiration time.Time `json:"expiration"`
}

func (t *TokenManager) NewToken(user core.User) (core.Token, error) {
	tokenStruct := tokenInternal{user.Name, time.Now().Add(time.Hour)}

	tokenPlain, err := json.Marshal(tokenStruct)
	if err != nil {
		return core.Token{}, err
	}

	gcm, err := cipher.NewGCM(t.c)
	if err != nil {
		return core.Token{}, err
	}

	nonce := make([]byte, gcm.NonceSize())

	bytes := gcm.Seal(nonce, nonce, tokenPlain, nil)

	return core.Token{
		Token:      hex.EncodeToString(bytes),
		Expiration: tokenStruct.Expiration,
		User:       user.Name,
	}, nil
}

func (t *TokenManager) UsernameForToken(token string) (string, error) {
	ciphertext, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("invlaid size: %d (expected %d)", len(ciphertext), nonceSize)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	var it tokenInternal
	err = json.Unmarshal(plaintext, &it)
	if err != nil {
		return "", err
	}

	if it.Expiration.Before(time.Now()) {
		return "", ErrorTokenExpired
	}

	return it.Username, nil
}
