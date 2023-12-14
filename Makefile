.PHONY: wsclient
wsclient:
	go run ./cmd/ws

.PHONY: wsserver
wsserver:
	go run ./cmd/wsserver

.PHONY: h2client
h2client:
	go run ./cmd/h2

.PHONY: h2server
h2server:
	go run ./cmd/h2server
