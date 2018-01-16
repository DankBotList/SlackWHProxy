## SlackWHProxy
====

Proxy webhooks to multiple recipients.

### Running
- go run main.go
- Edit config file.
- go run main.go

### Connecting
By default this runs and creates a socket in the cwd as SlackProxy.sock with permissions 0770 (srwxrwx---).
This can be changed by supplying -listener or -sock.
