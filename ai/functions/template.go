package functions

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/sashabaranov/go-openai"
)

// function name
const (
	RecommendFunctionName           = "recommend"
	SetGenresFunctionName           = "setGenres"
	SetDanceabilityFunctionName     = "setDanceability"
	SetValenceFunctionName          = "setValence"
	SetPopularityFunctionName       = "setPopularity"
	SetAcousticnessFunctionName     = "setAcousticness"
	SetEnergyFunctionName           = "setEnergy"
	SetInstrumentalnessFunctionName = "setInstrumentalness"
	SetLivenessFunctionaName        = "setLiveness"
	SetSpeechinessFunctionName      = "setSpeechiness"

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
							A quantitative expression of the music valence the user wants.
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
							A quantitative expression of the music popularity the user wants.
							A value ranges from 0 to 100.
							Tracks with high popularity is more popular.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetAcousticnessFunctionName,
			Description: heredoc.Doc(`
				Save the acousticness value the user wants to listen to. 
				Acousticness describes how much the track is acoustic`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music acousticness the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music acousticness the user wants.
							A value ranges from 0 to 100.
							Tracks with high acousticness is more acoustic.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetAcousticnessFunctionName,
			Description: heredoc.Doc(`
				Save the acousticness value the user wants to listen to. 
				Acousticness describes how much the track is acoustic`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music acousticness the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music acousticness the user wants.
							A value ranges from 0 to 1.0.
							Tracks with high acousticness is more acoustic.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetEnergyFunctionName,
			Description: heredoc.Doc(`
				Save the energy value the user wants to listen to. 
				Energy describes how much the track has energy`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music energy the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music energy the user wants.
							A value ranges from 0 to 1.0.
							Tracks with high energy is more energy.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetInstrumentalnessFunctionName,
			Description: heredoc.Doc(`
				Save the instrumentalness value the user wants to listen to. 
				Instrumentalness describes how much the track is instrumental`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music instrumentalness the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music instrumentalness the user wants.
							A value ranges from 0 to 1.0.
							Tracks with high instrumentalness is more instrumental.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetLivenessFunctionaName,
			Description: heredoc.Doc(`
				Save the liveness value the user wants to listen to. 
				Liveness describes how much the track is live`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music liveness the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music liveness the user wants.
							A value ranges from 0 to 1.0.
							Tracks with high liveness is more live.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
		{
			Name: SetSpeechinessFunctionName,
			Description: heredoc.Doc(`
				Save the speechiness value the user wants to listen to. 
				Speechiness describes how much the track is speech-like`,
			),
			Parameters: Parameters{
				Type: ObjectType,
				Properties: SetProperties{
					QualitativeValue: Property{
						Type:        StringType,
						Description: "A qualitative expression of the music speechiness the user wants.",
					},
					QuantitativeValue: Property{
						Type: NumberType,
						Description: heredoc.Doc(`
							A quantitative expression of the music speechiness the user wants.
							A value ranges from 0 to 1.0.
							Tracks with high speechiness is more speech-like.`,
						),
					},
				},
				Required: []string{"qualitative_value", "quatitative_value"},
			},
		},
	}
}
