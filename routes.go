package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"net/http/pprof"
	"regexp"
)

type MagazineRequest struct {
	FileId    string `json:"file_id"`
	Token     string `json:"token"`
	Qualities string `json:"qualities"`
}

type MagazineResponse struct {
	FileUrl string `json:"url"`
	Ok      bool   `json:"ok"`
	Error   string `json:"description"`
}

func setupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("123"))
	})

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Post("/magazine", magazineHandler)

	return r
}

var splitRegex = regexp.MustCompile("\n+")

func magazineHandler(w http.ResponseWriter, r *http.Request) {
	var (
		mreq MagazineRequest
		mres MagazineResponse
		err  error
	)
	defer func() {
		if err != nil {
			mres.Error = err.Error()
			mres.Ok = false
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			if err = json.NewEncoder(w).Encode(&mres); err != nil {
				return
			}
		}
	}()
	if err = json.NewDecoder(r.Body).Decode(&mreq); err != nil {

		fmt.Println(mreq)
		return
	}

	qualitites := splitRegex.Split(mreq.Qualities, 2)
	if len(qualitites) < 2 {
		err = errors.New("too few params")
		return
	}

	job := ScreenshotJob{
		mreq,
		qualitites,
		make(chan string),
	}

	ScreenshotJobQueue <- job

	select {
	case urlstr := <-job.C:
		if urlstr == "" {
			err = errors.New("urlstring is empty")
			return
		}
		mres.Ok = true
		mres.FileUrl = urlstr
		if err = json.NewEncoder(w).Encode(&mres); err != nil {
			return
		}
	case <-r.Context().Done():
		return
	}

}
