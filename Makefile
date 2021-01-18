

.PHONY: run_to_do_example
run_to_do_example: generate
	cd examples/to-do && go run .


.PHONY: ./cave-js/bundle.js
./cave-js/bundle.js:
	cd cave-js && make ./bundle.js

.PHONY: generate
generate: ./cave-js/bundle.js ./cmd/include-bundle/main.go
	cd cmd/include-bundle && go run .


demos_build:
	cd cmd/cave-demos && docker build -f ./Dockerfile ../../ -t maxmcd/cave:demos

demos_push: demos_build
	docker push maxmcd/cave:demos

demos_deploy: demos_push
	cd cmd/cave-demos && fly deploy -v
