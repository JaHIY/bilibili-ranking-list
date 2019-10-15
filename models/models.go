package models

import (
    "github.com/jinzhu/gorm"
    "time"
)

type Media struct {
    gorm.Model
    Actors string `gorm:"type:text;not null"`
    Alias string `gorm:"type:text;not null"`
    Areas []*Area `gorm:"many2many:media_areas"`
    BilibiliMediaId int `gorm:"unique;not null"`
    Copyright string `gorm:"type:text;not null"`
    Cover string `gorm:"type:text;not null"`
    EpisodeIndex string `gorm:"type:text;not null"`
    Evaluate string `gorm:"type:text;not null"`
    OriginName string `gorm:"type:text;not null"`
    IsFinish int `gorm:"not null"`
    IsStarted int `gorm:"not null"`
    PubDate time.Time `xorm:"type:date;not null"`
    RatingCount int `gorm:"not null"`
    RatingScore int `gorm:"not null"`
    Staff string `gorm:"type:text;not null"`
    StatDanmakus int `gorm:"not null"`
    StatFavorites int `gorm:"not null"`
    StatViews int `gorm:"not null"`
    Styles []*Style `gorm:"many2many:media_styles"`
    Title string `gorm:"type:text;not null"`
}

type Area struct {
    gorm.Model
    Name string `gorm:"type:text;unique;not null"`
    Medias []*Media `gorm:"many2many:media_areas"`
}

type Style struct {
    gorm.Model
    Name string `gorm:"type:text;unique;not null"`
    Medias []*Media `gorm:"many2many:media_styles"`
}
