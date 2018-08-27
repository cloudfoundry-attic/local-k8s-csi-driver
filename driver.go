package main

import (
	"code.cloudfoundry.org/goshims/filepathshim"
	"code.cloudfoundry.org/goshims/osshim"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/local-controller-plugin/controller"
	"code.cloudfoundry.org/local-node-plugin/node"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type driver struct {
	logger        lager.Logger
	os            osshim.Os
	osHelper      node.OsHelper
	filepath      filepathshim.Filepath
	volumesRoot   string
	mountPathRoot string
	endpoint      string
	nodeId        string
}

func NewDriver(
	logger lager.Logger,
	os osshim.Os,
	osHelper node.OsHelper,
	filepath filepathshim.Filepath,
	volumesRoot string,
	mountPathRoot string,
	endpoint string,
	nodeId string,
) *driver {
	return &driver{
		logger:        logger,
		os:            os,
		osHelper:      osHelper,
		filepath:      filepath,
		volumesRoot:   volumesRoot,
		mountPathRoot: mountPathRoot,
		endpoint:      endpoint,
		nodeId:        nodeId,
	}
}

func (d *driver) Run() {
	var ns csi.NodeServer
	ns = node.NewLocalNode(d.os, d.osHelper, d.filepath, d.logger, d.volumesRoot, d.nodeId)

	var localController interface{}
	localController = controller.NewController(d.os, d.filepath, d.mountPathRoot)
	cs := localController.(csi.ControllerServer)
	ids := localController.(csi.IdentityServer)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.endpoint, ids, cs, ns)
	s.Wait()
}
