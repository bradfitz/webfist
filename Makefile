ssh:
	ssh -i ~/keys/webfist.pem ubuntu@webfist.org

runprod:
	cd $(GOPATH)/src/github.com/bradfitz/webfist/webfistd
	go build
	sudo ./webfistd  -web=80 -smtp=25
