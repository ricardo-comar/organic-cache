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
	aws --endpoint-url=http://localhost:4566 dynamodb batch-write-item --request-items file://localstack/dynamodb_user_discounts.json
	aws --endpoint-url=http://localhost:4566 dynamodb batch-write-item --request-items file://localstack/dynamodb_products.json

scan:
	aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name $(table)

update-lambda:

	rm bin/lambda_$(lambda)*
	echo "Updating lambda_$(lambda)" 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/lambda_$(lambda) lambda_$(lambda)/handler/handler.go 
	zip -j bin/lambda_$(lambda).zip bin/lambda_$(lambda)
	aws --endpoint-url http://localhost:4566 lambda update-function-code --function-name $(lambda) --zip-file fileb://bin/lambda_$(lambda).zip --output text


subscribe:
	curl -i -X PUT http://localhost:4566/restapis/$(shell aws --endpoint-url=http://localhost:4566 apigateway get-rest-apis | jq -r '.items[0].id')/v1/\_user_request_/subscribe \
	-H "Content-Type: application/json" \
   -d '{"user_id": "$(ID)"}' 

quotation:
	curl -s -X POST http://localhost:4566/restapis/$(shell aws --endpoint-url=http://localhost:4566 apigateway get-rest-apis | jq -r '.items[0].id')/v1/\_user_request_/quotation \
	-H "Content-Type: application/json" \
   -d '{"user_id": "$(ID)", "products": [ {"id": "P01", "qtd": 19}, {"id": "P02", "qtd": 30}, {"id": "P05", "qtd": 10} ] }' | jq
