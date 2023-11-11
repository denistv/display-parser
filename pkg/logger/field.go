package logger

import "time"

// Field Поле context-based логгера
type Field struct {
	Key   string
	Value any
}

func NewInt64Field(key string, value int64) Field {
	return Field{
		Key: key,
		Value: value,
	}
}

func NewIntField(key string, value int) Field {
	return Field{
		Key: key,
		Value: value,
	}
}

func NewStringField(key string, value string) Field {
	return Field{
		Key: key,
		Value: value,
	}
}

func NewBoolField(key string, value bool) Field {
	return Field{
		Key: key,
		Value: value,
	}
}

func NewTimeField(key string, value time.Time) Field {
	return Field{
		Key: key,
		Value: value,
	}
}

func NewErrorField(key string, value error) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

