# cloud

Sync your data with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)

version: **0.0.0**

## Install & Update

```bash
go get github.com/OShalakhin/cloud
go get -u github.com/OShalakhin/cloud
```

## CLI

```bash
$ cloud help

DESCRIPTION
    Sync your data with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)

VERSION
     0.1

COMMANDS
    init    initialize .cloudcore and .cloud files
    sync    synchronize folder with the cloud
    ignore  ignore particular file with .cloudignore
    clear   clear container
    help    show this message

CONTRIBUTORS
    Olexandr Shalakhin <olexandr@shalakhin.com>

```

## Docs

See [Godoc](http://godoc.org/github.com/OShalakhin/cloud)

## TODO

- &#10003; `cloud init `
- &#10003; initialize `~/.cloudcore`, `.cloud` and `.cloudignore` with samples
- &#10003; `cloud sync [container name]` (almost - create operation left)
- &#10003; support `.cloudignore` and support the same behaviour like `.gitignore` (regexps)
- &#10003; `cloud help`
- `cloud ignore [folder/file]`
- `cloud clear`
- `cloud add [folder|file]`
- `cloud rm [folder|file]`
- `cloud update [folder|file]`
- improve `cloud sync` verifying if file was changed to upload it or not
- `cloud -v` to show verbose information
