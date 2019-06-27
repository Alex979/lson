# lson
A command line program written in go to output directory structure in JSON format.

## Installation
```
go get github.com/Alex979/lson
```

## Usage
```console
$ lson --help
Usage:
  lson [file] [flags]

Flags:
  -h, --help      help for lson
      --version   version for lson
```

## Examples
```console
$ lson dir1
{
  "name": "dir1",
  "type": "directory",
  "size": 29,
  "children": [
    {
      "name": "dir2",
      "type": "directory",
      "size": 17,
      "children": [
        {
          "name": "resume.pdf",
          "type": "file",
          "size": 17
        }
      ]
    },
    {
      "name": "hello.txt",
      "type": "file",
      "size": 6
    },
    {
      "name": "world.sh",
      "type": "file",
      "size": 6
    }
  ]
}
```
