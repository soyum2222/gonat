# Gonat
Port Mapping

## What done
proxy the port to other network

## Supports protocol
TCP

## Download
https://github.com/soyum2222/gonat/releases

## How to use
#### Server
` ./gonat_server -c config.json `  

#### Client
` ./gonat_client -c config `  

### Start parameters
You can try run `./gonat help` get start param help

## About config.json
#### Server
```
{
  "port": 1024, // gonat server use this port listen to gonat client
  "debug": true, // if this is ture , then log will be print more info
  "crypt": "aes-128-cbc", 
  "crypt_key": "gonat" // password in gonat server to gonat client communication
}
```
#### Client
```
{
"remote_ip":"192.168.0.2:1024",  // gonat server
"server_ip":"127.0.0.1:8080",  // you wnat porxy server addr ,such as your want proxy MYSQL , here fill in 127.0.0.1:3306
"remote_port":8081,             // your want on the gonat server listen port for your `server`
"crypt": "aes-128-cbc",       
"crypt_key": "gonat"  
}
```
this config map IP 127.0.0.1:8080 to IP 192.168.0.2:8081



---
 ̶I̶f̶ ̶y̶o̶u̶r̶ ̶u̶s̶e̶ ̶g̶o̶n̶a̶t̶ ̶c̶l̶i̶e̶n̶t̶ ̶i̶n̶ ̶w̶i̶n̶d̶o̶w̶s̶ ̶d̶e̶s̶k̶t̶o̶p̶,̶ ̶r̶e̶c̶o̶m̶m̶e̶n̶d̶ ̶y̶o̶u̶ ̶u̶s̶e̶ ̶h̶t̶t̶p̶s̶:̶/̶/̶g̶i̶t̶h̶u̶b̶.̶c̶o̶m̶/̶s̶o̶y̶u̶m̶2̶2̶2̶2̶/̶g̶o̶n̶a̶t̶_̶c̶l̶i̶e̶n̶t̶_̶g̶u̶i̶
 
 
 Now the gonat client GUI has been migrated to this project


