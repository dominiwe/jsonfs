package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var jsonMap map[string]interface{}

type inMemoryFS struct {
	fs.Inode
}

var _ = (fs.NodeOnAdder)((*inMemoryFS)(nil))

func handleValue(i string, v interface{}, parent *fs.Inode, ctx context.Context) {
	var child *fs.Inode
	if v == nil {
		embedder := &fs.MemRegularFile{
			Data: []byte{},
		}
		child = parent.NewPersistentInode(ctx, embedder, fs.StableAttr{})
		parent.AddChild(i, child, true)
		return
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		child = parent.NewPersistentInode(ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})
		handleArray(ctx, child, v.([]interface{}))
	case reflect.Map:
		child = parent.NewPersistentInode(ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})
		handleObject(ctx, child, v.(map[string]interface{}))
	case reflect.Bool:
		embedder := &fs.MemRegularFile{
			Data: []byte(strconv.FormatBool(v.(bool))),
		}
		child = parent.NewPersistentInode(ctx, embedder, fs.StableAttr{})
	case reflect.String:
		embedder := &fs.MemRegularFile{
			Data: []byte(v.(string)),
		}
		child = parent.NewPersistentInode(ctx, embedder, fs.StableAttr{})
	case reflect.Float64:
		embedder := &fs.MemRegularFile{
			Data: []byte(strconv.FormatFloat(v.(float64), 'f', -1, 64)),
		}
		child = parent.NewPersistentInode(ctx, embedder, fs.StableAttr{})
	}
	parent.AddChild(i, child, true)
}

func handleObject(ctx context.Context, parent *fs.Inode, field map[string]interface{}) {
	for k, v := range field {
		handleValue(k, v, parent, ctx)
	}
}

func handleArray(ctx context.Context, parent *fs.Inode, field []interface{}) {
	for i, v := range field {
		handleValue(strconv.Itoa(i), v, parent, ctx)
	}
}

func (root *inMemoryFS) OnAdd(ctx context.Context) {
	handleObject(ctx, &root.Inode, jsonMap)
}

func flagUsage() {
	w := flag.CommandLine.Output()
	_, err := fmt.Fprintf(w,
		"usage: %s [options] filename mountpoint\noptions:\n", os.Args[0])
	if err != nil {
		panic(err)
	}
	flag.PrintDefaults()
}

func usageErrExit() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage()
	os.Exit(2)
}

func main() {

	debugFlag := flag.Bool("d", false, "enable debug output")
	helpFlag := flag.Bool("h", false, "show help")
	flag.Usage = flagUsage
	flag.Parse()

	if len(os.Args) < 3 {
		if *helpFlag {
			flag.Usage()
			os.Exit(0)
		} else {
			usageErrExit()
		}
	}

	if flag.NArg() != 2 {
		usageErrExit()
	}

	jsonFile := flag.Arg(0)
	mntDir := flag.Arg(1)

	fileContents, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Panic(err)
	}

	jsonMap = make(map[string]interface{})
	err = json.Unmarshal(fileContents, &jsonMap)
	if err != nil {
		log.Panic(err)
	}

	root := &inMemoryFS{}
	server, err := fs.Mount(mntDir, root, &fs.Options{
		MountOptions: fuse.MountOptions{Debug: *debugFlag},
	})
	if err != nil {
		log.Panic(err)
	}

	server.Wait()
}
