SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

DECRYPTED_FILE := config/property.yml
ENCRYPTED_FILE := config/encrypted-property.yml

encrypt:
	@echo "Encrypt $(DECRYPTED_FILE) -> $(ENCRYPTED_FILE)"
	sops --encrypt $(DECRYPTED_FILE) > $(ENCRYPTED_FILE)
	@echo "Encrypt Done."

decrypt:
	@echo "Decrypt $(ENCRYPTED_FILE) -> $(DECRYPTED_FILE)"
	export SOPS_PRIVATE_KEY_FILE=./agekey.txt
	sops --decrypt $(ENCRYPTED_FILE) > $(DECRYPTED_FILE)
	@echo "Decrypt Done."

SOPS_VERSION := 3.9.3
sops:
	@echo "Installing SOPS version $(SOPS_VERSION)..."
	wget "https://github.com/getsops/sops/releases/download/v$(SOPS_VERSION)/sops-v$(SOPS_VERSION).linux.amd64" -O sops
	chmod +x sops
	sudo mv sops /usr/local/bin/sops
	sops --version
	@echo "SOPS installed."

.PHONY: encrypt decrypt sops
