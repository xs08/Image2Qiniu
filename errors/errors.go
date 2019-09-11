package errors

import "errors"

// ErrConfigFileNotExits config file not exists
var ErrConfigFileNotExits = errors.New("Config file doesn't exists")

// ErrOpenConfig error for open config file
var ErrOpenConfig = errors.New("Can't open config file'")

// ErrLoadConfig error for load localconfig to AppConfig struct
var ErrLoadConfig = errors.New("Can't load config to AppConfig")

// ErrNoImageSpecify error if no speci image
var ErrNoImageSpecify = errors.New("image must specify")

// ErrNoAccessKey no access key
var ErrNoAccessKey = errors.New("No access key specify")

// ErrNoSecretKey no secret key
var ErrNoSecretKey = errors.New("No secret key specify")

// ErrNoBucketName no bucketName specify
var ErrNoBucketName = errors.New("No bucket name specify")
