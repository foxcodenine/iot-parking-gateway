package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

type EnvHandler struct {
	SecretKey string // Symmetric encryption key
}

func (e *EnvHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Initialize the map
	environmentVariables := make(map[string]string)

	// Lookup the environment variable
	googleApiKey, ok := os.LookupEnv("GOOGLE_API_KEY")
	if ok {
		encryptedKey, err := e.encrypt(googleApiKey)
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to encrypt API key", http.StatusInternalServerError)
			helpers.LogError(err, "Failed to encrypt API key")
			return
		}
		environmentVariables["GOOGLE_API_KEY"] = encryptedKey
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(environmentVariables)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		helpers.LogError(err, "Failed to encode response")
		return
	}
}

// Encrypt encrypts the plaintext using AES encryption
func (e *EnvHandler) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(e.SecretKey))
	if err != nil {
		return "", err
	}

	// Generate a new IV for each encryption
	iv := make([]byte, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// Encode the encrypted data and IV as base64 for safe transport
	encoded := base64.StdEncoding.EncodeToString(append(iv, ciphertext...))
	return encoded, nil
}
