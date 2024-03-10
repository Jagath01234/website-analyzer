package util

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid URL with http",
			url:      "http://www.google.com",
			expected: true,
		},
		{
			name:     "Valid URL with https",
			url:      "https://www.google.com",
			expected: true,
		},
		{
			name:     "Valid URL with path",
			url:      "https://www.mail.google.com/mail/u/0/#inbox",
			expected: true,
		},
		{
			name:     "Invalid URL - missing http/https",
			url:      "www.google.com",
			expected: false,
		},
		{
			name:     "Invalid URL - incorrect format",
			url:      "http://www..google.com",
			expected: false,
		},
		{
			name:     "Invalid URL - incorrect prefix",
			url:      "ftp://www.google.com",
			expected: false,
		},
		{
			name:     "Invalid URL - incorrect prefix",
			url:      "ftp://www.%google.com",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid := IsValidURL(test.url)
			if valid != test.expected {
				t.Errorf("Expected IsValidURL(%s) to be %v, but got %v", test.url, test.expected, valid)
			}
		})
	}
}
