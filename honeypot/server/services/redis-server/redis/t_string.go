package redis

const StringType = uint64(0)
const StringTypeFancy = "string"

var _ Item = (*String)(nil)

type String struct {
	value *string
}

func NewString(value *string) *String {
	return &String{value: value}
}

func (s *String) Value() interface{} {
	return s.value
}

func (s *String) Type() uint64 {
	return StringType
}

func (s *String) TypeFancy() string {
	return StringTypeFancy
}

func (s *String) OnDelete(key *string, db *RedisDb) {
}
