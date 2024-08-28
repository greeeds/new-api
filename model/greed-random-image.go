package model

import "one-api/common"

type GreedImage struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Url  string `json:"url" gorm:"size:500;unique;comment:url"`
	Nsfw uint8  `json:"nsfw" gorm:"type:tinyint;default:0;index;comment:Is Not Safe For Work?"`
}

func GetGreedRandomImageUrlByNum(num int, nsfw uint8) (greedImage *GreedImage, err error) {
	err = DB.Where("`nsfw` = ?", nsfw).Limit(1).Offset(num).Order("id ASC").Omit("id").Find(&greedImage).Error
	return greedImage, err
}

func GetGreedRandomImageTotal(nsfw uint8) (total int64, err error) {
	err = DB.Model(&GreedImage{}).Where("`nsfw` = ?", nsfw).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, err
}

func AddGreedRandomImage(url string, nsfw uint8) (err error) {
	greedImage := &GreedImage{
		Url:  url,
		Nsfw: nsfw,
	}
	err = DB.Create(greedImage).Error
	if err != nil {
		common.SysError("failed to add greed-random-image: " + err.Error())
	}
	return err
}