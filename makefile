# Public repository 

all : libstool.so	

clean :
	@rm -f libstool.so
	@rm -f ./bin/*

libstool.so : *.go
	go build -buildmode=plugin -o libstool.so

install : 
	@if [ ! -d ./bin ];then mkdir ./bin;fi 
	@mv ./libstool.so ./bin/libstool.so	

push :
	@git push origin --all