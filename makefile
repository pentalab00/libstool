all : libstool.so	

clean :
	@rm -f libstool.so

libstool.so : *.go
	go build -buildmode=plugin -o libstool.so


