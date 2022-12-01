package models

import "gorm.io/gorm"

type Subscribe struct {
	gorm.Model
	UrlSha256 string `gorm:"column:url_sha256;primary_key;type:varchar(64);not null"`
	Name      string `gorm:"column:name;type:varchar(100);not null"`
	Url       string `gorm:"column:url;type:varchar(200);not null"`
	Using     bool   `gorm:"column:using;type:tinyint(1);not null"`
}
