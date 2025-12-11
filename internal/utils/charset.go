package utils

import (
	"bytes"
	"io"
	"gomoco/internal/models"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// ConvertCharset converts string content to the specified charset
func ConvertCharset(content string, charset string) ([]byte, error) {
	switch charset {
	case models.CharsetGBK:
		// Convert UTF-8 to GBK
		encoder := simplifiedchinese.GBK.NewEncoder()
		reader := transform.NewReader(bytes.NewReader([]byte(content)), encoder)
		result, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return result, nil
	case models.CharsetUTF8:
		// Already UTF-8
		return []byte(content), nil
	default:
		return []byte(content), nil
	}
}
