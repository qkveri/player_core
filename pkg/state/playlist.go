package state

import (
	"sync"

	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/progress"
)

type playlist struct {
	sync.RWMutex

	items []*domain.PlaylistTrack
}

func newPlaylist() playlist {
	return playlist{
		items: make([]*domain.PlaylistTrack, 0),
	}
}

func (p *playlist) Get() []*domain.PlaylistTrack {
	return p.items
}

func (p *playlist) Append(item *domain.PlaylistTrack) {
	p.items = append(p.items, item)
}

func (p *playlist) Replace(item *domain.PlaylistTrack, index int) {
	p.items[index] = item
}

func (p *playlist) SetDownloadProgress(item *domain.PlaylistTrack, downloadProgress progress.Progress) {
	for i, el := range p.items {
		if el == item {
			p.items[i].DownloadProgress = downloadProgress

			break
		}
	}
}

func (p *playlist) SetFilePath(item *domain.PlaylistTrack, filePath string) {
	for i, el := range p.items {
		if el == item {
			p.items[i].FilePath = filePath

			break
		}
	}
}

func (p *playlist) FirstItemDownloadProgress() progress.Progress {
	if len(p.items) == 0 {
		return 0
	}

	return p.items[0].DownloadProgress
}
