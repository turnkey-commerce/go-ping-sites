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
SRC_PATH=$(GOPATH)$(PATHSEP)src$(PATHSEP)github.com$(PATHSEP)turnkey-commerce$(PATHSEP)go-ping-sites

default:
	go install -ldflags "-X main.version=0.3.0" -v
	-mkdir $(DIST_PATH)
	cp $(GOPATH)$(PATHSEP)bin$(PATHSEP)$(EXE_NAME) $(DIST_PATH)$(PATHSEP)$(EXE_NAME)
	cp -r $(SRC_PATH)$(PATHSEP)templates $(DIST_PATH)
	cp -r $(SRC_PATH)$(PATHSEP)public $(DIST_PATH)

clean:
	-rm -rf $(DIST_PATH)

run: default
	cd $(DIST_PATH)$(CMDSEP)go-ping-sites

distribute: clean default
	cp $(SRC_PATH)$(PATHSEP)config$(PATHSEP)config.toml $(DIST_PATH)$(PATHSEP)config.toml
	cp $(SRC_PATH)$(PATHSEP)database$(PATHSEP)db-seed.toml $(DIST_PATH)$(PATHSEP)db-seed.toml
