build:
	goos=windows goarch=amd64 go build -o gonat-client-windows-amd64.exe ./client/cmd/*.go
	goos=windows goarch=amd64 go build -o gonat-server-windows-amd64.exe ./server/*.go
	goos=linux goarch=amd64 go build -o gonat-client-linux-amd64 ./client/cmd/*.go
	goos=linux goarch=amd64 go build -o gonat-server-linux-amd64 ./server/*.go
	goos=windows goarch=386 go build -o gonat-client-windows-386.exe ./client/cmd/*.go
	goos=windows goarch=386 go build -o gonat-server-windows-386.exe ./server/*.go
	goos=linux goarch=arm go build -o gonat-client-linux-arm ./client/cmd/*.go
	goos=linux goarch=arm go build -o gonat-server-linux-arm ./server/*.go

gui:
	goos=windows goarch=amd64 go build -ldflags "-H windowsgui" -o gonat-client-winGUI-amd64.exe ./client/gui/*.go
