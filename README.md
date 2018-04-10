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
So, there are at least three ways to look it up.

First - take a look at swagg    er.yaml in the repository 
with your favourite text editor.

Second - open https://editor.swagger.io/, and there:
 
File -> Import URL and paste there the next URL: 
https://raw.githubusercontent.com/scukonick/geddis/master/swagger.yaml

Swagger would build API documentation and even browser GUI.

Third, auto generated api documentation is located here:
https://github.com/scukonick/geddis/tree/master/cli/lib

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
curl -XDELETE "http://127.0.0.1:8080/delete/abc"
```

##### Getting existing keys #####
Returns all the keys starting with j:
```bash 
curl "http://127.0.0.1:8080/keys/j"
```


### Configuration ###
For now configuration file should be located in 
the same directory from which application is started.

Self explained example of config.toml.
 
```toml
[Store]
Size = 10  # approximate size of database, only for optimization 
StoreInterval = 15  # how often store to disk (in seconds)
WorkDir = "/tmp"  # where to store database 

[ServerAPI]
ListenAddr = ":8080"  # Listen Addr in form ":9090" or "127.0.0.1:8080"
```

### Command-line client ###

#### build ####
To build command line client run:
```bash
make client
```

#### requests with client ####
Here are examples of requests:

##### Strings #####
Setting value:
```bash
./cli-client strings set --url "http://127.0.0.1:8080" --key "a"  --value "b" --ttl 50
```
Getting value:
```bash
./cli-client strings get --url "http://127.0.0.1:8080" --key "a"
```

##### Maps #####
Setting value:
```bash
./cli-client maps set --url "http://127.0.0.1:8080" --key "people"  --map '{"first": "Jack", "second": "John"}' --ttl 100
```
Getting value:
```bash
./cli-client maps get --url "http://127.0.0.1:8080" --key "people"
```
Getting element of value by subkey:
```bash
./cli-client maps getSubKey --url "http://127.0.0.1:8080" --key "people"  --subkey "first"
```

##### Arrays #####
Setting value:
```bash
./cli-client arrays set --url "http://127.0.0.1:8080" --key "fruits"  --value "apple" --value "orange" --ttl 100
```
Getting value:
```bash
./cli-client arrays get --url "http://127.0.0.1:8080" --key "fruits"
```
Getting element of value by index:
```bash
./cli-client arrays getIndex --url "http://127.0.0.1:8080" --key "fruits"  --index 1
```

##### Deleting #####
Deletes element by key 'abc'
```bash
./cli-client arrays getIndex --url "http://127.0.0.1:8080" --key "fruits"  --index 1
```

##### Getting existing keys #####
Returns all the keys starting with peo:
```bash 
./cli-client common keys --url "http://127.0.0.1:8080" --key "peo" 
```

To get all existing keys pass "\*" in the key parameter.
(yep, it conflicts with key "\*" but it's the easiest way to implement it.
