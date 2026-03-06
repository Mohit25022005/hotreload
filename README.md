# HotReload

HotReload is a lightweight CLI tool written in Go that automatically rebuilds and restarts a server whenever source files change.

## Problem

Developers often need to manually stop, rebuild, and restart servers after every code change.
HotReload automates this workflow.

## Features

* Recursive file watching
* Debounced rebuilds
* Automatic server restart
* Real-time log streaming
* Ignore filters (.git, node_modules, build artifacts)
* Initial build on startup

## Usage

Build the tool:

go build -o hotreload.exe ./cmd/hotreload

Run the tool:

.\hotreload.exe --root ./testserver --build "go build -o ./bin/server.exe ./testserver" --exec ./bin/server.exe

## Demo

1. Start hotreload.
2. Edit a Go file in `testserver`.
3. Save the file.
4. The server automatically rebuilds and restarts.

## Architecture

* **Watcher** – monitors filesystem events
* **Debounce** – prevents excessive rebuilds
* **Builder** – executes build command
* **Runner** – manages server process lifecycle
