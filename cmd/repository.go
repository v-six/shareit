package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"time"

	"github.com/3di-clockwork/devops-test/cmd/types"
	"github.com/flytam/filenamify"
	"github.com/google/uuid"
	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

const (
	maxSize                = 128 << 20
	defaultRetentionPolicy = 24 * time.Hour
)

type ContentRepository struct {
	bucket *blob.Bucket
}

func (c *ContentRepository) CreateContentFromFile(ctx context.Context, file io.Reader, content types.ContentMeta) (*types.Content, error) {
	if content.Size > maxSize {
		return nil, errContentTooLarge
	} else if content.Size < 0 {
		return nil, errInvalidSize
	}

	var err error
	content.Filename, err = filenamify.FilenamifyV2(content.Filename)
	if err != nil {
		return nil, errInvalidFilename
	}

	id := types.ContentID("f-" + uuid.NewString())
	w, err := c.bucket.NewWriter(ctx, string(id), &blob.WriterOptions{
		ContentDisposition: fmt.Sprintf("inline; filename=\"%s\"", content.Filename),
		Metadata: map[string]string{
			"x-expiration-date": strconv.FormatInt(content.Expiry.UnixMilli(), 10),
			"x-filename":        content.Filename,
		},
	})
	if err != nil {
		return nil, err
	}

	defer w.Close()
	_, err = io.Copy(w, file)

	if err != nil {
		return nil, err
	}

	return &types.Content{
		ID:          id,
		ContentMeta: content,
	}, nil
}

func (c *ContentRepository) GetContentReaderFromContentID(ctx context.Context, id types.ContentID) (io.ReadCloser, error) {
	content, err := c.GetContentFromContentID(ctx, id)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, nil
	}
	return c.bucket.NewReader(ctx, string(id), nil)
}

func (c *ContentRepository) GetContentFromContentID(ctx context.Context, id types.ContentID) (*types.Content, error) {
	attribs, err := c.bucket.Attributes(ctx, string(id))
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			return nil, nil
		}
		return nil, err
	}

	expTimestamp, err := strconv.ParseInt(attribs.Metadata["x-expiration-date"], 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errContentCorrupted
	}
	exp := time.UnixMilli(expTimestamp)
	if time.Since(exp) > 0 {
		return nil, nil
	}

	return &types.Content{ID: id,
		ContentMeta: types.ContentMeta{
			Filename: attribs.Metadata["x-filename"],
			Expiry:   exp,
			Size:     attribs.Size,
		},
	}, nil
}
