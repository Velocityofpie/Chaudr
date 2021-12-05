# chaudr

Link to deployed project: http://ec2-18-232-171-115.compute-1.amazonaws.com:8080/

Unfortunatly, we were unable to figure out how to deploy the app with an HTTPS certificate in time, but it was our intention. Currently, the server is running on plain HTTP. We also had plans to deploy AWS WAF (Web Application Firewall) to protect against DDOS attacks, and a load balacer that terminated idle connections so that the server wouldn't be overloaded, but we simply didn't have enough time.

# Prequisites

You need to have a couple of tools installed to run the project.

## Node and NPM

First you need node to build the frontend. As such, you need `node` and `npm` installed, prefereably the latest LTS version, which is Node 16 and its associated NPM version. There are many ways to install this, and it differs per OS, so I will leave it up to user to install. However, there are giuides you can follow at the link here: https://nodejs.org/en/download/package-manager/#windows, links to download here: https://nodejs.org/en/download/. For the more advanced, I suggest using [`fnm`](https://github.com/Schniz/fnm), which I found to be quite useful, **although this is not mandatory at all**.

## Golang

The second thing you need installed is the Golang programming language. Again, the installion process differs from OS to OS, but there are guides for each OS here: https://go.dev/doc/install.

# Building the project

## Building the client

After you have the above tools installed, you can build the project. The first thing you have to do is go into the [`client`](https://github.com/Velocityofpie/chaudr/tree/main/client) folder and run `npm install`, followed by `npm run build`. After you have done that, go back the a folder, so you should be back where you started when you first cloned the project.

## Building the Golang program

Again, **make sure you are in the repository root**. Now run `go build -o chaudr`. After the command is complete, you should see a file called `chaudr`. This is an executable which you can run. This executable contains the bundled client code as well, so you can copy this single executable to wherever you want and run it without worrying about moving the client code to the same location.

# Running

On \*nix systems (Linux and Mac OS), you can run the program with `./chaudr`. This will start the server on the port `8080`, but you can control that by passing a flag to the program. If you want to run another port, for example `7070`, you can run `./chaudr -addr :7070`. Note that you must put a `:` before the port number. Navigate to [localhost:8080](http://localhost:8080)

# Issues

If you are a Window's user, you will most likely have issues running `go build`. There is an issue with a library we used (a sqlite library). However, if you have Docker Desktop installed with linux containers, you can run the app using `docker build . -t chaudr` and `docker run -p 8080:8080 chaudr`. Or if you have WSL, you can build and run the application inside of a WSL instance.
