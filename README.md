# Gonat
Port Mapping

## What done
Map one port of a different network to another port

## Supports protocol
TCP

## Download
https://github.com/soyum2222/gonat/releases

## How to use
#### Server
`./gonat_server -c config.json `  

OR  

`./gonat_server  -port=1024
-crypt="aes-128-cbc"
-crypt_key="gonat"
-debug=true

#### Client
`./gonat_client -c config`  
OR  
` ./gonat_client-remote_ip="127.0.0.1:1024"
-crypt="aes-128-cbc"
-crypt_key="gonat"
-server_ip="127.0.0.1:8080"
-debug=true
-remote_port=8880`

## About config.json
#### Server
```
{
  "port": 1024,
  "debug": true,
  "crypt": "aes-128-cbc",
  "crypt_key": "gonat"
}
```
#### Client
```
{
"remote_ip":"192.168.0.2:1024",  
"server_ip":"127.0.0.1:8080",  
"remote_port":8081,  
"crypt": "aes-128-cbc",  
"crypt_key": "gonat"  
}
```
this config map IP 127.0.0.1:8080 to IP 192.168.0.2:8081

---
If your use gonat client in windows desktop, recommend you use https://github.com/soyum2222/gonat_client_gui

