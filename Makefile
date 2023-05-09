.PHONY: deploy build run compress

deploy: build compress

build:
	go build -o goip.exe goip

run:
	go run goip

compress:
	upx -9 goip.exe