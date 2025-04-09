# 假的 LSP IT 管理机器人

没写完的实验性项目，不推荐使用

## 安装

1. 克隆仓库：
   ```bash
   git clone https://github.com/real-LiHua/fakelspitmgrbot
   cd fakelspitmgrbot
   ```

2. 在项目根目录创建 `.env` 文件，添加以下变量：
   ```
   TOKEN=<你的 Telegram 机器人 Token>
   URL=<你的 Web 应用地址>
   WEBHOOK_SECRET=<你的 Webhook 密钥>
   CHAT_ID=<你的聊天 ID>
   NAMESPACE=<用SSH验证签名的命名空间，`ssh-keygen` 命令 `-n` 参数的值>
   LISTEN_ADDR=<服务器监听地址，默认0.0.0.0:8080>
   ```

3. 运行机器人：
   ```bash
   go run .
   ```

## 使用方法

### 命令

| 命令            | 描述                      |
|-----------------|---------------------------|
| `/dl`           | 下载资源                  |
| `/dl_debug`     | 下载资源（调试版）        |
| `/ban`          | 根据 Telegram ID 封禁用户 |
| `/unban`        | 根据 Telegram ID 解封用户 |
| `/ban_github`   | 根据 GitHub ID 封禁用户   |
| `/unban_github` | 根据 GitHub ID 解封用户   |

### Web 接口

- `/` - Web 应用首页。
- `/validate` - 验证用户凭据
- `/submit` - 提交用户验证数据

## 数据库

机器人使用 SQLite 存储用户数据，数据库结构包括：

| 字段名            | 描述                    |
|-------------------|-------------------------|
| `telegram_id`     | Telegram 用户 ID (主键) |
| `github_id`       | GitHub 用户 ID          |
| `challenge_code`  | 用户验证的挑战码        |
| `github_username` | GitHub 用户名           |
| `flag`            | 用户状态（如封禁、成员）|

# TODO
- [x] 命令监听  
- [ ] 日志频道  
- [ ] 处理入群申请  
- [ ] 风控逻辑，写完也不会公开  
- [ ] 封禁  
- [ ] 解封  
- [ ] 重打包逻辑  
