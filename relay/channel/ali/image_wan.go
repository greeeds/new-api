package ali

import (
	"fmt"
	"strings"

	"baipiao-api/common"
	"baipiao-api/dto"
	relaycommon "baipiao-api/relay/common"

	"github.com/gin-gonic/gin"
)

func oaiFormEdit2WanxImageEdit(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (*AliImageRequest, error) {
	var err error
	var imageRequest AliImageRequest
	imageRequest.Model = request.Model
	imageRequest.ResponseFormat = request.ResponseFormat
	wanInput := WanImageInput{
		Prompt: request.Prompt,
	}

	if err := common.UnmarshalBodyReusable(c, &wanInput); err != nil {
		return nil, err
	}
	if wanInput.Images, err = getImageBase64sFromForm(c, "image"); err != nil {
		return nil, fmt.Errorf("get image base64s from form failed: %w", err)
	}
	wanParams := WanImageParameters{
		N: int(request.N),
	}
	imageRequest.Input = wanInput
	imageRequest.Parameters = wanParams
	return &imageRequest, nil
}

func isWanModel(modelName string) bool {
	return strings.Contains(modelName, "wan")
}
