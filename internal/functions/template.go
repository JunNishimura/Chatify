package functions

import (
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	genresDescTemplate = "Information on music genres. Multiple elements can be entered. e.x) j-pop, k-pop, chill."
)

type FunctionParameters struct {
	Type       string   `json:"type"`
	Properties any      `json:"properties"`
	Required   []string `json:"required"`
}

type GenresProperties struct {
	Genres Genres `json:"genres"`
}

type Genres struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func getGenresDescSentence(genres []string) string {
	joinedGenres := strings.Join(genres, ", ")
	return "The genres that can be taken as arguments are as follows: " + joinedGenres
}

func GetTemplate(genres []string) []openai.FunctionDefinition {
	genresDesc := genresDescTemplate + " " + getGenresDescSentence(genres)

	return []openai.FunctionDefinition{
		{
			Name:        "recommend",
			Description: "get information about recommended music tracks",
			Parameters: FunctionParameters{
				Type: "object",
				Properties: GenresProperties{
					Genres: Genres{
						Type:        "string",
						Description: genresDesc,
					},
				},
				Required: []string{"genres"},
			},
		},
	}
}
