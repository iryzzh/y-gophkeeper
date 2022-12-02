package store

import "github.com/pkg/errors"

var (
	// ErrUserNotFound returns when the user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists returns when the user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserCreateFailed returns when the user creation failed.
	ErrUserCreateFailed = errors.New("user creation failed")
	// ErrItemInvalidID returns when the item id is invalid
	ErrItemInvalidID = errors.New("item id is invalid")
	// ErrItemDataInvalidID returns when the item data id is invalid
	ErrItemDataInvalidID = errors.New("item data id is invalid")
	// ErrItemExists returns when the item already exists.
	ErrItemExists = errors.New("item already exists")
	// ErrItemNotFound is returned when item not found.
	ErrItemNotFound = errors.New("item not found")
	// ErrItemDataNotFound is returned when item data not found.
	ErrItemDataNotFound = errors.New("item data not found")
	// ErrItemCreateFailed returns when the item creation failed.
	ErrItemCreateFailed = errors.New("item creation failed")
	// ErrItemDataCreateFailed returns when the item data creation failed.
	ErrItemDataCreateFailed = errors.New("item data creation failed")
	// ErrItemUpdateFailed returns when the item update failed.
	ErrItemUpdateFailed = errors.New("item update failed")
	// ErrItemDataUpdateFailed returns when the item data update failed.
	ErrItemDataUpdateFailed = errors.New("item data update failed")
	// ErrItemDeleteFailed returns when the item delete failed.
	ErrItemDeleteFailed = errors.New("item delete failed")
	// ErrItemDataDeleteFailed returns when the item data delete failed.
	ErrItemDataDeleteFailed = errors.New("item data delete failed")
	// ErrItemMetaIsRequired returns when the item meta is nil.
	ErrItemMetaIsRequired = errors.New("meta is required")
)
