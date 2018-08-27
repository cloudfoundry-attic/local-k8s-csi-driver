package main

import (
	"flag"
	"os"

	"code.cloudfoundry.org/goshims/filepathshim"
	"code.cloudfoundry.org/goshims/osshim"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/local-node-plugin/oshelper"
)

var (
	endpoint      = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	nodeId        = flag.String("nodeId", "", "ID of the Kubernetes node")
	volumesRoot   = flag.String("volumesRoot", "/tmp/_volumes", "path to directory where node plugin mount point start with")
	mountPathRoot = flag.String("mountPathRoot", "", "path to directory where controller plugin mount point start with")
)

func main() {
	flag.Parse()

	handle()
	os.Exit(0)
}

func handle() {
	logger := lager.NewLogger("local-k8s-csi-driver")
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.DEBUG)
	logger.RegisterSink(sink)

	os := &osshim.OsShim{}
	driver := NewDriver(logger, os, oshelper.NewOsHelper(os), &filepathshim.FilepathShim{}, *volumesRoot, *mountPathRoot, *endpoint, *nodeId)
	driver.Run()
}
