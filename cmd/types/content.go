package types

import "time"

type ContentMeta struct {
	Filename string `validate:"filepath"`
	Size     int64  `validate:"min=1,max=2"`
	Expiry   time.Time
}

type Content struct {
	ID ContentID
	ContentMeta
}

type ContentID string
