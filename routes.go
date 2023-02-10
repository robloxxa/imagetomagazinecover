package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

type MagazineRequest struct {
	FileId string `json:"file_id"`
}

type MagazineResponse struct {
	FileUrl string `json:"url"`
	Name    string
}

func setupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Post("/magazine", func(w http.ResponseWriter, r *http.Request) {
		var (
			mreq MagazineRequest
			mres MagazineResponse
			err  error
		)
		defer func() {
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(400)
			}
		}()
		err = json.NewDecoder(r.Body).Decode(&mreq)
		if err != nil {
			return
		}

		link, err := getLinkToPhoto(mreq.FileId)
		if err != nil {
			return
		}
		url := fmt.Sprintf("http://localhost:3000/static/magazine.html?img_url=%s&name&qualities=", link)
		buf, err := processScreenshotMessage(url)
		if err != nil {
			return
		}

		file, err := os.Create("./static/images/123.img")
		if err != nil {
			return
		}
		defer file.Close()

		_, err = file.Write(buf)
		if err != nil {
			return
		}

		mres.FileUrl = "http://localhost:3000/static/images/123.img"

		if err = json.NewEncoder(w).Encode(mres); err != nil {
			return
		}
	})

	return r
}
