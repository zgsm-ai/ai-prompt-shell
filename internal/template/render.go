package template

import (
	"ai-prompt-shell/internal/utils"
	"bytes"
	"context"
	"text/template"
)

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer() *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplate 加载模板到内存
func (r *Renderer) LoadTemplate(id, content string) error {
	t := template.New(id)
	t, err := t.Parse(content)
	if err != nil {
		return err
	}
	r.templates[id] = t
	return nil
}

// Render 渲染模板
func (r *Renderer) Render(ctx context.Context, id string, data interface{}) (string, error) {
	t, ok := r.templates[id]
	if !ok {
		return "", utils.ErrTemplateNotFound
	}

	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
