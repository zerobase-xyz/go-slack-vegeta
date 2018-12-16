AWS_ACCOUNT := 858026153602
AWS_REGION := ap-northeast-1
SQS_QUEUE_NAME := vegeta-queue
SQS_URL := https://sqs.$(AWS_REGION).amazonaws.com/$(AWS_ACCOUNT)/$(SQS_QUEUE_NAME)
TEST_MESSAGE := 'hoge hoge "https://zerobase.xyz" GET 100 60'

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/slack

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose

.PHONY: test
test: test
	aws sqs send-message --queue-url $(SQS_URL) --message-body $(TEST_MESSAGE) --delay-seconds 10 --message-attributes file://testdata/sqs_attributes.json
