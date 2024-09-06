package util

import "testing"

func TestCheckIsUrl(t *testing.T) {
	type testData struct {
		url      string
		expected bool
	}

	tests := []testData{
		{
			url:      "https://www.google.com",
			expected: true,
		},
		{
			url:      "http://www.google.com",
			expected: true,
		},
		{
			url:      "www.google.com",
			expected: false,
		},
		{
			url:      "google.com",
			expected: false,
		},
		{
			url:      "google",
			expected: false,
		},
		{
			url:      "https://",
			expected: false,
		},
		{
			url:      "http://www",
			expected: true,
		},
	}

	for _, test := range tests {
		if !checkAsExpected(CheckIsUrl(test.url), test.expected) {
			t.Fatalf("expected %v, got %v", test.expected, CheckIsUrl(test.url))
		}
	}
}

func TestCheckIsImageUrl(t *testing.T) {
	type testData struct {
		url      string
		expected bool
	}

	tests := []testData{
		{
			url:      "https://www.google.com",
			expected: false,
		},
		{
			url:      "image.png",
			expected: false,
		},
		{
			url:      "https://www.google.com/image.jpg",
			expected: true,
		},
		{
			url:      "https://www.google.com/image.jpeg",
			expected: true,
		},
		{
			url:      "https://www.google.com/image.png",
			expected: true,
		},
		{
			url:      "https://www.google.com/gif.image",
			expected: false,
		},
	}

	for _, test := range tests {
		if !checkAsExpected(CheckIsImageUrl(test.url), test.expected) {
			t.Fatalf("expected %v, got %v", test.expected, CheckIsImageUrl(test.url))
		}
	}
}

func TestCheckIsContainsChinese(t *testing.T) {
	type testData struct {
		str      string
		expected bool
	}

	tests := []testData{
		{
			str:      "你好",
			expected: true,
		},
		{
			str:      "hello",
			expected: false,
		},
		{
			str:      "你好hello",
			expected: true,
		},
		{
			str:      "hello你好",
			expected: true,
		},
		{
			str:      "hello こんにちは world",
			expected: false,
		},
		{
			str:      "hello こんにちは 你好 world",
			expected: true,
		},
	}

	for _, test := range tests {
		if !checkAsExpected(CheckIsContainsChinese(test.str), test.expected) {
			t.Fatalf("expected %v, got %v", test.expected, CheckIsContainsChinese(test.str))
		}
	}
}

func TestCheckIsContainsEnglish(t *testing.T) {
	type testData struct {
		str      string
		expected bool
	}

	tests := []testData{
		{
			str:      "你好",
			expected: false,
		},
		{
			str:      "hello",
			expected: true,
		},
		{
			str:      "你好hello",
			expected: true,
		},
		{
			str:      "hello你好",
			expected: true,
		},
		{
			str:      "hello こんにちは world",
			expected: true,
		},
		{
			str:      "hello こんにちは 你好 world",
			expected: true,
		},
	}

	for _, test := range tests {
		if !checkAsExpected(CheckIsContainsEnglish(test.str), test.expected) {
			t.Fatalf("expected %v, got %v", test.expected, CheckIsContainsEnglish(test.str))
		}
	}
}

func TestCheckIsPinyin(t *testing.T) {
	type testData struct {
		str      string
		expected bool
	}

	tests := []testData{
		{
			str:      "你好",
			expected: false,
		},
		{
			str:      "hello",
			expected: true,
		},
		{
			str:      "你好hello",
			expected: false,
		},
		{
			str:      "xǔ",
			expected: true,
		},
		{
			str:      "fàng hello",
			expected: true,
		},
		{
			str:      "hello こんにちは fàng",
			expected: false,
		},
	}

	for _, test := range tests {
		if !checkAsExpected(CheckIsPinyin(test.str), test.expected) {
			t.Fatalf("str %s expected %v, got %v", test.str, test.expected, CheckIsPinyin(test.str))
		}
	}
}

func TestCheckIsHash(t *testing.T) {
	type testData struct {
		str      string
		expected bool
	}

	tests := []testData{
		{
			str:      "你好",
			expected: false,
		},
		{
			str:      "hello",
			expected: false,
		},
		{
			str:      "spjj47cdb84968e22b69ba874fde4d3cfdae",
			expected: false,
		},
		{
			str:      "spjj47cdb84968e22b69ba874fde4d3cfda",
			expected: false,
		},
		{
			str:      "e99a18c428cb38d5f260853678922e03",
			expected: true,
		},
		{
			str:      "6dcd4ce23d88d1f7f9d453b78b9cde65e3f5a6d8499f2c1c3d49a72b0c8e8d2",
			expected: false,
		},
	}

	for _, test := range tests {
		if !checkAsExpected(CheckIsHash(test.str), test.expected) {
			t.Fatalf("str %s expected %v, got %v", test.str, test.expected, CheckIsHash(test.str))
		}
	}
}

func checkAsExpected(err error, expected bool) bool {
	return (err == nil) == expected
}
