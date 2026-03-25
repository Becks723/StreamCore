
# 环境变量
DIR=$(shell pwd)
CMD_PATH=$(DIR)/cmd
CONFIG_PATH=$(DIR)/config
IDL_PATH=$(DIR)/idl
OUTPUT_PATH=$(DIR)/output
API_PATH=$(DIR)/cmd/api
# 自动获取 go.mod 的 module 名
MODULE := $(shell grep "^module" go.mod | awk '{print $$2}')

# 服务名
SERVICES=api user video interaction social chat group
service=$(word 1, $@)

PREFIX="[MAKEFILE]"

# 构建指定的服务并运行，如果仅构建需要加上 BUILD_ONLY 参数
.PHONY: $(SERVICES)
$(SERVICES):
	  @echo "$(PREFIX) Building $(service)..."
	  @mkdir -p output
	  @bash $(DIR)/docker/script/build.sh $(service)
	  @echo "$(PREFIX) Build $(service) completed"
ifndef BUILD_ONLY
	  @echo "$(PREFIX) Starting $(service)..."
	  @export SERVICE=$(service); bash ./docker/script/entrypoint.sh
endif

# 启动必要的环境，如 mysql和etcd 等
.PHONY: env-up
env-up:
	@docker compose -f ./docker/docker-compose.yml up -d

# 关闭环境
.PHONY: env-down
env-down:
	@cd ./docker && docker compose down

# 生成指定微服务的 kitex 脚手架代码
.PHONY: kitex-gen-%
kitex-gen-%:
	@service_name=$*; \
	echo "$(PREFIX) Generating kitex code for $$service_name..."; \
	mkdir -p cmd/$$service_name; \
	cd cmd/$$service_name && kitex -module $(MODULE) -service $$service_name \
		-gen-path ../../kitex_gen ../../idl/$$service_name.thrift; \
	echo "$(PREFIX) kitex generate completed for $$service_name"

# 更新api网关的 hertz 脚手架代码
.PHONY: hz-update-api
hz-update-api:
	@hz update --idl idl/api.thrift

# 运行单元测试
# -gcflags="all=-N -l" 表示禁用内联优化，否则 monkey patch 函数替换可能不成功
.PHONY: test
test:
	@echo "$(PREFIX) Running go test..."; \
	go test -gcflags="all=-N -l" -v ./...; \