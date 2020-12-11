package persistence

type MockPersistence struct {
	MockCreate            func(interface{}) error
	MockFind              func(interface{}) error
	MockFindAll           func(interface{}, ...interface{}) error
	MockUpdate            func(interface{}) error
	MockDelete            func(interface{}) error
	MockLoadRelationships func() Persistence
}

func (p MockPersistence) LoadRelationships() Persistence {
	return p.MockLoadRelationships()
}

func (p MockPersistence) Create(item interface{}) error {
	return p.MockCreate(item)
}

func (p MockPersistence) Find(item interface{}) error {
	return p.MockFind(item)
}

func (p MockPersistence) FindAll(items interface{}, conds ...interface{}) error {
	return p.MockFindAll(items, conds...)
}

func (p MockPersistence) Update(item interface{}) error {
	return p.MockUpdate(item)
}

func (p MockPersistence) Delete(item interface{}) error {
	return p.MockDelete(item)
}
