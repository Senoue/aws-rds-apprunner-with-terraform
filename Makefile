# 変数定義
include .env.development
export $(shell sed 's/=.*//' .env)

# コンテナ関連の操作
build-image:
	docker build --platform linux/amd64 -t ${ECR_REPO}:latest .

push-image:
	docker push ${ECR_REPO}:latest

auth-ecr:
	aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com

# command 関連の操作
generate-container:
	docker build -t ${CMD_CONTAINER} docker/cmd/.

generate-swagger:
	docker run --rm -v ./:/app ${MGT_CONTAINER} init -g ./main.go

# migrate 関連の操作
.PHONY: migrate-up migrate-down migrate-create
migrate-create:
	migrate create -ext sql -dir ./migrations -seq create_posts

migrate-up:
	migrate -path=./migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" up

migrate-down:
	migrate -path=./migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" down
