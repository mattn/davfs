# davfs

WebDAV filesystem

## Usage

```
$ davfs
```

## Supported Drivers

|Driver    |Options to be specified           |
|----------|----------------------------------|
|file      |-driver=file -source=/path/to/root|
|memory    |-driver=memory                    |
|sqlite3   |-driver=sqlite3 -source=fs.db     |
|mysql     |-driver=mysql -source=blah...     |
|postgresql|-driver=postgres -source=blah...  |


## Installation

```
$ go get github.com/mattn/davfs/cmd/davfs
```

At the first time, you need to create filesystem for database drivers you specified like below.

```
$ davfs -driver=sqlite3 -source=fs.db -create
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
