package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand/v2"
	"net/http"
	"one-api/common"
	"one-api/model"
	"os"
	"strconv"
)

func GetGreedRandomImageUrlByNum(c *gin.Context) {
	num, _ := strconv.Atoi(c.Query("num"))
	if num < 0 {
		num = 0
	}
	nsfw64, err := strconv.ParseUint(c.Query("nsfw"), 10, 8)
	nsfw := uint8(nsfw64)
	var total = int64(0)
	if num == 0 {
		total, _ = model.GetGreedRandomImageTotal(nsfw)
		num = rand.IntN(int(total))
	}
	showTotal, _ := strconv.ParseBool(c.Query("total"))
	if showTotal {
		if total == 0 {
			total, _ = model.GetGreedRandomImageTotal(nsfw)
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data":    total,
		})
		return
	}
	greedImage, err := model.GetGreedRandomImageUrlByNum(num, nsfw, c.Query("keyword"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 下载图片
	resp, err := http.Get(greedImage.Url)
	if resp == nil || err != nil {
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
	if contentType == "" && greedImage.ContentType != "" {
		contentType = greedImage.ContentType
	}
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

func GetGreedRandomImageUrlByNumRedirect(c *gin.Context) {
	num, _ := strconv.Atoi(c.Query("num"))
	if num < 0 {
		num = 0
	}
	nsfw64, err := strconv.ParseUint(c.Query("nsfw"), 10, 8)
	nsfw := uint8(nsfw64)
	var total = int64(0)
	if num == 0 {
		total, _ = model.GetGreedRandomImageTotal(nsfw)
		num = rand.IntN(int(total))
	}
	showTotal, _ := strconv.ParseBool(c.Query("total"))
	if showTotal {
		if total == 0 {
			total, _ = model.GetGreedRandomImageTotal(nsfw)
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data":    total,
		})
		return
	}
	greedImage, err := model.GetGreedRandomImageUrlByNum(num, nsfw, c.Query("keyword"))
	if greedImage == nil || err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to get image",
		})
		return
	}
	// 使用 http.Redirect 重定向到新的 URL
	http.Redirect(c.Writer, c.Request, greedImage.Url, http.StatusMovedPermanently)
}

func AddGreedRandomImage(c *gin.Context) {
	hehedaKey := os.Getenv("hehedaKey")
	if hehedaKey != "" {
		accessToken := c.Request.Header.Get("Heheda")
		if accessToken != hehedaKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
	}
	var greedImage model.GreedImage
	err := json.NewDecoder(c.Request.Body).Decode(&greedImage)
	if err != nil || greedImage.Url == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "invalid parameter",
		})
		return
	}
	if err := common.Validate.Struct(&greedImage); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Input not legal: " + err.Error(),
		})
		return
	}
	if err := model.AddGreedRandomImage(greedImage.Url, greedImage.Nsfw, greedImage.ContentType, greedImage.Keyword); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "failed to add greed-random-image: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}
