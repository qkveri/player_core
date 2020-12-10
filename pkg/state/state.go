package state

type State struct {
	PlayerInfo playerInfo
	MusicData  musicData
}

func NewState() *State {
	return &State{
		PlayerInfo: newPlayerInfo(),
		MusicData:  newMusicData(),
	}
}
