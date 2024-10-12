Day 2 stack
===========

We want a ready for production (Kubernetes compatible) Golang app skeleton.
This means that a certain number of prerequisite have to be met:

1. DB connection + Routing + Rendering (the basic hello world feature)
1. 12 factor app:
	1. Configuration is loaded as follow:
		```
		File < Env variables < CLI args
		```
	1. Logs are outputed to stdout
1. DB operations
	1. Migrations
	1. Backup + Restore + Test
1. Kubernetes ready:
	1. Containerized
	1. Liveness/Readiness probes
	1. Graceful shutdown
1. Monitoring:
	1. Logs
		1. Structured
		1. Stack trace
	1. Metrics
	1. Tracing
	1. 5 Golden signals
1. Caching

1. **Bonus:** Application can be run locally without going through a container.
This ease the debugging process by simplifying the build phase
(no build container + remote debug anymore)

Decision logs
-------------

**Why not slog, why not zerolog?**

Development
-----------


```sh
go install github.com/a-h/templ/cmd/templ@v0.2.778
go install github.com/air-verse/air@latest
```

