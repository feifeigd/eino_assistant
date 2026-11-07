package einoagent

import (
	"context"

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

func GetTools() ([]tool.BaseTool, error) {
	return nil, nil
}
