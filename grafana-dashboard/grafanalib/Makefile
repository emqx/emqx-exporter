# Makefile to run Python scripts with customizable arguments
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

# Phony target to handle all tasks
.PHONY: all
all: emqx_4 emqx_ee_4 emqx_5 emqx_ee_5

.PHONY: emqx_4
emqx_4: EDITION_ARG=ce
emqx_4: VERSION_ARG=4
emqx_4: JSON_DIR = $(PROJECT_DIR)/../template/emqx-4
emqx_4: $(JSON_DIR) overview.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json

.PHONY: emqx_ee_4
emqx_ee_4: EDITION_ARG=ee
emqx_ee_4: VERSION_ARG=4
emqx_ee_4: JSON_DIR = $(PROJECT_DIR)/../template/emqx-enterprise-4
emqx_ee_4: $(JSON_DIR) overview.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json

.PHONY: emqx_5
emqx_5: EDITION_ARG=ce
emqx_5: VERSION_ARG=5
emqx_5: JSON_DIR = $(PROJECT_DIR)/../template/emqx-5
emqx_5: $(JSON_DIR) overview.json authentication.json authorization.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json

.PHONY: emqx_ee_5
emqx_ee_5: EDITION_ARG=ee
emqx_ee_5: VERSION_ARG=5
emqx_ee_5: JSON_DIR = $(PROJECT_DIR)/../template/emqx-enterprise-5
emqx_ee_5: $(JSON_DIR) overview.json authentication.json authorization.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json

# Create json directory if it doesn't exist
$(JSON_DIR):
	mkdir -p $(JSON_DIR)

# Rule for generating overview.json with customizable arguments
overview.json:
	./overview.py -e $(EDITION_ARG) -v $(VERSION_ARG) > $(JSON_DIR)/overview.json

authentication.json:
	./authentication.py > $(JSON_DIR)/authentication.json

authorization.json:
	./authorization.py > $(JSON_DIR)/authorization.json

client-events.json:
	./client-events.py -v $(VERSION_ARG) > $(JSON_DIR)/client-events.json

messages.json:
	./messages.py > $(JSON_DIR)/messages.json

rule-engine-rate.json:
	./rule-engine-rate.py > $(JSON_DIR)/rule-engine-rate.json

rule-engine-count.json:
	./rule-engine-count.py > $(JSON_DIR)/rule-engine-count.json
