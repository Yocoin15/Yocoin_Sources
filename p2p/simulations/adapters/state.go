// Authored and revised by YOC team, 2017-2018
// License placeholder #1
package adapters

type SimStateStore struct {
	m map[string][]byte
}

func (self *SimStateStore) Load(s string) ([]byte, error) {
	return self.m[s], nil
}

func (self *SimStateStore) Save(s string, data []byte) error {
	self.m[s] = data
	return nil
}

func NewSimStateStore() *SimStateStore {
	return &SimStateStore{
		make(map[string][]byte),
	}
}
