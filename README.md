# ai-prompt-shell

## Introduction

ai-prompt-shell is an interpreter for Prompt-type extensions.

Using the text/template engine, it processes Prompt templates into LLM requests compatible with the OpenAI format, calling AI capabilities provided by LLMs. During template processing, ai-prompt-shell can access shared variables set by system modules, call tool interfaces exposed by various modules, and obtain context information to extend LLM capabilities.

## Installation

```shell
go install github.com/zgsm-ai/ai-prompt-shell@latest
```

## Usage

1. Start ai-prompt-shell

```shell
ai-prompt-shell
```

2. Register Prompt-type extensions

```shell
smc extension add "zgsm.translator" -d "examples/extension/translator/package.json"
```

3. Register tools available to ai-prompt-shell

```shell
smc tool add "codebase.lookup_reference" -d "examples/tool/lookup_reference.json"
```

4. Register Prompt templates available to ai-prompt-shell

```shell
smc prompt add "translator.translate_zh_en" -d "examples/prompt/translate_zh_en.json"
```

5. Set shared variables available to ai-prompt-shell

```shell
smc variable set "completion.model" -v "deepseek-codelite-v3"
```

6. Call interfaces provided by ai-prompt-shell to interact with LLMs using specified Prompt templates

```shell
smc prompt render "agent.code_review" -v "{\"language\": \"cpp\", \"code\": \"int main(){\n}\"}"
smc prompt chat "agent.code_review" -m "deepseek-v3" -v "{\"language\": \"cpp\", \"code\": \"int main(){\n}\"}"
```

## Examples
