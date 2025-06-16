# AI-Prompt-Shell 设计文档

## 概述

AI-Prompt-Shell是执行Prompt类型扩展的引擎，负责获取各类上下文信息，利用这些信息将Prompt模板渲染成最终Prompt，发送给LLM进行交互。

这是一个常驻进程，采用go语言实现，对外提供RESTful API接口。

## 背景

Prompt类型扩展，指一种通过定义Prompt模板，供其它AI程序调用大语言模型的扩展形式。

Prompt类型扩展，给了用户一种简单快捷的方法，充分利用AI-Prompt-Shell获取上下文知识的能力，来调用及扩展LLM的智能。

AI-Prompt-Shell具有获取共享变量，调用工具的能力，它利用text/template的扩展变量和外部函数的机制，在渲染模板过程中应用这些能力，动态构造发送给LLM的请求内容，令LLM变相具有访问共享变量，调用工具的扩展能力。

## 关键需求

AI-Prompt-Shell的核心功能特性有：

1. Prompt模板管理：加载、缓存和更新Prompt扩展中的模板，并根据ID查找调用用户需要的Prompt模板
2. 支持请求参数注入: 模板渲染时可以访问用户请求时传递的参数，作为模板变量使用
3. 支持共享变量注入：模板渲染时，可以从redis获取共享变量，作为模板变量使用
4. 支持扩展工具：从redis读取工具定义，构建成模板渲染过程中可用的函数表，作为模板扩展函数使用。这些工具，采用RESTful API/GRPC/MCP等协议调用系统内部应用提供的接口，执行结果作为模板函数的返回值，输出到渲染后的模板中，替换模板中的函数调用标记
5. 模板渲染：使用go 'text/template'的模板引擎将Prompt模板渲染成最终Prompt/LLM调用请求
6. 获取渲染结果：用户调用`/api/render/prompts/{prompt_id}`，可获取指定Prompt模板的渲染结果
7. 调用LLM：用户调用`/api/chat/prompts/{prompt_id}`，可将Prompt模板渲染完毕后作为参数发送给指定LLM，然后接收LLM返回的请求结果。

## 对外接口

AI-Prompt-Shell支持的RESTful API接口如下：

| 功能 | 接口 | 说明|
|------|------|----|
| 列出Prompt类型扩展 | `GET /api/extensions` | 列出系统有哪些Prompt类型扩展可用 |
| 获取Prompt类型扩展的详情 | `GET /api/extensions/{extension_id}`| 获取指定Prompt类型扩展的详情 |
| 列出Prompt模板 | `GET /api/prompts` | 列出系统有哪些Prompt模板可用 |
| 获取Prompt模板详情 | `GET /api/prompts/{prompt_id}` | 获取指定Prompt模板的详情 |
| 列出共享变量 | `GET /api/environs` | 列出系统有哪些共享变量可用 |
| 获取共享变量值 | `GET /api/environs/{environ_id}` | 获取共享变量的值|
| 列出Tool定义 | `GET /api/tools` | 列出系统有哪些工具可用 |
| 获取Tool定义详情 | `GET /api/tools/{tool_id}` | 获取指定工具的定义详情|
| 获取渲染后的Prompt | `POST /api/render/prompts/{prompt_id}` | 获取指定Prompt模板的渲染结果 |
| 调用LLM | `POST /api/chat/prompts/{prompt_id}` | 采用指定的Prompt模板，使用渲染结果调用LLM，获取LLM的输出结果|

详情请参考下述章节。

### 渲染Prompt

```bash
POST /api/render/prompts/{prompt_id}
```

请求参数：

```json
{
    "variables": {
        "key1": "value1",
        "key2": "value2"
    }
}
```

响应格式：

```json
{
    "rendered_prompt": "渲染后的Prompt", 
    "status": "success"
}
```

返回的rendered_prompt字段可以是文本或JSON值，根据Prompt扩展定义而定。

### 调用LLM

```bash
POST /api/chat/prompts/{prompt_id}
```

请求参数：

```json
{
  "model": "deepseek-v3",
  "variables": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

响应格式是openai chat接口的标准回复格式：

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-3.5-turbo",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I assist you today?"
      },
      "finish_reason": "stop"
    }
  ]
}
```

### 错误处理

| 错误码 | 说明 |
|--|--|
| 404 | Prompt ID不存在 |
| 400 | 缺少必要变量 |
| 500 | 模板渲染错误 |

## 原理

系统提供两大类机制将特定知识嵌入LLM调用请求中，扩展LLM能力。

这两大类机制，分别是“共享变量(SharedVariable)”，“扩展工具(Tool)”。

共享变量，即系统各模块把可供LLM使用的上下文信息，以键值对的形式写到redis作为共享变量。
扩展工具，即以RESTful API/GRPC/MCP等协议接口方式访问各应用提供的能力，完成特定某个工作，类似FunctionCall。
把可供LLM访问的扩展工具的接口定义，以元数据的方式写入到redis，则AI-Prompt-Shell/ChatRAG等需要调用扩展工具的程序，可以通过遍历读取“扩展工具”元数据方式，获得动态调用这些扩展工具的能力。

## 数据结构

### 总述

AI-Prompt-Shell从redis获取其它程序注册的四类信息：

- 遍历redis的'shenma:extensions:'目录，加载扩展定义。
- 遍历redis的'shenma:templates:'目录，加载Prompt类型扩展所定义的Prompt模板。
- 遍历redis的'shenma:environs:'目录，加载共享变量，构建成给Prompt模板使用的共享变量查找表。
- 遍历redis的'shenma:tools:'目录，加载'扩展工具'的元数据，构建给Prompt模板使用的函数查找表。

这四类信息的定义如下所述。

### 扩展

redis 'shenma:extensions:'目录下，存储若干扩展的定义。

Prompt类型扩展是扩展的一种，即PromptExtension，它是一个JSON对象，包含如下字段：

```json
{
  "name": "evaluator",
  "publisher": "zgsm-ai",
  "displayName": "项目质量评估器",
  "icon": "images/evaluator.png",
  "description": "评估项目代码的质量，给出各个质量维度的评估报告",
  "version": "1.0.0",
  "extensionType": "prompt",
  "license": "Apache-2.0",
  "engines": {
    "name": "ai-prompt-shell",
    "version": ">=1.0.0"
  },
  "contributes": {
    "prompts": [
      {
        "name": "evaluate_quality",
        "messages": [
          {
            "role": "system",
            "content": "You are a code review assistant who can evaluate the quality of the project."
          },
          {
            "role": "user",
            "content": "Evaluate the project code from the following dimensions and give a quality assessment report.\nEvaluation dimensions:\n{{.vscode.rules}}\nCode context:\n{{.codebase.current_project}}"
          }
        ],
        "supports": [
          "chat",
          "codereview"
        ],
        "parameters": [
          {
            "name": "repo",
            "type": "string",
            "default": "",
            "description": "仓库地址"
          }
        ],
        "returns": []
      }
    ],
    "languages": ["c++", "c", "lua", "python"],
    "dependences": [
      {
        "name": "scan-gitrepo",
        "version": "^1.0.0",
        "failStrategy": "abort"
      },
      {
        "name": "context",
        "version": "^1.2.0",
        "failStrategy": "ignore"
      }
    ]
  }
}
```

扩展定义的JSONSCHEMA请参考`jsonschema/extension.json`。

说明：

| JsonPath | 说明 |
|---------|-----|
|name | 扩展名称 |
|publisher | 扩展发布者 |
|displayName | 扩展显示名称 |
|icon | 扩展图标 |
|description | 扩展描述 |
|version | 扩展版本 |
|extensionType | 扩展类型，目前只支持prompt |
|license | 扩展许可证 |
|engines | 扩展引擎信息 |
|contributes | 扩展提供的能力 |
|contributes.prompts| 扩展提供的Prompt模板形式的接口 |
|contributes.prompts.name | 模板名称 |
|contributes.prompts.messages | 模板需要渲染的内容，这个字段说明要渲染的是消息列表，这个字段和userPrompt互斥 |
|contributes.prompts.userPrompt | 模板需要渲染的内容，这个字段说明要渲染的是用户请求文本串，这个字段和messages互斥 |
|contributes.prompts.supports | 模板支持的场景，目前支持chat、completion、codereview |
|contributes.prompts.parameters | 用户向本接口发送请求的参数列表定义 |
|contributes.prompts.returns | 本接口返回值列表定义 |

### Prompt模板

redis 'shenma:templates:'目录下，存储若干Prompt模板定义。

Prompt模板即Prompt类型扩展中定义的'contributes.prompts'字段。狭义上，也可以特指'contributes.prompts.messages'字段和'contributes.prompts.userPrompt'字段。

### 共享变量

redis 'shenma:environs:'目录下，存储其它服务提供的共享变量。

共享变量的值，是一个JSON值，可以是文本，数字，bool值，对象，数组。

### 扩展工具

redis 'shenma:tools:'目录下，存储若干工具定义。

工具定义，即ToolCall，是一个JSON对象，包含如下字段：

```json
{
  "name": "translate_zh_en",
  "module": "translator",
  "type": "restful", //grpc, mcp
  "url": "/translate/zh/en",
  "description": "将输入代码中的中文注释、字符串翻译为英文",
  "supports": [
    "chat",
    "completion",
    "codereview"
  ],
  "parameters": {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
      "code": {
        "type": "string",
        "description": "代码文本"
      }
    },
    "required": [
      "code"
    ]
  },
  "returns": {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
      "translated_code": {
        "type": "string",
        "description": "翻译后的代码文本"
      }
    }
  }
}
```

对应的jsonschema请参考`jsonschema/tool.json`：

说明：

| JsonPath | 说明 |
|--------|-----|
|name | 扩展工具名称 |
|module | 扩展工具所属模块 |
|type | 扩展工具接口类型，目前支持restful。后续支持grpc、mcp |
|url| 扩展工具接口地址 |
|description | 扩展工具描述 |
|supports | 扩展工具支持的场景，目前支持chat、completion、codereview |
|parameters | 扩展工具参数列表定义 |
|returns | 扩展工具返回值列表定义 |

## 结构设计

| 目录/文件  | 职责 | 说明 |
|------|------|----|
| api | 接口 | 提供RESTful API接口 |
| internal| 内部机制 | 实现通用的机制性代码，如核心逻辑 |
| service | 服务 | 实现接口所依赖的业务逻辑 |
| internal/extension| 扩展管理 | 扩展的管理，包括加载、查找等|
| internal/template | 模板管理 | 模板的管理，包括加载、查找、渲染、准备渲染模板需要的数据对象，函数表等|
| internal/variable | 变量管理 | 共享变量管理，包括遍历、加载、构建共享变量表等 |
| internal/tool| 扩展工具 | 扩展工具的管理，包括遍历，加载，构建工具查找表等 |
| internal/llm | 大模型调用 | |
| internal/redis | redis | 实现redis操作 |
| internal/utils | 工具函数 | 实现工具函数 |
| internal/config | 配置 | 实现配置加载 |
| internal/logger | 日志 | 实现日志记录 |
| internal/http | HTTP客户端 | 实现RESTful类型扩展工具的访问逻辑 |

## 流程或机制

### 模板语法

Prompt类型扩展中的Prompt模板，采用go语言text/template库使用的模板语法，支持的特性包括：

- 变量注入：`{{.variable}}`
- 条件判断：`{{if .condition}}...{{end}}`
- 循环：`{{range .items}}...{{end}}`
- 作用域： `{{with .object}}...{{end}}`
- 函数调用：`{{codebase_look_ref "CreateObject"}}`
- 调用其它模板(相当于调用子过程)：`{{template "name" .}}`

AI-Prompt-Shell使用go的text/template完成模板的实例化。实例化前需要先构建数据对象，以及函数查找表。

### 扩展加载

AI-Prompt-Shell从redis中加载所有Prompt类型扩展，获取扩展定义的Prompt模板，缓存在Prompt模板查找表中。

AI-Prompt-Shell从redis中遍历'shenma:templates:'目录下的所有内容。

并根据遍历得到的KEY(如'shenma:templates:{extension-name}:{prompt-name}')获取对应的值，即Prompt扩展，内容Prompt模板定义。

Prompt模板查找表是`map[string]PromptTemplate`类型，map的键是调用keyToJsonPath(KEY)得到的结果，
PromptTemplate即从Redis中加载的Prompt模板。

### 构建数据对象

AI-Prompt-Shell从多种途径获取变量，并按照下述规则构建为一个叫做context的数据对象，提供给text/template做Prompt生成。

AI-Prompt-Shell可用的变量，包括两大部分：从用户请求中获取的'variables'字段，从redis中读取的共享变量。

AI-Prompt-Shell渲染模板时，从redis中获取需要的变量，根据下述规则组织成JSON值，传递给模板渲染引擎，作为上下文变量。

redis中所有前缀为'shenma:environs:'的KEY，其值都是一个JSON值(包括字符串，数字、对象、数组等)。

系统其它模块，将可供AI-Prompt-Shell进行渲染的变量统一写入到redis的'shenma:environs:'下，供
AI-Prompt-Shell渲染模板使用。

AI-Prompt-Shell遍历'shenma:environs:'为前缀的所有KEY，读取其JSON值。

对于每个读取到的JSON值，获取其KEY，通过KeyToJsonPath(KEY)计算得到一个路径Path，然后将该JSON值保存到JSON对象context的中Path路径下。

keyToJsonPath的计算方法为：去掉KEY的头部'shenma:environs:'，然后把剩余部分的':'替换成'.'，得到的就是JsonPath。

比如，redis有如下键值对：

1. KEY: 'shenma:environs:vscode:programming_language', VALUE: "go"
2. KEY: 'shenma:environs:vscode:frameworks', VALUE: ["gin", "gorm", "gin-swagger"]

用户请求中的variables字段如下：

```json
"variables": {
    "key1": "value1",
    "key2": "value2"
  }
```

根据这两部分内容，构造出来的JSON对象context如下：

```json
{
  "vscode": {
    "programming_language": "go",
    "frameworks": ["gin","gorm", "gin-swagger"]
  },
  "variables": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

各模块注册给AI-Prompt-Shell的变量KEY，格式通常为：'shenma:environs:{module}:{subject}'。

目前，确定会生成上下文信息到redis，给AI-Prompt-Shell使用的模块如下：

| 模块 | Key前缀 | 说明 |
|------|--------|-----|
| vscode | 'shenma:environs:vscode:*' | vscode编辑上下文 |
| codebase | 'shenma:environs:codebase:*' | 代码索引的文件列表、目录结构等信息 |
| chat | 'shenma:environs:chat:*' | 对话历史 |
| completion | 'shenma:environs:completion:*' | 补全历史 |
| model | 'shenma:environs:model:*' | 提供模型列表、模型特征描述等信息 |
| modules | 'shenma:environs:modules:*' | 给AI-Prompt-Shell提供信息的各模块的元信息 |

### 构建函数查找表

AI-Prompt-Shell从redis中获取工具定义，并构建为工具查找表(`type ToolRegistry=map[string]Tool`)。
使用text/template渲染Prompt模板时，会根据工具查找表构建模板函数表Funcs，作为生成Prompt过程中可用的外部函数。

redis中所有前缀为'shenma:tools'的KEY，其值都是工具类型。AI-Prompt-Shell定时遍历redis的'shenma:tools:'目录，刷新工具查找表。

工具查找表的键名toolName是keyToFunc(KEY)的计算结果。

keyToFunc(KEY)的逻辑：

- 去掉 KEY 的前缀 'shenma:tools:'
- 将剩余的路径中的':'替换为'_'（如 'codebase:lookup_ref' → 'codebase_lookup_ref'）

keyToFunc(KEY)计算样例：

- keyToFunc('shenma:tools:codebase:lookup_ref') = 'codebase_lookup_ref'
- keyToFunc('shenma:tools:codebase:caller') = 'codebase_caller'
- keyToFunc('shenma:tools:mcp:chrome:xx') = 'mcp_chrome_xx'

渲染Prompt模板时，需要遍历工具查找表，构建模板函数表Funcs，Funcs的键名和ToolRegistry键名(即toolName)一致，键值为函数`TemplateFunc(toolName)`返回的函数对象。

### 模板渲染错误处理

1. 模板语法错误：
   - 在加载模板时进行检查，如果有语法错误直接返回400错误

2. 变量未定义：
   - 检查所有引用变量是否存在于上下文，缺失任何变量返回400错误

3. 工具调用失败：
   - 记录错误日志，但不中断渲染
   - 在模板中使用默认值或空值替代

4. 渲染超时：
   - 内置500ms超时控制
   - 超时后返回503服务不可用错误

### LLM调用重试策略

1. 首次调用失败后，等待100ms进行第一次重试
2. 第二次失败后，等待300ms进行第二次重试
3. 最多重试2次
4. 如果所有重试都失败，返回502 Bad Gateway错误

### 性能优化策略

1. 模板缓存：
   - 自动缓存最近使用的100个模板
   - LRU算法淘汰不常用模板

2. 变量缓存：
   - 高频访问的redis变量本地缓存5秒

3. 批量调用工具：
   - 使用goroutine并发执行不相关工具调用

4. 内存池：
   - 重用渲染结果缓冲区
   - 减少GC压力
