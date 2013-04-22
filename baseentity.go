package chatservice

import (
	"errors"
	"sync"
)

//BaseEntity is not a Entity. It's just a helper class for other real Entities
//It become this because I didn't relize what the Go's inheritance really is.
type BaseEntity struct {
	id            idType
	store         *BaseStorage
	newEntityFunc func() Entity
	cloneFunc     func(dst Entity, src Entity)
}

type BaseStorage struct {
	basemap map[idType]Entity
	idindex idType
	mutex   *sync.Mutex
}

// var (
// 	messagestore = BaseStorage{make(map[idType]*BaseEntity), idType(0), sync.Mutex{}}
// )

func NewBaseEntity(store *BaseStorage, newentityfunc func() Entity, clonefunc func(dst, src Entity)) BaseEntity {
	return BaseEntity{idType(0), store, newentityfunc, clonefunc}
}

func NewBaseStore() *BaseStorage {
	return &BaseStorage{make(map[idType]Entity), idType(0), new(sync.Mutex)}
}

func (e *BaseEntity) SaveEntity(self Entity) error {
	store := e.store
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if v, ok := store.basemap[e.id]; ok {
		e.cloneFunc(v, self)
	} else {
		//need a new id.
		store.idindex++
		newid := store.idindex
		toadd := e.newEntityFunc()
		toadd.SetId(newid)
		e.id = newid
		e.cloneFunc(toadd, self)
		store.basemap[e.id] = toadd
	}
	return nil
}

func (e *BaseEntity) LoadEntity(self Entity, id idType) error {
	store := e.store
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if v, ok := store.basemap[id]; ok {
		e.id = id
		e.cloneFunc(self, v)
		return nil
	}
	return errors.New("BaseEntity.Load failed for id did not exist!")
}

func (e *BaseEntity) Delete() error {
	store := e.store
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if _, ok := store.basemap[e.id]; ok {
		delete(store.basemap, e.id)
		return nil
	}
	return errors.New("BaseEntity.Delete failed for id did not exist!")
}

func (e *BaseEntity) Id() idType {
	return e.id
}

func (e *BaseEntity) SetId(id idType) {
	e.id = id
}
