package util

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	imageRegex  = regexp.MustCompile(`\.(jpg|jpeg|png|gif|bmp)$`)
	pinyinRegex = regexp.MustCompile(`^[a-zA-Zāáǎàēéěèīíǐìōóǒòūúǔùǖǘǚǜü]+$`)
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
