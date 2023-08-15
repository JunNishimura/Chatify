package functions

import (
	"github.com/JunNishimura/Chatify/util"
)

func SetGenres(genresInfo *util.Info[[]string], value []string) {
	util.SetInfo(genresInfo, value)
}

func SetDanceability(danceabilityInfo *util.Info[float64], value float64) {
	util.SetInfo(danceabilityInfo, value)
}

func SetValence(valence *util.Info[float64], value float64) {
	util.SetInfo(valence, value)
}

func SetPopularity(popularity *util.Info[int], value int) {
	util.SetInfo(popularity, value)
}
