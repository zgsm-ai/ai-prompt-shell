# 编程规范

## 编程语言

使用GO语言进行开发,要求支持go 1.21版本。

## 框架

采用如下框架：

| 框架名 | 说明 |
|---|---|
| gin | 轻量级web框架 |
| gorm | 轻量级ORM框架 |
| github.com/go-redis/redis/v8 | 键值数据库,可用作缓存及消息队列 |
| logrus |  日志 |
| spf13/viper | 配置文件 |
| spf13/cobra | 命令行工具 |
| text/template | 文本模板 |
| gin-swagger | swagger文档生成 |

## 注释

所有注释内容，要求使用英文进行描述。

RESTful API实现函数，采用swagger注释标准进行注释，保证能生成所有API的swagger文档。

范例：

```go
// RenderPrompt 渲染Prompt模板
// @Summary 渲染Prompt模板
// @Description 根据Prompt ID获取模板，使用输入变量渲染生成最终Prompt
// @Tags Render
// @Accept json
// @Produce json
// @Param prompt_id path string true "Prompt模板ID"
// @Param variables body string false "模板变量" SchemaExample({"variables":{"text":"单例模式实现"}})
// @Success 200 {object} map[string]interface{} "渲染结果"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "模板不存在"
// @Failure 500 {object} map[string]interface{} "渲染失败"
// @Router /api/render/prompts/{prompt_id} [post]
func (pc *PromptController) RenderPrompt(c *gin.Context) {
}
```

其它函数，采用jsdoc风格给每个函数进行注释，说明函数的功能，参数，返回值，用法，注意事项等。

范例:

```go
/**
 * 上传文件到服务器
 * @param {string} serverPath - 文件在服务器上的目标存储路径
 * @param {*resource.AI_File} file - 要上传的文件对象，包含文件大小(Size)等元数据
 * @returns {error} 返回错误对象，成功时返回nil
 * @description
 * - 自动处理文件大小：小文件直接上传，大文件(>DEF_PART_SIZE)调用PostHugeFile
 * - 设置HTTP请求头部：Content-Type、Cookie和Accept
 * - 处理服务器响应并更新上传进度条
 * @throws
 * - 文件流转失败(createFileBuffer错误)
 * - POST请求创建失败(http.NewRequest错误)
 * - HTTP请求发送错误(client.Do错误)
 * - 服务器返回非200状态码(statusToError)
 * @example
 * err := session.PostFile("/upload/path", file)
 * if err != nil {
 *     log.Fatal(err)
 * }
 */
func (ss *AI_Session) PostFile(serverPath string, file *resource.AI_File) error {
...
}
```

## 单元测试

每个函数都需要编写单元测试，单元测试需覆盖函数主体逻辑。

## 支持swagger文档

需要支持swagger文档。

## 指标监控

需要支持prometheus监控，向prometheus报告必要的指标数据。

## 结构

采用分层设计方式，从上到下分为四层：

- 交互层: 与用户进行交互的逻辑
- 业务层：处理用户的业务逻辑，也可以称为策略层。
- 机制层：业务逻辑依赖的基础机制，公共机制。
- 数据层/IO层：负责处理数据，调用数据库的方法，返回结果给机制层。

各层逻辑之间分离到不同的代码文件中，必要时保存在不同目录。

其中RESTful API接口实现部分放在controllers目录下；
业务层放在services目录下；
机制层放在internal目录下；
数据IO层放在dao目录下；

