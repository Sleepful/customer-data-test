package datastore

import (
	_ "encoding/json"
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
var Store = Datastore{listmap.New()}
var eventIds = make(map[string]struct{})

func PrintStore() {
	bytes, err := Store.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
	fmt.Println("Size: ", Store.Size())
}

func (d Datastore) MapGet(id int) (*serve.Customer, bool) {
	customer, exists := Store.Map.Get(id)
	if !exists {
		return nil, exists
	}
	// casting runtime errors for nil values
	return customer.(*serve.Customer), exists
}
func (d Datastore) MapValues() []*serve.Customer {
	arr := Store.Map.Values()
	result := make([]*serve.Customer, len(arr))
	for index, val := range arr {
		// yawn
		result[index] = val.(*serve.Customer)
	}
	return result
}
func (d Datastore) MapPut(id int, c *serve.Customer) {
	Store.Map.Put(id, c)
}
func (d Datastore) MapRemove(id int) {
	Store.Map.Remove(id)
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
	if old, exists := Store.MapGet(userId); exists {
		switch rec.Type {
		case "attributes":
			customer.Attributes = rec.Data
			customer.LastUpdated = rec.Timestamp
			mergeAttributes(old, &customer)
			Store.MapPut(userId, &customer)
			return
		case "event":
			mergeEvent(old, &customer, rec.Name)
			Store.MapPut(userId, &customer)
			return
		}
	} else {
		// first record for a user
		switch rec.Type {
		case "attributes":
			customer.Attributes = rec.Data
			customer.LastUpdated = rec.Timestamp
			Store.MapPut(userId, &customer)
			return
		case "event":
			customer.Events[rec.Name] = 1
			Store.MapPut(userId, &customer)
			return
		}
	}
}

func (d Datastore) Get(id int) (*serve.Customer, error) {
	if customer, exists := Store.MapGet(id); exists {
		return customer, nil
	} else {
		return nil, errors.New("Not found")
	}
}

func (d Datastore) List(page, count int) ([]*serve.Customer, error) {
	firstIndex := (page - 1) * count
	lastIndex := firstIndex + count
	values := Store.MapValues()
	if firstIndex > len(values) {
		return nil, errors.New("Out of bounds")
	}
	if lastIndex > len(values) {
		lastIndex = len(values)
	}
	result := values[firstIndex:lastIndex]
	return result, nil
}

func (m Datastore) Create(id int, attributes map[string]string) (*serve.Customer, error) {
	var customer = serve.Customer{
		ID:         id,
		Attributes: attributes,
		Events:     make(map[string]int),
	}
	Store.MapPut(id, &customer)
	created, _ := Store.MapGet(id)
	return created, nil
}

func (m Datastore) Update(id int, attributes map[string]string) (*serve.Customer, error) {
	old, _ := Store.MapGet(id)
	var new = serve.Customer{
		ID:         id,
		Attributes: attributes,
	}
	mergeAttributes(old, &new)
	Store.MapPut(id, &new)
	result, _ := Store.MapGet(id)
	return result, nil
}

func (m Datastore) Delete(id int) error {
	Store.MapRemove(id)
	return nil
}

func (m Datastore) TotalCustomers() (int, error) {
	return Store.Map.Size(), nil
}
