package imp

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Importer interface {
	Import() error
}

// Here is an implementation that talks over RPC
type ImporterRPC struct{ client *rpc.Client }

func (g *ImporterRPC) Import() error {
	var resp error
	return g.client.Call("Plugin.Import", new(interface{}), &resp)
}

// Here is the RPC server that ImporterRPC talks to, conforming to
// the requirements of net/rpc
type ImporterRPCServer struct {
	// This is the real implementation
	Impl Importer
}

func (s *ImporterRPCServer) Import(args interface{}, resp *error) error {
	*resp = s.Impl.Import()
	return nil
}

type ImporterPlugin struct {
	Impl Importer
}

func (p *ImporterPlugin) Server(*plugin.MuxBroker) (any, error) {
	return &ImporterRPCServer{Impl: p.Impl}, nil
}

func (p *ImporterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &ImporterRPC{client: c}, nil
}
