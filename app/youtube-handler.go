package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscribers int    `json:"subscribers"`
	ChannelName string `json:"channel_name"`
	TotalVideos int    `json:"total_videos"`
	View        int    `json:"view"`
}

func getChannelStats(k, channelID string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		yt := YoutubeStats{
			Subscribers: 100,
			ChannelName: "Golang",
			TotalVideos: 1000,
			View:        10000,
		}
		ctx := context.Background()
		yts, err := youtube.NewService(ctx, option.WithAPIKey(k))
		if err != nil {
			fmt.Println("error creating new youtube service: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		call := yts.Channels.List([]string{"snippet, contentDetails, statistics"})
		response, err := call.Id(channelID).Do()
		if err != nil {
			fmt.Println("error calling youtube channels list api: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(response.Items[0].Snippet.Title)

		yt = YoutubeStats{}
		if len(response.Items) > 0 {
			val := response.Items[0]
			yt = YoutubeStats{
				Subscribers: int(val.Statistics.SubscriberCount),
				ChannelName: val.Snippet.Title,
				TotalVideos: int(val.Statistics.VideoCount),
				View:        int(val.Statistics.ViewCount),
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(yt); err != nil {
			panic(err)
		}
	}
}
