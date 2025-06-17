# AI-Prompt-Shell Design Document

## Overview

AI-Prompt-Shell is an engine for executing Prompt-type extensions, responsible for acquiring various context information and using this information to render Prompt templates into final Prompts, which are then sent to LLMs for interaction.

This is a resident process implemented in Go, providing RESTful API interfaces externally.

## Background

Prompt-type extensions refer to a form of extension that defines Prompt templates for other AI programs to call large language models.

Prompt-type extensions provide users with a simple and fast method to fully utilize AI-Prompt-Shell's ability to acquire contextual knowledge, thereby invoking and extending LLM capabilities.

AI-Prompt-Shell has the capabilities to acquire shared variables and call tools. It leverages text/template's variable expansion and external function mechanisms to apply these capabilities during template rendering, dynamically constructing requests to be sent to LLMs. This indirectly gives LLMs the ability to access shared variables and call tools.

## Key Requirements

The core features of AI-Prompt-Shell include:

1. Prompt template management: Load, cache, and update templates in Prompt extensions, and locate Prompt templates needed by users based on IDs
2. Support for request parameter injection: Access parameters passed during user requests during template rendering, using them as template variables
3. Support for shared variable injection: During template rendering, retrieve shared variables from Redis for use as template variables
4. Support for extension tools: Read tool definitions from Redis and construct them into function tables available during template rendering, serving as template extension functions. These tools use protocols such as RESTful API/GRPC/MCP to call interfaces provided by internal system applications, with execution results serving as return values for template functions, outputting to the rendered templates to replace function call markers
5. Template rendering: Use Go's 'text/template' template engine to render Prompt templates into final Prompts/LLM call requests
6. Obtaining rendering results: Users can call `/api/prompts/{prompt_id}/render` to obtain the rendering results of specified Prompt templates
7. Calling LLMs: Users can call `/api/prompts/{prompt_id}/chat` to send Prompt templates as parameters after rendering to specified LLMs, then receive the request results returned by LLMs.

## External Interfaces

The RESTful API interfaces supported by AI-Prompt-Shell are as follows:

| Functionality | Interface | Description |
|------|------|----|
| List Prompt-type extensions | `GET /api/extensions` | List available Prompt-type extensions in the system |
| Get details of a Prompt-type extension | `GET /api/extensions/{extension_id}` | Get details of a specified Prompt-type extension |
| List Prompt templates | `GET /api/prompts` | List available Prompt templates in the system |
| Get details of a Prompt template | `GET /api/prompts/{prompt_id}` | Get details of a specified Prompt template |
| Get rendered Prompt | `POST /api/prompts/{prompt_id}/render` | Get rendering results of a specified Prompt template |
| Call LLM | `POST /api/prompts/{prompt_id}/chat` | Use specified Prompt template, call LLM with rendering results, and get output from LLM |
| List shared variables | `GET /api/environs` | List available shared variables in the system |
| Get value of a shared variable | `GET /api/environs/{environ_id}` | Get the value of a shared variable |
| List tool definitions | `GET /api/tools` | List available tools in the system |
| Get details of a tool definition | `GET /api/tools/{tool_id}` | Get definition details of a specified tool |

For details, please refer to the following sections.

### Render Prompt
POST /api/prompts/{prompt_id}/render
```

Request parameters:

```json
{
    "args": {
        "key1": "value1",
        "key2": "value2"
    }
}
```

Response format:

```json
{
    "rendered_prompt": "Rendered Prompt", 
    "status": "success"
}
```

The rendered_prompt field can be text or JSON value, depending on the Prompt extension definition.

### Call LLM

```bash
POST /api/prompts/{prompt_id}/chat
```

Request parameters:

```json
{
  "model": "deepseek-v3",
  "args": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

Response follows the standard OpenAI chat interface format:

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

### Error Handling

| Error Code | Description |
|--|--|
| 404 | Prompt ID does not exist |
| 400 | Missing required args |
| 500 | Template rendering error |

## Principles

The system provides two main mechanisms to embed specific knowledge into LLM request calls and extend LLM capabilities.

These two mechanisms are "Shared Variables" and "Extension Tools".

Shared Variables: Various system modules write context information available for LLM usage into Redis as key-value pairs to serve as shared variables.

Extension Tools: Access capabilities provided by various applications through interface protocols such as RESTful API/GRPC/MCP to complete specific tasks, similar to FunctionCall. By writing interface definitions of extension tools available for LLM access into Redis as metadata, programs like AI-Prompt-Shell/ChatRAG that need to call extension tools can dynamically invoke these extension tools by traversing and reading the "Extension Tools" metadata.

## Data Structures

### Overview

AI-Prompt-Shell retrieves four types of information registered by other programs from Redis:

- Traverse Redis's 'shenma:extensions:' directory to load extension definitions.
- Traverse Redis's 'shenma:templates:' directory to load Prompt templates defined by Prompt-type extensions.
- Traverse Redis's 'shenma:environs:' directory to load shared variables and construct a shared variable lookup table for Prompt templates.
- Traverse Redis's 'shenma:tools:' directory to load metadata of 'extension tools' and construct a function lookup table for Prompt templates.

The definitions of these four types of information are described below.

### Extensions

Under Redis's 'shenma:extensions:' directory, definitions of various extensions are stored.

Prompt-type extensions are a type of extension, specifically PromptExtension, which is a JSON object containing the following fields:
```json
{
  "name": "evaluator",
  "publisher": "zgsm-ai",
  "displayName": "Project Quality Evaluator",
  "icon": "images/evaluator.png",
  "description": "Evaluates the quality of project code and provides an assessment report across various quality dimensions",
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
            "description": "Repository URL"
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

The JSONSchema for extension definitions can be found in `jsonschema/extension.json`.

Description:

| JsonPath | Description |
|---------|-----|
|name | Extension name |
|publisher | Extension publisher |
|displayName | Extension display name | 
|icon | Extension icon |
|description | Extension description |
|version | Extension version |
|extensionType | Extension type, currently only supports prompt |
|license | Extension license |
|engines | Extension engine information |
|contributes | Capabilities provided by the extension |
|contributes.prompts| Prompt template interfaces provided by the extension |
|contributes.prompts.name | Template name |
|contributes.prompts.messages | Content to be rendered, indicates a message list to render (mutually exclusive with prompt) |
|contributes.prompts.prompt | Content to be rendered, indicates a user request text string to render (mutually exclusive with messages) |
|contributes.prompts.supports | Supported scenarios, currently supports chat, completion, codereview |
|contributes.prompts.parameters | Parameter definitions for requests to this interface | 
|contributes.prompts.returns | Return value definitions for this interface |

### Prompt Templates

Under Redis's 'shenma:templates:' directory, various Prompt template definitions are stored.

Prompt templates refer to the 'contributes.prompts' field defined in Prompt-type extensions. In a narrow sense, they can also specifically refer to the 'contributes.prompts.messages' field and 'contributes.prompts.prompt' field.

### Shared Variables

Under Redis's 'shenma:environs:' directory, shared variables provided by other services are stored.

The value of a shared variable is a JSON value, which can be text, numbers, boolean values, objects, or arrays.

### Extension Tools
Under Redis's 'shenma:tools:' directory, various tool definitions are stored.

A tool definition, or ToolCall, is a JSON object containing the following fields:

```json
{
  "name": "translate_zh_en",
  "module": "translator",
  "type": "restful", //grpc, mcp
  "restful": {
    "url": "/translate/zh/en",
    "method": "GET" //PUT,DELETE,POST,GET
  },
  "description": "Translates Chinese comments and strings in input code to English",
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
        "description": "Code text"
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
        "description": "Translated code text"
      }
    }
  }
}
```

The corresponding JSONSchema can be found in `jsonschema/tool.json`:

Description:

| JsonPath | Description |
|--------|-----|
|name | Tool name |
|module | Module the tool belongs to |
|type | Tool interface type, currently supports restful. Will support grpc, mcp in the future |
|restful.url| http url for restful api |
|restful.method| http method for restful api |
|description | Tool description |
|supports | Supported scenarios, currently supports chat, completion, codereview |
|parameters | Parameter list definition for the tool |
|returns | Return value list definition for the tool |

## Architecture Design

| Directory/File | Responsibility | Description |
|------|------|----|
| api | Interfaces | Provides RESTful API interfaces |
| internal| Internal mechanisms | Implements common mechanism code, such as core logic |
| service | Services | Implements business logic required by interfaces |
| internal/extension| Extension management | Manages extensions, including loading and lookup |
| internal/template | Template management | Manages templates, including loading, lookup, rendering, and preparing data objects and function tables needed for rendering |
| internal/variable | Variable management | Manages shared variables, including traversal, loading, and building shared variable lookup tables |
| internal/tool| Tool management | Manages extension tools, including traversal, loading, and building tool lookup tables |
| internal/llm | LLM calls | |
| internal/redis | Redis | Implements Redis operations |
| internal/utils | Utility functions | Implements utility functions |
| internal/config | Configuration | Implements configuration loading |
| internal/logger | Logging | Implements logging |
| internal/http | HTTP client | Implements access logic for RESTful type extension tools |

## Processes and Mechanisms

### Template Syntax

The Prompt templates in Prompt-type extensions use the template syntax of Go's text/template library, supporting features including:

- Variable injection: `{{.variable}}`
- Conditional statements: `{{if .condition}}...{{end}`
- Loops: `{{range .items}}...{{end}`
- Scoping: `{{with .object}}...{{end}`
- Function calls: `{{codebase_look_ref "CreateObject"}}
- Calling other templates (equivalent to calling subroutines): `{{template "name" .}}`

AI-Prompt-Shell uses Go's text/template to instantiate templates. Before instantiation, it needs to build data objects and function lookup tables.

### Extension Loading

AI-Prompt-Shell loads all Prompt-type extensions from Redis, obtains the Prompt templates defined by these extensions, and caches them in the Prompt template lookup table.

AI-Prompt-Shell traverses all contents under Redis's 'shenma:templates:' directory.

Based on the traversed KEY (such as 'shenma:templates:{extension-name}:{prompt-name}'), it obtains the corresponding value, which is the Prompt extension containing the Prompt template definition.

The Prompt template lookup table is of type `map[string]PromptTemplate`, where the map key is the result of calling keyToJsonPath(KEY), 
and PromptTemplate is the Prompt template loaded from Redis.
### Building Data Objects

AI-Prompt-Shell acquires variables from various sources and constructs them into a data object called context according to the following rules, which is then provided to text/template for Prompt generation.

Variables available to AI-Prompt-Shell consist of two main parts: the 'args' field obtained from user requests, and shared variables read from Redis.

When rendering templates, AI-Prompt-Shell retrieves required variables from Redis and organizes them into JSON values according to the following rules, passing them to the template rendering engine as context variables.

All KEYs in Redis with the prefix 'shenma:environs:' have values that are JSON values (including strings, numbers, objects, arrays, etc.).

Other system modules uniformly write variables available for AI-Prompt-Shell rendering under Redis's 'shenma:environs:' for template rendering.

AI-Prompt-Shell traverses all KEYs prefixed with 'shenma:environs:' and reads their JSON values.

For each JSON value read, it obtains its KEY, calculates a path Path through KeyToJsonPath(KEY), and then saves the JSON value under the Path in the JSON object context.

The keyToJsonPath calculation method is: remove the header 'shenma:environs:' from the KEY, then replace the remaining ':' with '.' to get the JsonPath.

For example, Redis has the following key-value pairs:

1. KEY: 'shenma:environs:vscode:programming_language', VALUE: "go"
2. KEY: 'shenma:environs:vscode:frameworks', VALUE: ["gin", "gorm", "gin-swagger"]

The args field in the user request is as follows:

```json
"args": {
    "key1": "value1",
    "key2": "value2"
  }
```

Based on these two parts of content, the constructed JSON object context is as follows:

```json
{
  "vscode": {
    "programming_language": "go",
    "frameworks": ["gin","gorm", "gin-swagger"]
  },
  "args": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

The KEY format for variables registered by modules to AI-Prompt-Shell is typically: 'shenma:environs:{module}:{subject}'.

Currently, modules that will generate context information to Redis for AI-Prompt-Shell use include:

| Module | Key Prefix | Description |
|------|--------|-----|
| vscode | 'shenma:environs:vscode:*' | VSCode editing context |
| codebase | 'shenma:environs:codebase:*' | File lists, directory structures and other information from code indexing |
| chat | 'shenma:environs:chat:*' | Conversation history |
| completion | 'shenma:environs:completion:*' | Completion history |
| model | 'shenma:environs:model:*' | Provides model lists, model feature descriptions and other information |
| modules | 'shenma:environs:modules:*' | Metadata of modules providing information to AI-Prompt-Shell |

### Building Function Lookup Tables

AI-Prompt-Shell retrieves tool definitions from Redis and builds them into tool lookup tables (`type ToolRegistry=map[string]Tool`). 

When rendering Prompt templates using text/template, it builds template function tables Funcs based on the tool lookup tables, serving as external functions available during Prompt generation.

All KEYs in Redis with the prefix 'shenma:tools' have values that are tool types. AI-Prompt-Shell periodically traverses Redis's 'shenma:tools:' directory to refresh the tool lookup tables.

The key name toolName in the tool lookup table is the result of calculating keyToFunc(KEY).

The logic of keyToFunc(KEY):

- Remove the prefix 'shenma:tools:' from KEY
- Replace the remaining ':' with '_' in the path (e.g. 'codebase:lookup_ref' → 'codebase_lookup_ref')

Examples of keyToFunc(KEY) calculations:

- keyToFunc('shenma:tools:codebase:lookup_ref') = 'codebase_lookup_ref'
- keyToFunc('shenma:tools:codebase:caller') = 'codebase_caller'
- keyToFunc('shenma:tools:mcp:chrome:xx') = 'mcp_chrome_xx'

When rendering Prompt templates, it needs to traverse the tool lookup table and build the template function table Funcs, where the Funcs key names match the ToolRegistry key names (i.e. toolName), and the values are the function objects returned by `TemplateFunc(toolName)`.

### Template Rendering Error Handling

1. Template syntax errors:
   - Check during template loading, directly return 400 error if any syntax errors exist

2. Undefined variables:
   - Check if all referenced variables exist in context, return 400 error if any variables are missing

3. Tool call failures:
   - Log error but do not interrupt rendering 
   - Use default or empty values as substitutes in templates

4. Rendering timeout:
   - Built-in 500ms timeout control
   - Return 503 Service Unavailable error after timeout

### LLM Call Retry Strategy

1. After first failure, wait 100ms for first retry
2. After second failure, wait 300ms for second retry
3. Maximum of 2 retries
4. If all retries fail, return 502 Bad Gateway error

### Performance Optimization Strategies

1. Template caching:
   - Automatically cache the most recently used 100 templates
   - Evict less commonly used templates with LRU algorithm

2. Variable caching:
   - Locally cache frequently accessed Redis variables for 5 seconds

3. Batch tool calls:
   - Use goroutines to concurrently execute unrelated tool calls

4. Memory pool:
   - Reuse rendering result buffers
   - Reduce GC pressure