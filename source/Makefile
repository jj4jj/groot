SHELL := /bin/bash
BASEDIR = $(shell pwd)
RUNDIR = $(BASEDIR)/../run


server_crt=${RUNDIR}/keys/server.crt
server_key=${RUNDIR}/keys/server.key
rsa_prk=${RUNDIR}/keys/rsa_pri.pem
rsa_pbk=${RUNDIR}/keys/rsa_pub.pem
sym_key=${RUNDIR}/keys/sym.key

# share  proto

#pyhon check
PYTHON=~/.grootdev/venv/bin/python


# server proto
pssgenpy:=$(shell find ./proto -type d | awk '{print $$1"/*.py"}')
ssgenpy=$(wildcard $(pssgenpy))
ssgenpb=$(patsubst %.py,%.proto,$(ssgenpy))


# csproto proto
pcsgenpy:=$(shell find ./groot/proto -type d | awk '{print $$1"/*.py"}')
csgenpy=$(wildcard $(pcsgenpy))
csgenpb=$(patsubst %.py,%.proto,$(csgenpy))
csgengo=$(patsubst %.py,%.pb.go,$(csgenpy))
csgengo+=$(BASEDIR)/groot/proto/cserr/cserr.errno.go



# server proto
pssgenpb:=$(shell find proto -type d | awk '{print $$1"/*.proto"}')
pssgenpb+=$(shell find service -type d | awk '{print $$1"/*.proto"}')
pssgenpb+=$(shell find sfw -type d | awk '{print $$1"/*.proto"}')
ssmanpb:=$(wildcard $(pssgenpb))


ssgengo:=$(patsubst %.proto,%.pb.go,$(ssgenpb))
ssgengo+=$(patsubst %.proto,%.pb.go,$(ssmanpb))



build:
	cd cmd && go build -v -x -work -ldflags ${ldflags}
	mv cmd/cmd ${RUNDIR}/bin/groot

all: clean gen build

key: ca rsakey symkey

mkdirs:
	mkdir -p ${RUNDIR}/logs ${RUNDIR}/db ${RUNDIR}/bin ${RUNDIR}/conf ${RUNDIR}/keys
	mkdir -p ${RUNDIR}/upload

redis:
	echo "redis server starting ..."
	nohup redis-server &

initdev:
	cd ${BASEDIR}/tools && sh initenv.sh

initenv: initdev mkdirs key pbdc gen redis build genconf
	[ -f ${RUNDIR}/conf/config.json ] || mv ${RUNDIR}/conf/config.json.default ${RUNDIR}/conf/config.json

run: build
	cd ${RUNDIR}/bin && ./groot --config=../conf/config.json

syncdb:
	cd ${RUNDIR}/bin && ./groot --syncdb --config=../conf/config.json

genconf:
	cd ${RUNDIR}/bin && ./groot --genconf --config=../conf/config.json

#avoiding make auto delete proto files autogen
.PRECIOUS: %.proto

%.proto: %.py
	$(PYTHON) $< --outdir=$(dir $@) --codegen=all

%.pb.go: %.proto
	protoc --proto_path=${GOPATH}/src:$(dir $<):. --go_out=paths=source_relative:$(dir $<) $<


#$(info gen proto go path: $(csgenpb))

gen: $(csgenpb) $(csgengo) $(ssgengo)
	$(info generating cs:$(csgengo))
	$(info generating ss:$(ssgengo))


pbdc:
	cd groot/tools && bash initenv.sh && cd ../.. && make clean && make gen


clean:
	rm -f groot
	rm -f $(csgengo) $(csgenpb) $(ssgengo) $(ssgenpb)

fmt:
	gofmt -w .

ca:
	openssl req -new -nodes -x509 -out ${server_crt} -keyout ${server_key} -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=123@123.com"

rsakey:
	openssl genrsa -out ${rsa_prk} 2048
	openssl rsa -in ${rsa_prk} -pubout -out ${rsa_pbk}
	#convert to pkcs8
	openssl pkcs8 -topk8 -inform PEM -in ${rsa_prk} -outform PEM -nocrypt > ${rsa_prk}.tmp
	rm -f ${rsa_prk}
	mv ${rsa_prk}.tmp ${rsa_prk}


symkey:
	openssl rand -base64 16 > ${sym_key}

test:
	cd tests && go test -v -run GetCode account_test.go


help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"

.PHONY: clean gotool ca help



