package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/linuxerwang/dgraph-bench/tasks"
)

func insertPersonConflict(dgraphCli *dgo.Dgraph) error {
	var start int64 = 1000

	person := &Person{
		Name:      "tom1",
		CreatedAt: start,
		UpdatedAt: start,
	}

	txn := dgraphCli.NewTxn()
	defer txn.Discard(context.Background())

	randomName := "tom4"
	friends := []string{"tom2", "tom3", randomName}
	old, err := getFriendList(txn, friends)
	if err != nil {
		return err
	}

	for _, f := range friends {
		if fo, ok := old[f]; ok {
			person.FriendOf = append(person.FriendOf, fo)
		} else {
			person.FriendOf = append(person.FriendOf, &Person{
				Name:      f,
				CreatedAt: start,
				UpdatedAt: start,
			})
		}
	}
	// Insert person node
	payload, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("Failed to marshal person object: %v\n", err)
		return err
	}

	mu := &api.Mutation{
		SetJson: payload,
	}

	fmt.Printf("\n ...  payload: %+v\n", string(payload))
	as, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		fmt.Printf("Failed to call mutate: %v\n", err)
		return err
	}

	go insertPersonConcurrent(dgraphCli, randomName)
	time.Sleep(time.Second * 3)

	if err = txn.Commit(context.Background()); err != nil {
		fmt.Printf("Failed to commit: %v\n", err)
		return err
	}

	fmt.Printf("as: %+v\n", as)
	fmt.Printf("InsertPerson: %v,  uid: %v \n", string(payload), as.Uids)
	_ = as

	return nil
}

func insertPersonConcurrent(dgraphCli *dgo.Dgraph, name string) error {
	var start int64 = 2000

	person := &tasks.Person{
		Name:      name,
		CreatedAt: start,
		UpdatedAt: start,
	}

	txn := dgraphCli.NewTxn()
	defer txn.Discard(context.Background())

	// Insert person node
	payload, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("insertPersonConcurrent Failed to marshal person object: %v\n", err)
		return err
	}

	mu := &api.Mutation{
		SetJson: payload,
	}

	fmt.Printf("\n ...  insertPersonConcurrent payload: %+v\n", string(payload))
	as, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		fmt.Printf("insertPersonConcurrent Failed to call mutate: %v\n", err)
		return err
	}

	if err = txn.Commit(context.Background()); err != nil {
		fmt.Printf("insertPersonConcurrent Failed to commit: %v\n", err)
		return err
	}

	fmt.Printf("insertPersonConcurrent as: %+v\n", as)
	fmt.Printf("insertPersonConcurrent InsertPerson: %v,  uid: %v \n", string(payload), as.Uids)
	_ = as

	return nil
}
