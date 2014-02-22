// Package cloud is a command line interface to sync files with cloud storages
// Just make
//
//     go get github.com/OShalakhin/cloud
//
// and try
//
//     cloud help
//
// `~/.cloudcore` stores credentials
//
//     {
//         "providers": [
//             {
//                 "provider": "CloudFiles",
//                 "name": "mynamehere",
//                 "key": "myapikey",
//                 "auth_url": "LON"
//             }
//         ]
//     }
//
// `.cloud` must be in the root folder you want to sync. Here you can define
// which containers to synchronize and define storage (CloudFiles, S3). More
// than one container is useful if you test with one container (even other storage)
// and deploy in another container.
//
//     {
//         "containers": [
//             {
//                 "provider": "CloudFiles",
//                 "name": "mycontainer1"
//             },
//             {
//                 "provider": "CloudFiles",
//                 "name": "mycontainer2"
//             }
//         ]
//     }
//
// `.cloudignore` is just like `.gitignore` file where you can add regexps which
// files or folders to ignore.
//
//     // Put here what to ignore. Syntax like .gitignore
//     .cloud
//     .cloudignore
//
package main
