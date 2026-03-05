# Our Records Makefile

.PHONY: all build build-linux build-frontend package clean help

# 项目配置
PROJECT := our_records
BUILD_DIR := build
FRONTEND_DIR := frontend

# 默认目标
all: build

# 帮助
help:
	@echo "Our Records Makefile (Linux Only)"
	@echo ""
	@echo "Targets:"
	@echo "  build          - Build all (backend + frontend)"
	@echo "  build-linux    - Build Linux backend"
	@echo "  build-frontend - Build frontend"
	@echo "  package        - Package Linux version"
	@echo "  clean          - Clean build files"
	@echo ""

# 构建所有
build: build-linux build-frontend
	@echo "Build complete!"

# 构建 Linux 后端
build-linux:
	@echo "Building Linux backend..."
	@mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(PROJECT) main.go

# 构建前端
build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm run build

# 打包 Linux
package: build-linux
	@echo "Packaging Linux version..."
	@mkdir -p $(BUILD_DIR)/linux/dist
	cp -r $(FRONTEND_DIR)/dist/* $(BUILD_DIR)/linux/dist/
	cp application.yaml $(BUILD_DIR)/linux/
	@echo ""
	@echo "========================================"
	@echo "  Package Complete!"
	@echo "========================================"
	@echo ""
	@echo "Linux package: $(BUILD_DIR)/linux/"
	@echo ""
	@echo "Files:"
	@echo "  - our_records       (executable)"
	@echo "  - application.yaml  (config)"
	@echo "  - dist/             (static files)"
	@echo "  - deploy.sh         (deploy script)"
	@echo "  - stop.sh           (stop script)"
	@echo ""

# 清理
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f $(PROJECT)
	@echo "Clean complete!"
