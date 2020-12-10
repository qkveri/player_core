package state

import (
	"sync"

	"github.com/qkveri/player_core/pkg/domain"
)

type musicDataSubscriberFunc = func(musicData *domain.MusicData)

type musicData struct {
	sync.Mutex

	lastID      int
	subscribers map[int]musicDataSubscriberFunc
	musicData   *domain.MusicData
}

func newMusicData() musicData {
	return musicData{
		subscribers: make(map[int]musicDataSubscriberFunc),
	}
}

func (m *musicData) Set(musicData *domain.MusicData) {
	m.Lock()
	m.musicData = musicData
	m.Unlock()

	m.dispatch(musicData)
}

func (m *musicData) Subscribe(cb func(musicData *domain.MusicData)) func() {
	if cb == nil {
		return func() {}
	}

	var currentID int

	m.Lock()
	currentID = m.lastID
	m.subscribers[currentID] = cb
	m.lastID++
	m.Unlock()

	return m.unsubscribe(currentID)
}

func (m *musicData) unsubscribe(id int) func() {
	return func() {
		m.Lock()
		delete(m.subscribers, id)
		m.Unlock()
	}
}

func (m *musicData) dispatch(musicData *domain.MusicData) {
	m.Lock()
	for _, cb := range m.subscribers {
		cb(musicData)
	}
	m.Unlock()
}
