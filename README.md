# Hello World App

## Introduction

This is simple Hello World Golang application with following three endpoints.

- **/hello** is a sample endpoint

- **/health**  is an endpoint for an external agent to monitor the app's liveliness

- **/metrics** presents number of successful and failed invocations
of the '/hello' endpoint

## Build

Export the environment variable to build the binary.
defines the operating system for which to compile the golang code.
Examples are linux, windows, darwin etc.

```bash
export GOOS="linux"
```

Run the following command for creating the binary file. You can specify the
executable file by using [-o] flag. And in our example the executable is
target/eric-oss-hello-world-app.

**mod=vendor** flag constructs a directory named vendor in the root
directory that contains copies of all packages needed to support builds
and tests of packages in the src module.

```bash
go build -mod=vendor -o target/eric-oss-hello-world-go-app ./src
```

## Run

**Note:** The following procedure might not work on Windows.
If you have issues, try using a Linux terminal with WSL to execute the binary file.

Provide Execute Permissions to the Binary File

```bash
sudo chmod +x target/eric-oss-hello-world-go-app
```

Execute the Hello World Application

```bash
./target/eric-oss-hello-world-go-app
```

Make a Request to **/hello** Endpoint

```bash
curl -is localhost:8050/hello
```

Example Output:

```python
HTTP/1.1 200 OK
Date: Thu, 17 Jun 2021 14:46:46 GMT
Content-Length: 13
Content-Type: text/plain; charset=utf-8

Hello World!!
```

## Test

Execute the unit testcases

```bash
go test -v ./src
```

For specific unit test case command:

```bash
go test -v --run [TESTCASE_NAME] .\src\
```

Hello World APP SDK Documentation [Here](https://arm1s11-eiffel004.eiffel.gic.ericsson.se:8443/nexus/content/sites/tor/idun-sdk/latest/index.html#getting-started).
