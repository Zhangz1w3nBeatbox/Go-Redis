package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

//实现dict接口

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	val, ok := dict.m.Load(key)
	return val, ok
}

func (dict *SyncDict) Len() int {
	len := 0
	//range方法是把func施加到每个kv上
	dict.m.Range(func(key, value interface{}) bool {
		len++
		return true
	})
	return len
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if existed {
		return 0
	}
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		return 0
	}

	dict.m.Store(key, val)

	return 1
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		dict.m.Store(key, val)
		return 1

	}
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exited := dict.m.Load(key)
	dict.m.Delete(key)
	if exited {
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	res := make([]string, dict.Len())

	i := 0

	dict.m.Range(func(key, value interface{}) bool {
		res[i] = key.(string)
		i++
		return true
	})

	return res
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	res := make([]string, dict.Len())

	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value interface{}) bool {
			res[i] = key.(string)
			return false
		})
	}

	return res
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	res := make([]string, dict.Len())

	i := 0

	dict.m.Range(func(key, value interface{}) bool {
		res[i] = key.(string)
		i++
		if i == limit {
			return false
		}
		return true
	})

	return res
}

func (dict *SyncDict) clear() {
	*dict = *MakeSyncDict()
}
