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

### 2. API 接口说明

#### 认证相关

| 方法 | 路径 | 说明 | 需要认证 |
|------|------|------|----------|
| POST | `/api/auth/register` | 用户注册 | ❌ |
| POST | `/api/auth/login` | 用户登录 | ❌ |
| GET | `/api/user` | 获取当前用户信息 | ✅ |
| PUT | `/api/user/password` | 修改密码 | ✅ |

#### 记录相关

| 方法 | 路径 | 说明 | 需要认证 |
|------|------|------|----------|
| POST | `/api/records` | 创建记录 | ✅ |
| GET | `/api/records` | 获取记录列表 | ✅ |
| GET | `/api/records/:id` | 获取单条记录 | ✅ |
| PUT | `/api/records/:id` | 更新记录 | ✅ |
| DELETE | `/api/records/:id` | 删除记录 | ✅ |
| POST | `/api/upload` | 文件上传 | ✅ |

#### 请求示例

**注册：**
```bash
curl -X POST http://localhost:8088/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456","email":"admin@example.com"}'
```

**登录：**
```bash
curl -X POST http://localhost:8088/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

**访问受保护接口（携带 Token）：**
```bash
curl -X GET http://localhost:8088/api/user \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. 安装前端依赖

```bash
cd frontend
npm install
```

### 4. 启动服务

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

### 5. Windows 构建脚本

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

### 6. 部署到服务器

1. 上传 `build/linux/` 目录到服务器
2. 执行部署脚本：
```bash
chmod +x deploy.sh stop.sh our_records
sudo ./deploy.sh
```

## 📄 License

MIT
