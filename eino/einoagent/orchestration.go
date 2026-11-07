package einoagent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildEinoAgent(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToHistory = "InputToHistory"
		ChatTemplate   = "ChatTemplate"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambda(newLambda))
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	_ = g.AddEdge(compose.START, compose.END)
	_ = g.AddEdge(compose.START, InputToHistory)
	_ = g.AddEdge(InputToHistory, ChatTemplate)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
