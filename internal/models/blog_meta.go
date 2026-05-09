package models

import "time"

type BlogMeta struct {
	Description string
	Url         string
	Title       string
	ImageUri    string
	Published   time.Time
}
