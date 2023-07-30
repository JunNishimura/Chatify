package object

import (
	"fmt"
	"strings"
)

type info struct {
	value      any
	hasChanged bool
}

type MusicOrientationInfo struct {
	Genres       info
	Danceability info
	Valence      info
	Popularity   info
}

func NewMusicOrientationInfo() *MusicOrientationInfo {
	return &MusicOrientationInfo{
		Genres: info{
			value:      []string{},
			hasChanged: false,
		},
		Danceability: info{
			value:      0.0,
			hasChanged: false,
		},
		Valence: info{
			value:      0.0,
			hasChanged: false,
		},
		Popularity: info{
			value:      0,
			hasChanged: false,
		},
	}
}

func (i *MusicOrientationInfo) String() string {
	var output string
	output += fmt.Sprintf("genres: %s\n", strings.Join(i.Genres.value.([]string), ","))
	output += fmt.Sprintf("danceability: %f\n", i.Danceability.value)
	output += fmt.Sprintf("valence: %f\n", i.Valence.value)
	output += fmt.Sprintf("popularity: %d\n", i.Popularity.value)
	return output
}

func (i *MusicOrientationInfo) IsSet() bool {
	return i.Genres.hasChanged && i.Danceability.hasChanged && i.Valence.hasChanged && i.Popularity.hasChanged
}

func (i *MusicOrientationInfo) SetGenres(genres []string) {
	i.Genres.value = genres
	i.Genres.hasChanged = true
}

func (i *MusicOrientationInfo) SetDanceability(danceability float64) {
	i.Danceability.value = danceability
	i.Danceability.hasChanged = true
}

func (i *MusicOrientationInfo) SetValence(valence float64) {
	i.Valence.value = valence
	i.Valence.hasChanged = true
}

func (i *MusicOrientationInfo) SetPopularity(popularity int) {
	i.Popularity.value = popularity
	i.Popularity.hasChanged = true
}
