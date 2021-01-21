# Hearsay

Listens for requests and passes them along, using a defined upstream proxy.

Useful for concentrating egress from multiple internal sources.

## Usage

```
Usage of hearsay:
  -dest string
        Final destination URL (default "https://example.com")
  -port string
        Listen port (default "8080")
  -proxy string
        Upstream proxy URI (default "http://localhost:8181")
  -v string
        Verbose
```



## Compiling

You can "bake in" any arguments during the build process so they don't need to
be passed in on the command-line. You'll need to have `make` installed.

```bash
PORT=1080 DEST=google.com PROXY=http://test:test@proxy.local make
```

When flags are passed to make-built binaries, they will override the built-in arguments.
