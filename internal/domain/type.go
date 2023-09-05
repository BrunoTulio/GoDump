package domain

import (
	"errors"
	"strings"
)

type (
	Type string
)

const (
	DUMP Type = "dump"
	SQL  Type = "sql"
)

func (t Type) String() string {
	return string(t)
}

var (
	TypeExtensions       = []string{DUMP.Extension(), SQL.Extension()}
	DefaultTypeExtension = DUMP.Extension()
)

func ReplaceFileExtension(fileName string) string {
	for _, t := range TypeExtensions {
		fileName = strings.ReplaceAll(fileName, t, "")
	}

	return fileName
}

func (c Type) Extension() string {
	switch c {
	case DUMP:
		return ".dump"
	case SQL:
		return ".sql"
	default:
		return ".dump"
	}
}

func ParseFileType(value string) (Type, error) {
	switch value {
	case "dump":
		return DUMP, nil
	case "sql":
		return SQL, nil
	default:
		return DUMP, errors.New("Type not supported")
	}
}
