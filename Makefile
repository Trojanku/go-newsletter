PHONY: build cover deploy start test test-integration

export image := `aws lightsail get-container-images --service-name goo | jq -r '.containerImages[0].image'`

deploy:
	aws lightsail push-container-image --service-name goo --label app --image goo
	jq <containers.json ".app.image=\"$(image)\"" >containers2.json
	mv containers2.json containers.json
	aws lightsail create-container-service-deployment --service-name goo \
		--containers file://containers.json \
		--public-endpoint '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}'
		
build:
	docker build -t goo .

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...
