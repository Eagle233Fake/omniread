innospark-run:
	@echo "初始化 innospark 数据库和缓存服务"
	docker-compose --project-name innospark -f ./docker-compose.yml up -d innospark-mongodb innospark-redis
	@echo "初始化完成"

innospark-clean:
	@echo "删除 innospark 数据库和缓存服务"
	docker-compose --project-name innospark -f ./docker-compose.yml down --remove-orphans
	@echo "删除完成"

psych-run:
	@echo "初始化 psych 数据库和缓存服务 (仅限 psych-mongodb 和 psych-redis)"
	docker-compose --project-name psych -f ./docker-compose.yml up -d psych-mongodb psych-redis
	@echo "初始化完成"

psych-clean:
	@echo "删除 psych 数据库和缓存服务"
	docker-compose --project-name psych -f ./docker-compose.yml down --remove-orphans
	@echo "删除完成"

meowpick-run:
	@echo "初始化 meowpick 数据库和缓存服务"
	docker-compose --project-name meowpick -f ./docker-compose.yml up -d meowpick-mongodb meowpick-redis
	@echo "初始化完成"

meowpick-clean:
	@echo "删除 meowpick 数据库和缓存服务"
	docker-compose --project-name meowpick -f ./docker-compose.yml down --remove-orphans
	@echo "删除完成"

synapse-run:
	@echo "启动 Synapse MySQL 服务"
	docker compose --project-name synapse -f ./docker-compose.yml up -d synapse-mysql synapse-redis
	@echo "启动完成"

synapse-clean:
	@echo "删除 Synapse MySQL 服务"
	docker compose --project-name synapse -f ./docker-compose.yml down --remove-orphans
	@echo "删除完成"

omniread-run:
	@echo "初始化 omniread 数据库和缓存服务"
	docker-compose --project-name omniread -f ./docker-compose.yml up -d omniread-mongodb omniread-redis
	@echo "初始化完成"

omniread-clean:
	@echo "删除 omniread 数据库和缓存服务"
	docker-compose --project-name omniread -f ./docker-compose.yml down --remove-orphans
	@echo "删除完成"
