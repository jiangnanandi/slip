# Slip - 简单的笔记应用

## 概述

Slip 即「小纸条系统」是一个使用 Go 和 Gin 框架构建的简单笔记发布系统。它通过提供一系列 API 允许用户创建、查看和管理笔记。该应用以 Markdown 格式存储笔记，并生成索引页面以便于导航。

## 开发状态

请注意，当前代码仍在开发中，某些功能尚未完善。我们欢迎社区的参与和反馈，希望能够共同维护和改进这个项目。如果你对项目有任何建议、问题或想要贡献代码，请随时提交拉取请求或提出问题。

## 特性

- 创建带有标题和正文的笔记。
- 在网页界面中查看笔记。
- 为所有笔记生成索引页面。
- 使用 Go 和 Gin 框架进行高效的网页处理。
- 使用 Markdown 格式进行笔记排版。

## 快速开始

### 前提条件

- Go 1.19 或更高版本
- Docker（可选，用于容器化部署）

### 安装

1. 克隆仓库：

   ```bash
   git clone https://github.com/yourusername/slip.git
   cd slip
   ```

2. 安装依赖：

   ```bash
   go mod download
   ```

3. 准备必需文件：

   a. 配置文件：
   在 `configs` 目录下，复制 `config.example.yaml` 文件并重命名为 `config.yaml`，然后根据你的实际情况修改内容：

   ```yaml
   keys:
       client_id: "your_client_id" # 客户端ID
       secret_key: "your_secret_key1" # 密钥 16字节
   ```

   b. 模板文件：
   在项目根目录下创建 `templates` 目录，并从示例中复制模板文件：
   ```bash
   mkdir -p templates
   cp templates.example/index.html.tmpl templates/index.html.tmpl
   ```

4. 构建应用程序：

   ```bash
   go build -o slip main.go
   ```

5. 运行应用程序：

   ```bash
   ./slip
   ```

   应用程序将在 `http://localhost:8084` 启动。

### Docker 部署

要使用 Docker 运行应用程序，可以构建并运行 Docker 容器：

1. 构建 Docker 镜像：

   ```bash
   docker build -t slip .
   ```

2. 运行 Docker 容器：

   ```bash
   docker run -p 8084:8084 slip
   ```

   访问应用程序：`http://localhost:8084`。

## 使用方法

- 登录获取 Token，发送一个 GET 请求到 `/login?encrypted_string=&client_id=`，获取到生成的 `tokenstr`。
  - 两个入参分别代表「密钥字符串」和「客户端ID」。
  - 每个客户端注册的时候提供一个「加密字符串」，服务端会提供一个 `clientId` ，此即为「客户端ID」
  - 客户端须用 `examples/auth.js` 代码中的「加密函数」生成「密钥字符串」。
- 要创建笔记，发送一个 POST 请求到 `/send-notes`，请求体为包含标题和正文的 JSON，请求中的 Header 包含 `Authorization: token tokenstr`

  示例：

  ```curl
  curl -X "POST" "http://127.0.0.1:8084/send-notes" \
       -H 'Authorization: token eyJhbGciOiJIUzI1NiIsInR5cC' \
       -H 'Content-Type: application/json; charset=utf-8' \
       -d $'{
    "title": "Slip 是什么",
    "Body": "Slip 是一个简单易用的笔记应用，允许用户快速创建、查看和管理笔记。它支持 Markdown 格式，提供直观的网页界面，用户可以轻松访问和组织自己的笔记。无论是工作记录还是生活感悟，Slip 都能帮助用户高效地整理思路。"
  }'
  ```

- 要查看笔记索引，导航到 `/index`。
- 要查看特定笔记，访问 `/notes/:title`，将 `:title` 替换为实际的笔记标题。

## 许可证

本项目根据 GNU 通用公共许可证 v3.0 进行许可。有关详细信息，请参见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎贡献！请随时提交拉取请求或提出任何建议或改进的问题。

## 联系方式

如有任何疑问，请联系：

- 姓名：xzsj.wang
- 邮箱：melody.wang1984@gmail.com
