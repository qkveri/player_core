package state

type State struct {
	PlayerInfo playerInfo
	MusicData  musicData
	Playlist   playlist
}

func NewState() *State {
	return &State{
		PlayerInfo: newPlayerInfo(),
		MusicData:  newMusicData(),
		Playlist:   newPlaylist(),
	}
}
