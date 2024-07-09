package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"
)

func init() {
	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func main() {
	// Load the app config from env
	err := Config.ParseFromEnv()
	if err != nil {
		slog.Error("cannot load config", slog.Any("err", err))
		os.Exit(CONFIG_ERROR)
	}

	bkt, err := blob.OpenBucket(context.Background(), Config.BlobStorageURL)
	if err != nil {
		slog.Error("cannot open blob storage", slog.Any("err", err))
		os.Exit(BLOB_STORAGE_ERROR)
	}

	// Health check endpoint, for k8s for instance
	http.Handle("GET /healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessible, err := bkt.IsAccessible(r.Context())

		w.Header().Add("Content-Type", "application/json")
		if err != nil || !accessible {
			if err != nil {
				slog.ErrorContext(r.Context(), "storage is unavailable", slog.Any("err", err))
			} else {
				slog.ErrorContext(r.Context(), "storage is unavailable")
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"storage":"unavailable"}`))
			return
		}

		_, _ = w.Write([]byte(`{"storage":"available"}`))
	}))

	// Storage
	repo := ContentRepository{bucket: bkt}

	// HTTP controller (sometimes referred to as the view)
	ctrl := Controller{repo: &repo}

	// Mount the endpoints
	ctrl.Mount(http.DefaultServeMux)

	// Home
	http.Handle("/{$}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _ = Home().Render(w) }))

	// Fallback to home
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/", http.StatusSeeOther) }))

	slog.Info("starting shareit app on port 8080")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
}
