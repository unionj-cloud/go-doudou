package caches

import (
	"errors"
	"sync"
)

type cacherMock struct {
	store *sync.Map
}

func (c *cacherMock) Delete(tag string, tags ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *cacherMock) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *cacherMock) Get(key string) *Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*Query)
}

func (c *cacherMock) Store(key string, val *Query) error {
	c.init()
	c.store.Store(key, val)
	return nil
}

type cacherStoreErrorMock struct {
	store *sync.Map
}

func (c *cacherStoreErrorMock) Delete(tag string, tags ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *cacherStoreErrorMock) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *cacherStoreErrorMock) Get(key string) *Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*Query)
}

func (c *cacherStoreErrorMock) Store(string, *Query) error {
	return errors.New("store-error")
}
