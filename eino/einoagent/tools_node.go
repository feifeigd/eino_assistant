package einoagent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &duckduckgo.Config{}

	bt, err = duckduckgo.NewTextSearchTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

type Tool1Impl struct {
	config *Tool1Config
}

type Tool1Config struct {
}

func newTool1(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &Tool1Config{}
	bt = &Tool1Impl{config: config}
	return bt, nil
}

func (impl *Tool1Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}

func (impl *Tool1Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	panic("implement me")
}
