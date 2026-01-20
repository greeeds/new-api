package service

import (
	"baipiao-api/service/openaicompat"
	"baipiao-api/setting/model_setting"
)

func ShouldChatCompletionsUseResponsesPolicy(policy model_setting.ChatCompletionsToResponsesPolicy, channelID int, model string) bool {
	return openaicompat.ShouldChatCompletionsUseResponsesPolicy(policy, channelID, model)
}

func ShouldChatCompletionsUseResponsesGlobal(channelID int, model string) bool {
	return openaicompat.ShouldChatCompletionsUseResponsesGlobal(channelID, model)
}
