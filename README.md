# Our Records - 记了么

> 日常生活记录项目 - 支持图片、文字、语音、视频等信息的记录与管理

## 🌟 功能特性

- ✅ 图文混排记录
- ✅ 时间轴展示
- ✅ 多图片上传（最多 9 张）
- ✅ 移动端响应式适配
- ✅ 图片预览
- ✅ 标签分类
- ✅ 登录验证

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
├── build.bat                # Windows 构建脚本
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

**注意：** 请根据实际环境修改配置。

### 2. 安装前端依赖

```bash
cd frontend
npm install
```

### 3. 启动服务

**启动后端：**
```bash
cd ..
go run main.go
```

**启动前端：**
```bash
cd frontend
npm run dev
```

访问 http://localhost:3000

## 📦 打包部署

### Windows 构建脚本

```cmd
# 查看所有命令
build.bat help

# 打包 Linux 版本
build.bat package

# 构建所有
build.bat build

# 清理
build.bat clean
```

### 部署到服务器

1. 上传 `build/linux/` 目录到服务器
2. 执行部署脚本：
```bash
chmod +x deploy.sh stop.sh our_records
sudo ./deploy.sh
```

## 📄 License

MIT
