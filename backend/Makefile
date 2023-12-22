.PHONY: create-keypair

# Makefiles will not work on Windows OS unless the Developer has downloaded GCC (GNU C Compiler), or some other Open-Source tool
# Windows does not natively support Makefiles >:(

# Database migrations directory path and other variables needed for DB migrations.
PWD = $(shell pwd)
SERVPATH = $(PWD)
MPATH = $(PWD)/migrations
PORT = 5432

# Default number of migrations executed when running up or down
N = 1

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(SERVPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(SERVPATH)/rsa_private_$(ENV).pem -pubout -out $(SERVPATH)/rsa_public_$(ENV).pem

# We use the Go golang-migrate package to run these commands | https://github.com/golang-migrate/migrate
# golang-migrate must be installed before being able to run these cmds: go install github.com/golang-migrate/migrate
# review the golang-migrate docs for help installing this on linux, MacOS, or Windows.
# Run these cmds in a terminal like so: "make migrate-up"

# "make migrate-create"
migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)

# "make migrate-up"
migrate-up:
	migrate -source file://$(MPATH) -database postgres://postgres:DukeRedrix44@localhost:$(PORT)/colony-office?sslmode=disable up $(N)

# "make migrate-down"
migrate-down:
	migrate -source file://$(MPATH) -database postgres://postgres:DukeRedrix44@localhost:$(PORT)/colony-office?sslmode=disable down $(N)

# "make migrate-force"
migrate-force:
	migrate -source file://$(MPATH) -database postgres://postgres:DukeRedrix44@localhost:$(PORT)/colony-office?sslmode=disable force $(VERSION)

# "make migrate-drop"
migrate-drop:
	migrate -source file://$(MPATH) -database postgres://postgres:DukeRedrix44@localhost:$(PORT)/colony-office?sslmode=disable drop

