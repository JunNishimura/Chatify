package functions

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/sashabaranov/go-openai"
)

// function name
const (
	RecommendFunctionName       = "recommend"
	SetGenresFunctionName       = "setGenres"
	SetDanceabilityFunctionName = "setDanceability"
	SetValenceFunctionName      = "setValence"
	SetPopularityFunctionName   = "setPopularity"

	ObjectType = "object"
	StringType = "string"
	NumberType = "number"
)

type Parameters struct {
	Type       string   `json:"type"`
	Properties any      `json:"properties"`
	Required   []string `json:"required"`
}

type SetProperties struct {
	QualitativeValue  Property `json:"qualitative_value,omitempty"`
	QuantitativeValue Property `json:"quantitative_value,omitempty"`
}

type Property struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

func GetFunctionDefinitions(genres []string) []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name:        SetGenresFunctionName,
			Description: "Save the genre of the music the user wants to listent to.",
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type: StringType,
						Description: heredoc.Docf(`
							Information on music genres.
							Multiple elements can be entered. e.g. j-pop,k-pop,chill.
							The genres that can be taken as arguments are as follows: %s"`,
							strings.Join(genres, ","),
						),
					},
				},
				Required: []string{"qualitative_value"},
			},
		},
		{
			Name: SetDanceabilityFunctionName,
			Description: heredoc.Doc(`
				Save the danceability value the user wants to listen to. 
				Danceability describes how suitable a track is for dancing based on a combination of musical elements including tempo, rhythm stability, beat strength, and overall regularity.`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music danceability the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music danceability the user wants.
							A value ranges from 0.0 to 1.0.
							A value of 0.0 is least danceable and 1.0 is most danceable.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetValenceFunctionName,
			Description: heredoc.Doc(`
				Save the valence value the user wants to listen to. 
				Valence describes how much the musical positiveness conveyed by a track.`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music valence the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music danceability the user wants.
							A value ranges from 0.0 to 1.0.
							Tracks with high valence sound more positive, while tracks with low valence sound more negative.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetPopularityFunctionName,
			Description: heredoc.Doc(`
				Save the popularity value the user wants to listen to. 
				Popularity describes how much the track is popular`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music popularity the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music danceability the user wants.
							A value ranges from 0 to 100.
							Tracks with high popularity is more popular.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
	}
}
