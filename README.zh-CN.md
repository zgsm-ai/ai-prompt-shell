# ai-prompt-shell

## 介绍

ai-prompt-shell是Prompt类型扩展的解释器。

使用text/template引擎，将Prompt模板处理成OpenAI兼容格式的LLM请求，调用LLM提供的AI能力。在模板处理过程中，ai-prompt-shell能够获取系统各模块设置的共享变量，调用各模块开放的工具接口，获取上下文信息，以扩展LLM的能力。

## 安装

```shell
go install github.com/zgsm-ai/ai-prompt-shell@latest
```

## 使用

1. 启动ai-prompt-shell

```shell
ai-prompt-shell
```

2. 注册Prompt类型扩展

```shell
smc extension add "zgsm.translator" -d "examples/extension/translator/package.json"
```

3. 注册ai-prompt-shell可用的工具

```shell
smc tool add "codebase.lookup_reference" -d "examples/tool/lookup_reference.json"
```

4. 注册ai-prompt-shell可用的Prompt模板

```shell
smc prompt add "translator.translate_zh_en" -d "examples/prompt/translate_zh_en.json"
```

5. 写入ai-prompt-shell可用的共享变量

```shell
smc variable set "completion.model" -v "deepseek-codelite-v3"
```

6. 调用ai-prompt-shell提供的接口，使用指定Prompt模板与LLM交互

```shell
ai-prompt-shell -p "translator.translate_zh_en" -i "你好，世界"
```

## 示例
