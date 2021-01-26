package playlister

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/rs/zerolog"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"

	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/state"
)

const (
	trackCount              = 5
	updateDuration          = time.Second
	numberSelectionAttempts = 30

	secondsInDay = 86400
)

type service struct {
	musicData *domain.MusicData

	state      *state.State
	logger     zerolog.Logger
	clock      clockwork.Clock
	trackCount int
}

func NewService(state *state.State, logger zerolog.Logger, clock clockwork.Clock) *service {
	return &service{
		state:      state,
		logger:     logger.With().Str("service", "playlister").Logger(),
		clock:      clock,
		trackCount: trackCount,
	}
}

func (s *service) Run(ctx context.Context) error {
	s.logger.Debug().Msg("starts up")
	defer s.logger.Debug().Msg("stopped")

	rand.Seed(time.Now().UnixNano())

	t := time.NewTicker(updateDuration)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-t.C:
			s.state.MusicData.RLock()

			md := s.state.MusicData.Get()

			if md != s.musicData {
				s.musicData = md

				s.updateList(true)
			} else {
				s.updateList(false)
			}

			s.state.MusicData.RUnlock()
		}
	}
}

func (s *service) updateList(force bool) {
	s.logger.Debug().Bool("force", force).Msg("starting update list...")

	if s.musicData == nil {
		s.logger.Debug().Msg("update list skipped (no musicData)")
		return
	}

	// lock playlist...
	s.state.Playlist.Lock()
	defer s.state.Playlist.Unlock()

	now := s.clock.Now()
	seconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
	addSeconds := func(d time.Duration) { seconds = (seconds + int(d.Seconds())) % secondsInDay }
	playlist := s.state.Playlist.Get()
	updated := 0

	for i := 0; i < s.trackCount; i++ {
		trackIndex := i
		intervalIndex := s.intervalIndexBySeconds(seconds)

		if !force {
			// если трек на нужном месте - пропускам
			if len(playlist) > trackIndex && playlist[trackIndex].BackgroundIntervalIndex == intervalIndex {
				addSeconds(playlist[trackIndex].Track.Duration)
				continue
			}
		}

		newTrack, err := s.getPlaylistTrackByIntervalIndex(intervalIndex)

		if err != nil {
			s.logger.Err(err).
				Int("trackIndex", trackIndex).
				Int("intervalIndex", intervalIndex).
				Msg("getPlaylistTrackByIntervalIndex fail")

			continue
		}

		if len(playlist) > trackIndex {
			s.state.Playlist.Replace(newTrack, trackIndex)
		} else {
			s.state.Playlist.Append(newTrack)
		}

		addSeconds(newTrack.Track.Duration)

		updated++
	}

	s.logger.Debug().Int("updated", updated).Msg("updated playlist")
}

func (s *service) intervalIndexBySeconds(seconds int) int {
	for index, interval := range s.musicData.Intervals {
		// переходящий интервал (со дня в другой день)
		if interval.Start > interval.End {
			if interval.Start <= seconds && seconds < 86400 || 0 <= seconds && seconds < interval.End {
				return index
			}
		}

		if interval.Start <= seconds && seconds < interval.End {
			return index
		}
	}

	return 0
}

var intervalsEmpty = errors.New("musicData.Interval not exists")

func (s *service) getTrackIdByIntervalIndex(index int) (int, error) {
	if index >= len(s.musicData.Intervals) {
		return 0, intervalsEmpty
	}

	playlist := s.state.Playlist.Get()
	currentTrackIDs := make(map[int]struct{}, len(playlist))

	for _, t := range playlist {
		currentTrackIDs[t.Track.ID] = struct{}{}
	}

	// Проверяем есть ли трек в списке.
	// Треки берутся рандомно, если добавляемый трек есть в списке, то мы рекурсивно пробуем другой.
	var randTrackID int

	for ai := 0; ai < numberSelectionAttempts; ai++ {
		//nolint:gosec
		randTrackIndex := rand.Intn(len(s.musicData.Intervals[index].TrackIDs)-0) + 0
		randTrackID = s.musicData.Intervals[index].TrackIDs[randTrackIndex]

		if _, has := currentTrackIDs[randTrackID]; !has {
			break
		}
	}

	return randTrackID, nil
}

func (s *service) getPlaylistTrackByIntervalIndex(index int) (*domain.PlaylistTrack, error) {
	if index >= len(s.musicData.Intervals) {
		return nil, errors.New("musicData.Interval not exists")
	}

	playlist := s.state.Playlist.Get()
	currentTrackIDs := make(map[int]struct{}, len(playlist))

	for _, t := range playlist {
		currentTrackIDs[t.Track.ID] = struct{}{}
	}

	// Проверяем есть ли трек в списке.
	// Треки берутся рандомно, если добавляемый трек есть в списке, то мы рекурсивно пробуем другой.
	var randTrackID int

	for ai := 0; ai < numberSelectionAttempts; ai++ {
		//nolint:gosec
		randTrackIndex := rand.Intn(len(s.musicData.Intervals[index].TrackIDs)-0) + 0
		randTrackID = s.musicData.Intervals[index].TrackIDs[randTrackIndex]

		if _, has := currentTrackIDs[randTrackID]; !has {
			break
		}
	}

	var track *domain.Track

	for _, t := range s.musicData.Tracks {
		if t.ID == randTrackID {
			track = t
			break
		}
	}

	if track == nil {
		//nolint:goerr113
		return nil, fmt.Errorf("musicData.Tracks not exists track with ID %d", randTrackID)
	}

	return &domain.PlaylistTrack{
		Track:                   track,
		Type:                    domain.PlaylistTrackTypeBackground,
		BackgroundIntervalIndex: index,
	}, nil
}
