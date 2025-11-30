run:
	templ generate && go run .

build:
	templ generate
	rsync -r ../weblights lights:/home/erikrodabaugh/

clean:
	go clean
	rm -f weblights

install: install-bin install-service enable-start

install-bin:
	@echo "-- Creating /usr/local/bin/weblights"
	sudo mkdir -p /usr/local/bin/weblights
	@echo "-- Copying .env and weblights binary"
	sudo cp .env weblights /usr/local/bin/weblights/

install-service:
	@echo "-- Installing weblights.service"
	sudo cp weblights.service /etc/systemd/system/
	@echo "-- Reloading systemd"
	sudo systemctl daemon-reload

enable-start: install-service
	@echo "-- Enabling weblights.service"
	sudo systemctl enable weblights.service
	@echo "-- Starting weblights.service"
	sudo systemctl start weblights.service
