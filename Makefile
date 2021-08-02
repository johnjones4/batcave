dependencies:
	mkdir -p vendor
	wget https://github.com/johnjones4/mqtt2kasa/archive/refs/heads/main1.zip
	unzip main1.zip
	mv mqtt2kasa-main1 vendor/mqtt2kasa
	rm main1.zip

build:
	docker-compose build

clean:
	rm -rf vendor
