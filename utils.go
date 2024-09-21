package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"hash/crc32"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

func (a *App) FormatJSON(input string, compact bool) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", err
	}

	var output []byte
	if compact {
		output, err = json.Marshal(data)
	} else {
		output, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (a *App) TimestampToDate(timestamp int64, isMilliseconds bool) string {
	if isMilliseconds {
		timestamp /= 1000
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func (a *App) DateToTimestamp(date string, toMilliseconds bool) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return 0, err
	}
	timestamp := t.Unix()
	if toMilliseconds {
		timestamp *= 1000
	}
	return timestamp, nil
}

func (a *App) HashString(input string, algorithm string) string {
	var result []byte
	switch strings.ToLower(algorithm) {
	case "crc32":
		crc := crc32.ChecksumIEEE([]byte(input))
		result = []byte(hex.EncodeToString([]byte{byte(crc >> 24), byte(crc >> 16), byte(crc >> 8), byte(crc)}))
	case "md5":
		hash := md5.Sum([]byte(input))
		result = hash[:]
	case "sha1":
		hash := sha1.Sum([]byte(input))
		result = hash[:]
	case "sha256":
		hash := sha256.Sum256([]byte(input))
		result = hash[:]
	case "sha512":
		hash := sha512.Sum512([]byte(input))
		result = hash[:]
	case "sha3-256":
		hash := sha3.Sum256([]byte(input))
		result = hash[:]
	case "sha3-512":
		hash := sha3.Sum512([]byte(input))
		result = hash[:]
	default:
		return "Unsupported algorithm"
	}
	return hex.EncodeToString(result)
}

func (a *App) EncodeString(input string, encoding string) string {
	switch strings.ToLower(encoding) {
	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(input))
	case "url":
		return url.QueryEscape(input)
	default:
		return "Unsupported encoding"
	}
}

func (a *App) DecodeString(input string, encoding string) (string, error) {
	switch strings.ToLower(encoding) {
	case "base64":
		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	case "url":
		decoded, err := url.QueryUnescape(input)
		if err != nil {
			return "", err
		}
		return decoded, nil
	default:
		return "Unsupported decoding", nil
	}
}
