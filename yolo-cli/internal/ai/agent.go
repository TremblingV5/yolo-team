package ai

import (
	"context"
	"fmt"
	"log"

	"yolo-team/yolo-cli/internal/ai/tools"
	"yolo-team/yolo-cli/internal/handler"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// Agent wraps an LLM with tools and session management for conversational AI.
type Agent struct {
	runner     *adk.Runner
	chatModel  *openai.ChatModel
	cfg        *Config
	sessionMgr *SessionManager
}

// taggedTool wraps a tool to add <tool_call> and <tool_result> tags around its I/O.
type taggedTool struct {
	inner tool.InvokableTool
	name  string
}

func (t *taggedTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return t.inner.Info(ctx)
}

func (t *taggedTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	log.Printf("[tool] InvokableRun: tool=%s params_len=%d", t.name, len(arguments))
	result, err := t.inner.InvokableRun(ctx, arguments, opts...)
	if err != nil {
		log.Printf("[tool] InvokableRun error: tool=%s err=%v", t.name, err)
		return "", err
	}
	tagged := fmt.Sprintf("<tool_call data-tool-name=\"%s\">\n%s\n</tool_call>\n<tool_result>\n%s\n</tool_result>\n\n", t.name, arguments, result)
	log.Printf("[tool] InvokableRun result: tool=%s result_len=%d tagged_len=%d", t.name, len(result), len(tagged))
	return tagged, nil
}

// NewAgent creates a new Agent with the given configuration, tools, and session manager.
func NewAgent(cfg *Config, sessionMgr *SessionManager, ph *handler.ProjectHandler, eh *handler.ExecutorHandler, ih *handler.IssueHandler, th *handler.IssueHandler, dh *handler.DocumentHandler) (*Agent, error) {
	log.Printf("[AI Agent] creating agent: model=%s base_url=%s", cfg.Model, cfg.BaseURL)

	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
	})
	if err != nil {
		return nil, fmt.Errorf("create chat model: %w", err)
	}
	log.Printf("[AI Agent] openai chat model created")

	toolRegistry := tools.NewToolRegistry(ph, eh, ih, th, dh)
	allTools := toolRegistry.AllTools()

	// Add session tools (managed in the ai package to avoid circular dependency).
	sessionTools := []tool.InvokableTool{
		&sessionListTool{sm: sessionMgr},
		&sessionGetTool{sm: sessionMgr},
		&sessionDeleteTool{sm: sessionMgr},
	}
	allTools = append(allTools, sessionTools...)

	baseTools := make([]tool.BaseTool, 0, len(allTools))
	for _, t := range allTools {
		info, _ := t.Info(context.Background())
		name := ""
		if info != nil {
			name = info.Name
		}
		baseTools = append(baseTools, &taggedTool{inner: t, name: name})
	}
	log.Printf("[AI Agent] registered %d tools", len(baseTools))

	agent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Model:         chatModel,
		Instruction:   cfg.SystemPrompt,
		MaxIterations: 10,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: baseTools,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create chat model agent: %w", err)
	}
	log.Printf("[AI Agent] chat model agent created")

	runner := adk.NewRunner(context.Background(), adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	log.Printf("[AI Agent] runner created, streaming enabled")

	return &Agent{
		runner:     runner,
		chatModel:  chatModel,
		cfg:        cfg,
		sessionMgr: sessionMgr,
	}, nil
}

// ChatStream sends a user message and streams the AI response via callback.
func (a *Agent) ChatStream(ctx context.Context, sessionID, message string, callback func(chunk string)) (string, string, error) {
	log.Printf("[AI] ChatStream: session=%q msg_len=%d", sessionID, len(message))

	if sessionID == "" {
		var err error
		sessionID, err = a.sessionMgr.CreateSession()
		if err != nil {
			return "", "", fmt.Errorf("create session: %w", err)
		}
	}

	if err := a.sessionMgr.AddMessage(sessionID, "user", message); err != nil {
		return sessionID, "", fmt.Errorf("save user message: %w", err)
	}

	// Load history (last 50 rounds).
	msgs, _ := a.sessionMgr.GetHistory(sessionID, 50)
	log.Printf("[AI] built %d history messages for session=%s", len(msgs), sessionID)

	iter := a.runner.Run(ctx, msgs)

	var fullResponse string
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("[AI] runner error: %v", event.Err)
			return sessionID, "", fmt.Errorf("agent error: %w", event.Err)
		}

		msg, _, err := adk.TypedGetMessage(event)
		if err != nil {
			continue
		}
		if msg == nil {
			continue
		}

		if msg.Content != "" {
			fullResponse += msg.Content
			callback(msg.Content)
		}
	}

	log.Printf("[AI] ChatStream done: session=%s response_len=%d", sessionID, len(fullResponse))

	if fullResponse != "" {
		a.sessionMgr.AddMessage(sessionID, "assistant", fullResponse)
	} else {
		log.Printf("[AI] WARNING: empty response for session=%s", sessionID)
	}
	return sessionID, fullResponse, nil
}
