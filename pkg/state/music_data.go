package state

import (
	"sync"

	"github.com/qkveri/player_core/pkg/domain"
)

type musicData struct {
	sync.RWMutex

	musicData *domain.MusicData
}

func newMusicData() musicData {
	return musicData{}
}

func (m *musicData) Set(musicData *domain.MusicData) {
	m.musicData = musicData
}

func (m *musicData) Get() *domain.MusicData {
	return m.musicData
}
