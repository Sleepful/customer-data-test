package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "golang.org/x/exp/slices"
	"log"
	"strconv"

	"github.com/customerio/homework/serve"
	"github.com/customerio/homework/stream"
	listmap "github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/imdario/mergo"
)

type Datastore struct {
	*listmap.Map
}

var _ serve.Datastore = Datastore{}
var store = Datastore{listmap.New()}
var eventIds = make(map[string]struct{})

func PrintStore() {
	bytes, err := store.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func mergeEvent(old *serve.Customer, customer *serve.Customer, eventName string) {
	if event, exists := old.Events[eventName]; exists {
		customer.Events[eventName] = event + 1
	} else {
		customer.Events[eventName] = 1
	}
	// merge to retain all previous Events
	if err := mergo.Merge(customer, *old); err != nil {
		log.Fatal(err)
	}
}

func mergeAttributes(old *serve.Customer, customer *serve.Customer) {
	// empty override function to have a valid
	// patameter for mergo.Merge
	var override = func(_ *mergo.Config) {}
	if old.LastUpdated > customer.LastUpdated {
		// override merge because
		// values from `old` are the
		// more recent values and must prevail
		override = mergo.WithOverride
	}
	if err := mergo.Merge(customer, *old, override); err != nil {
		log.Fatal(err)
	}
}

func InsertRecord(rec *stream.Record) {
	if _, exists := eventIds[rec.ID]; exists {
		return
	}
	eventIds[rec.ID] = struct{}{}
	if rec.UserID == "" {
		return
	}
	userId, err := strconv.Atoi(rec.UserID)
	if err != nil {
		log.Fatal(err)
	}
	// empty customer meant to represent
	// new values from record
	var customer = serve.Customer{
		ID:         userId,
		Attributes: make(map[string]string),
		Events:     make(map[string]int),
	}
	if old, exists := store.Map.Get(userId); exists {
		oldCasted := old.(serve.Customer)
		switch rec.Type {
		case "attributes":
			customer.Attributes = rec.Data
			customer.LastUpdated = rec.Timestamp
			mergeAttributes(&oldCasted, &customer)
			store.Put(userId, customer)
			return
		case "event":
			mergeEvent(&oldCasted, &customer, rec.Name)
			store.Put(userId, customer)
			return
		}
	} else {
		// first record for a user
		switch rec.Type {
		case "attributes":
			customer.Attributes = rec.Data
			customer.LastUpdated = rec.Timestamp
			store.Put(userId, customer)
			return
		case "event":
			customer.Events[rec.Name] = 1
			store.Put(userId, customer)
			return
		}
	}
}

func (d Datastore) Get(id int) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (d Datastore) List(page, count int) ([]*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (m Datastore) Create(id int, attributes map[string]string) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (m Datastore) Update(id int, attributes map[string]string) (*serve.Customer, error) {
	return nil, errors.New("unimplemented")
}

func (m Datastore) Delete(id int) error {
	return errors.New("unimplemented")
}

func (m Datastore) TotalCustomers() (int, error) {
	return 0, errors.New("unimplemented")
}
