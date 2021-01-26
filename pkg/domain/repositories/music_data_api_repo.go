package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/domain"
)

type musicDataApiRepo struct {
	client api.Client
}

func NewMusicDataApiRepo(client api.Client) *musicDataApiRepo {
	return &musicDataApiRepo{
		client: client,
	}
}

func (m *musicDataApiRepo) Get(ctx context.Context) (*domain.MusicData, error) {
	resRaw, err := m.client.GET(ctx, "/music-data")

	if err != nil {
		return nil, err
	}

	type resMusicDataTrack struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Artist struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"artist"`
		Duration     float64 `json:"duration"`
		ImagePreview *string `json:"imagePreview"`
		URL          string  `json:"url"`
	}

	var resMusicDataTrackToTrack = func(resTrack resMusicDataTrack) *domain.Track {
		// nolint:gomnd
		return &domain.Track{
			ID:    resTrack.ID,
			Title: resTrack.Title,
			Artist: domain.Artist{
				ID:   resTrack.Artist.ID,
				Name: resTrack.Artist.Name,
			},
			Duration:        time.Millisecond * time.Duration(resTrack.Duration*1000),
			ImagePreviewURL: resTrack.ImagePreview,
			MP3URL:          resTrack.URL,
		}
	}

	var resMusicData struct {
		Hash string `json:"hash"`
		CDN  string `json:"cdn"`
		Ads  []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			Times struct {
				Days  []int `json:"days"`
				Times []int `json:"times"`
			} `json:"times"`
			Track resMusicDataTrack `json:"track"`
		} `json:"ads"`
		Intervals []struct {
			Start    int   `json:"start"`
			End      int   `json:"end"`
			TrackIDs []int `json:"trackIds"`
		} `json:"intervals"`
		Tracks []resMusicDataTrack `json:"tracks"`
	}

	if err := json.Unmarshal(resRaw, &resMusicData); err != nil {
		return nil, fmt.Errorf("music data unmarshall fail: %w", err)
	}

	musicData := &domain.MusicData{
		Hash:      resMusicData.Hash,
		CdnURL:    resMusicData.CDN,
		Ads:       make([]*domain.Ad, len(resMusicData.Ads)),
		Intervals: make([]*domain.MusicDataInterval, len(resMusicData.Intervals)),
		Tracks:    make([]*domain.Track, len(resMusicData.Tracks)),
	}

	for i, ad := range resMusicData.Ads {
		musicData.Ads[i] = &domain.Ad{
			ID:    ad.ID,
			Title: ad.Title,
			Times: domain.AdTimes{
				Days:  ad.Times.Days,
				Times: ad.Times.Times,
			},
			Track: resMusicDataTrackToTrack(ad.Track),
		}
	}

	for i, interval := range resMusicData.Intervals {
		musicData.Intervals[i] = &domain.MusicDataInterval{
			Start:    interval.Start,
			End:      interval.End,
			TrackIDs: interval.TrackIDs,
		}
	}

	for i, track := range resMusicData.Tracks {
		musicData.Tracks[i] = resMusicDataTrackToTrack(track)
	}

	return musicData, nil
}
