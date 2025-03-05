package util

import (
	"fmt"
	"os"

	"github.com/A11Might/lark-alert/model"
	"resty.dev/v3"
)

// deepseek-chat
// 乌克兰总统泽连斯基在与特朗普的激烈交锋后表示，尽管此次冲突“不好”，但两国关系仍可修复。泽连斯基强调乌克兰需要美国的支持来对抗俄罗斯，并感谢美国人民的帮助。民主党人批评特朗普的行
// 为“可耻”，认为其助长了普京的气焰。泽连斯基在采访中表示，乌克兰人民渴望和平，但需要安全保障。
//
// deepseek-reasoner
// 泽连斯基与特朗普在椭圆形办公室会晤中爆发激烈争执，民主党指责特朗普和万斯为普京服务，损害美国领导力。泽连斯基会后受访称冲突“不愉快”，但相信两国关系可修复，强调乌克兰需美支持且无
// 法单方面停火。特朗普称泽连斯基“想继续战斗”，主张立即停火但暗示可能缩减援助。加拿大总理特鲁多重申支持乌克兰，俄军则袭击哈尔科夫致医院起火。事件凸显美乌关系紧张及国内对援乌分歧。

func CallOpenAIAPI(prompt, content string) (string, error) {
	apiKey := os.Getenv("API_KEY")
	llmModel := os.Getenv("MODEL")
	endpoint := os.Getenv("ENDPOINT")

	c := resty.New()
	defer c.Close()

	res, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
		SetBody(map[string]interface{}{
			"model": llmModel,
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": prompt,
				},
				{
					"role":    "user",
					"content": content,
				},
			},
			"stream": false,
		}).
		SetResult(&model.OpenAIResponse{}).
		Post(endpoint)
	if err != nil {
		return "", err
	}
	return res.Result().(*model.OpenAIResponse).Choices[0].Message.Content, nil
}
