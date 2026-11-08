package main

import (
	"context"
	"eino_assistant/eino/knowledge_indexing"
	"eino_assistant/pkg/env"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	clc "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/document"
	"github.com/coze-dev/cozeloop-go"
)

func init() {
	env.MustHasEnvs("ARK_API_KEY", "ARK_EMBEDDING_MODEL")
}

func main() {
	repoDir := flag.String("dir", "./eino-docs", "thie directory of markdown files")
	flag.Parse()

	finalPath := *repoDir
	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		fixedPath := filepath.Join("cmd/knowledge_indexing", *repoDir)
		if _, err := os.Stat(fixedPath); err == nil {
			finalPath = fixedPath
		}
	}

	cozeloopApiToken := os.Getenv("COZELOOP_API_TOKEN")
	cozeloopWorkspaceID := os.Getenv("COZELOOP_WORKSPACE_ID")

	ctx := context.Background()
	var handlers []callbacks.Handler
	if cozeloopApiToken != "" && cozeloopWorkspaceID != "" {
		client, err := cozeloop.NewClient(
			cozeloop.WithAPIToken(cozeloopApiToken),
			cozeloop.WithWorkspaceID(cozeloopWorkspaceID),
		)
		if err != nil {
			panic(err)
		}
		defer client.Close(ctx)
		handlers = append(handlers, clc.NewLoopHandler(client))

	}
	callbacks.AppendGlobalHandlers(handlers...)

	if err := indexMarkdownFiles(ctx, finalPath); err != nil {
		panic(err)
	}
	fmt.Println("index success")

}

func indexMarkdownFiles(ctx context.Context, dir string) error {
	runner, err := knowledge_indexing.BuildKnowledgeIndexing(ctx)
	if err != nil {
		return fmt.Errorf("build index graph failed: %w", err)
	}

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("work dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}
		fmt.Printf("[start] indexing file: %s\n", path)

		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			return fmt.Errorf("invoke index graph failed: %w", err)
		}
		fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))
		return nil
	})
}
