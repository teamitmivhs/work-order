package utils

import "os"

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
