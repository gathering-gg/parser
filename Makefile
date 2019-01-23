build:
	go build -ldflags \
	  "-X 'gitlab.com/gathering-gg/gathering.root=http://localhost:8600' -X 'gitlab.com/gathering-gg/gathering.version=0.0.3'" \
	  -o gathering \
	  ./cli

test:
	go test .

cov:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html
