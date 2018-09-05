NAME         =chainstack
MAIN_PATH    =server/
MAIN_FILE    =server/main.go
BIN_DIR      =bin

default: 
	@echo "USAGE: make <command>"
	@echo ""
	@echo "    dev: Running with gin hot-reload"
	@echo "    run: Execute build file"
	@echo ""
	@echo "    build: Build new bin file"
	@echo ""
	
dev: 
	gin --bin $(BIN_DIR)/$(NAME) --path . --build $(MAIN_PATH) -i run $(MAIN_FILE)

install:
	@echo "Installing..."
	@go install $(MAIN_FILE)

run: 
	$(BIN_DIR)/$(NAME)

build:
	@echo "STEP: BUILD"
	@echo "   1. create dir: $(BIN_DIR)" \
		&& mkdir -p $(BIN_DIR)\
		&& echo "   ==> ok"
	@echo "   2. build: $(MAIN_FILE)" \
		&& go build -o $(BIN_DIR)/$(NAME) $(MAIN_FILE) \
		&& echo "   ==> ok: SERVICE=$(BIN_DIR)/$(NAME)"
