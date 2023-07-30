package functions

import (
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	genresDescTemplate = "Information on music genres. Multiple elements can be entered. e.g. j-pop,k-pop,chill."

	TypeObject = "object"
	TypeString = "string"
	TypeNumber = "number"
)

var RecommendFunctionDefinition = openai.FunctionDefinition{
	Name:        RecommendFunction,
	Description: "get information about recommended music tracks",
	Parameters: FunctionParameters{
		Type: TypeObject,
		Properties: RecommendProperties{
			Genres: Property{
				Type:        TypeString,
				Description: "Information on music genres. Multiple elements can be entered. e.g. j-pop, k-pop, chill.",
			},
			Danceability: Property{
				Type:        TypeNumber,
				Description: "Information on how suitable a track is for dancing based on a combination of musical elements including tempo, rhythm stability, beat strength, and overall regularity. Range from 0.0 to 1.0.",
			},
			Valence: Property{
				Type:        TypeNumber,
				Description: "Information on the musical positiveness conveyed by a track. Range from 0.0 to 1.0.",
			},
			Popularity: Property{
				Type:        TypeNumber,
				Description: "Integer. Information on how much the track is popular. Range from 0 to 100.",
			},
		},
		Required: []string{"genres", "danceability", "valence", "popularity"},
	},
}

type FunctionName = string

const (
	RecommendFunction       FunctionName = "recommend"
	SetGenresFunction       FunctionName = "setGenres"
	SetDanceabitliyFunction FunctionName = "setDanceability"
	SetValenceFunction      FunctionName = "setValence"
	SetPopularityFunction   FunctionName = "setPopularity"
)

type FunctionParameters struct {
	Type       string   `json:"type"`
	Properties any      `json:"properties"`
	Required   []string `json:"required"`
}

type RecommendProperties struct {
	Genres       Property `json:"genres"`
	Danceability Property `json:"danceability"`
	Valence      Property `json:"valence"`
	Popularity   Property `json:"popularity"`
}

type MusicElementProperties struct {
	Value    Property `json:"value"`
	SubValue Property `json:"subvalue,omitempty"`
}

type Property struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

func getGenresDescSentence(genres []string) string {
	joinedGenres := strings.Join(genres, ",")
	return "The genres that can be taken as arguments are as follows: " + joinedGenres
}

func GetTemplate(genres []string) []openai.FunctionDefinition {
	genresDesc := genresDescTemplate + " " + getGenresDescSentence(genres)

	return []openai.FunctionDefinition{
		{
			Name:        SetGenresFunction,
			Description: "get music genres and return object storing data about the music orientation.",
			Parameters: FunctionParameters{
				Type: TypeObject,
				Properties: MusicElementProperties{
					Value: Property{
						Type:        TypeString,
						Description: genresDesc,
					},
				},
				Required: []string{"value"},
			},
		},
		{
			Name:        SetDanceabitliyFunction,
			Description: "get music danceability and return object storing data about the music orientation. Danceability describes how suitable a track is for dancing based on a combination of musical elements including tempo, rhythm stability, beat strength, and overall regularity.",
			Parameters: FunctionParameters{
				Type: TypeObject,
				Properties: MusicElementProperties{
					Value: Property{
						Type:        TypeNumber,
						Description: "The value ranges from 0.0 to 1.0. A value of 0.0 is least danceable and 1.0 is most danceable.",
					},
					SubValue: Property{
						Type:        TypeString,
						Description: "How much danceability is wanted",
					},
				},
				Required: []string{"value", "subvalue"},
			},
		},
		{
			Name:        SetValenceFunction,
			Description: "get music valence and return object storing data about the music orientation. A measure from 0.0 to 1.0 describing the musical positiveness conveyed by a track.",
			Parameters: FunctionParameters{
				Type: TypeObject,
				Properties: MusicElementProperties{
					Value: Property{
						Type:        TypeNumber,
						Description: "The value ranges from 0.0 to 1.0. Tracks with high valence sound more positive, while tracks with low valence sound more negative.",
					},
					SubValue: Property{
						Type:        TypeString,
						Description: "How much valence is wanted",
					},
				},
				Required: []string{"value", "subvalue"},
			},
		},
		{
			Name:        SetPopularityFunction,
			Description: "get music popularity and return object storing data about the music orientation. A measure from 0 to 100 describing how much the track is popular.",
			Parameters: FunctionParameters{
				Type: TypeObject,
				Properties: MusicElementProperties{
					Value: Property{
						Type:        TypeNumber,
						Description: "The value ranges from 0 to 100. Tracks with high popularity is more popular.",
					},
					SubValue: Property{
						Type:        TypeString,
						Description: "How much popularity is wanted",
					},
				},
				Required: []string{"value", "subvalue"},
			},
		},
	}
}
