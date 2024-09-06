package util

import (
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
)

// CheckIsUrl checks if the given string is a URL.
func CheckIsUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	// Check if the scheme or host is valid.
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// CheckIsImageUrl checks if the given string is a URL of an image.
func CheckIsImageUrl(str string) bool {
	// Check if the URL is valid.
	if !CheckIsUrl(str) {
		return false
	}

	// Check if the URL is an image.
	return imageRegex.MatchString(str)
}

// CheckIsContainsChinese checks if the given string contains Chinese characters.
func CheckIsContainsChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}

	return false
}

// CheckIsContainsEnglish checks if the given string contains English characters.
func CheckIsContainsEnglish(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Latin, r) {
			return true
		}
	}

	return false
}

// CheckIsPinyin checks if the given string is a Pinyin.
func CheckIsPinyin(str string) bool {
	str = strings.Replace(str, " ", "", -1)
	return pinyinRegex.MatchString(str)
}

// CheckIsHash checks if the given string is a hash.
// Supported hash types: MD5, SHA-1, SHA-256.
func CheckIsHash(str string) bool {
	if len(str) == 0 {
		return false
	}

	return md5Regex.MatchString(str) || sha1Regex.MatchString(str) || sha256Regex.MatchString(str)
}

// CheckIsInt checks if the given string is an integer.
func CheckIsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

// CheckIsFloat checks if the given string is a float.
func CheckIsFloat(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
