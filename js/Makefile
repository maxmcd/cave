



./src/index.js: ./src/index.ts
	npx tsc

./test/test.js: ./test/test.ts ./src/index.js

.PHONY: test
test: ./src/index.js ./test/test.js
	npm run test
