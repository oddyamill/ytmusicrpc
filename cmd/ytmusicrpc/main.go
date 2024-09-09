package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oddyamill/ytmusicrpc/discord"
)

const clientId = "1142840799738986556"
const authKey = "e8ab39d4b23d2877af508538de8424fd7c8ea4734870f462591b759acdf07199"

var activityType int

func init() {
	discord.SendHandshake(clientId)
}

func main() {
	flag.IntVar(&activityType, "type", 2, "discord activity type")
	flag.Parse()

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request %s %s", r.Method, r.URL.Path)

		if r.Header.Get("Authorization") != authKey {
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case "POST":
			updatePresence(w, r)
		case "DELETE":
			deletePresence(w)
		default:
			http.Error(w, "405 invalid request method", http.StatusMethodNotAllowed)
		}
	})

	log.Panic(http.ListenAndServe("127.0.0.1:32484", nil))
}

type updatePresenceBody struct {
	TrackId string `json:"trackId"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Artwork string `json:"artwork"`
	Album   string `json:"album"`
	Current *int64 `json:"current"`
	End     *int64 `json:"end"`
}

func updatePresence(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "415 invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var body updatePresenceBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, "400 bad json", http.StatusBadRequest)
		return
	}

	if body.Artist == "" || body.Artwork == "" || body.Title == "" || body.TrackId == "" {
		http.Error(w, "400 invalid request body", http.StatusBadRequest)
		return
	}

	activity := discord.Activity{
		Type:    activityType,
		Details: body.Title,
		State:   body.Artist,
		Assets: discord.Assets{
			LargeImage: body.Artwork,
		},
		Buttons: []discord.Button{
			{
				Label: "Слушать",
				Url:   "https://music.youtube.com/watch?v=" + body.TrackId,
			},
		},
	}

	if body.Album != "" && body.Title != body.Album && body.Artist != body.Album {
		activity.Assets.LargeText = body.Album
	}

	if body.Current != nil && body.End != nil {
		timestamp := time.Now().UnixMilli()

		activity.Timestamps = &discord.Timestamps{
			Start: timestamp - *body.Current,
			End:   timestamp - *body.Current + *body.End,
		}
	}

	discord.UpdatePresence(activity)
	fmt.Fprintf(w, "201 created")
}

func deletePresence(w http.ResponseWriter) {
	discord.DeletePresence()
	fmt.Fprintf(w, "200 ok")
}
