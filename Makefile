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

.PHONY: run
run: ready air
	echo "Successfully run the goranchise application"

.PHONY: ready
ready:
	@docker compose up -d

.PHONY: down
down:
	@docker compose down -v