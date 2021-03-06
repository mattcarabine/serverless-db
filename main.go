package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	badger "github.com/dgraph-io/badger/v3"
)

type MyEvent struct {
	ID string `json:"id"`
}

var db badger.DB

func init() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		panic(err)
	}
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte("answer"), []byte("42"))
		return err
	})
	if err != nil {
		panic(err)
	}
}

func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	var valCopy []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(event.ID))
		if err != nil {
			panic(err)
		}

		err = item.Value(func(val []byte) error {
			// Copying or parsing val is valid.
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Hello %s! The result was %s", event.ID, valCopy), nil
}

func main() {
	lambda.Start(HandleRequest)
}
