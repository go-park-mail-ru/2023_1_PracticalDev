package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrBadTokenData = errors.New("bad token data")
	ErrBadTokenTime = errors.New("bad token time")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenDecode  = errors.New("can't hex decode token")
)

type SessionParams struct {
	Token      string
	LivingTime time.Duration
}

type HashToken struct {
	Secret []byte
}

func NewHMACHashToken(secret string) *HashToken {
	return &HashToken{Secret: []byte(secret)}
}

func (tk *HashToken) Create(s *SessionParams, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s$%d", s.Token, tokenExpTime)
	h.Write([]byte(data))
	token := fmt.Sprintf("%s$%s", hex.EncodeToString(h.Sum(nil)), strconv.FormatInt(tokenExpTime, 10))
	return token, nil
}

func (tk *HashToken) Check(s *SessionParams, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, "$")
	if len(tokenData) != 2 {
		return false, ErrBadTokenData
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, ErrBadTokenTime
	}

	if tokenExp < time.Now().Unix() {
		return false, ErrTokenExpired
	}

	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s$%d", s.Token, tokenExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, ErrTokenDecode
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
