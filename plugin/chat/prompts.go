// Package chat use act to chat with GPT-3.5-Turbo.
package chat

import (
	"encoding/json"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"os"
	"strings"
)

var (
	chatGPTPrompts []map[string]string
	prompts        string
)

func init() {
	// initialize the prompts first.
	getPromptsJson, err := os.ReadFile(engine.DataFolder() + "prompts-zh.json")
	if err != nil {
		panic(err)
	}
	_ = json.Unmarshal(getPromptsJson, &chatGPTPrompts)
	engine.OnRegex(`!test\schatgpt\sact\s(.*)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getAct := ctx.State["regex_matched"].([]string)[1]
		getLength := len(chatGPTPrompts)
		for i := 0; i < getLength; i++ {
			if chatGPTPrompts[i]["act"] == getAct {
				prompts = chatGPTPrompts[i]["prompt"]
				break
			}
		}
		if prompts == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有获得预设内容哦~"))
			return
		}
		key := sessionKey{
			group: ctx.Event.GroupID,
			user:  ctx.Event.UserID,
		}
		// 使用预设时，先清理掉之前的所有list。
		cache.Delete(key)
		messages := cache.Get(key)
		messages = append(messages, chatMessage{
			Role:    "user",
			Content: prompts,
		})
		resp, err := completions(messages, os.Getenv("gptkey"))
		if err != nil {
			ctx.SendChain(message.Text("Some errors occurred when requesting :( : ", err))
			return
		}
		reply := resp.Choices[0].Message
		reply.Content = strings.TrimSpace(reply.Content)
		messages = append(messages, reply)
		cache.Set(key, messages)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经使用预设 ", getAct, " ~返回内容如下:\n", reply.Content))
	})
}