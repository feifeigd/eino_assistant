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
		InputToQuery   = "InputToQuery"
		ReactAgent     = "ReactAgent"
		RedisRetriever = "RedisRetriever"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(newLambda))
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(newLambda1))
	reactAgentKeyOfLambda, err := newLambda2(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, reactAgentKeyOfLambda)
	redisRetrieverKeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(RedisRetriever, redisRetrieverKeyOfRetriever, compose.WithOutputKey("documents"))
	_ = g.AddEdge(compose.START, InputToHistory)
	_ = g.AddEdge(compose.START, InputToQuery)
	_ = g.AddEdge(ReactAgent, compose.END)
	_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(RedisRetriever, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	_ = g.AddEdge(InputToQuery, RedisRetriever)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
