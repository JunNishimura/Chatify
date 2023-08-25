package functions

import (
	"github.com/JunNishimura/Chatify/utils"
)

func SetGenres(genresInfo *utils.Info[[]string], value []string) {
	utils.SetInfo(genresInfo, value)
}

func SetDanceability(danceabilityInfo *utils.Info[float64], value float64) {
	utils.SetInfo(danceabilityInfo, value)
}

func SetValence(valence *utils.Info[float64], value float64) {
	utils.SetInfo(valence, value)
}

func SetPopularity(popularity *utils.Info[int], value int) {
	utils.SetInfo(popularity, value)
}

func SetAcousticness(acousticness *utils.Info[float64], value float64) {
	utils.SetInfo(acousticness, value)
}

func SetEnergy(energy *utils.Info[float64], value float64) {
	utils.SetInfo(energy, value)
}

func SetInstrumentalness(instrumentalness *utils.Info[float64], value float64) {
	utils.SetInfo(instrumentalness, value)
}

func SetLiveness(liveness *utils.Info[float64], value float64) {
	utils.SetInfo(liveness, value)
}

func SetSpeechiness(speechiness *utils.Info[float64], value float64) {
	utils.SetInfo(speechiness, value)
}
