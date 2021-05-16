.PHONY: protos clean infrastructure

PROTOS_DIR := ./
QUALIFIED_DIR = $(PROTOS_DIR)github.com/dalmarcogd/mobstore

TARGETS := products users discounts

protos:
	-cd protos && \
	protoc -I=$(PROTOS_DIR) --go_out=./  domains.proto && \
	protoc -I=$(PROTOS_DIR) --go-grpc_out=./ discounts.proto && \
	mv $(QUALIFIED_DIR)/products/internal/domains/domainsgrpc/domains.pb.go ../products/internal/domains/domainsgrpc/domains.pb.go && \
	mv $(QUALIFIED_DIR)/products/internal/discounts/discountsgrpc/discounts_grpc.pb.go ../products/internal/discounts/discountsgrpc/discounts_grpc.pb.go && \
	mockgen -source=../products/internal/discounts/discountsgrpc/discounts_grpc.pb.go -destination=../products/internal/discounts/discountsgrpc/discounts_mock.pb.go -package=discountsgrpc && \
	rm -rf $(PROTOS_DIR)github.com

clean:
	@echo "\nRemoving localstack, mysql"
	@docker-compose down -v | true
	sudo rm -rf .localstack | true
	sudo rm -rf infrastructure/terraform/local/products/.terraform* infrastructure/terraform/local/products/*tfstate* | true
	sudo rm -rf infrastructure/terraform/local/discounts/.terraform* infrastructure/terraform/local/discounts/*tfstate*  | true
	sudo rm -rf infrastructure/terraform/local/users/.terraform* infrastructure/terraform/local/users/*tfstate* | true

terraform_apply:
	for target in $(TARGETS) ; do \
  		echo "terraform apply ($$target)" && \
			cd infrastructure/terraform/local/$$target/ && \
			terraform init && \
			terraform destroy -auto-approve && \
			terraform plan && \
			terraform apply -auto-approve && \
			cd ../../../../ ; \
	done

database_migration:
	@echo "\nMigrating database products"
	cd products && go build -a -tags netgo -o products_migration cmd/migration/main.go && ./products_migration && rm products_migration
	cd users && go build -a -tags netgo -o users_migration cmd/migration/main.go && ./users_migration && rm users_migration
    #cd discounts && go build -a -tags netgo -o discounts_migration cmd/migration/main.go && ./discounts_migration && rm discounts_migration


infrastructure: clean
	@echo "\nStarting localstack container and creating AWS local resources"
	@docker-compose up -d --build --force-recreate
	@echo "\nWaiting until localstack be ready"
	@until docker inspect --format='{{json .State.Health}}' localstack | grep -o healthy; do sleep 1; done
	@echo "\nCreating AWS resources locally"
	$(MAKE) terraform_apply
	$(MAKE) database_migration
