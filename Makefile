export devel=0

all: installgobindata generatestatic install

devmode: devel = 1
devmode: all

install:
	go install github.com/defcube/webservice-proxy

generatestatic:
	$(MAKE) -C server all

installgobindata:
	go install github.com/jteeuwen/go-bindata/go-bindata
