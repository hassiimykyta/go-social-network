package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

func BuildPostKey(userID int64, filename string, now time.Time) string {
	ext := strings.ToLower(filename)
	if ext == "" {
		ext = ".bin"
	}
	id := ulid.Make().String()
	return fmt.Sprintf("posts/%04d/%02d/%d/%s%s", now.Year(), int(now.Month()), userID, id, ext)
}
