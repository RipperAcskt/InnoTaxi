package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
)

type Token struct {
	Access           string `json:"access_token"`
	RT               string `json:"refresh_token"`
	AccessExpiration time.Time
	RTExpiration     time.Time
}

func NewToken(id uint64, cfg *config.Config) (*Token, error) {
	jwtExp := time.Now().Add(time.Duration(cfg.ACCESS_TOKEN_EXP) * time.Minute)

	access, accessExp, err := newJwt(jwtExp, id, cfg)
	if err != nil {
		return nil, fmt.Errorf("new jwt failed: %w", err)
	}

	jwtExp = time.Now().Add(time.Duration(cfg.REFRESH_TOKEN_EXP) * 24 * time.Hour)

	rt, rtExp, err := newJwt(jwtExp, id, cfg)
	if err != nil {
		return nil, fmt.Errorf("new rt failed: %w", err)
	}

	return &Token{access, rt, accessExp, rtExp}, nil
}

func newJwt(jwtExp time.Time, id uint64, cfg *config.Config) (string, time.Time, error) {
	header := make(map[string]string, 2)
	header["typ"] = "JWT"
	header["alg"] = "HS256"

	headerEncoded, err := json.Marshal(header)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal header failed: %w", err)
	}

	b64Header := base64.RawURLEncoding.EncodeToString(headerEncoded)

	playload := make(map[string]uint64, 2)

	playload["user_id"] = id
	playload["exp"] = uint64(jwtExp.UTC().Unix())

	playloadEncoded, err := json.Marshal(playload)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal playload failed: %w", err)
	}

	b64Playload := base64.RawURLEncoding.EncodeToString(playloadEncoded)

	data := fmt.Sprintf("%s.%s", b64Header, b64Playload)

	secret := []byte(cfg.HS256_SECRET)
	h := hmac.New(sha256.New, secret)

	_, err = h.Write([]byte(data))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sha256 write failed: %w", err)
	}

	sig := h.Sum(nil)
	b64Sig := base64.RawURLEncoding.EncodeToString(sig)

	jwt := fmt.Sprintf("%s.%s", data, b64Sig)

	return jwt, jwtExp, nil
}

func Verify(token string, cfg *config.Config) (bool, uint64, error) {
	rawSegs := strings.Split(token, ".")
	if len(rawSegs) != 3 {
		return false, 0, nil
	}

	b64header, err := base64.RawURLEncoding.DecodeString(rawSegs[0])
	if err != nil {
		return false, 0, fmt.Errorf("decode string failed: %w", err)
	}

	header := make(map[string]string, 2)
	err = json.Unmarshal(b64header, &header)
	if err != nil {
		return false, 0, fmt.Errorf("json unmarshal header error: %w", err)
	}

	if header["typ"] != "JWT" {
		return false, 0, nil
	}

	alg, ok := header["alg"]
	if !ok {
		return false, 0, nil
	}

	if alg != "HS256" {
		return false, 0, nil
	}

	b64playload, err := base64.RawURLEncoding.DecodeString(rawSegs[1])
	if err != nil {
		return false, 0, fmt.Errorf("decode string failed: %w", err)
	}

	playload := make(map[string]uint64, 2)
	err = json.Unmarshal(b64playload, &playload)
	if err != nil {
		return false, 0, fmt.Errorf("json unmarshal playload error: %w", err)
	}

	if playload["exp"] < uint64(time.Now().UTC().Unix()) {
		return false, 0, ErrTokenExpired
	}

	body := fmt.Sprintf("%s.%s", rawSegs[0], rawSegs[1])

	secret := []byte(cfg.HS256_SECRET)
	h := hmac.New(sha256.New, secret)

	_, err = h.Write([]byte(body))
	if err != nil {
		return false, 0, fmt.Errorf("sha256 write failed: %w", err)
	}

	sig, err := base64.RawURLEncoding.DecodeString(rawSegs[2])
	if err != nil {
		return false, 0, fmt.Errorf("decode string failed: %w", err)
	}

	expectedSig := h.Sum(nil)
	return hmac.Equal([]byte(sig), expectedSig), playload["user_id"], nil
}
