package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/3di-clockwork/devops-test/cmd/types"
	"github.com/ggicci/httpin"
	"github.com/justinas/alice"
)

const (
	maxTTL = 24 * 30 // in hours = 30 days
)

type Controller struct {
	repo *ContentRepository
}

func (ctrl *Controller) Mount(mux *http.ServeMux) {
	mux.Handle("POST /files", alice.New(httpin.NewInput(filesCreateInput{})).ThenFunc(ctrl.filesCreate))
	mux.Handle("GET /files/{id}", http.HandlerFunc(ctrl.filesGet))
	mux.Handle("GET /files/{id}/raw", http.HandlerFunc(ctrl.filesGetRaw))
}

type filesCreateInput struct {
	File       *httpin.File `in:"form=file;required" validate:"filepath"`
	TTLInHours int          `in:"form=ttl;default=1" validate:"min=1,max=720"`
}

func (c *Controller) filesCreate(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*filesCreateInput)
	input.TTLInHours = min(input.TTLInHours, maxTTL)

	f, err := input.File.OpenReceiveStream()
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to open file stream", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	content, err := c.repo.CreateContentFromFile(r.Context(), f, types.ContentMeta{
		Filename: input.File.Filename(),
		Size:     input.File.Size(),
		Expiry:   time.Now().Add(time.Duration(input.TTLInHours) * time.Hour),
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create content from file", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "content created", slog.String("id", string(content.ID)))

	http.Redirect(w, r, "/files/"+string(content.ID), http.StatusSeeOther)
}

func (c *Controller) filesGet(w http.ResponseWriter, r *http.Request) {
	id := types.ContentID(r.PathValue("id"))
	content, err := c.repo.GetContentFromContentID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get content", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if content == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = FileNotFound().Render(w)
		return
	}
	_ = FileDetail(r, content).Render(w)
}

func (c *Controller) filesGetRaw(w http.ResponseWriter, r *http.Request) {
	id := types.ContentID(r.PathValue("id"))
	content, err := c.repo.GetContentFromContentID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get content", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if content == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = FileNotFound().Render(w)
		return
	}

	rdr, err := c.repo.GetContentReaderFromContentID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get content", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rdr.Close()

	w.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", content.Filename))
	w.Header().Add("Content-Length", strconv.FormatInt(content.Size, 10))
	io.Copy(w, rdr)
}
