# コンテナ関連の操作
build-image:
	docker build --platform linux/amd64 -t 637423403983.dkr.ecr.us-west-2.amazonaws.com/example/app:latest .

push-image:
	docker push 637423403983.dkr.ecr.us-west-2.amazonaws.com/example/app:latest

auth-ecr:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 637423403983.dkr.ecr.us-west-2.amazonaws.com

# Swagger関連の操作
generate-conteiner:
	docker build -t cmd-container docker/cmd/.

generate-swagger:
	docker run --rm -v ./:/app cmd-container init -g ./main.go

go-fmt:
	docker run --rm -v ./:/app cmd-container fmt