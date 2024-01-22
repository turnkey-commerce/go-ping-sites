ifdef ComSpec
	PATHSEP2=\\
	EXE_NAME=go-ping-sites.exe
	CMDSEP=&
else
	PATHSEP2=/
	EXE_NAME=go-ping-sites
	CMDSEP=;
endif
PATHSEP=$(strip $(PATHSEP2))
DIST_PATH=$(GOPATH)$(PATHSEP)dist$(PATHSEP)go-ping-sites
SRC_PATH=.

default:
	$(eval export GO15VENDOREXPERIMENT = 1)
	go install -ldflags "-X main.version=1.3.6" -v
	-mkdir -p $(DIST_PATH)
	cp $(GOPATH)$(PATHSEP)bin$(PATHSEP)$(EXE_NAME) $(DIST_PATH)$(PATHSEP)$(EXE_NAME)

clean:
	-rm -rf $(DIST_PATH)

run: default
	cd $(DIST_PATH)$(CMDSEP)go-ping-sites

distribute: clean default
	cp $(SRC_PATH)$(PATHSEP)config$(PATHSEP)config.toml $(DIST_PATH)$(PATHSEP)config_sample.toml
	cp $(SRC_PATH)$(PATHSEP)database$(PATHSEP)db-seed.toml $(DIST_PATH)$(PATHSEP)db-seed_sample.toml
