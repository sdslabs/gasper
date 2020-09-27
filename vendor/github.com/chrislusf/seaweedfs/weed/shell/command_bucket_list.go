package shell

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math"

	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
)

func init() {
	Commands = append(Commands, &commandBucketList{})
}

type commandBucketList struct {
}

func (c *commandBucketList) Name() string {
	return "bucket.list"
}

func (c *commandBucketList) Help() string {
	return `list all buckets

`
}

func (c *commandBucketList) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	bucketCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	if err = bucketCommand.Parse(args); err != nil {
		return nil
	}

	_, parseErr := commandEnv.parseUrl(findInputDirectory(bucketCommand.Args()))
	if parseErr != nil {
		return parseErr
	}

	var filerBucketsPath string
	filerBucketsPath, err = readFilerBucketsPath(commandEnv)
	if err != nil {
		return fmt.Errorf("read buckets: %v", err)
	}

	err = filer_pb.List(commandEnv, filerBucketsPath, "", func(entry *filer_pb.Entry, isLast bool) error {
		if entry.Attributes.Replication == "" || entry.Attributes.Replication == "000" {
			fmt.Fprintf(writer, "  %s\n", entry.Name)
		} else {
			fmt.Fprintf(writer, "  %s\t\t\treplication: %s\n", entry.Name, entry.Attributes.Replication)
		}
		return nil
	}, "", false, math.MaxUint32)
	if err != nil {
		return fmt.Errorf("list buckets under %v: %v", filerBucketsPath, err)
	}

	return err

}

func readFilerBucketsPath(filerClient filer_pb.FilerClient) (filerBucketsPath string, err error) {
	err = filerClient.WithFilerClient(func(client filer_pb.SeaweedFilerClient) error {

		resp, err := client.GetFilerConfiguration(context.Background(), &filer_pb.GetFilerConfigurationRequest{})
		if err != nil {
			return fmt.Errorf("get filer configuration: %v", err)
		}
		filerBucketsPath = resp.DirBuckets

		return nil

	})

	return filerBucketsPath, err
}
