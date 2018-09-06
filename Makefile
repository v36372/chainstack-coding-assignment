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

deploy: backup build

run: 
	$(BIN_DIR)/$(NAME)

backup:
	@echo "STEP: BACKUP"
	@echo "   1. backup: binary file"
	@[ ! -f $(BIN_DIR)/$(NAME) ] && \
		 echo "   => skip: SERVICE=bin/$(NAME) DOES NOT EXIST" || ( \
		 cp $(BIN_DIR)/$(NAME) $(BIN_DIR)/$(NAME)_$(GIT_TAG) && \
	  	 echo "   => ok: SERVICE=$(BIN_DIR)/$(NAME)_$(GIT_TAG)" )
	@echo "   2. backup: git commit"
	@git tag -f $(GIT_TAG) && \
		echo  "   => ok: TAG=$(GIT_TAG)"


build:
	@echo "STEP: BUILD"
	@echo "   1. create dir: $(BIN_DIR)" \
		&& mkdir -p $(BIN_DIR)\
		&& echo "   ==> ok"
	@echo "   2. build: $(MAIN_FILE)" \
		&& go build -o $(BIN_DIR)/$(NAME) $(MAIN_FILE) \
		&& echo "   ==> ok: SERVICE=$(BIN_DIR)/$(NAME)"

clean:
	@echo "STEP: CLEAN"
	@echo "   1. remove dir: $(BIN_DIR)"
	@rm -rf bin \
	 	&& echo "   ==> ok"

restart: 
ifneq ("$(wildcard $(PID_API_FILE))","") 
	@$(MAKE) -f $(THIS_FILE) stop
	@$(MAKE) -f $(THIS_FILE) start
else
	@$(MAKE) -f $(THIS_FILE) start
endif

start:
	@echo "Starting..."
	@mkdir -p $(LOG_DIR)
ifneq ("$(wildcard $(PID_API_FILE))","") 
	@echo "[FAIL] A processing is running. Stop it first or restart"
else
	@pid= nohup bin/$(NAME) >> $(LOG_FILE) 2>&1 & echo "$$!" > $(PID_API_FILE)
	@echo "Done!"
endif

stop:
	@echo "Stopping..."
	@kill -USR2 $(PID_API)
	@rm $(PID_API_FILE)
	@echo "Done!"

status:
ifeq ("$(wildcard $(PID_API_FILE))","") 
	@echo "There is no running process"
else
	@echo "Service is running at PID=$(PID_API)"
	@echo "---"
	@ps aux | grep $(PID_API)
endif
