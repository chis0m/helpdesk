package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2IDParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var params = Argon2IDParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("password cannot be empty")
	}

	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	if strings.TrimSpace(password) == "" {
		return false, errors.New("password cannot be empty")
	}
	if strings.TrimSpace(encodedHash) == "" {
		return false, errors.New("encoded hash cannot be empty")
	}

	cfg, saltBytes, hashBytes, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	hashToCompare := argon2.IDKey(
		[]byte(password),
		saltBytes,
		cfg.Iterations,
		cfg.Memory,
		cfg.Parallelism,
		uint32(len(hashBytes)),
	)

	return subtle.ConstantTimeCompare(hashToCompare, hashBytes) == 1, nil
}

func decodeHash(encodedHash string) (Argon2IDParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("invalid hash format")
	}
	if parts[1] != "argon2id" {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("unsupported algorithm: %s", parts[1])
	}

	var version uint32
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("invalid hash version: %w", err)
	}
	if version != argon2.Version {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("incompatible hash version: %d", version)
	}

	var cfg Argon2IDParams
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &cfg.Memory, &cfg.Iterations, &cfg.Parallelism); err != nil {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("invalid hash parameters: %w", err)
	}

	saltBytes, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("invalid salt: %w", err)
	}
	cfg.SaltLength = uint32(len(saltBytes))

	hashBytes, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Argon2IDParams{}, nil, nil, fmt.Errorf("invalid hash: %w", err)
	}
	cfg.KeyLength = uint32(len(hashBytes))

	return cfg, saltBytes, hashBytes, nil
}