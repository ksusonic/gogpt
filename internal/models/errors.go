package models

import "errors"

type GenerationError error

var (
	GenerationTimeoutErr  GenerationError = errors.New("generation timeout")
	GenerationNotReadyErr GenerationError = errors.New("not ready")
)
