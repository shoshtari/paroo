package cache

type InmemoryCache[KeyType comparable, ValType any] struct {
	Data map[KeyType]ValType
}

func (i InmemoryCache[KeyType, ValType]) Get(key KeyType) (ValType, error) {
	if val, exists := i.Data[key]; !exists {
		return val, NotfoundError
	} else {
		return val, nil
	}
}

func (i InmemoryCache[KeyType, ValType]) Exists(key KeyType) (bool, error) {
	_, exists := i.Data[key]
	return exists, nil
}

func (i InmemoryCache[KeyType, ValType]) Set(key KeyType, val ValType) error {
	i.Data[key] = val
	return nil
}

func (i InmemoryCache[KeyType, ValType]) Delete(key KeyType) error {
	if _, exists := i.Data[key]; !exists {
		return NotfoundError
	}

	delete(i.Data, key)
	return nil
}

func NewInmemoryCache[keytype comparable, valtype any]() Cache[keytype, valtype] {
	c := InmemoryCache[keytype, valtype]{
		Data: make(map[keytype]valtype),
	}
	return c

}
