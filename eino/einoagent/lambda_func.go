package einoagent

import (
	"context"
	"time"
)

// lambda 节点的实现

// newLambda component initialization function of node 'InputToHistory' in graph 'EinoAgent'
func newLambda(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// newLambda1 component initialization function of node 'InputToQuery' in graph 'EinoAgent'
func newLambda1(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	return input.Query, nil
}
