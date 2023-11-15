# Makefile to run Python scripts with customizable arguments

# Default values for arguments
# EDITION_ARG can be ee or ce
EDITION_ARG ?= ee
# VERSION_ARG can be 4 or 5
VERSION_ARG ?= 5
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
JSON_DIR = $(PROJECT_DIR)/json

# Phony target to handle all tasks
.PHONY: all clean
ifeq ($(VERSION_ARG),4)
all: $(JSON_DIR) overview.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json
else
all: $(JSON_DIR) overview.json authentication.json authorization.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json
endif

# Create json directory if it doesn't exist
$(JSON_DIR):
	mkdir -p $(JSON_DIR)

# Rule for generating overview.json with customizable arguments
overview.json:
	./overview.py -e $(EDITION_ARG) -v $(VERSION_ARG) > $(JSON_DIR)/overview-$(EDITION_ARG)-$(VERSION_ARG).json

authentication.json:
	./authentication.py > $(JSON_DIR)/authentication.json

authorization.json:
	./authorization.py > $(JSON_DIR)/authorization.json

client-events.json:
	./client-events.py -v $(VERSION_ARG) > $(JSON_DIR)/client-events-$(VERSION_ARG).json

messages.json:
	./messages.py > $(JSON_DIR)/messages.json

rule-engine-rate.json:
	./rule-engine-rate.py > $(JSON_DIR)/rule-engine-rate.json

rule-engine-count.json:
	./rule-engine-count.py > $(JSON_DIR)/rule-engine-count.json

# Clean up by removing the json directory
clean:
	rm -rf $(JSON_DIR)

