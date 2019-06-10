package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

type jsonResponse struct {
	Status string   `json:"status"`
	Fucks  []string `json:"fucks"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	r.Route("/v1", func(r chi.Router) {
		r.Use(crs.Handler)
		r.Get("/give/{num}/fucks", giveFucks)
	})

	r.Handle("/*", StaticContentCacher(http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":3333", r)
}

func giveFucks(w http.ResponseWriter, r *http.Request) {
	numFucksString := chi.URLParam(r, "num")
	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Set("Cache-Control", "no-cache")

	numFucks, err := strconv.Atoi(numFucksString)
	if err != nil {
		w.WriteHeader(400)
		errResponse := errorResponse{"error", fmt.Sprintf("What the fuck kind of number is %s?", numFucksString)}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	if numFucks > 1000 {
		w.WriteHeader(410)
		errResponse := errorResponse{"error", "No one has that many fucking fucks to give."}
		json.NewEncoder(w).Encode(errResponse)
	} else if numFucks > 20 {
		w.WriteHeader(410)
		errResponse := errorResponse{"error", "Out of fucks to give."}
		json.NewEncoder(w).Encode(errResponse)
	} else if numFucks < 0 {
		w.WriteHeader(400)
		errResponse := errorResponse{"error", "Negative fucks? Are you fucking kidding me? Real cute... asshole."}
		json.NewEncoder(w).Encode(errResponse)
	} else {
		fucks := make([]string, numFucks)
		for i := 0; i < numFucks; i++ {
			fucks[i] = "fuck"
		}

		resp := jsonResponse{"ok", fucks}
		json.NewEncoder(w).Encode(resp)
	}
}

func StaticContentCacher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Do your magic here.
		rec := cacheControlWriter{w}

		next.ServeHTTP(&rec, req)
	})
}

type cacheControlWriter struct {
	http.ResponseWriter
}

func (w *cacheControlWriter) WriteHeader(code int) {
	content_type := w.Header().Get("Content-Type")

	if strings.HasPrefix(content_type, "text/html") {
		w.Header().Set("Cache-Control", "max-age=3600")
	} else {
		w.Header().Set("Cache-Control", "max-age=1209600")
	}

	w.ResponseWriter.WriteHeader(code)
}
