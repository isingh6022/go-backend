This is a microservices project I am working on. It will consist of multiple servers in multiple modules, based on `coreserver` module. It uses go workspaces.

## `coreserver`

Every other server depends on `coreserver` module that provides a basic struct `App` which has methods for starting a server. The server provided by `App` listens to termination signals such as `syscall.SIGTERM` and shuts down gracefully. It also provides methods for graceful shutdown.

However the `App` should generally be created with `CreateApp` function because it provides reasonable defaults.

#### `demoserver`

Creates a sever based on `coreserver` for sanity and testing.
