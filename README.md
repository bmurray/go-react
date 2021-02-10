# go-react: A quick and usable development method for integrating Go and React

## Purpose

go-react is a simple method for developing and deploying React apps that rely on a go backend. I needed a method to be able to easily develop for scenarios where the React code was served from the same end point as a Go based API. Fiddling with CORS policies seemed prone to error as I went front dev to prod. I was also finding myself deploying one server based apps to various servers, and found it tedious to ensure the React code was synchronized with the Go binary. Thus, this method was born. This isn't for larger deployments where you may have frontend application servers or reverse proxies, like nginx or Amazon ALB, but it can give you an easy way to develop for those environments. 

All the functionality is in the [main.go](./main.go) file. The key areas are documented for your convenience. The react app is in the [react-app](./react-app) directory, but would work in any directory with minor changes. 

This repository also serves as an easy location for me to grab the code from when starting new projects. 

## Usage

The built binary operates in one of three modes: proxy, dir, and embed mode. The produced binary is for demonstration purposes, and should be easily integrated into your own code. 

```
Usage of go-react:
  -dir string
        Directory where the static built app resides (default "./react-app/build/")
  -embed string
        Directory where the static built embeded app resides (1.16+) (default "react-app/build")
  -listen string
        Listen on address (default ":8080")
  -mode string
        Mode to serve REACT site: proxy, dir, pkger, embed  (default "proxy")
  -pkger string
        Directory where the static built embeded app resides (default "/react-app/build/")
  -proxy string
        Address to proxy requests to (default "http://localhost:3000/")
```

### Proxy mode

The default mode is proxy mode, and is most useful for development of React apps with a Go backend. The built-in proxy happily proxies websocket connections to the Node environment, which preserves live-updates. 

Start the npm server. Be careful; this often opens a browsers to the wrong port! Your API will appear broken if you use port 3000.
```
cd react-app/
npm start
```

Build the go code, and run it in proxy mode.
```
go run main.go

OR

go run main.go -mode proxy
```

Open your browser and point to http://localhost:8080/

You can edit your react code and benefit from live updates. However, if you restart the go server, you will need to refresh the page as the live-update socket will drop.

### Dir mode

The classic server mode; this serves files directly from the filesystem. Works as a nice intermediate. There is no live-reload. If you are deploying a larger package, such as via Docker, this method can sometimes be preferred as it keeps the react code out of the binary. 

First, build the react code into a deployment/production version:

```
cd react-app/
npm build
```

Run the go server in dir mode
```
go run main.go -mode dir
```


### Packaged mode (pre 1.16)

This method will build the react code into a deployment/production version, then embed it in the binary. This means you only need to produce a single binary for distribution without needing to deploy the react code as well. I find this very useful for programs that need a user interface but dont want the heft of a full update system. 

First, build the react code into a deployment/production version

```
cd react-app/
npm build
```

Next, package the contents of the build directory. Requires [github.com/markbates/pkger](github.com/markbates/pkger)

Install pkger (only once)
```
go get github.com/markbates/pkger/cmd/pkger
```

Package the react files into the binary; creates the pkged.go file in the root (./) directory
```
pkger -o ./ -include /react-app/build
```

Build and run the binary
```
go build -o go-react main.go pkged.go

./go-react -mode pkger

```

Alternatively, the Makefile has all these steps:
```
make run
```

### Embed mode (1.16+)

This provides the same methods as the packager method, but uses the new 1.16+ interface. 

First, build the react code into a deployment/production version

```
cd react-app/
npm build
```

Next, run the code. The files are automatically embeded. 

```
go run main.go -mode embed
```



