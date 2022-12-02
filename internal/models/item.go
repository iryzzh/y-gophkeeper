package models

import (
	"encoding/base64"
	"time"

	jsoniter "github.com/json-iterator/go"
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

type Items struct {
	Meta Meta    `json:"meta,omitempty"`
	Data []*Item `json:"data"`
}

type Meta struct {
	TotalItems int `json:"totalItems"`
}

func (i *Item) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigFastest
	return json.Marshal(&i)
}

func (id *ItemData) DecodeDataToString() (string, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(id.Data)))
	_, err := base64.StdEncoding.Decode(buf, id.Data)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (id *ItemData) DecodeDataToBytes() ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(id.Data)))
	_, err := base64.StdEncoding.Decode(buf, id.Data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
