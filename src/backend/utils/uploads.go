package utils

import (
	"fmt"
	"html"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// PublicUploadDir returns the filesystem directory used for files served under
// /static/public. In Docker this is mounted at /static/public; in local Go runs
// the backend serves ../static, so ../static/public is the correct fallback.
func PublicUploadDir() string {
	if dir := os.Getenv("PUBLIC_UPLOAD_DIR"); dir != "" {
		return dir
	}
	if info, err := os.Stat("/static"); err == nil && info.IsDir() {
		return "/static/public"
	}
	return "../static/public"
}

func SanitizeText(value string) string {
	return html.EscapeString(strings.TrimSpace(value))
}

func ValidateImageUpload(file multipart.File, allowedTypes map[string]bool) error {
	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read uploaded file: %w", err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset uploaded file: %w", err)
	}
	contentType := http.DetectContentType(head[:n])
	if !allowedTypes[contentType] {
		return fmt.Errorf("invalid image content type: %s", contentType)
	}
	return nil
}
