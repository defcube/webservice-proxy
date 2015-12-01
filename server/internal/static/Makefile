all:
ifeq ($(devel),1)
	go-bindata-assetfs -pkg static -ignore bindata\\.go -debug --nometadata .
else
	go-bindata-assetfs -pkg static -ignore bindata\\.go .
endif

