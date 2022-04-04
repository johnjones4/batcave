package api

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"main/core"
	"os"
	"time"
)

var (
	errorTokenExpired = errors.New("token is expired")
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

type token struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

type tokenInternal struct {
	Username   string    `json:"username"`
	Expiration time.Time `json:"expiration"`
}

func (t *TokenManager) newToken(user core.User) (token, error) {
	tokenStruct := tokenInternal{user.Name, time.Now().Add(time.Hour)}

	tokenPlain, err := json.Marshal(tokenStruct)
	if err != nil {
		return token{}, err
	}

	gcm, err := cipher.NewGCM(t.c)
	if err != nil {
		return token{}, err
	}

	nonce := make([]byte, gcm.NonceSize())

	bytes := gcm.Seal(nonce, nonce, tokenPlain, nil)

	return token{hex.EncodeToString(bytes), tokenStruct.Expiration}, nil
}

func (t *TokenManager) usernameForToken(token string) (string, error) {
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
		return "", errorTokenExpired
	}

	return it.Username, nil
}
