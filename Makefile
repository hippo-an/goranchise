.PHONY: run
run: ready wait air
	@echo "Successfully run the goranchise application"


.PHONY: install-tailwindcss
install-tailwindcss:
	@sh < ./scripts/install_tailwindcss.sh
	@./bin/twc init --full

.PHONY: css-build
css-build:
	@./bin/twc -i ./static/css/app.css -o ./public/style.css

.PHONY: air
air:
	@air -c .air.toml

# ex) make gen-model MODEL=Color
.PHONY: gen-model
gen-model:
	@go run -mod=mod entgo.io/ent/cmd/ent new $(MODEL)

.PHONY: ent
ent:
	@go generate ./ent

.PHONY: ready
ready:
	@docker compose up -d

.PHONY: down
down:
	@docker compose down -v

.PHONY: ent-install
ent-install:
	@go get entgo.io/ent/cmd/ent

wait:
	@echo "Waiting for database to be ready..."; \
	for i in $$(seq 1 10); do \
		docker container exec goranchise-db-1 pg_isready -U admin -d goranchise; \
		if [ $$? -eq 0 ]; then \
			echo "Database is ready!"; \
			break; \
		fi; \
		echo "Database is not ready yet. Waiting..."; \
		sleep 5; \
	done; \
	if [ $$i -eq 10 ]; then \
		echo "Failed to start the database"; \
		exit 1; \
	fi; \
	echo "Database goranchise-db successfully started!"; \
