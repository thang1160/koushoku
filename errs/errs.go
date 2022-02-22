package errs

import "errors"

var ErrUnknown = errors.New("Unknown error")

var (
	ErrArchiveNotFound  = errors.New("Archive does not exist")
	ErrArtistNotFound   = errors.New("Artist does not exist")
	ErrCircleNotFound   = errors.New("Circle does not exist")
	ErrMagazineNotFound = errors.New("Magazine does not exist")
	ErrTagNotFound      = errors.New("Tag does not exist")
	ErrParodyNotFound   = errors.New("Parody does not exist")
)

var (
	ErrArchivePathRequired  = errors.New("Archive path is required")
	ErrArtistNameRequired   = errors.New("Artist name is required")
	ErrArtistNameTooLong    = errors.New("Artist name must be at most 128 characters")
	ErrCircleNameRequired   = errors.New("Circle name is required")
	ErrCircleNameTooLong    = errors.New("CIrcle name must be at most 128 characters")
	ErrMagazineNameRequired = errors.New("Magazine name is required")
	ErrMagazineNameTooLong  = errors.New("Magazine name must be at most 128 characters")
	ErrParodyNameRequired   = errors.New("Parody name is required")
	ErrParodyNameTooLong    = errors.New("Parody name must be at most 128 characters")
	ErrTagNameRequired      = errors.New("Tag name is required")
	ErrTagNameTooLong       = errors.New("Tag name must be at most 128 characters")
)
