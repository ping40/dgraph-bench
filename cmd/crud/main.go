package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

var (
	flagServers = flag.String("servers", "", "Comma separated dgraph server endpoints")

	dgraphCli *dgo.Dgraph
)

func connect(servers string) *dgo.Dgraph {
	clis := make([]api.DgraphClient, 0, 5)
	for _, s := range strings.Split(strings.Replace(servers, " ", "", -1), ",") {
		if len(s) > 0 {
			fmt.Printf("Connect to server %s\n", s)
			conn, err := grpc.Dial(s, grpc.WithInsecure())
			if err != nil {
				panic(err)
			}
			clis = append(clis, api.NewDgraphClient(conn))
		}
	}
	return dgo.NewDgraphClient(clis...)
}

func main() {
	flag.Parse()

	if *flagServers == "" {
		fmt.Println("Flag --servers is required.")
		os.Exit(2)
	}
	dgraphCli = connect(*flagServers)
	insertPersonConflict(dgraphCli)

	fmt.Println("\n game2 is over \n ")
}
