package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"eino_assistant/pkg/env"

	"github.com/cloudwego/eino-ext/callbacks/apmplus"
	clc "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino-ext/callbacks/langfuse"
	"github.com/cloudwego/eino-ext/devops"
	"github.com/coze-dev/cozeloop-go"

	"github.com/cloudwego/eino/callbacks"

	"eino_assistant/pkg/mem"
)

type ()

var id = flag.String("id", "", "conversation id")

var memory = mem.GetDefaultMemory()

var cbHandler callbacks.Handler

func main() {
	flag.Parse()

	ctx := context.Background()

	// 开启Eino 的可视化调试能力
	if err := devops.Init(ctx); err != nil {
		log.Printf("[eino dev] init failed, err=%v", err)
		return
	}

	if *id == "" {
		*id = strconv.Itoa(rand.Intn(1000000))
	}

	if err := Init(ctx); err != nil {
		log.Printf("[eino agent] init failed, err=%v", err)
		return
	}

	// Start interactive dialogue
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[eino agent] read input failed, err=%v", err)
			break
		}
		input = strings.TrimSpace(input)
		if input == "" || input == "exit" || input == "quit" {
			break
		}
		// Call the RunAgent with the input
		if stream, err := RunAgent(ctx, *id, input); err != nil {
			log.Printf("[eino agent] run agent failed, err=%v", err)
			continue
		} else {
			log.Printf(": ")
			for {
				msg, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("[eino agent] receive message failed, err=%v", err)
					break
				}
				log.Print(msg.Content)
			}
			fmt.Println("\n")
		}
	}
	log.Printf("[eino agent] ended")
}

func Init(ctx context.Context) error {
	env.MustHasEnvs("OPENAI_API_KEY", "OPENAI_API_BASE", "OPENAI_API_MODEL")

	os.MkdirAll("log", 0755)
	f, err := os.OpenFile("log/eino.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	cbConfig := &LogCallbackConfig{
		Detail: true,
		Writer: f,
		Debug:  os.Getenv("DEBUG") == "true",
	}

	// this is for invoke option of WithCallback
	cbHandler = LogCallback(cbConfig)

	callbackHandlers := make([]callbacks.Handler, 0)

	// 火山引擎：全栈可观测平台
	if os.Getenv("APMPLUS_APP_KEY") != "" {
		region := os.Getenv("APMPLUS_REGION")
		if region == "" {
			region = "cn-beijing"
		}
		fmt.Println("[eino agent] INFO: use apmplus as callback, watch at: https://console.volcengine.com/apmplus-server")
		cbh, _, err := apmplus.NewApmplusHandler(&apmplus.Config{
			Host:        fmt.Sprintf("apmplus-%s.volces.com:4317", region),
			AppKey:      os.Getenv("APMPLUS_APP_KEY"),
			ServiceName: "eino-assistant",
			Release:     "release/v0.1.0",
		})
		if err != nil {
			log.Fatal(err)
		}
		callbackHandlers = append(callbackHandlers, cbh)
	}

	if os.Getenv("LANGFUSE_PUBLIC_KEY") != "" && os.Getenv("LANGFUSE_SECRET_KEY") != "" {
		fmt.Println("[eino agent] INFO: use langfuse as callback, watch at: https://cloud.langfuse.com")
		cbh, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
			Name:      "eino-assistant",
			PublicKey: os.Getenv("LANGFUSE_PUBLIC_KEY"),
			SecretKey: os.Getenv("LANGFUSE_SECRET_KEY"),
			Host:      "https://cloud.langfuse.com",
			Public:    true,
			Release:   "release/v0.1.0",
			UserID:    "user-123",
			Tags:      []string{"env:dev", "team:d7kj"},
		})

		callbackHandlers = append(callbackHandlers, cbh)
	}

	// https://loop.coze.cn/open/docs/cozeloop/go-sdk#4a8c980e
	cozeloopApiToken := os.Getenv("COZELOOP_API_TOKEN")
	if cozeloopApiToken != "" {
		cozeloopWorkspaceID := os.Getenv("COZELOOP_WORKSPACE_ID")
		if cozeloopWorkspaceID != "" {
			fmt.Println("[eino agent] INFO: use cozeloop as callback, watch at: https://cozeloop.com")
			client, err := cozeloop.NewClient(
				cozeloop.WithAPIToken(cozeloopApiToken),
				cozeloop.WithWorkspaceID(cozeloopWorkspaceID),
			)

			if err != nil {
				log.Fatal(err)
			}
			defer client.Close(ctx)
			callbackHandlers = append(callbackHandlers, clc.NewLoopHandler(client))
		}
	}

	if len(callbackHandlers) > 0 {
		callbacks.InitCallbackHandlers(callbackHandlers)
	}
	return nil
}
