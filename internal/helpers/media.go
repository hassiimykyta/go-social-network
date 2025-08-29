package helpers

import "strings"

func SanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "file"
	}
	return name
}

func InferKind(ct string) string {
	ct = strings.ToLower(ct)
	switch {
	case strings.HasPrefix(ct, "image/"):
		if ct == "image/svg+xml" {
			return ""
		}
		return "image"
	case strings.HasPrefix(ct, "video/"):
		return "video"
	case strings.HasPrefix(ct, "audio/"):
		return "audio"
	default:
		return ""
	}
}
