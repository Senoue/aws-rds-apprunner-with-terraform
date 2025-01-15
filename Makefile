# 変数定義
REGION=us-west-2
ACCOUNT_ID=<you_account_id>
ECR_REPO=${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/example/app
CMD_CONTAINER=cmd-container

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
	docker run --rm -v ./:/app ${CMD_CONTAINER} init -g ./main.go

go-fmt:
	docker run --rm -v ./:/app ${CMD_CONTAINER} fmt
