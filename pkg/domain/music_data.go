package domain

import "context"

type (
	MusicData struct {
		Hash      string
		CdnURL    string
		Ads       []Ad
		Intervals []MusicDataInterval
		Tracks    []Track
	}

	MusicDataInterval struct {
		Start    int
		End      int
		TrackIDs []int
	}

	MusicDataRepository interface {
		Get(ctx context.Context) (*MusicData, error)
	}
)
