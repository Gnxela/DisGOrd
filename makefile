GOCMD := go
GOBUILD := $(GOCMD) build
GOBUILD_PLUGIN := $(GOBUILD) --buildmode=plugin -i

all: fmt bin/disGOrd modules

fmt:
	gofmt -w *.go common/*.go modules/*.go

bin/disGOrd: dirs disGOrd.go common/*
	$(GOBUILD) -o bin/disGOrd .

dirs:
	mkdir -p bin/modules
	mkdir -p config

MODULE_DEPS := bin/modules/command_avatar.so
bin/modules/command_avatar.so: modules/command_avatar.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)

bin/modules/command_call.so: modules/command_call.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_call.so

bin/modules/command_dad.so: modules/command_dad.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_dad.so

bin/modules/command_dota.so: modules/command_dota.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_dota.so

bin/modules/command_flip.so: modules/command_flip.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_flip.so

bin/modules/command_ping.so: modules/command_ping.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_ping.so

bin/modules/command_roll.so: modules/command_roll.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_roll.so

bin/modules/command_timer.so: modules/command_timer.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_timer.so

bin/modules/command_restrict.so: modules/command_restrict.go common/*
	$(GOBUILD_PLUGIN) -o $@ $(word 1,$^)
MODULE_DEPS += bin/modules/command_restrict.so

modules: dirs $(MODULE_DEPS)

clean:
	rm -rf bin/disGOrd bin/modules/*
