run:
	templ generate && go run .

build:
	~/go/bin/templ generate && go build .

clean:
	go clean
	rm -f weblights

install: install-bin install-service enable-start

install-bin:
	@echo "-- Stopping weblights.service (if running)"
	sudo systemctl stop weblights.service || true
	@echo "-- Creating /usr/local/bin/weblights"
	sudo mkdir -p /usr/local/bin/weblights
	@echo "-- Copying .env and weblights binary"
	sudo cp -f .env weblights /usr/local/bin/weblights/
	@echo "-- Setting ownership to root:root (optional)"
	sudo chown -R root:root /usr/local/bin/weblights

install-service:
	@echo "-- Installing weblights.service"
	cp weblights.service /etc/systemd/system/
	@echo "-- Reloading systemd"
	systemctl daemon-reload

enable-start: install-service
	@echo "-- Enabling weblights.service"
	systemctl enable weblights.service
	@echo "-- Starting weblights.service"
	systemctl start weblights.service
