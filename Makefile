dependencies:
	mkdir -p vendor
	wget https://github.com/johnjones4/mqtt2kasa/archive/refs/heads/main1.zip
	unzip main1.zip
	mv mqtt2kasa-main1 vendor/mqtt2kasa
	rm main1.zip

backup-secrets:
	cp vendor/mqtt2kasa/data/config.yaml ${dest}/mqtt2kasa-config.yaml
	cp -R api/data ${dest}/
	cp api/.env ${dest}/api.env