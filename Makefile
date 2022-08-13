MODULE_DIR=user-service

PROTO_DIR=proto
STUBS_DIR=internal/pb

REPO_DIR=internal/repo
DBMOCK_DIR=$(REPO_DIR)/mock

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
	rm -rf $(DBMOCK_DIR)
	mkdir $(DBMOCK_DIR)
	mockgen -package mockdb\
		-destination $(DBMOCK_DIR)/db.go\
	       $(MODULE_DIR)/$(REPO_DIR) DB

.PHONY: stubs mockdb
       	

