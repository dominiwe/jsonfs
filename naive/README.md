# Naive implementation

This is a first proof of concept implementation of the idea.

It can mount a [json](https://www.json.org/json-en.html) file as a read-only fuse filesystem and unmounts it again when a `SIGINT` [signal](https://man7.org/linux/man-pages/man7/signal.7.html) is received.

## Usage

Mount a json file like this:

```shell
jsonfs-naive <file> <mountpoint>
```

Send a `SIGINT` signal (Ctrl + C) to unmount the filesystem or unmount it with the `fusermount` utility.

## Drawbacks

Some drawbacks of this implementation are:

- It unmarsalls the json into a go struct and then turns the values back into text which is superfluous.
- It uses reflection to figure out the types of the fields.
- It is read-only.