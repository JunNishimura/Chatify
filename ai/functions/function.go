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
