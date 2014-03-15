# cloud [![GoDoc](http://godoc.org/github.com/OShalakhin/cloud?status.png)](http://godoc.org/github.com/OShalakhin/cloud)

version: **0.1.0**

Sync your files with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)

I decided it would be great to have simple abstaction layer for cloud storage
to sync data with CLI and to integrate it inside the app.

That's why `github.com/OShalakhin/cloud` is CLI and packages inside are API.

## Usage

```bash
$ cloud init # edit created files
$ cloud sync # synchronize files
```

Yep. So simple.

## Config files

`github.com/OShalakhin/cloud/storage` and other parts are for integration into
the apps. Tests will be supplied later.

- supported `init`, `sync` and `help` commands
- included documentation
- abstract structure

`~/.cloudcore` stores credentials

```json
{
    "providers": [
        {
            "provider": "CloudFiles",
            "name": "mynamehere",
            "key": "myapikey",
            "auth_url": "LON"
        }
    ]
}
```

`.cloud` must be in the root folder you want to sync. Here you can define
which containers to synchronize and define storage (CloudFiles, S3). More
than one container is useful if you test with one container (even other storage)
and deploy in another container.

```json
{
    "containers": [
        {
            "provider": "CloudFiles",
            "name": "mycontainer1"
        },
        {
            "provider": "CloudFiles",
            "name": "mycontainer2"
        }
    ]
}
```

`.cloudignore` is just like `.gitignore` file where you can add regexps which
files or folders to ignore.

```
// Put here what to ignore. Syntax like .gitignore
.cloud
.cloudignore
```

## Install & Update

```bash
# Install
go get github.com/OShalakhin/cloud
# Update
go get -u github.com/OShalakhin/cloud
```

## CLI

```bash
$ cloud help

DESCRIPTION
    Sync your data with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)

VERSION
     0.1.0

COMMANDS
    init    initialize .cloudcore and .cloud files
    sync    synchronize folder with the cloud
    ignore  ignore particular file with .cloudignore
    clear   clear container
    help    show this message

CONTRIBUTORS
    Olexandr Shalakhin <olexandr@shalakhin.com>

```

## TODO

- &#10003; `cloud init `
- &#10003; initialize `~/.cloudcore`, `.cloud` and `.cloudignore` with samples
- &#10003; `cloud sync [container]` (almost - create operation left)
- &#10003; support `.cloudignore` and support the same behaviour like `.gitignore` (regexps)
- &#10003; `cloud help`
- `cloud ignore [folder/file]`
- `cloud clear`
- `cloud ls`
- `cloud info [container] [file]`
- `cloud add [folder|file]`
- `cloud rm [folder|file]`
- `cloud update [folder|file]`
- improve `cloud sync` verifying if file was changed to upload it or not
- `cloud -v` to show verbose information
- support specific features like CDN operation for CloudFiles
