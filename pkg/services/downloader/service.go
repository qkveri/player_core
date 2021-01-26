package downloader

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog"

	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/progress"
	"github.com/qkveri/player_core/pkg/state"
)

const (
	checkDuration = time.Second
)

type service struct {
	cm      sync.Mutex
	current *current

	state      *state.State
	logger     zerolog.Logger
	errCb      func(error)
	mp3RootDir string
}

func NewService(state *state.State, logger zerolog.Logger, errCb func(error), mp3RootDir string) *service {
	return &service{
		state:      state,
		logger:     logger.With().Str("service", "downloader").Logger(),
		errCb:      errCb,
		mp3RootDir: mp3RootDir,
	}
}

func (s *service) Run(ctx context.Context) error {
	s.logger.Debug().Msg("starts up")
	defer s.logger.Debug().Msg("stopped")

	t := time.NewTicker(checkDuration)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-t.C:
			s.checkAndDownload(ctx)
		}
	}
}

func (s *service) checkAndDownload(ctx context.Context) {
	s.logger.Debug().Msg("starting checkAndDownload...")

	s.state.Playlist.RLock()
	defer s.state.Playlist.RUnlock()

	playlist := s.state.Playlist.Get()

	if len(playlist) == 0 {
		s.logger.Debug().Msg("download skipped (playlist empty)")
		return
	}

	var playlistTrack *domain.PlaylistTrack

	for _, pt := range playlist {
		if pt.FilePath == "" {
			playlistTrack = pt
			break
		}
	}

	if playlistTrack == nil {
		s.logger.Debug().Msg("download skipped (all tracks loaded)")
		return
	}

	s.cm.Lock()
	defer s.cm.Unlock()

	if s.current != nil {
		if s.current.playlistTrack == playlistTrack {
			return
		}

		s.logger.Debug().Msg("current cancel")
		s.current.cancel()
	}

	s.logger.Debug().Interface("playlistTrack", playlistTrack).Msg("start download...")

	cur := &current{
		playlistTrack: playlistTrack,
		mp3RootDir:    s.mp3RootDir,
	}

	s.current = cur

	go s.download(ctx, cur)
}

func (s *service) download(ctx context.Context, cur *current) {
	progressCh := make(chan progress.Progress)
	defer close(progressCh)

	doneCh := make(chan bool)
	defer close(doneCh)

	g := run.Group{}

	// progress handler...
	g.Add(func() error {
		for {
			select {
			case <-doneCh:
				return nil

			case val := <-progressCh:
				s.logger.Debug().Int("trackId", cur.playlistTrack.Track.ID).
					Stringer("progress", val).
					Msg("download progress")

				s.state.Playlist.Lock()
				s.state.Playlist.SetDownloadProgress(cur.playlistTrack, val)
				s.state.Playlist.Unlock()
			}
		}
	}, func(err error) {
		select {
		case <-ctx.Done():
			break

		case doneCh <- true:
			break
		}
	})

	// download...
	g.Add(func() error {
		filePath, err := cur.download(ctx, progressCh)

		if err != nil {
			s.logger.Err(err).Interface("playlistTrack", cur.playlistTrack).
				Msg("mp3 download error")

			return fmt.Errorf("mp3 download error: %w", err)
		}

		s.logger.Debug().Interface("playlistTrack", cur.playlistTrack).
			Str("filePath", filePath).
			Msg("mp3 downloaded")

		s.state.Playlist.Lock()
		s.state.Playlist.SetDownloadProgress(cur.playlistTrack, progress.Passed)
		s.state.Playlist.SetFilePath(cur.playlistTrack, filePath)
		s.state.Playlist.Unlock()

		s.cm.Lock()
		if s.current == cur {
			s.current = nil
		}
		s.cm.Unlock()

		return nil
	}, func(err error) {})

	if err := g.Run(); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		s.errCb(err)
	}
}
