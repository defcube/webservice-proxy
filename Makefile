export devel=0

all: installgobindata generatestatic
	go install github.com/defcube/webservice-proxy

dev: devel = 1
dev: all

generatestatic:
	$(MAKE) -C server all

installgobindata:
	go install github.com/jteeuwen/go-bindata/go-bindata
