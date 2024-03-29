APP = hearsay
OUT = release

PROXY ?= http://localhost:8181
PORT ?= 8080
DEST ?= https://example.com

GARBLE=${GOPATH}/bin/garble
BUILD=garble -tiny build

LD.windows=-ldflags "${STRIP} -X main.port=${PORT} -X main.dest=${DEST} -X main.proxy=${PROXY} -H windowsgui"
LD.linux=-ldflags "${STRIP} -X main.port=${PORT} -X main.dest=${DEST} -X main.proxy=${PROXY}"
LD.darwin=${LD.linux}

PLATFORMS=windows linux darwin
OS=$(word 1, $@)

all: ${PLATFORMS}

${PLATFORMS}: $(GARBLE)
	GOOS=${OS} ${BUILD} ${LD.${OS}} -o ${OUT}/${APP}_${OS}

release: all
	@tar caf ${APP}.tar.gz ${OUT}
	@rm -rf ${OUT}

clean:
	rm -rf ${OUT} ${APP}*

$(GARBLE):
	go install mvdan.cc/garble@latest
