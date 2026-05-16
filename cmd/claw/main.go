// cmd/claw/main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/llxxgg/go-tiny-claw/internal/engine"
	"github.com/llxxgg/go-tiny-claw/internal/provider"
	"github.com/llxxgg/go-tiny-claw/internal/tools"
)

func main() {
	if os.Getenv("MINIMAX_API_KEY") == "" {
		log.Fatal("请先导出 MINIMAX_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()
	workDir = workDir + "/workdir/"

	// 2. 初始化真实的大脑 (指向MiniMax-M2.7，使用上一讲的 OpenAI 适配器)
	llmProvider := provider.NewMinimaxOpenAIProvider("MiniMax-M2.7")
	registry := tools.NewRegistry()

	// 挂载极简工具集
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))

	// 实例化核心引擎，关闭慢思考阶段，享受 YOLO 急速模式
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, false)

	// 发起一个需要局部修改的指令
	prompt := ` 我当前目录下有一个 server.go 文件。 请帮我把里面 "TODO: 增加鉴权逻辑" 下面的那个 if 语句，整个替换为： if user == nil { fmt.Println("Forbidden!") return } `
	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
