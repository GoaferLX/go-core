package http

import "errors"

// ErrNotFound is a standard error to return with 404 status.
var ErrNotFound error = errors.New("Hello, is it me you're looking for? Because the page you requested doesn't exist.")
