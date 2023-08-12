package model

import (
	"fmt"
	"strings"

	"github.com/JunNishimura/Chatify/util"
)

type User struct {
	name             string
	MusicOrientation *MusicOrientation
}

func NewUser(name string) *User {
	return &User{
		name:             name,
		MusicOrientation: newMusicOrientation(),
	}
}

type OrientationElement int

const (
	GenresElement OrientationElement = iota
	DanceabilityElement
	ValenceElement
	PopularityElement
)

type MusicOrientation struct {
	Genres       util.Info[[]string]
	Danceability util.Info[float64]
	Valence      util.Info[float64]
	Popularity   util.Info[int]
}

func newMusicOrientation() *MusicOrientation {
	return &MusicOrientation{
		Genres: util.Info[[]string]{
			Value:      []string{},
			HasChanged: false,
		},
		Danceability: util.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Valence: util.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Popularity: util.Info[int]{
			Value:      0,
			HasChanged: false,
		},
	}
}

func (m *MusicOrientation) String() string {
	var output string
	output += fmt.Sprintf("genres: %s\n", strings.Join(m.Genres.Value, ","))
	output += fmt.Sprintf("danceability: %f\n", m.Danceability.Value)
	output += fmt.Sprintf("valence: %f\n", m.Valence.Value)
	output += fmt.Sprintf("popularity: %d\n", m.Popularity.Value)
	return output
}

func (m *MusicOrientation) HasAllSet() bool {
	return m.Genres.HasChanged &&
		m.Danceability.HasChanged &&
		m.Valence.HasChanged &&
		m.Popularity.HasChanged
}
