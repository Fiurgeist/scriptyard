.PHONY: backend-run backend-test backend-fmt backend-vet node-deps frontend-test frontend-watch frontend-build frontend-clean

ROOT = $(abspath $(CURDIR))
BACKEND_ROOT = $(ROOT)/backend
FRONTEND_ROOT = $(ROOT)/frontend
PUBLIC_ROOT = $(ROOT)/public

GO = cd $(BACKEND_ROOT) && go
YARN = yarn --cwd $(FRONTEND_ROOT)

backend-run:
	$(GO) run -race ./cmd/server/main.go

backend-test:
	$(GO) test -v -race ./...

backend-fmt:
	$(GO) fmt ./...

backend-vet:
	$(GO) vet ./...

node-deps:
	$(YARN) install

frontend-test: node-deps
	$(YARN) run test:dev-watch

frontend-watch: node-deps
	$(YARN) run watch-develop

frontend-build: node-deps
	$(YARN) run build

frontend-clean:
	rm -rf $(FRONTEND_ROOT)/node_modules
	rm -rf $(FRONTEND_ROOT)/coverage
	rm -rf $(PUBLIC_ROOT)/static/js/*.js.map

help:
	@echo "  backend-run              to run the server"
	@echo "  backend-test             to test the backend code"
	@echo "  backend-fmt              to run go fmt"
	@echo "  backend-vet              to run go vet"
	@echo "  node-deps                to install frontend dependencies"
	@echo "  frontend-test            to test the frontend code"
	@echo "  frontend-watch           to build dev frontend and hot-reload changes"
	@echo "  frontend-build           to build prod frontend"
	@echo "  frontend-clean           to remove build dependencies"
