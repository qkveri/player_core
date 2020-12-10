package repositories

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/qkveri/player_core/pkg/domain"
)

type authFileRepo struct {
	filePath string
	key      []byte
}

func NewAuthFileRepo(filePath string, key string) *authFileRepo {
	return &authFileRepo{
		filePath: filePath,
		key:      []byte(key),
	}
}

func (a *authFileRepo) Set(_ context.Context, auth *domain.Auth) error {
	rawData, err := json.Marshal(auth)

	if err != nil {
		return fmt.Errorf("connot marshal auth: %w", err)
	}

	gcm, err := a.createGCM()

	if err != nil {
		return fmt.Errorf("connot create gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = rand.Read(nonce); err != nil {
		return fmt.Errorf("rand.Read failed: %w", err)
	}

	cipherData := gcm.Seal(nonce, nonce, rawData, nil)

	if err := ioutil.WriteFile(a.filePath, cipherData, 0600); err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	return nil
}

func (a *authFileRepo) Get(_ context.Context) (*domain.Auth, error) {
	rawData, err := ioutil.ReadFile(a.filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}

	gcm, err := a.createGCM()

	if err != nil {
		return nil, fmt.Errorf("connot create gcm: %w", err)
	}

	nonce, ciphertext := rawData[:gcm.NonceSize()], rawData[gcm.NonceSize():]

	plainData, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, fmt.Errorf("gcm.Open failed: %w", err)
	}

	auth := &domain.Auth{}

	if err := json.Unmarshal(plainData, auth); err != nil {
		return nil, fmt.Errorf("connot unmarshal auth: %w, plainData: %s", err, plainData)
	}

	return auth, nil
}

func (a *authFileRepo) createGCM() (cipher.AEAD, error) {
	blockCipher, err := aes.NewCipher(a.key)

	if err != nil {
		return nil, fmt.Errorf("cannot create chipher: %w", err)
	}

	gcm, err := cipher.NewGCM(blockCipher)

	if err != nil {
		return nil, fmt.Errorf("cannot create GCM: %w", err)
	}

	return gcm, nil
}
