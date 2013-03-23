penguin-go
==========

A 'go' based version of the penguin Rest APIs.

Install 'go' from http://golang.org/doc/install
Make sure you have a GOPATH setup correctly so you can install go packages and build and run go programs.

Install pengui-go:
go get -u github.com/simonajones/penguin-go

This should download several dependent packages too. Some use the Bazaar VCS, so you'll need to have 'bzr' installed too.

Build and Run the server:
go run github.com/simonajones/penguin-go/main/penguin.go
