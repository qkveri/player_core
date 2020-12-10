package state

import (
	"sync"

	"github.com/qkveri/player_core/pkg/domain"
)

type playerInfo struct {
	sync.RWMutex

	playerInfo *domain.PlayerInfo
}

func newPlayerInfo() playerInfo {
	return playerInfo{}
}

func (p *playerInfo) Set(playerInfo *domain.PlayerInfo) {
	p.Lock()
	defer p.Unlock()

	p.playerInfo = playerInfo
}

func (p *playerInfo) Get() *domain.PlayerInfo {
	p.RLock()
	defer p.RUnlock()

	return p.playerInfo
}
