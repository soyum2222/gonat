build:
	go build -o gonat-client ./client/cmd/*.go
	go build -o gonat-server ./server/*.go
	go build -ldflags "-H windowsgui" -o gonat-client-gui ./client/gui/*.go

