package relay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"baipiao-api/common"
	"baipiao-api/dto"
	"baipiao-api/logger"
	relaycommon "baipiao-api/relay/common"
	"baipiao-api/relay/helper"
	"baipiao-api/service"
	"baipiao-api/types"

	"github.com/gin-gonic/gin"
)

func EmbeddingHelper(c *gin.Context, info *relaycommon.RelayInfo) (newAPIError *types.NewAPIError) {
	info.InitChannelMeta(c)

	embeddingReq, ok := info.Request.(*dto.EmbeddingRequest)
	if !ok {
		newApiErr := types.NewErrorWithStatusCode(fmt.Errorf("invalid request type, expected *dto.EmbeddingRequest, got %T", info.Request), types.ErrorCodeInvalidRequest, http.StatusBadRequest, types.ErrOptionWithSkipRetry())
		postConsumeQuota(c, info, nil, newApiErr.Err.Error(), "")
		return newAPIError
	}

	request, err := common.DeepCopy(embeddingReq)
	if err != nil {
		postConsumeQuota(c, info, nil, err.Error(), "")
		return types.NewError(fmt.Errorf("failed to copy request to EmbeddingRequest: %w", err), types.ErrorCodeInvalidRequest, types.ErrOptionWithSkipRetry())
	}

	err = helper.ModelMappedHelper(c, info, request)
	if err != nil {
		postConsumeQuota(c, info, nil, err.Error(), "")
		return types.NewError(err, types.ErrorCodeChannelModelMappedError, types.ErrOptionWithSkipRetry())
	}

	adaptor := GetAdaptor(info.ApiType)
	if adaptor == nil {
		newApiErr := types.NewError(fmt.Errorf("invalid api type: %d", info.ApiType), types.ErrorCodeInvalidApiType, types.ErrOptionWithSkipRetry())
		postConsumeQuota(c, info, nil, newApiErr.Err.Error(), "")
		return newAPIError
	}
	adaptor.Init(info)

	convertedRequest, err := adaptor.ConvertEmbeddingRequest(c, info, *request)
	if err != nil {
		postConsumeQuota(c, info, nil, err.Error(), "")
		return types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}
	jsonData, err := json.Marshal(convertedRequest)
	if err != nil {
		postConsumeQuota(c, info, nil, err.Error(), "")
		return types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
	}
	logger.LogDebug(c, fmt.Sprintf("converted embedding request body: %s", string(jsonData)))
	requestBody := bytes.NewBuffer(jsonData)
	statusCodeMappingStr := c.GetString("status_code_mapping")
	resp, err := adaptor.DoRequest(c, info, requestBody)
	if err != nil {
		postConsumeQuota(c, info, nil, err.Error(), "")
		return types.NewOpenAIError(err, types.ErrorCodeDoRequestFailed, http.StatusInternalServerError)
	}

	var bodyContent string
	bodyContent = string(jsonData)
	var httpResp *http.Response
	if resp != nil {
		httpResp = resp.(*http.Response)
		if httpResp.StatusCode != http.StatusOK {
			newAPIError = service.RelayErrorHandler(c.Request.Context(), httpResp, false)
			// reset status code 重置状态码
			service.ResetStatusCode(newAPIError, statusCodeMappingStr)
			postConsumeQuota(c, info, nil, newAPIError.Err.Error(), bodyContent)
			return newAPIError
		}
	}

	usage, newAPIError := adaptor.DoResponse(c, httpResp, info)
	if newAPIError != nil {
		// reset status code 重置状态码
		service.ResetStatusCode(newAPIError, statusCodeMappingStr)
		postConsumeQuota(c, info, nil, newAPIError.Err.Error(), bodyContent)
		return newAPIError
	}
	postConsumeQuota(c, info, usage.(*dto.Usage), "", bodyContent)
	return nil
}
