# Makefile to run Python scripts with customizable arguments

# Default values for arguments
E_ARG ?= ee
V_ARG ?= 5
JSON_DIR = json

# Phony target to handle all tasks
.PHONY: all clean
ifeq ($(V_ARG),4)
all: $(JSON_DIR) overview.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json
else
all: $(JSON_DIR) overview.json authentication.json authorization.json client-events.json messages.json rule-engine-rate.json rule-engine-count.json
endif

# Create json directory if it doesn't exist
$(JSON_DIR):
	mkdir -p $(JSON_DIR)

# Rule for generating overview.json with customizable arguments
overview.json:
	./overview.py -e $(E_ARG) -v $(V_ARG) > $(JSON_DIR)/overview-$(E_ARG)-$(V_ARG).json

authentication.json:
	./authentication.py > $(JSON_DIR)/authentication.json

authorization.json:
	./authorization.py > $(JSON_DIR)/authorization.json

client-events.json:
	./client-events.py -v $(V_ARG) > $(JSON_DIR)/client-events-$(V_ARG).json

messages.json:
	./messages.py > $(JSON_DIR)/messages.json

rule-engine-rate.json:
	./rule-engine-rate.py > $(JSON_DIR)/rule-engine-rate.json

rule-engine-count.json:
	./rule-engine-count.py > $(JSON_DIR)/rule-engine-count.json

# Clean up by removing the json directory
clean:
	rm -rf $(JSON_DIR)

