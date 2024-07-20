package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"math/rand/v2"
	"net/http"
	"one-api/model"
	"os"
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
	// 修正：将greedImage 的url对应图片下载下来返回给前端
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
	// 将图片保存到临时文件
	tmpFile, err := os.CreateTemp("", "downloaded-*.jpg")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to get temp file",
		})
		return
	}
	defer func(tmpFile *os.File) {
		err := tmpFile.Close()
		if err != nil {

		}
	}(tmpFile)

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to resp image",
		})
		return
	}

	// 将图片内容返回给前端
	_, err = tmpFile.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	c.Header("Content-Type", contentType)
	if _, err := io.Copy(c.Writer, tmpFile); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to send image",
		})
		return
	}
	return
}
