package openaicompat

import "baipiao-api/setting/model_setting"

func ShouldChatCompletionsUseResponsesPolicy(policy model_setting.ChatCompletionsToResponsesPolicy, channelID int, model string) bool {
	if !policy.IsChannelEnabled(channelID) {
		return false
	}
	return matchAnyRegex(policy.ModelPatterns, model)
}

func ShouldChatCompletionsUseResponsesGlobal(channelID int, model string) bool {
	return ShouldChatCompletionsUseResponsesPolicy(
		model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy,
		channelID,
		model,
	)
}
