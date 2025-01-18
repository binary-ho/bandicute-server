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
	export SOPS_AGE_KEY_FILE=./agekey.txt
	sops --decrypt $(ENCRYPTED_FILE) > $(DECRYPTED_FILE)
	@echo "Decrypt Done."

.PHONY: encrypt decrypt
