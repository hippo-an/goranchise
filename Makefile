
.PHONY: install-tailwindcss
install-tailwindcss:
	@sh < ./scripts/install_tailwindcss.sh
	@./bin/twc init --full


.PHONY: css-build
css-build:
	@./bin/twc -i ./static/css/app.css -o ./public/style.css