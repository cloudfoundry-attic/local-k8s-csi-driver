package main

import (
	"code.cloudfoundry.org/goshims/filepathshim"
	"code.cloudfoundry.org/goshims/osshim"
	"code.cloudfoundry.org/lager"
	csi "github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/jeffpak/local-controller-plugin/controller"
	"github.com/jeffpak/local-node-plugin/node"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type driver struct {
	logger        lager.Logger
	os            osshim.Os
	filepath      filepathshim.Filepath
	volumesRoot   string
	mountPathRoot string
	endpoint      string
	nodeId        string
}

func NewDriver(
	logger lager.Logger,
	os osshim.Os, filepath filepathshim.Filepath,
	volumesRoot string,
	mountPathRoot string,
	endpoint string,
	nodeId string,
) *driver {
	return &driver{
		logger:        logger,
		os:            os,
		filepath:      filepath,
		volumesRoot:   volumesRoot,
		mountPathRoot: mountPathRoot,
		endpoint:      endpoint,
		nodeId:        nodeId,
	}
}

func (d *driver) Run() {
	var ns csi.NodeServer
	ns = node.NewLocalNode(d.os, d.filepath, d.logger, d.volumesRoot, d.nodeId)

	var localController interface{}
	localController = controller.NewController(d.os, d.filepath, d.mountPathRoot)
	cs := localController.(csi.ControllerServer)
	ids := localController.(csi.IdentityServer)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.endpoint, ids, cs, ns)
	s.Wait()
}
