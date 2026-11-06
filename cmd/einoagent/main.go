package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino/compose"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	graph := compose.NewGraph[string, string]()
	if err := graph.AddLambdaNode("assistant", compose.InvokableLambda(
		func(_ context.Context, input string) (string, error) {
			return fmt.Sprintf("收到你的消息：%s\n请配置 Eino ChatModel 以启用 AI 回复。", input), nil
		},
	)); err != nil {
		return fmt.Errorf("add assistant node: %w", err)
	}
	if err := graph.AddEdge(compose.START, "assistant"); err != nil {
		return fmt.Errorf("connect assistant graph start: %w", err)
	}
	if err := graph.AddEdge("assistant", compose.END); err != nil {
		return fmt.Errorf("connect assistant graph end: %w", err)
	}

	assistant, err := graph.Compile(ctx)
	if err != nil {
		return fmt.Errorf("compile assistant graph: %w", err)
	}

	fmt.Println("Eino 助手已启动。输入 /exit 退出。")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "/exit" {
			break
		}
		if input == "" {
			continue
		}

		output, err := assistant.Invoke(ctx, input)
		if err != nil {
			return fmt.Errorf("invoke assistant: %w", err)
		}
		fmt.Println(output)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	return nil
}
