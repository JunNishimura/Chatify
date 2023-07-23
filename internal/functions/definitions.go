package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/zmb3/spotify/v2"
)

const (
	RecommendCount = 5
)

func Recommend(ctx context.Context, client *spotify.Client, genres string) (string, error) {
	// genres length needs to be less than 5
	if len(genres) > 5 {
		genres = genres[:5]
	}
	genresSlice := strings.Split(genres, ",")

	seeds := spotify.Seeds{
		Genres: genresSlice,
	}

	trackAttrib := spotify.NewTrackAttributes()

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
