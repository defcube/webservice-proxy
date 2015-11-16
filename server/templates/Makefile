all:
ifeq ($(devel),1)
	go-bindata -pkg templates -ignore bindata\\.go -debug --nometadata .
else
	go-bindata -pkg templates -ignore bindata\\.go .
endif

