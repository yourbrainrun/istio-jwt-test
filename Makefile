.PHONY: cert

VERSION=v0.0.14
PROJECTNAME=istio-jwt-test

cert:
	openssl genpkey -algorithm RSA -out private-key.pem
    #openssl rsa -pubout -in private-key.pem -out public-key.pem

build:
	docker build --platform linux/amd64 -t $(PROJECTNAME):$(VERSION) .

clean:
	docker rmi $(PROJECTNAME):$(VERSION)