package model

type GreedImage struct {
	Id  int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Url string `json:"url" gorm:"size:500;unique;comment:url"`
}

func GetGreedRandomImageUrlByNum(num int) (greedImage *GreedImage, err error) {
	err = DB.Limit(1).Offset(num).Omit("id").Find(&greedImage).Error
	return greedImage, err
}

func GetGreedRandomImageTotal() (total int64, err error) {
	err = DB.Model(&GreedImage{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, err
}
