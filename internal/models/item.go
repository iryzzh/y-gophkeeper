package models

import (
	"time"
)

// Item is the model of the item. The fields `Item.ID`,
// 'Item.DataID', 'Item.CreatedAt', 'Item.UpdatedAt' are filled
// by the database service after creating or updating. Field
// 'Item.DataID' must correspond to field 'Item.ID' of
// `models.ItemData` struct.
type Item struct {
	ID        int        `json:"id,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	Meta      string     `json:"meta"`
	DataID    int        `json:"data_id,omitempty"`
	DataType  string     `json:"data_type,omitempty"`
	ItemData  *ItemData  `json:"item_data,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// ItemData is a data model. The field `ItemData.ID` are
// filled in by the database service
// after creating or updating.
type ItemData struct {
	ID   int    `json:"data_id,omitempty"`
	Data []byte `json:"data"`
}

type ItemResponse struct {
	Meta meta    `json:"meta,omitempty"`
	Data []*Item `json:"data"`
}

type meta struct {
	TotalPages int `json:"totalPages"`
}
