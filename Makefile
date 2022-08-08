MODULE_DIR=user-service

PROTO_DIR=proto
STUBS_DIR=pb

REPO_DIR=repo
DBMOCK_DIR=mock

stubs:
	rm -rf $(STUBS_DIR)
	mkdir $(STUBS_DIR)
	protoc -I $(PROTO_DIR)\
	       	--go_out=$(STUBS_DIR)\
	       	--go_opt=paths=source_relative\
	       	--go-grpc_out=$(STUBS_DIR)\
	       	--go-grpc_opt=paths=source_relative\
		$(PROTO_DIR)/*.proto

mockdb:
	rm -rf $(REPO_DIR)/$(DBMOCK_DIR)
	mkdir $(REPO_DIR)/$(DBMOCK_DIR)
	mockgen -package mockdb\
		-destination $(REPO_DIR)/$(DBMOCK_DIR)/db.go\
	       $(MODULE_DIR)/$(REPO_DIR) DB

.PHONY: stubs mockdb
       	

