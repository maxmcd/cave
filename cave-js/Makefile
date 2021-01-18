# SHELL=/usr/bin/env bash



./src/index.js: ./src/*.ts ./test/*.ts ./Makefile
	npx tsc

.PHONY: test
test: ./src/index.js ./test/test.js
	npm run test

bundle.js: ./src/index.js ./index.ts
	npx tsc index.ts
	npx browserify index.js | npx uglifyjs --source-map --output bundle.js
	sed -i '1s;^;//# sourceMappingURL=bundle.js.map\n;' bundle.js
.PHONY: clean
clean:
	rm ./*.js{,.map} || true
	rm ./{src,test}/*.js{,.map} || true
