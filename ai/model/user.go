package model

import (
	"fmt"
	"strings"

	"github.com/JunNishimura/Chatify/utils"
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
	AcousticnessElement
	EnergyElement
	InstrumentalnessElement
	LivenessElement
	SpeechinessElement
)

type MusicOrientation struct {
	Genres           utils.Info[[]string]
	Danceability     utils.Info[float64]
	Valence          utils.Info[float64]
	Popularity       utils.Info[int]
	Acousticness     utils.Info[float64]
	Energy           utils.Info[float64]
	Instrumentalness utils.Info[float64]
	Liveness         utils.Info[float64]
	Speechiness      utils.Info[float64]
}

func newMusicOrientation() *MusicOrientation {
	return &MusicOrientation{
		Genres: utils.Info[[]string]{
			Value:      []string{},
			HasChanged: false,
		},
		Danceability: utils.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Valence: utils.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Popularity: utils.Info[int]{
			Value:      0,
			HasChanged: false,
		},
		Acousticness: utils.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Energy: utils.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Instrumentalness: utils.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Liveness: utils.Info[float64]{
			Value:      0.0,
			HasChanged: false,
		},
		Speechiness: utils.Info[float64]{
			Value:      0.0,
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
	output += fmt.Sprintf("acousticness: %f\n", m.Acousticness.Value)
	output += fmt.Sprintf("energy: %f\n", m.Energy.Value)
	output += fmt.Sprintf("instrumentalness: %f\n", m.Instrumentalness.Value)
	output += fmt.Sprintf("liveness: %f\n", m.Liveness.Value)
	output += fmt.Sprintf("speechiness: %f\n", m.Speechiness.Value)
	return output
}

func (m *MusicOrientation) HasAllSet() bool {
	return m.Genres.HasChanged &&
		m.Danceability.HasChanged &&
		m.Valence.HasChanged &&
		m.Popularity.HasChanged &&
		m.Acousticness.HasChanged &&
		m.Energy.HasChanged &&
		m.Instrumentalness.HasChanged &&
		m.Liveness.HasChanged &&
		m.Speechiness.HasChanged
}
