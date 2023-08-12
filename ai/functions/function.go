package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/JunNishimura/Chatify/ai/model"
	"github.com/JunNishimura/Chatify/util"
	"github.com/JunNishimura/spotify/v2"
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

const (
	RecommendCount = 5
)

func Recommend(ctx context.Context, client *spotify.Client, musicOrientation *model.MusicOrientation) (string, error) {
	// genres length needs to be less than 5
	genres := musicOrientation.Genres.Value
	if len(genres) > 5 {
		genres = genres[:5]
	}

	seeds := spotify.Seeds{
		Genres: genres,
	}

	trackAttrib := spotify.NewTrackAttributes().
		TargetDanceability(musicOrientation.Danceability.Value).
		TargetValence(musicOrientation.Valence.Value).
		TargetPopularity(musicOrientation.Popularity.Value)

	recommendations, err := client.GetRecommendations(ctx, seeds, trackAttrib, spotify.Limit(RecommendCount))
	if err != nil {
		return "", fmt.Errorf("fail to get recommendations: %v", err)
	}

	var output string
	for i, track := range recommendations.Tracks {
		output += fmt.Sprintf("[%d]\n", i+1)
		output += fmt.Sprintf("album: %s\n", track.Album.Name)
		var artists []string
		for _, artist := range track.Artists {
			artists = append(artists, artist.Name)
		}
		output += fmt.Sprintf("artists: %s\n", strings.Join(artists, ","))
		output += fmt.Sprintf("url: %s\n", track.ExternalURLs["spotify"])
		output += "\n"
	}

	return output, nil
}
