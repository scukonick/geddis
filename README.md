## Geddis ##
Geddis is a simple KV storage for golang with 
REST API. Also it's possible to use it as embedded storage.

### Installation ###
#### Docker ####
The easiest way to run geddis is to run it's docker
image:
```bash
docker run -d -p 8080:8080 scukonick/geddis
```
Geddis would be accessible on http://127.0.0.1:8080:
```bash
curl -X POST "http://127.0.0.1:8080/strings/jack" -d '{ "value": "It works!", "ttl": 100}'
curl "http://127.0.0.1:8080/strings/jack"
It works!
```

#### Build and run ####
Another way is to check out repo and build geddis locally.
```bash
go get -u  github.com/scukonick/geddis

mkdir app
cp GOPATH/bin/geddis app
cp GOPATH/src/github.com/scukonick/geddis/config.example.toml app/config.toml
cd app
./geddis
```
Where GOPATH is your GOPATH, usually `~/go`.

### API ###

#### Swagger ####
Geddis has REST API done with help of swagger.
So, there are at least two ways to look it up.

First - take a look at swagg    er.yaml in the repository 
with your favourite text editor.

Second - open https://editor.swagger.io/, and there:
 
File -> Import URL and paste there the next URL: 
https://raw.githubusercontent.com/scukonick/geddis/master/swagger.yaml

Swagger would build API documentation and even browser GUI.

#### Curl ####
