package models

type Tag struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}
