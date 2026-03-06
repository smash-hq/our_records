# Our Records - 记了么

> 日常生活记录项目 - 支持图片、文字、语音、视频等信息的记录与管理

## 🌟 功能特性

- ✅ 图文混排记录
- ✅ 时间轴展示
- ✅ 多图片上传（最多 9 张）
- ✅ 移动端响应式适配
- ✅ 图片预览
- ✅ 标签分类
- ✅ 用户注册/登录（JWT 认证）
- ✅ 密码修改
- ✅ 个人中心
- ✅ 可见性控制（公开/私密）
- ✅ 群组功能
- ✅ 评论互动

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.25
- **框架**: Gin
- **数据库**: PostgreSQL + GORM
- **对象存储**: MinIO

### 前端
- **框架**: Vue 3 (Composition API)
- **UI 组件**: Element Plus
- **构建工具**: Vite

## 📁 项目结构

```
our_records/
├── frontend/                # 前端项目
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── vite.config.js
├── internal/                # 后端内部包
├── pkg/                     # 后端公共包
├── application.yaml         # 配置文件
├── main.go                  # 后端入口
├── Makefile                 # 构建脚本
└── README.md                # 项目说明
```

## 🚀 快速开始

### 环境要求

- Go 1.25+
- Node.js 18+
- PostgreSQL 12+
- MinIO

### 1. 配置应用

编辑 `application.yaml`:

```yaml
server:
  port: "8088"
  mode: "debug"
  jwt_secret: "your-secret-key-change-in-production"  # 生产环境请修改

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  dbname: "our_records"
  sslmode: "disable"

minio:
  endpoint: "localhost:9000"
  access_key_id: "minioadmin"
  secret_access_key: "minioadmin"
  bucket: "our-records"
  use_ssl: false
```

**注意：** 请根据实际环境修改配置，**生产环境务必修改 `jwt_secret`**。

### 2. 安装依赖

**后端：**
```bash
go mod download
```

**前端：**
```bash
cd frontend
npm install
```

### 3. 启动服务

**启动后端：**
```bash
go run main.go
```

**启动前端（开发模式）：**
```bash
cd frontend
npm run dev
```

访问 http://localhost:3000

## 📦 打包部署

### 方式一：Windows 一键打包（推荐）

在 Windows 环境下，使用 PowerShell 脚本一键打包：

```powershell
# 执行打包脚本
.\build.ps1
```

脚本会自动完成以下操作：
1. 构建前端到 `frontend/dist/`
2. 交叉编译后端为 Linux amd64 可执行文件
3. 复制配置文件 `application.yaml`
4. 生成启动脚本 `start.sh` 和 `stop.sh`

打包完成后，文件位于 `build/windows/` 目录。

### 方式二：使用 Makefile（Linux/macOS）

```bash
# 查看帮助
make help

# 构建所有（后端 + 前端）
make build

# 仅构建 Linux 后端
make build-linux

# 构建前端
make build-frontend

# 打包 Linux 版本（包含后端、前端静态文件和配置文件）
make package

# 清理构建文件
make clean
```

### 方式三：手动构建

#### 1. 构建前端

```bash
cd frontend
npm run build
```

构建完成后，静态文件将输出到 `frontend/dist/` 目录。

#### 2. 交叉编译后端

**Windows (PowerShell)：**
```powershell
$env:GOOS='linux'
$env:GOARCH='amd64'
$env:CGO_ENABLED='0'
go build -o build/linux/our_records main.go
```

**Linux/macOS：**
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/linux/our_records main.go
```

### 部署到 Linux 服务器

#### 1. 上传文件

将打包目录（`build/windows/` 或 `build/linux/`）下的所有文件上传到服务器：

```bash
# 假设上传到 /opt/our_records
scp -r build/windows/* root@your_server:/opt/our_records/
```

目录结构：
```
/opt/our_records/
├── our_records        # 后端可执行文件
├── application.yaml   # 配置文件
├── dist/              # 前端静态文件
│   ├── index.html
│   ├── assets/
│   └── ...
├── start.sh           # 启动脚本
└── stop.sh            # 停止脚本
```

#### 2. 配置文件权限

```bash
# 设置执行权限
chmod +x /opt/our_records/start.sh
chmod +x /opt/our_records/stop.sh
chmod +x /opt/our_records/our_records

# 设置配置文件权限
chmod 644 /opt/our_records/application.yaml
```

#### 3. 启动服务

使用启动脚本：
```bash
cd /opt/our_records
./start.sh
```

或手动启动：
```bash
cd /opt/our_records
./our_records
```

#### 4. 停止服务

使用停止脚本：
```bash
./stop.sh
```

或手动停止：
```bash
# 查找进程
ps aux | grep our_records

# 杀死进程
kill <PID>
```

#### 5. 查看日志

```bash
# 查看运行日志
tail -f our_records.log
```

#### 6. 配置 Nginx 反向代理（可选）

创建 Nginx 配置文件 `/etc/nginx/sites-available/our-records`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/our_records/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8088/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/our-records /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 📄 License

MIT

## 👤 Author

smash_hq
