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
		if CheckIsUrl(test.url) != test.expected {
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
		if CheckIsImageUrl(test.url) != test.expected {
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
		if CheckIsContainsChinese(test.str) != test.expected {
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
		if CheckIsContainsEnglish(test.str) != test.expected {
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
		if CheckIsPinyin(test.str) != test.expected {
			t.Fatalf("str %s expected %v, got %v", test.str, test.expected, CheckIsPinyin(test.str))
		}
	}
}
