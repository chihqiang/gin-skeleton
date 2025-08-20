# 版本号变量
# 使用git命令获取当前标签或提交哈希作为版本号
version := $(shell git describe --tags --always)
# 项目Makefile文件
# 用于自动化构建、测试和代码检查等任务

# 变量定义
OUTPUT := skeleton  # 编译输出的可执行文件名
MAIN := main.go         # 主程序入口文件

# 检查目标
# 用于运行代码检查工具和格式化代码
check:
	@echo "🔍 Running linters..."
	# 检查并安装golangci-lint (代码静态分析工具)
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4; \
	else \
		echo "golangci-lint already installed, skipping..."; \
	fi
	golangci-lint run ./...  # 运行代码检查
	@echo "✅ Linting passed"

	# 检查并安装errcheck (检查未处理的错误)
	@if ! command -v errcheck >/dev/null 2>&1; then \
		echo "Installing errcheck..."; \
		go install github.com/kisielk/errcheck@latest; \
	else \
		echo "errcheck already installed, skipping..."; \
	fi
	errcheck ./...  # 检查代码中的未处理错误
	@echo "✅ Error checks passed"

	@find . -name "*.go" -exec go fmt {} \;  # 格式化所有Go文件
	@go mod tidy  # 整理Go模块依赖

# 测试目标
# 用于运行单元测试和检查代码状态
test:
	@echo "🧪 Running tests..."
	# 运行单元测试，带详细输出、覆盖率和竞争检测，超时2分钟
	go test -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.txt ./... -run . -timeout=2m

	@echo "🚀 Running benchmark tests (stress/performance)..."
	# 运行所有基准测试，基准运行时间5秒，单独报告
	go test -bench=. -benchtime=5s -run=^$$ ./...

	@echo "🔍 Checking git status..."
	# 检查工作目录是否有未提交的更改
	@git diff --quiet || (echo "❌ Uncommitted changes detected in working directory!" && git status && exit 1)
	# 检查暂存区是否有未提交的更改
	@git diff --cached --quiet || (echo "❌ Staged but uncommitted changes detected!" && git status && exit 1)
	@echo "✅ Git status clean"


# 构建目标
# 用于编译项目生成可执行文件
build:
	@echo "🔧 Building $(OUTPUT) with version $(version)..."
	GO111MODULE=on go build -ldflags "-s -w -X main.version=$(version)" -o $(OUTPUT) $(MAIN)
	@echo "✅ Build complete: $(OUTPUT)"

