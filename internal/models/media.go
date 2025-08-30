package models

import "time"

type Media struct {
	ID         int64      `json:"id"`
	OwnerID    int64      `json:"owner_id"`
	Kind       string     `json:"kind"`
	StorageKey string     `json:"storage_key"`
	MimeType   string     `json:"mime_type"`
	SizeBytes  int64      `json:"size_bytes"`
	Width      *int32     `json:"width,omitempty"`
	Height     *int32     `json:"height,omitempty"`
	DurationMs *int32     `json:"duration_ms,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type MediaPublic struct {
	ID         int64  `json:"id"`
	Kind       string `json:"kind"`
	MimeType   string `json:"mime_type"`
	URL        string `json:"url"`
	Width      *int32 `json:"width,omitempty"`
	Height     *int32 `json:"height,omitempty"`
	DurationMs *int32 `json:"duration_ms,omitempty"`
}

func (m Media) PublicWithURL(u string) MediaPublic {
	return MediaPublic{
		ID:       m.ID,
		Kind:     m.Kind,
		MimeType: m.MimeType,
		URL:      u,
		Width:    m.Width,
		Height:   m.Height,
	}
}
