package client // import "github.com/sdslabs/docker/client"

import (
	"context"
	"net/url"
	"strconv"

	"github.com/sdslabs/docker/api/types"
)

// PluginEnable enables a plugin
func (cli *Client) PluginEnable(ctx context.Context, name string, options types.PluginEnableOptions) error {
	query := url.Values{}
	query.Set("timeout", strconv.Itoa(options.Timeout))

	resp, err := cli.post(ctx, "/plugins/"+name+"/enable", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
