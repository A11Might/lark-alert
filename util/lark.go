package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var client = lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"))

func SendMsg(msgType, content string) (string, error) {
	log.Default().Println(content)
	// 创建 Client
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId("oc_5af19f112b4e71be0fac728dc5e155fe").
			MsgType(msgType).
			Content(content).
			// Uuid(`选填，每次调用前请更换，如a0d69e20-1dd1-458b-k525-dfeca4015204`).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Im.V1.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		return "", fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
	return *resp.Data.MessageId, nil
}

func ReplayMsg(messageId, msgType, content string) {
	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			Content(content).
			MsgType(msgType).
			// ReplyInThread(true).
			// Uuid(`选填，每次调用前请更换，如a0d69e20-1dd1-458b-k525-dfeca4015204`).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Im.V1.Message.Reply(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return
	}
}

type Card struct {
	Type string            `json:"type,omitempty"`
	Data *TemplateCardData `json:"data,omitempty"`
}

type TemplateCardData struct {
	TemplateID          string                 `json:"template_id,omitempty"`
	TemplateVariable    map[string]interface{} `json:"template_variable,omitempty"`
	TemplateVersionName string                 `json:"template_version_name,omitempty"`
}

// https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/feishu-card-cardkit/configure-card-variables#a6abb723
func BuildTemplateCard(list []map[string]string) string {
	data := &TemplateCardData{
		TemplateID:          "AAqBH01pfadns",
		TemplateVersionName: "1.0.2",
		TemplateVariable: map[string]interface{}{
			"object_img": list,
		},
	}
	card := &Card{
		Type: "template",
		Data: data,
	}
	b, _ := json.Marshal(card)
	return string(b)
}

type AuditData struct {
	FileKey string `json:"file_key,omitempty"`
}

func BuildAuidtMsg(fileKey string) string {
	content := &AuditData{
		FileKey: fileKey,
	}
	b, _ := json.Marshal(content)
	return string(b)
}

func UploadFile(filename string, duration int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建请求对象
	req := larkim.NewCreateFileReqBuilder().
		Body(larkim.NewCreateFileReqBodyBuilder().
			FileType("opus").
			FileName(filename).
			Duration(duration).
			File(file).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Im.V1.File.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		return "", fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
	return *resp.Data.FileKey, nil
}
