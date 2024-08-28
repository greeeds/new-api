package model

import (
	"one-api/common"
)

type GreedImage struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Url         string `json:"url" gorm:"size:500;unique;comment:url"`
	Nsfw        uint8  `json:"nsfw" gorm:"type:tinyint;default:0;index;comment:Is Not Safe For Work? 0-涩图 1-R18 2-妹子图 3-jk 4-bs 5-hs 6-luoli"`
	ContentType string `json:"content_type" gorm:"size:50;comment:content-type"`
	Keyword     string `json:"keyword" gorm:"size:50;index;comment:关键词"`
}

func GetGreedRandomImageUrlByNum(num int, nsfw uint8, keyword string) (greedImage *GreedImage, err error) {
	tx := DB
	if nsfw > 0 {
		tx = tx.Where("`nsfw` = ?", nsfw)
	}
	if keyword != "" {
		tx = tx.Where("`key_word` = ?", keyword)
	}
	err = tx.Order("id ASC").Limit(1).Offset(num).Find(&greedImage).Error
	return greedImage, err
}

func GetGreedRandomImageTotal(nsfw uint8) (total int64, err error) {
	tx := DB
	if nsfw > 0 {
		tx = tx.Where("`nsfw` = ?", nsfw)
	}
	err = tx.Model(&GreedImage{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, err
}

func AddGreedRandomImage(url string, nsfw uint8, contentType string, keyword string) (err error) {
	greedImage := &GreedImage{
		Url:         url,
		Nsfw:        nsfw,
		ContentType: contentType,
		Keyword:     keyword,
	}
	err = DB.Create(greedImage).Error
	if err != nil {
		common.SysError("failed to add greed-random-image: " + err.Error())
	}
	return err
}
