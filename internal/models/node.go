package models

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	SubscribeUrlSha256 string `gorm:"column:subscribe_url_sha256;type:varchar(64);not null"`
	Data               string `gorm:"column:data;type:longtext;not null"`
}
