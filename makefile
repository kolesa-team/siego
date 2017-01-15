export GOPATH=$(CURDIR)/.go

APP_NAME = siego
DEBIAN_TMP = $(CURDIR)/deb
VERSION = `$(CURDIR)/out/$(APP_NAME) -v | cut -d ' ' -f 3`

$(CURDIR)/out/$(APP_NAME): $(CURDIR)/src/main.go
	go build -o $(CURDIR)/out/$(APP_NAME) $(CURDIR)/src/main.go

dep-install:
	go get github.com/codegangsta/cli
	
install:
	cp $(CURDIR)/out/$(APP_NAME) /usr/local/bin/$(APP_NAME)

fmt:
	gofmt -s=true -w $(CURDIR)/src

run:
	go run $(CURDIR)/src/main.go ${ARGS}
	
run-dev:
	go run --race $(CURDIR)/src/main.go ${ARGS}

keys:
	openssl genrsa -out $(CURDIR)/out/rsakey 512
	openssl rsa -in $(CURDIR)/out/rsakey -pubout > $(CURDIR)/out/rsakey.pub

strip: $(CURDIR)/out/$(APP_NAME)
	strip $(CURDIR)/out/$(APP_NAME)

deb: $(CURDIR)/out/$(APP_NAME)
	mkdir $(DEBIAN_TMP)
	mkdir -p $(DEBIAN_TMP)/usr/local/bin
	install -m 755 $(CURDIR)/out/$(APP_NAME) $(DEBIAN_TMP)/usr/local/bin
	fpm -n $(APP_NAME) \
		-v $(VERSION) \
		-t deb \
		-s dir \
		-C $(DEBIAN_TMP) \
		.
	rm -fr $(DEBIAN_TMP)

clean:
	rm -f $(CURDIR)/out/*

clean-deb:
	rm -fr $(DEBIAN_TMP)
	rm -f $(CURDIR)/*.deb

debug:
	echo $(GOPATH)
