# HotReload 🔥

HotReload is a lightweight CLI tool written in Go that automatically rebuilds and restarts a server whenever source files change.

It eliminates the need for developers to manually stop, rebuild, and restart services during development.

---

# Problem

During development, engineers often follow this workflow:

1. Edit code
2. Stop server
3. Rebuild project
4. Restart server

This process slows down development and interrupts flow.

HotReload automates this workflow by **watching files and restarting the server automatically when code changes.**

---

# Features

* Recursive file watching
* Automatic rebuild on file changes
* Automatic server restart
* Debounced file events to avoid redundant rebuilds
* Build cancellation when new changes occur
* Crash-loop protection
* Real-time log streaming
* File filtering to ignore unnecessary directories
* Dynamic detection of newly created folders

---

# Architecture

The system is organized into modular components.

```
            ┌─────────────┐
            │   Watcher   │
            │ (fsnotify)  │
            └──────┬──────┘
                   │
             file events
                   │
                   ▼
            ┌─────────────┐
            │  Debouncer  │
            │  (500 ms)   │
            └──────┬──────┘
                   │
             trigger rebuild
                   │
                   ▼
            ┌─────────────┐
            │   Builder   │
            │ executes    │
            │ build cmd   │
            └──────┬──────┘
                   │
           build successful
                   │
                   ▼
            ┌─────────────┐
            │   Runner    │
            │ stops old   │
            │ server and  │
            │ starts new  │
            └─────────────┘
```

---

# Project Structure

```
hotreload
│
├ cmd/hotreload
│  └ main.go
│
├ internal
│  ├ watcher
│  │   watcher.go
│  ├ debounce
│  │   debounce.go
│  ├ builder
│  │   builder.go
│  └ runner
│      runner.go
│
├ testserver
│  └ main.go
│
├ Makefile
├ go.mod
└ README.md
```

---

# Usage

Build the tool:

```
go build -o hotreload.exe ./cmd/hotreload
```

Run the tool:

```
.\hotreload.exe --root ./testserver --build "go build -o ./bin/server.exe ./testserver" --exec ./bin/server.exe
```

Parameters:

| Flag      | Description                         |
| --------- | ----------------------------------- |
| `--root`  | directory to watch                  |
| `--build` | command used to rebuild the project |
| `--exec`  | command used to start the server    |

---

# Demo

1. Start the tool.

```
.\hotreload.exe --root ./testserver --build "go build -o ./bin/server.exe ./testserver" --exec ./bin/server.exe
```

2. Open browser

```
http://localhost:9090
```

3. Edit `testserver/main.go`

4. Save the file.

HotReload automatically:

* detects file change
* rebuilds the project
* restarts the server

Logs will show:

```
INFO change detected
INFO building project
INFO stopping server
INFO starting server
```

---

# File Filtering

The watcher ignores unnecessary files to improve performance.

Ignored directories:

```
.git
node_modules
bin
tmp
```

Ignored temporary files:

```
*.swp
*.tmp
*~
```

Only the following files trigger rebuilds:

```
.go
.mod
.sum
```

---

# Build Cancellation

If a rebuild is already running and a new file change occurs, the previous build is cancelled.

This ensures the tool always builds the **latest version of the project**.

---

# Crash Loop Protection

If the server crashes immediately after starting, HotReload delays the restart to prevent a rapid restart loop.

---

# Logging

HotReload uses Go's `log/slog` package for structured logging.

Example logs:

```
time=2026-03-06T17:21:46 level=INFO msg="change detected"
time=2026-03-06T17:21:46 level=INFO msg="building project"
time=2026-03-06T17:21:48 level=INFO msg="starting server"
```

Server logs are streamed in real time.

---

# Tests

Unit tests cover tricky components such as the debounce logic.

Run tests:

```
go test ./...
```

---

# Demo Server

The repository includes a simple HTTP server inside `testserver/`.

```
testserver/main.go
```

This server is used to demonstrate the hot reload functionality.

---

# Design Decisions

**fsnotify**

Used as the event source for filesystem monitoring.

**Debouncing**

Editors often trigger multiple file events when saving files.
A debounce mechanism prevents excessive rebuilds.

**Process Management**

The runner ensures that the previous server instance is terminated before starting a new one.

**Modular Architecture**

Each responsibility is separated into its own package to maintain clean and testable code.

---

# Future Improvements

Possible improvements include:

* cross-platform process tree termination
* distributed file watching for very large projects
* configurable debounce duration
* plugin system for language-specific builds

---

# Author

Mohit Swarnkar
