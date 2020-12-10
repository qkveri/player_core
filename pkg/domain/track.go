package domain

import "time"

type Track struct {
	ID              int
	Title           string
	Artist          Artist
	Duration        time.Duration
	ImagePreviewURL *string
	MP3URL          string
}
