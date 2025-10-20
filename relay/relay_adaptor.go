package relay

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"baipiao-api/constant"
	"baipiao-api/relay/channel"
	"baipiao-api/relay/channel/ali"
	"baipiao-api/relay/channel/aws"
	"baipiao-api/relay/channel/baidu"
	"baipiao-api/relay/channel/baidu_v2"
	"baipiao-api/relay/channel/claude"
	"baipiao-api/relay/channel/cloudflare"
	"baipiao-api/relay/channel/cohere"
	"baipiao-api/relay/channel/coze"
	"baipiao-api/relay/channel/deepseek"
	"baipiao-api/relay/channel/dify"
	"baipiao-api/relay/channel/gemini"
	"baipiao-api/relay/channel/jimeng"
	"baipiao-api/relay/channel/jina"
	"baipiao-api/relay/channel/mistral"
	"baipiao-api/relay/channel/mokaai"
	"baipiao-api/relay/channel/moonshot"
	"baipiao-api/relay/channel/ollama"
	"baipiao-api/relay/channel/openai"
	"baipiao-api/relay/channel/palm"
	"baipiao-api/relay/channel/perplexity"
	"baipiao-api/relay/channel/siliconflow"
	"baipiao-api/relay/channel/submodel"
	taskdoubao "baipiao-api/relay/channel/task/doubao"
	taskjimeng "baipiao-api/relay/channel/task/jimeng"
	"baipiao-api/relay/channel/task/kling"
	tasksora "baipiao-api/relay/channel/task/sora"
	"baipiao-api/relay/channel/task/suno"
	taskvertex "baipiao-api/relay/channel/task/vertex"
	taskVidu "baipiao-api/relay/channel/task/vidu"
	"baipiao-api/relay/channel/tencent"
	"baipiao-api/relay/channel/vertex"
	"baipiao-api/relay/channel/volcengine"
	"baipiao-api/relay/channel/xai"
	"baipiao-api/relay/channel/xunfei"
	"baipiao-api/relay/channel/zhipu"
	"baipiao-api/relay/channel/zhipu_4v"
)

func GetAdaptor(apiType int) channel.Adaptor {
	switch apiType {
	case constant.APITypeAli:
		return &ali.Adaptor{}
	case constant.APITypeAnthropic:
		return &claude.Adaptor{}
	case constant.APITypeBaidu:
		return &baidu.Adaptor{}
	case constant.APITypeGemini:
		return &gemini.Adaptor{}
	case constant.APITypeOpenAI:
		return &openai.Adaptor{}
	case constant.APITypePaLM:
		return &palm.Adaptor{}
	case constant.APITypeTencent:
		return &tencent.Adaptor{}
	case constant.APITypeXunfei:
		return &xunfei.Adaptor{}
	case constant.APITypeZhipu:
		return &zhipu.Adaptor{}
	case constant.APITypeZhipuV4:
		return &zhipu_4v.Adaptor{}
	case constant.APITypeOllama:
		return &ollama.Adaptor{}
	case constant.APITypePerplexity:
		return &perplexity.Adaptor{}
	case constant.APITypeAws:
		return &aws.Adaptor{}
	case constant.APITypeCohere:
		return &cohere.Adaptor{}
	case constant.APITypeDify:
		return &dify.Adaptor{}
	case constant.APITypeJina:
		return &jina.Adaptor{}
	case constant.APITypeCloudflare:
		return &cloudflare.Adaptor{}
	case constant.APITypeSiliconFlow:
		return &siliconflow.Adaptor{}
	case constant.APITypeVertexAi:
		return &vertex.Adaptor{}
	case constant.APITypeMistral:
		return &mistral.Adaptor{}
	case constant.APITypeDeepSeek:
		return &deepseek.Adaptor{}
	case constant.APITypeMokaAI:
		return &mokaai.Adaptor{}
	case constant.APITypeVolcEngine:
		return &volcengine.Adaptor{}
	case constant.APITypeBaiduV2:
		return &baidu_v2.Adaptor{}
	case constant.APITypeOpenRouter:
		return &openai.Adaptor{}
	case constant.APITypeXinference:
		return &openai.Adaptor{}
	case constant.APITypeXai:
		return &xai.Adaptor{}
	case constant.APITypeCoze:
		return &coze.Adaptor{}
	case constant.APITypeJimeng:
		return &jimeng.Adaptor{}
	case constant.APITypeMoonshot:
		return &moonshot.Adaptor{} // Moonshot uses Claude API
	case constant.APITypeSubmodel:
		return &submodel.Adaptor{}
	}
	return nil
}

func GetTaskPlatform(c *gin.Context) constant.TaskPlatform {
	channelType := c.GetInt("channel_type")
	if channelType > 0 {
		return constant.TaskPlatform(strconv.Itoa(channelType))
	}
	return constant.TaskPlatform(c.GetString("platform"))
}

func GetTaskAdaptor(platform constant.TaskPlatform) channel.TaskAdaptor {
	switch platform {
	//case constant.APITypeAIProxyLibrary:
	//	return &aiproxy.Adaptor{}
	case constant.TaskPlatformSuno:
		return &suno.TaskAdaptor{}
	}
	if channelType, err := strconv.ParseInt(string(platform), 10, 64); err == nil {
		switch channelType {
		case constant.ChannelTypeKling:
			return &kling.TaskAdaptor{}
		case constant.ChannelTypeJimeng:
			return &taskjimeng.TaskAdaptor{}
		case constant.ChannelTypeVertexAi:
			return &taskvertex.TaskAdaptor{}
		case constant.ChannelTypeVidu:
			return &taskVidu.TaskAdaptor{}
		case constant.ChannelTypeDoubaoVideo:
			return &taskdoubao.TaskAdaptor{}
		case constant.ChannelTypeSora, constant.ChannelTypeOpenAI:
			return &tasksora.TaskAdaptor{}
		}
	}
	return nil
}
