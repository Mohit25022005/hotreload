# HotReload 🔥

HotReload is a lightweight CLI tool written in **Go** that automatically rebuilds and restarts a server whenever source files change.

It eliminates the need for developers to manually stop, rebuild, and restart services during development.

---

# Problem

During development, engineers often follow this repetitive workflow:

1. Edit code
2. Stop server
3. Rebuild project
4. Restart server

This process slows down development and interrupts flow.

HotReload solves this problem by **watching files and restarting the server automatically when code changes.**

---

# Features

* Recursive project directory watching
* Automatic rebuild on file changes
* Automatic server restart
* Debounced file events to avoid redundant rebuilds
* Build cancellation when new changes occur
* Real-time server log streaming
* Dynamic detection of newly created directories
* Ignore filters for irrelevant files and directories
* Cross-platform CLI tool

---

# Architecture

The system is organized into modular components:

```
Filesystem Events
        ↓
     Watcher
        ↓
    Event Queue
        ↓
     Debounce
        ↓
      Builder
        ↓
      Runner
```

### Components

**Watcher**

* Uses `fsnotify` to monitor filesystem events
* Watches directories recursively
* Detects new directories dynamically

**Event Queue**

* Prevents event storms caused by editors
* Ensures only one rebuild is queued at a time

**Debouncer**

* Collapses multiple rapid file events into a single rebuild trigger

**Builder**

* Executes the build command
* Cancels previous builds if new changes occur
* Streams build logs in real time

**Runner**

* Manages the server process lifecycle
* Stops the old server before starting a new one
* Protects against crash restart loops

---

# Project Structure

```
hotreload
│
├ cmd
│   └ hotreload
│       └ main.go
│
├ internal
│   ├ builder
│   ├ debounce
│   ├ event
│   ├ logx
│   ├ runner
│   └ watcher
│
├ testserver
│   └ main.go
│
├ Makefile
├ go.mod
└ README.md
```

---

# Installation

Clone the repository:

```
git clone <repo-url>
cd hotreload
```

Build the CLI tool:

```
go build -o hotreload ./cmd/hotreload
```

---

# Usage

Run HotReload with the following parameters:

```
hotreload --root <project-folder> --build "<build-command>" --exec "<run-command>"
```

Example:

```
./hotreload \
--root ./testserver \
--build "go build -o ./bin/server ./testserver" \
--exec "./bin/server"
```

---

# Demo

A sample HTTP server is included in the repository under:

```
testserver/
```

Run the demo using:

```
make run
```

Open the server:

```
http://localhost:9090
```

Edit `testserver/main.go` and save the file.
HotReload will automatically:

* detect the change
* rebuild the project
* restart the server

---

# Logging

HotReload uses structured logging with Go's `log/slog` package.

Example logs:

```
[WATCH] file change detected
[BUILD] building project
[BUILD] build completed
[SERVER] restarting server
```

---

# File Filtering

The watcher ignores unnecessary files to reduce noise and improve performance.

Ignored directories:

```
.git
node_modules
bin
tmp
.vscode
.idea
```

Ignored temporary files:

```
*.swp
*.tmp
*~
```

---

# Tests

Unit tests cover critical components such as:

* debounce timing logic
* build execution behavior

Run tests using:

```
go test ./...
```

---

# Design Decisions

**Debouncing**

Editors often generate multiple file system events when saving files.
Debouncing prevents redundant rebuilds.

**Event Queue**

An event queue ensures that multiple rapid events trigger only a single rebuild.

**Build Cancellation**

If a rebuild is already in progress and a new change occurs, the previous build is cancelled so the latest state is always built.

**Process Lifecycle Management**

The runner ensures that the previous server instance is terminated before starting a new one.

---

# Author

Mohit Swarnkar
