package service

import (
	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/utils"
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

type RenderStats struct {
	Total         int64
	Success       int64
	Failures      int64
	TotalDuration time.Duration
	AvgDuration   time.Duration
	CacheHits     int64
	CacheMisses   int64
	RecentErrors  []error
	mu            sync.RWMutex
}

type ToolExecutor func(args ...interface{}) (interface{}, error)

type Renderer struct {
	refreshInterval time.Duration
	templates       map[string]*template.Template
	funcMap         template.FuncMap
	stats           *RenderStats
}

var renderer *Renderer = NewRenderer()

/**
 * Create new renderer instance with default settings
 * @return initialized renderer with 10s refresh interval
 */
func NewRenderer() *Renderer {
	return &Renderer{
		refreshInterval: 10 * time.Second,
		templates:       make(map[string]*template.Template),
		funcMap:         make(template.FuncMap),
		stats:           &RenderStats{},
	}
}

/**
 * Construct context data by combining environment variables with template args
 * @param args additional args to include in context
 * @return combined map containing all variables
 */
func constructContextData(args map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for k, v := range environs.All() {
		data[k] = v
	}
	data["args"] = args
	return data
}

/**
 * Convert tool/prompt ID to template function/variable name
 * @param id input ID string
 * @return converted variable name (lowercase with dots replaced)
 */
func idToVariable(id string) string {
	return strings.ReplaceAll(strings.ToLower(id), ".", "_")
}

/**
 * Update template functions when tools are refreshed
 */
func onRefreshTools() {
	newFuncs := make(template.FuncMap)
	for k, v := range tools.All() {
		newFuncs[idToVariable(k)] = newToolExecutor(&v)
	}
	renderer.funcMap = newFuncs
}

/**
 * Create executor function for tool
 * @param t tool definition to create executor for
 * @return executor function that calls the tool
 */
func newToolExecutor(t *dao.Tool) ToolExecutor {
	return func(args ...interface{}) (interface{}, error) {
		return Call(context.Background(), t, args)
	}
}

/**
 * Update templates when prompts are refreshed
 */
func onRefreshPrompts() {
	for key, content := range prompts.All() {
		if content.Prompt.Prompt != "" {
			key = key + ".prompt"
			t, err := template.New(key).Funcs(renderer.funcMap).Parse(content.Prompt.Prompt)
			if err != nil {
				continue
			}
			renderer.templates[key] = t
		} else if content.Messages != nil {
			for i, _ := range content.Messages {
				msgkey := fmt.Sprintf("%s.messages.%d", key, i)
				t, err := template.New(msgkey).Funcs(renderer.funcMap).Parse(content.Messages[i].Content)
				if err != nil {
					continue
				}
				renderer.templates[msgkey] = t
			}
			continue
		} else {
			logrus.Errorf("prompt %s is invalid", key)
		}
	}
}

/**
 * Execute template rendering with given arguments
 * @param templateKey identifier for template to render
 * @param args input values for template
 * @return rendered template as string
 * @return error if template not found or execution fails
 */
func renderTemplate(templateKey string, args map[string]interface{}) (string, error) {
	t, ok := renderer.templates[templateKey]
	if !ok {
		return "", utils.ErrBug
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, constructContextData(args))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

/**
 * Render all messages in conversation
 * @param prompt_id ID of the message template set
 * @param messages message templates to render
 * @param args input values for template
 * @return fully rendered messages
 * @return error if rendering fails
 */
func renderMessages(prompt_id string, messages []dao.Message, args map[string]interface{}) ([]dao.Message, error) {
	var results []dao.Message
	for i, message := range messages {
		content, err := renderTemplate(fmt.Sprintf("%s.messages.%d", prompt_id, i), args)
		if err != nil {
			return []dao.Message{}, err
		}
		results = append(results, dao.Message{
			Content: content,
			Role:    message.Role,
		})
	}
	return results, nil
}

/**
 * Render prompt with args
 * @param prompt_id ID of prompt to render
 * @param args input args for template
 * @return type of rendered content ("prompt" or "messages")
 * @return rendered content or messages
 * @return error if rendering fails
 */
func RenderPrompt(prompt_id string, args map[string]interface{}) (string, interface{}, error) {
	prompt, origin := prompts.Get(prompt_id)
	if origin == dao.PromptOrigin_Notexist {
		return "", "", utils.ErrPromptNotFound
	}
	if prompt.Prompt != "" {
		text, err := renderTemplate(prompt_id+".prompt", args)
		return "prompt", text, err
	} else if prompt.Messages != nil {
		messages, err := renderMessages(prompt_id, prompt.Messages, args)
		return "messages", messages, err
	}
	return "", "", utils.ErrPromptInvalid
}
