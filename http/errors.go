package http

import "errors"

// ErrNotFound is a standard error to return with 404 status code.
var ErrNotFound error = errors.New("http: Hello, is it me you're looking for? Because the page you requested doesn't exist")

// ErrNotAllowed is a standard error to return with 405 status.
var ErrNotAllowed error = errors.New("http: we don't do that here")

// ErrInvalidRequest is a standard error to return with 400 status code.
var ErrInvalidRequest error = errors.New("http: invalid request or payload")

// ErrUnsupportedMedia is a standard error to return with 415 status code.
var ErrUnsupportedMedia error = errors.New("http: unsupported media type")
