package errors

import "errors"

// ErrOpenConfig error for open config file
var ErrOpenConfig = errors.New("Can't open config file'")

// ErrLoadConfig error for load localconfig to AppConfig struct
var ErrLoadConfig = errors.New("Can't load config to AppConfig")

// ErrNoImageSpecify error if no speci image
var ErrNoImageSpecify = errors.New("image must specify")
