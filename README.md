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
Here are examples of curl requests to geddis
##### Strings #####
Setting value:
```bash
curl -X POST "http://127.0.0.1:8080/strings/jack" -d '{ "value": "It works!", "ttl": 200000}'
```
Getting value:
```bash
curl "http://127.0.0.1:8080/strings/jack"
```

##### Maps #####
Setting value:
```bash
curl -X POST "http://127.0.0.1:8080/maps/y" -d '{"value": {"a":"1","b":"2"}, "ttl": 0}'
```
Getting value:
```bash
curl "http://127.0.0.1:8080/maps/y"
```
Getting element of value by subkey:
```bash
curl "http://127.0.0.1:8080/maps/y/a"
```

##### Arrays #####
Setting value:
```bash
curl -X POST "http://127.0.0.1:8080/arrays/abc" -d '{ "values": [ "qqq", "www", "eeee" ], "ttl": 0}'
```
Getting value:
```bash
curl "http://127.0.0.1:8080/arrays/abc"
```
Getting element of value by index:
```bash
curl "http://127.0.0.1:8080/arrays/abc/2"
```

##### Deleting #####
Deletes element by key 'abc'
```bash
curl -XDELETE "http://127.0.0.1:8080/common/abc"
```

##### Getting existing keys #####
Returns all the keys starting with j:
```bash 
curl "http://127.0.0.1:8080/keys/j"
```


### Configuration ###
Self explained example of config.toml
```toml
[Store]
Size = 10  # approximate size of database, only for optimization 
StoreInterval = 15  # how often store to disk (in seconds)
WorkDir = "/tmp"  # where to store database 

[ServerAPI]
ListenAddr = ":8080"  # Listen Addr in form ":9090" or "127.0.0.1:8080"
```