# cloud

Sync your data with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)

version: **0.1**

## Install & Update

```bash
go get github.com/OShalakhin/cloud
go get -u github.com/OShalakhin/cloud
```

## Command line interface

```bash
$ cloud help
DESCRIPTION
    Sync your data with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)

VERSION
     0.1

COMMANDS
    init    initialize .cloudrc and .cloud files
    sync    synchronize folder with the cloud
    help    show this message

```

## Docs

See [Godoc](http://godoc.org/github.com/OShalakhin/cloud)

## TODO

- `cloud sync [container name]`
- `cloud add [folder|file]`
- `cloud update [folder|file]`
- improve `cloud sync` verifying if file was changed to upload it or not
- add `.cloudignore` and support the same behaviour like `.gitignore` (regexps)
- add `cloud ignore [folder/file]`
