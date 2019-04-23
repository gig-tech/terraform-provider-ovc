package ovc

import (
	"crypto"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwtLib "github.com/dgrijalva/jwt-go"
)

const (
	iyoPublicKeyStr = `
-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2
7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6
6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv
-----END PUBLIC KEY-----
`
	iyoRefreshURL = "https://itsyou.online/v1/oauth/jwt/refresh"
)

var (
	// ErrExpiredJWT represents an expired JWT error
	ErrExpiredJWT = fmt.Errorf("JWT is expired")
	// ErrInvalidJWT represents an invalid JWT error
	ErrInvalidJWT = fmt.Errorf("invalid JWT token")
	// ErrClaimNotPresent represents an error where a claim was not found in the token
	ErrClaimNotPresent = fmt.Errorf("claim was not found in the JWT token")

	jwtPublicKey        crypto.PublicKey
	expirationBuffer, _ = time.ParseDuration("5m")
)

func init() {
	err := SetJWTPublicKey(iyoPublicKeyStr)
	if err != nil {
		log.Fatalf("failed to parse pub key:%v", err)
	}
}

// SetJWTPublicKey configure the public key used to verify JWT token
func SetJWTPublicKey(key string) error {
	var err error
	jwtPublicKey, err = jwtLib.ParseECPublicKeyFromPEM([]byte(key))
	if err != nil {
		return err
	}
	return nil
}

// NewJWT returns a new JWT type
// supported identity providers:
// IYO (itsyou.online)
func NewJWT(jwtStr string, idProvider string) (*JWT, error) {
	token, err := parseJWT(jwtStr)
	if err != nil {
		return nil, err
	}

	var jwt = &JWT{}
	jwt.original = token

	switch idProvider {
	case "IYO":
		jwt.refreshFunc = getIYORefreshedJWT
	default:
		return nil, fmt.Errorf("unsupported identity provider. Supported providers are: IYO")
	}

	refreshable, err := isRefreshable(token)
	if err != nil {
		return nil, err
	}
	if refreshable {
		jwt.refreshable = true
	}

	return jwt, nil
}

func parseJWT(jwtStr string) (*jwtLib.Token, error) {
	log.Printf("[DEBUG] parsing JWT")
	parser := new(jwtLib.Parser)
	parser.SkipClaimsValidation = true
	return parser.Parse(jwtStr, func(token *jwtLib.Token) (interface{}, error) {
		if token.Method != jwtLib.SigningMethodES384 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtPublicKey, nil
	})
}

func isRefreshable(token *jwtLib.Token) (bool, error) {
	log.Printf("[DEBUG] checking if JWT is refreshable")
	claims, ok := token.Claims.(jwtLib.MapClaims)
	if !ok {
		return false, ErrInvalidJWT
	}
	_, ok = claims["refresh_token"]
	return ok, nil
}

// JWT represents a JWT
type JWT struct {
	original    *jwtLib.Token
	current     *jwtLib.Token
	refreshable bool
	refreshFunc func(jwtStr string) (string, error)
}

// Get returns the JWT
// If the JWT is expired (or nearly so) and refreshable, a refreshed token is returned
func (j *JWT) Get() (string, error) {
	err := j.refresh()
	return j.current.Raw, err
}

// Claim returns the value of Claim
func (j *JWT) Claim(key string) (interface{}, error) {
	log.Printf("[DEBUG] Checking for claim %s", key)
	token := j.current
	if token == nil {
		token = j.original
	}

	claims, ok := token.Claims.(jwtLib.MapClaims)
	if !ok {
		return nil, ErrClaimNotPresent
	}

	val, ok := claims[key]
	if !ok {
		return nil, ErrClaimNotPresent
	}

	return val, nil
}

// refresh refreshes the current JWT if expired (or nearly so) and the original JWT is refreshable
func (j *JWT) refresh() error {
	log.Printf("[DEBUG] checking to refresh the JWT")
	if j.current == nil {
		j.current = j.original
	}
	if isExpired(j.current) {
		log.Printf("[DEBUG] refreshing JWT")
		if !j.refreshable {
			return ErrExpiredJWT
		}

		newJWTStr, err := j.refreshFunc(j.original.Raw)
		if err != nil {
			return fmt.Errorf("Something went wrong refreshing the JWT: %s", err)
		}
		newToken, err := parseJWT(newJWTStr)
		if err != nil {
			return fmt.Errorf("Something went wrong parsing the refreshed JWT: %s", err)
		}
		j.current = newToken
	}

	return nil
}

func isExpired(token *jwtLib.Token) bool {
	exp, err := isExpiredWithErr(token)
	if err != nil {
		log.Printf("Something went wrong checking experation of JWT: %s\n", err)
	}

	return exp
}

func isExpiredWithErr(token *jwtLib.Token) (bool, error) {
	exp, err := getJWTExpiration(token)
	if err != nil {
		return false, err
	}
	if time.Until(time.Unix(exp, 0)) <= expirationBuffer {
		return true, nil
	}
	return false, nil
}

func getJWTExpiration(token *jwtLib.Token) (int64, error) {
	claims, ok := token.Claims.(jwtLib.MapClaims)
	if !(ok && token.Valid) {
		return 0, ErrInvalidJWT
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid expiration claim in token")
	}

	return int64(expFloat), nil
}

func getIYORefreshedJWT(token string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", iyoRefreshURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	return string(body), nil
}
