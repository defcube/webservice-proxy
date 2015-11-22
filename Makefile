export devel=0
.DEFAULT_GOAL=install

publish: installgobindata generatestatic install

devmode: devel = 1
devmode: publish

entrdev: devmode
	find . -iname "*.go" -or -iname "*.html"  | entr -r make installAndRun

install:
	go install github.com/defcube/webservice-proxy

installAndRun: install
	-@pkill webservice-proxy
	webservice-proxy

generatestatic:
	$(MAKE) -C server all

installgobindata:
	go install github.com/jteeuwen/go-bindata/go-bindata
	go install github.com/elazarl/go-bindata-assetfs/...

