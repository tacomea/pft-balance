package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strings"
)

func CreateToken(sessionId string) string {
	hash := hmac.New(sha256.New, []byte("this is the key"))
	hash.Write([]byte(sessionId))

	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	log.Println(signature)
	return signature + "|" + sessionId
}

func ParseToken(token string) (string, error) {
	ss := strings.SplitN(token, "|", 2)
	if len(ss) != 2 {
		return "", errors.New("cookie was changed")
	}

	decoded, err := base64.StdEncoding.DecodeString(ss[0])
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte("this is the key"))
	h.Write([]byte(ss[1]))

	if !hmac.Equal(decoded, h.Sum(nil)) {
		return "", errors.New("cookie was changed")
	}
	return ss[1], nil
}