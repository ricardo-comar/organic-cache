build: compile

clean:
		@echo "Cleaning bin folder"
		@rm -rf bin

prepare:
		@echo "Creating bin folder"
		@mkdir -p bin

test: clean prepare		
		@echo "Testing lambdas"
		@find . -maxdepth 1 -mindepth 1 -type d  -name 'lambda_*'| while read dir; do\
			echo "Testing $$dir"; cd $$dir; go test -coverprofile=../bin/$$dir-coverage.out ./...; cd ..; \
		done

compile: clean prepare		
		@echo "Compiling lambdas"
		@find . -maxdepth 1 -mindepth 1 -type d -name 'lambda_*'| while read dir; do\
			echo "Compiling $$dir"; cd $$dir; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../bin/$$dir handler/handler.go; cd ..; \
		done

package: build 
		@find . -maxdepth 1 -mindepth 1 -type d  -name 'lambda_*'| while read dir; do \
			zip -j bin/$$dir.zip bin/$$dir; \
		done

localstack:
	aws --endpoint http://localhost:4566 iam create-user --user-name test
	
load:
	aws --endpoint-url=http://localhost:4566 dynamodb batch-write-item --request-items file://dynamodb_user_discounts.json
	aws --endpoint-url=http://localhost:4566 dynamodb batch-write-item --request-items file://dynamodb_products.json

subscribe:
	curl -i -X PUT http://localhost:4566/restapis/$(shell aws --endpoint-url=http://localhost:4566 apigateway get-rest-apis | jq -r '.items[0].id')/v1/\_user_request_/subscribe \
	-H "Content-Type: application/json" \
   -d '{"id": "$(ID)"}' 

quotation:
	curl -s -X POST http://localhost:4566/restapis/$(shell aws --endpoint-url=http://localhost:4566 apigateway get-rest-apis | jq -r '.items[0].id')/v1/\_user_request_/quotation \
	-H "Content-Type: application/json" \
   -d '{"id": "$(ID)", "products": [ {"id": "P01", "qtd": 19}, {"id": "P02", "qtd": 30}, {"id": "P05", "qtd": 10} ] }' | jq
