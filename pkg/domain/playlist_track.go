package domain

import "github.com/qkveri/player_core/pkg/progress"

type PlaylistTrackType string

const (
	PlaylistTrackTypeBackground = "background"
	PlaylistTrackTypeAd         = "ad"
)

type PlaylistTrack struct {
	Track                   *Track
	Type                    PlaylistTrackType
	BackgroundIntervalIndex int

	DownloadProgress progress.Progress
	FilePath         string
}
