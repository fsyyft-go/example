# 编码输出的应用程序名称。
APP_OUTPUT=./bin
APP_LOGS=./logs

# 读取环境变量。
# 完整工作方式，示例：make MAKE_FULL=1。
EVN_FULL=${MAKE_FULL}
# 更新软件包，示例：make MAKE_GO_GET=1 bindata。
ENV_GOGET=${MAKE_GO_GET}
# 压缩包，示例：make MAKE_COMPRESSION=1
EVN_COMPRESSION=${MAKE_COMPRESSION}
# 打包，示例：make MAKE_PACK=1
ENV_PACK=${MAKE_PACK}
# 打开 CGO，开启后只针对非交叉编译的有效，示例：make MAKE_ENABLE_CGO=1。
ENV_ENABLE_CGO=${MAKE_ENABLE_CGO}
# 是否后台工作状态。
ENV_ENABLE_DAEMON=${MAKE_ENABLE_DAEMON}

ifeq ($(EVN_FULL),1)
	ENV_GOGET=1
	EVN_COMPRESSION=1
	ENV_PACK=1
endif

# 设置命令路径。
CMD_CAT=cat
CMD_CD=cd
CMD_CP=cp
CMD_CLEAR=clear
CMD_DATE=date
CMD_DIRNAME=dirname
CMD_GREP=grep
CMD_ECHO=echo
CMD_FILD=find
CMD_GO=go
CMD_GOLDS=golds
# 安装方式：https://golangci-lint.run/usage/install/#local-installation 。
CMD_GOLANGCI-LINT=golangci-lint
CMD_GOVULNCHECK=govulncheck
CMD_GO_ASSETS_BUILDER=go-assets-builder
CMD_GO_BINDATA=go-bindata
CMD_GO_CALLVIS=go-callvis
CMD_GIT=git
CMD_GMCHART=gmchart
CMD_LS=ls
CMD_MAKE=make
CMD_MKDIR=mkdir
CMD_MV=mv
CMD_PROTOC=protoc
CMD_PWD=pwd
CMD_RM=rm
CMD_SED=sed
CMD_SWAG=swag
CMD_TINYGO=tinygo
CMD_UPX=upx
CMD_UNZIP=unzip
CMD_ZIP=zip

.PHONY: test
test:
	# 准备测试所有用例。
	@$(CMD_MAKE) -f Makefile_Test test-all
	# 完成测试所有用例。
	# 准备单元测试后置工作。
	-@$(CMD_RM) -rf *.test
	@$(CMD_GO) mod tidy
	# 准备单元测试后置工作。

.PHONY: test-pre
test-pre: clean
	# 准备单元测试前置工作。
	# 准备对依赖包进行梳理。
	@$(CMD_GO) mod tidy
	# 完成对依赖包进行梳理。
	# 准备静态代码检查。
	@$(CMD_GOLANGCI-LINT) run --verbose ./...
	# 结束静态代码检查。
	# 准备代码安全检查。
	@${CMD_GOVULNCHECK} ./...
	# 完成代码安全检查。
	# 结束单元测试前置工作。

.PHONY: clean
clean:
	# 准备清理已编译文件。
ifeq ($(ENV_ENABLE_DAEMON),1)
	# 后台工作，无需清屏。
else
	@$(CMD_CLEAR)
endif
	-@$(CMD_RM) -rf $(APP_OUTPUT)
	-@$(CMD_RM) -rf $(APP_LOGS)
	@$(CMD_MKDIR) $(APP_OUTPUT)
	@$(CMD_MKDIR) $(APP_OUTPUT)/out
	@$(CMD_MKDIR) $(APP_OUTPUT)/test
	@$(CMD_MKDIR) $(APP_OUTPUT)/darwin_amd64
	@$(CMD_MKDIR) $(APP_OUTPUT)/linux_amd64
	@$(CMD_MKDIR) $(APP_OUTPUT)/linux_arm64
	@$(CMD_MKDIR) $(APP_OUTPUT)/windows_amd64
	@$(CMD_MKDIR) $(APP_LOGS)
	@$(CMD_GO) mod tidy
	# 完成清理已编译文件。

.PHONY: tidy
tidy:
	# 准备清理无效引用。
	@$(CMD_CAT) go.mod | $(CMD_GREP) -v "// indirect" > go.mod.tmp
	@$(CMD_RM) -rf go.mod
	@$(CMD_MV) go.mod.tmp go.mod
	@$(CMD_GO) mod tidy
	# 完成清理无效引用。