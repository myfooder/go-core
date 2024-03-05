package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var (
	SigningMethodES256 = jwt.SigningMethodES256
	SigningMethodES384 = jwt.SigningMethodES384
	SigningMethodES512 = jwt.SigningMethodES512
	SigningMethodEdDSA = jwt.SigningMethodEdDSA
	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS384 = jwt.SigningMethodHS384
	SigningMethodHS512 = jwt.SigningMethodHS512
	SigningMethodRS256 = jwt.SigningMethodRS256
	SigningMethodRS384 = jwt.SigningMethodRS384
	SigningMethodRS512 = jwt.SigningMethodRS512
)

type TokenManager struct {
	method     jwt.SigningMethod
	publicKey  interface{}
	privateKey interface{}
}

func NewTokenManager(method jwt.SigningMethod, cfg *Config) (*TokenManager, error) {

	// parse keys
	publicKey, privateKey := parseKeys(cfg)

	// parse keys
	var parsedPublicKey interface{}
	var parsedPublicKeyErr error
	var parsedPrivateKey interface{}
	var parsedPrivateKeyErr error
	switch method.(type) {
	case *jwt.SigningMethodECDSA:
		if len(publicKey) > 0 {
			parsedPublicKey, parsedPublicKeyErr = jwt.ParseECPublicKeyFromPEM(publicKey)
		}
		if len(privateKey) > 0 {
			parsedPrivateKey, parsedPrivateKeyErr = jwt.ParseECPrivateKeyFromPEM(privateKey)
		}
		break
	case *jwt.SigningMethodEd25519:
		if len(publicKey) > 0 {
			parsedPublicKey, parsedPublicKeyErr = jwt.ParseEdPublicKeyFromPEM(publicKey)
		}
		if len(privateKey) > 0 {
			parsedPrivateKey, parsedPrivateKeyErr = jwt.ParseEdPrivateKeyFromPEM(privateKey)
		}
		break
	case *jwt.SigningMethodHMAC:
		if len(publicKey) > 0 {
			parsedPublicKey = publicKey
		}
		if len(privateKey) > 0 {
			parsedPrivateKey = privateKey
		}
	case *jwt.SigningMethodRSA:
		if len(publicKey) > 0 {
			parsedPublicKey, parsedPublicKeyErr = jwt.ParseRSAPublicKeyFromPEM(publicKey)
		}
		if len(privateKey) > 0 {
			parsedPrivateKey, parsedPrivateKeyErr = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
		}
		break
	case *jwt.SigningMethodRSAPSS:
		if len(publicKey) > 0 {
			parsedPublicKey, parsedPublicKeyErr = jwt.ParseRSAPublicKeyFromPEM(publicKey)
		}
		if len(privateKey) > 0 {
			parsedPrivateKey, parsedPrivateKeyErr = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
		}
		break
	default:
		return nil, errors.New("unsupported signing method")
	}
	if parsedPublicKeyErr != nil {
		return nil, errors.Join(errors.New("error while parsing public key"), parsedPublicKeyErr)
	}
	if parsedPrivateKeyErr != nil {
		return nil, errors.Join(errors.New("error while parsing private key"), parsedPrivateKeyErr)
	}

	// return helper
	return &TokenManager{
		method:     method,
		publicKey:  parsedPublicKey,
		privateKey: parsedPrivateKey,
	}, nil

}

func parseKeys(cfg *Config) (publicKey []byte, privateKey []byte) {

	if len(cfg.PublicKeyFile) > 0 {
		publicKeyBytes, err := os.ReadFile(cfg.PublicKeyFile)
		if err == nil {
			publicKey = publicKeyBytes
		}
	}

	if len(cfg.PrivateKeyFile) > 0 {
		privateKeyBytes, err := os.ReadFile(cfg.PrivateKeyFile)
		if err == nil {
			privateKey = privateKeyBytes
		}
	}

	return publicKey, privateKey

}

func (m *TokenManager) SignedToken(claims Claims) (string, error) {
	token := m.NewToken(claims)
	return m.SignedString(token)
}

func (m *TokenManager) NewToken(claims Claims) *jwt.Token {
	return jwt.NewWithClaims(m.method, claims)
}

func (m *TokenManager) SignedString(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *TokenManager) Parse(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// TODO verify that token.Method matches JWT.method
		return m.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil

}

func (m *TokenManager) GetClaims(token *jwt.Token) (claims Claims) {
	return token.Claims.(Claims)
}
