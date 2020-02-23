# Hot

- Hot help

```shell
> hot -h

NAME:
   Hot - Generate a generate.go file from the raml files in the api directory.

USAGE:
   hot [global options] command [command options] [arguments...]

VERSION:
   v0.2.0

COMMANDS:
   server   Generate a server according to a RAML specification
   client   Create a client for a RAML specification
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

- Hot server help

```shell
> hot server -h

NAME:
   hot server - Generate a server according to a RAML specification

USAGE:
   hot server [command options] [arguments...]

OPTIONS:
   --language value, -l value  Language to construct a server for (default: "go")
   --kind value                Kind of server to generate (gorestful) (default: "gorestful")
   --module value              Module name for go mod (default: "gitlab.com/hotbug/clients/gitproxy")
```

- Hot client help

```shell
> hot client -h

NAME:
   hot client - Create a client for a RAML specification

USAGE:
   hot client [command options] [arguments...]

OPTIONS:
   --language value, -l value  Language to construct a client for (default: "go")
   --kind value                Kind of client to generate (requests,grequests) (default: "requests")
   --package value             package name (default: "gitproxy")
   --module value              Module name for go mod (default: "gitlab.com/hotbug/clients/gitproxy")
```
