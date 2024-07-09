package main

import "errors"

var (
	errContentTooLarge  = errors.New("content is too large")
	errInvalidFilename  = errors.New("invalid filename")
	errInvalidSize      = errors.New("invalid size")
	errContentCorrupted = errors.New("content is corrupted")
)
