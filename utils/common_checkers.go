package util

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	imageRegex  = regexp.MustCompile(`\.(jpg|jpeg|png|gif|bmp)$`)
	pinyinRegex = regexp.MustCompile(`^[a-zA-Zāáǎàēéěèīíǐìōóǒòūúǔùǖǘǚǜü]+$`)
	md5Regex    = regexp.MustCompile(`^[a-f0-9]{32}$`)
	sha1Regex   = regexp.MustCompile(`^[a-f0-9]{40}$`)
	sha256Regex = regexp.MustCompile(`^[a-f0-9]{64}$`)

	ErrInvalidURL      = errors.New("invalid URL")
	ErrInvalidImageURL = errors.New("invalid image URL")
	ErrInvalidPinyin   = errors.New("invalid Pinyin")
	ErrInvalidHash     = errors.New("invalid hash")
	ErrInvalidInt      = errors.New("invalid integer")
	ErrInvalidFloat    = errors.New("invalid float")
	ErrInvalidChinese  = errors.New("invalid Chinese")
	ErrInvalidEnglish  = errors.New("invalid English")
)

// CheckIsUrl checks if the given string is a URL.
func CheckIsUrl(str string) error {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return ErrInvalidURL
	}

	// Check if the scheme or host is valid.
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

// CheckIsImageUrl checks if the given string is a URL of an image.
func CheckIsImageUrl(str string) error {
	// Check if the URL is valid.
	if CheckIsUrl(str) != nil {
		return ErrInvalidURL
	}

	// Check if the URL is an image.
	if !imageRegex.MatchString(str) {
		return ErrInvalidImageURL
	}

	return nil
}

// CheckIsContainsChinese checks if the given string contains Chinese characters.
func CheckIsContainsChinese(str string) error {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return nil
		}
	}

	return ErrInvalidChinese
}

// CheckIsContainsEnglish checks if the given string contains English characters.
func CheckIsContainsEnglish(str string) error {
	for _, r := range str {
		if unicode.Is(unicode.Latin, r) {
			return nil
		}
	}

	return ErrInvalidEnglish
}

// CheckIsPinyin checks if the given string is a Pinyin.
func CheckIsPinyin(str string) error {
	str = strings.Replace(str, " ", "", -1)
	if !pinyinRegex.MatchString(str) {
		return ErrInvalidPinyin
	}

	return nil
}

// CheckIsHash checks if the given string is a hash.
// Supported hash types: MD5, SHA-1, SHA-256.
func CheckIsHash(str string) error {
	if len(str) == 0 {
		return ErrInvalidHash
	}

	if md5Regex.MatchString(str) || sha1Regex.MatchString(str) || sha256Regex.MatchString(str) {
		return nil
	}

	return ErrInvalidHash
}

// CheckIsInt checks if the given string is an integer.
func CheckIsInt(str string) error {
	_, err := strconv.Atoi(str)
	return err
}

// CheckIsFloat checks if the given string is a float.
func CheckIsFloat(str string) error {
	_, err := strconv.ParseFloat(str, 64)
	return err
}
