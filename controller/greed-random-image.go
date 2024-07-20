package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"math/rand/v2"
	"net/http"
	"one-api/model"
	"strconv"
)

func GetGreedRandomImageUrlByNum(c *gin.Context) {
	num, _ := strconv.Atoi(c.Query("num"))
	if num < 0 {
		num = 0
	}
	if num == 0 {
		total, _ := model.GetGreedRandomImageTotal()
		num = rand.IntN(int(total))
	}
	greedImage, err := model.GetGreedRandomImageUrlByNum(num)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 下载图片
	resp, err := http.Get(greedImage.Url)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to get image",
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// 获取响应头中的 Content-Type
	contentType := resp.Header.Get("Content-Type")
	c.Header("Content-Type", contentType)
	// 将下载的文件流直接写入response流
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to send image",
		})
		return
	}
}
