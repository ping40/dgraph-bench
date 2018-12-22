package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgraph-io/dgo"
)

const (
	qlFriendList = `
query test($a:  string) {
  everyone(func: anyofterms(name, $a)) {
    uid
    expand(_all_) 
  }
}`
)

func getFriendList(txn *dgo.Txn, names []string) (map[string]*Person, error) {

	nameStr := strings.Join(names, ",")

	fmt.Println("nameStr:", nameStr)
	resp, err := txn.QueryWithVars(context.Background(), qlFriendList, map[string]string{"$a": nameStr})
	if err != nil {
		fmt.Printf("Failed to getFriendList for %v: %v\n", names, err)
		return nil, err
	}

	var qp QueryPerson

	json.Unmarshal(resp.Json, &qp)
	fmt.Println("resp.Json:", string(resp.Json))
	fmt.Println("resultMap:", &qp)
	maps := make(map[string]*Person, 0)
	for _, v := range qp.Everyone {
		fmt.Printf("resultMap v2 %+v. \n", v)

		maps[v.Name] = v
	}

	fmt.Printf("resultMap v3  %+v. \n", maps)

	_ = resp.Json

	return maps, nil
}
