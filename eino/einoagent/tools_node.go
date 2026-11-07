package einoagent

import (
	"context"

	"eino_assistant/pkg/tool/einotool"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ToolImpl struct {
	config *ToolConfig
}

// 实现接口
var _ tool.BaseTool = (*ToolImpl)(nil)

type ToolConfig struct {
}

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &ToolConfig{}
	bt = &ToolImpl{config: config}
	return bt, nil
}

func (impl *ToolImpl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}

// func (impl *ToolImpl) StreamableRun(ctx context.Context, input argumentsInJSON, opts tool.Option) (schema.StreamReader[string], error) {
// 	panic("implement me")
// }

func GetTools(ctx context.Context) ([]tool.BaseTool, error) {
	einoAssistantTool, err := newEinoAssistantTool(ctx)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{
		einoAssistantTool,
	}, nil
}

func newEinoAssistantTool(ctx context.Context) (tool.BaseTool, error) {
	return einotool.NewEinoAssistantTool(ctx, nil)
}
