MAIN_CLASS=imchatroomgateway.go
APP_NAME=gateway
CURRENT_PATH:=`pwd`

version := "1.0.0"
DATE := `date +%FT%T%z`
profile=test


.PHONY: help clean resources all windows linux darwin

help:
	@echo "usage: make <option>"
	@echo "options and effects:"
	@echo "    help   		: Show help"
	@echo "    clean  		: Clean the conf dictionary"
	@echo "    resource  	: Copy resource into conf"
	@echo "    all    		: Build multiple binary of this project"
	@echo "    build 		: Build the linux binary of this project"
	@echo "    windows		: Build the windows binary of this project"
	@echo "    linux  		: Build the linux binary of this project"
	@echo "    darwin 		: Build the darwin binary of this project"

all:clean resources windows linux darwin

build:clean resources linux

windows:
	@cd src && GOOS=windows go build -o ../bin/${APP_NAME}-windows.${version} ${MAINCLASS}
	@echo build  success


linux:
	@cd src && GOOS=linux go build -o ../bin/${APP_NAME}-linux.${version} ${MAINCLASS}

darwin:
	@cd src && GOOS=darwin go build -o ../bin/${APP_NAME}-darwin.${version} ${MAINCLASS}

clean:
	@echo clean ...
	@rm -rf ${CURRENT_PATH}/bin
	@rm -rf ${CURRENT_PATH}/conf
	@rm -rf ${CURRENT_PATH}/logs
	@mkdir ${CURRENT_PATH}/bin
	@mkdir ${CURRENT_PATH}/conf
	@mkdir ${CURRENT_PATH}/logs
	@echo clean ${APP_NAME} conf logs success


resources:
	@echo copy resources ...
	@cp ${CURRENT_PATH}/resource/conf-${profile}.json ${CURRENT_PATH}/conf/conf.json
	@echo copy resources \[conf-${profile}.json\] to ${CURRENT_PATH}/conf/conf.json success





