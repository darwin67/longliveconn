.PHONY: server
server:
	go run ./cmd/server

.PHONY: wsclient
wsclient:
	go run ./cmd/ws

.PHONY: h2client
h2client:
	go run ./cmd/h2
