package object

import (
	"fmt"
	"strings"
)

type MusicOrientationInfo struct {
	Genres       []string
	Danceability float64
	Valence      float64
	Popularity   int
}

func NewMusicOrientationInfo() *MusicOrientationInfo {
	return &MusicOrientationInfo{
		Genres:       []string{},
		Danceability: 0.0,
		Valence:      0.0,
		Popularity:   0,
	}
}

func (i *MusicOrientationInfo) String() string {
	var output string
	output += fmt.Sprintf("genres: %s\n", strings.Join(i.Genres, ","))
	output += fmt.Sprintf("danceability: %f\n", i.Danceability)
	output += fmt.Sprintf("valence: %f\n", i.Valence)
	output += fmt.Sprintf("popularity: %d\n", i.Popularity)
	return output
}

func (i *MusicOrientationInfo) IsSet() bool {
	return len(i.Genres) != 0 && i.Danceability != 0.0 && i.Valence != 0.0 && i.Popularity != 0
}
