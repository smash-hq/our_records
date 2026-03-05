# Our Records - 记了么

> 日常生活记录项目 - 支持图片、文字、语音、视频等信息的记录与管理

## 🌟 功能特性

- ✅ 图文混排记录
- ✅ 时间轴展示
- ✅ 多图片上传（最多 9 张）
- ✅ 文件 MinIO 对象存储
- ✅ PostgreSQL 数据库
- ✅ 移动端响应式适配
- ✅ 图片预览
- ✅ 标签分类
- ✅ 登录验证
- ✅ 软删除

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.25
- **框架**: Gin
- **数据库**: PostgreSQL + GORM
- **对象存储**: MinIO
- **配置**: Viper (YAML)

### 前端
- **框架**: Vue 3 (Composition API)
- **UI 组件**: Element Plus
- **路由**: Vue Router 4
- **HTTP**: Axios
- **构建工具**: Vite

## 📁 项目结构

```
our_records/
├── frontend/                # 前端项目
│   ├── src/
│   │   ├── api/             # API 接口
│   │   ├── views/           # 页面组件
│   │   ├── router/          # 路由配置
│   │   └── utils/           # 工具函数
│   ├── public/              # 静态资源
│   ├── package.json
│   └── vite.config.js
├── internal/                # 后端内部包
│   ├── config/              # 配置管理
│   ├── handlers/            # 请求处理
│   ├── middleware/          # 中间件
│   ├── models/              # 数据模型
│   └── routes/              # 路由配置
├── pkg/                     # 后端公共包
│   └── minio/               # MinIO 客户端
├── application.yaml         # 配置文件
├── go.mod                   # Go 模块定义
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

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd our_records
```

### 2. 安装前端依赖

```bash
cd frontend
npm install
```

### 3. 准备数据库

```sql
-- 创建数据库
CREATE DATABASE our_records;
```

### 4. 配置应用

编辑 `application.yaml`:

```yaml
server:
  port: "8088"

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

**注意：** 请根据你的实际环境修改数据库和 MinIO 配置。

### 5. 启动后端

```bash
cd ..
go mod tidy
go run main.go
```

后端服务：http://localhost:8088

### 6. 启动前端

```bash
cd frontend
npm run dev
```

前端开发：http://localhost:3000

## 📦 打包部署

### 使用 Makefile (推荐)

```bash
# 查看所有命令
make help

# 构建所有
make build

# 构建 Linux 后端
make build-linux

# 构建前端
make build-frontend

# 打包 Linux 版本
make package

# 清理
make clean
```

### 手动打包

#### 1. 安装前端依赖

```bash
cd frontend
npm install
```

#### 2. 编译后端

```bash
cd ..

# Linux
GOOS=linux GOARCH=amd64 go build -o build/linux/our_records main.go
```

#### 3. 构建前端

```bash
cd frontend
npm run build
```

#### 4. 复制文件

```bash
# 复制前端到构建目录
cp -r dist/* ../build/linux/dist/

# 复制配置文件
cp ../application.yaml ../build/linux/
```

### 打包结果

```
build/linux/
├── our_records          # Linux 可执行文件 (24MB)
├── application.yaml     # 配置文件
├── dist/                # 前端静态文件 (1.6MB)
│   ├── index.html
│   ├── logo.svg
│   └── assets/
├── deploy.sh            # 部署脚本
└── stop.sh              # 停止脚本
```

## 📝 部署到 Linux 服务器

### 1. 上传文件

```bash
scp -r build/linux/* root@your-server:/opt/our_records/
```

### 2. 执行部署

```bash
ssh root@your-server
cd /opt/our_records
chmod +x deploy.sh stop.sh our_records
sudo ./deploy.sh
```

### 3. 手动启动

```bash
# 启动
./our_records &

# 停止
./stop.sh
# 或
pkill -f our_records

# 查看日志
tail -f our_records.log
```

## 🔌 API 接口

### 记录管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/records` | 创建记录 |
| GET | `/api/records` | 获取记录列表 |
| GET | `/api/records/:id` | 获取单条记录 |
| PUT | `/api/records/:id` | 更新记录 |
| DELETE | `/api/records/:id` | 删除记录 |

### 文件上传

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/upload` | 上传文件到 MinIO |

## 📊 数据库设计

### Record 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| type | varchar(20) | 类型 (text/image/audio/video) |
| title | varchar(255) | 标题 |
| content | text | 内容 |
| media_path | text | 媒体路径 (MinIO 相对路径) |
| tags | varchar(500) | 标签 |
| deleted_at | timestamp | 软删除时间 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

## 🔐 登录验证

可用用户名：
- `huangqi` - 男主角
- `zhongyanling` - 女主角

登录信息存储在 `localStorage`：
- `isLoggedIn` - 登录状态
- `lastLoginUser` - 上次登录用户
- `lastLoginTime` - 上次登录时间

## 🌐 端口说明

| 服务 | 端口 | 说明 |
|------|------|------|
| 后端 API | 8088 | Go 后端（含前端静态文件） |
| MinIO | 9000 | 对象存储 API |
| MinIO 控制台 | 9001 | Web 管理界面 |
| PostgreSQL | 5432 | 数据库 |

## 📝 存储策略

所有文件（图片、音频、视频）统一存储到 MinIO：

| 文件类型 | MinIO 路径 | 访问方式 |
|---------|-----------|---------|
| 图片 | `images/` | 签名 URL（24 小时有效） |
| 音频 | `audios/` | 签名 URL（24 小时有效） |
| 视频 | `videos/` | 签名 URL（24 小时有效） |

**说明：**
- 数据库只存储相对路径
- 获取记录时自动生成签名 URL
- 签名 URL 有效期 24 小时

## 🔧 常用命令

### 后端

```bash
# 开发模式
go run main.go

# 编译
go build -o our_records main.go

# 清理
go clean
```

### 前端

```bash
# 安装依赖
cd frontend && npm install

# 开发模式
cd frontend && npm run dev

# 构建
cd frontend && npm run build
```

## 🐛 故障排查

### 后端启动失败

```bash
# 查看日志
cat our_records.log

# 检查端口占用
lsof -i :8088
```

### 数据库连接失败

```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 创建数据库
sudo -u postgres psql -c "CREATE DATABASE our_records;"
```

### MinIO 连接失败

```bash
# 检查 MinIO 服务
docker ps | grep minio

# 查看日志
docker logs minio
```

### 前端依赖问题

```bash
# 删除 node_modules
rm -rf frontend/node_modules

# 重新安装
cd frontend && npm install
```

## 📄 License

MIT

## 👨‍💻 开发

如有问题或建议，欢迎提交 Issue 或 Pull Request。
