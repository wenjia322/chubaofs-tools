package main

import (
	"flag"
	"fmt"
	"github.com/chubaofs/chubaofs-tools/audit-daemon/daemon"
	"github.com/chubaofs/chubaofs-tools/audit-daemon/gather"
	"github.com/chubaofs/chubaofs-tools/audit-daemon/server"
	"github.com/chubaofs/chubaofs-tools/audit-daemon/util"
)

var (
	module   = flag.String("module", "", "start module about 'gather', 'daemon' , 'server'")
	port     = flag.Int("port", 8080, "Port Settings for the service")
	logLevel = flag.String("log_level", "debug", "log level")
	dbAddr   = flag.String("db_addr", "", "chubaodb address to send")
	dbTable  = flag.String("db_table", "", "chubaodb table to send")

	//gather need config
	config = flag.String("gather_conf", "", "gather module config path")
	logDir = flag.String("gather_log", "", "gather module log dir for raft parse")
)

func main() {
	flag.Parse()

	util.ConfigLog(*module, *logLevel)

	switch *module {
	case "gather":
		if *config == "" {
			panic("must set '-gather_conf' in gather module")
		}

		if *logDir == "" {
			panic("must set '-gather_log' in gather module")
		}

		if *dbAddr == "" || *dbTable == "" {
			panic("must set '-db_addr' and '-db_table' in gather module")
		}

		gather.StartGather(*config)
		util.LOG.Fatal(gather.StartRaftParse(*logDir, *dbAddr, *dbTable))
	case "daemon":
		daemon.StartServer(*port)
	case "server":
		server.StartServer(*port)
	default:
		fmt.Println(fmt.Sprintf("module type has err: [not support `%s`] use 'audit-daemon -h' see more", *module))
	}
}
