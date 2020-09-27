package shell

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/util"
)

func init() {
	Commands = append(Commands, &commandFsMetaSave{})
}

type commandFsMetaSave struct {
}

func (c *commandFsMetaSave) Name() string {
	return "fs.meta.save"
}

func (c *commandFsMetaSave) Help() string {
	return `save all directory and file meta data to a local file for metadata backup.

	fs.meta.save /               # save from the root
	fs.meta.save -v -o t.meta /  # save from the root, output to t.meta file.
	fs.meta.save /path/to/save   # save from the directory /path/to/save
	fs.meta.save .               # save from current directory
	fs.meta.save                 # save from current directory

	The meta data will be saved into a local <filer_host>-<port>-<time>.meta file.
	These meta data can be later loaded by fs.meta.load command, 

`
}

func (c *commandFsMetaSave) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	fsMetaSaveCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	verbose := fsMetaSaveCommand.Bool("v", false, "print out each processed files")
	outputFileName := fsMetaSaveCommand.String("o", "", "output the meta data to this file")
	// chunksFileName := fsMetaSaveCommand.String("chunks", "", "output all the chunks to this file")
	if err = fsMetaSaveCommand.Parse(args); err != nil {
		return nil
	}

	path, parseErr := commandEnv.parseUrl(findInputDirectory(fsMetaSaveCommand.Args()))
	if parseErr != nil {
		return parseErr
	}

	fileName := *outputFileName
	if fileName == "" {
		t := time.Now()
		fileName = fmt.Sprintf("%s-%d-%4d%02d%02d-%02d%02d%02d.meta",
			commandEnv.option.FilerHost, commandEnv.option.FilerPort, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}

	dst, openErr := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if openErr != nil {
		return fmt.Errorf("failed to create file %s: %v", fileName, openErr)
	}
	defer dst.Close()

	err = doTraverseBfsAndSaving(commandEnv, writer, path, *verbose, func(outputChan chan interface{}) {
		sizeBuf := make([]byte, 4)
		for item := range outputChan {
			b := item.([]byte)
			util.Uint32toBytes(sizeBuf, uint32(len(b)))
			dst.Write(sizeBuf)
			dst.Write(b)
		}
	}, func(entry *filer_pb.FullEntry, outputChan chan interface{}) (err error) {
		bytes, err := proto.Marshal(entry)
		if err != nil {
			fmt.Fprintf(writer, "marshall error: %v\n", err)
			return
		}

		outputChan <- bytes
		return nil
	})

	if err == nil {
		fmt.Fprintf(writer, "meta data for http://%s:%d%s is saved to %s\n", commandEnv.option.FilerHost, commandEnv.option.FilerPort, path, fileName)
	}

	return err

}

func doTraverseBfsAndSaving(filerClient filer_pb.FilerClient, writer io.Writer, path string, verbose bool, saveFn func(outputChan chan interface{}), genFn func(entry *filer_pb.FullEntry, outputChan chan interface{}) error) error {

	var wg sync.WaitGroup
	wg.Add(1)
	outputChan := make(chan interface{}, 1024)
	go func() {
		saveFn(outputChan)
		wg.Done()
	}()

	var dirCount, fileCount uint64

	err := filer_pb.TraverseBfs(filerClient, util.FullPath(path), func(parentPath util.FullPath, entry *filer_pb.Entry) {

		protoMessage := &filer_pb.FullEntry{
			Dir:   string(parentPath),
			Entry: entry,
		}

		if err := genFn(protoMessage, outputChan); err != nil {
			fmt.Fprintf(writer, "marshall error: %v\n", err)
			return
		}

		if entry.IsDirectory {
			atomic.AddUint64(&dirCount, 1)
		} else {
			atomic.AddUint64(&fileCount, 1)
		}

		if verbose {
			println(parentPath.Child(entry.Name))
		}

	})

	close(outputChan)

	wg.Wait()

	if err == nil && writer != nil {
		fmt.Fprintf(writer, "total %d directories, %d files\n", dirCount, fileCount)
	}
	return err
}
