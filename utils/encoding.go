package utils

import (
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Function to convert ASCII codes to string
func AsciiToString(asciiCodes []byte) string {
	var result strings.Builder

	for _, code := range asciiCodes {
		// Convert each ASCII code to corresponding character and append to result
		result.WriteByte(byte(code))
	}

	return result.String()
}

func Utf8ToGBK(utf8Str string) (string, error) {
	// Create GBK Encoder
	encoder := simplifiedchinese.GBK.NewEncoder()

	// Transform the string
	gbkBytes, err := io.ReadAll(transform.NewReader(strings.NewReader(utf8Str), encoder))
	if err != nil {
		return "", err
	}

	// Convert the bits back to string
	return string(gbkBytes), nil
}

// GBKToUtf8 converts a GBK encoded string to a UTF-8 string
func GBKToUtf8(gbkStr string) (string, error) {
	// Create GBK Decoder
	decoder := simplifiedchinese.GBK.NewDecoder()

	// Transform the string
	utf8Bytes, err := io.ReadAll(transform.NewReader(strings.NewReader(gbkStr), decoder))
	if err != nil {
		return "", err
	}

	// Convert the bytes back to string
	return string(utf8Bytes), nil
}
