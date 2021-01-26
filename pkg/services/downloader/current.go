package downloader

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/cavaliercoder/grab"

	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/progress"
)

const durationProgressPoll = 200 * time.Millisecond

type current struct {
	ctxCancel context.CancelFunc

	playlistTrack *domain.PlaylistTrack
	mp3RootDir    string
}

func (c *current) cancel() {
	c.ctxCancel()
}

func (c *current) download(ctx context.Context, progressCh chan<- progress.Progress) (string, error) {
	ctx, c.ctxCancel = context.WithCancel(ctx)
	defer c.ctxCancel()

	filePath := path.Join(c.mp3RootDir, strconv.Itoa(c.playlistTrack.Track.ID))
	client := grab.NewClient()

	req, err := grab.NewRequest(filePath, c.playlistTrack.Track.MP3URL)

	if err != nil {
		return "", fmt.Errorf("cannot grab.NewRequest: %w, filePath: %s, mp3URL: %s",
			err, filePath, c.playlistTrack.Track.MP3URL)
	}

	req = req.WithContext(ctx)

	resp := client.Do(req)

	t := time.NewTicker(durationProgressPoll)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			// dispatch to progress chan...
			select {
			case <-ctx.Done():
				return "", ctx.Err()

			case progressCh <- progress.Progress(resp.Progress()):
				break
			}
		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	return filePath, nil
}
