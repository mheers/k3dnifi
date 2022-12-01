all: build

build:
	docker build -t mheers/k3dnifi:latest .

push:
	docker push mheers/k3dnifi:latest
