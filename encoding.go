package snowflake

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

func (id ID) Int64() int64 {
	return int64(id)
}

func ParseInt64(id int64) ID {
	return ID(id)
}

// String returns the decimal string representation of the ID
func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func ParseString(id string) (ID, error) {
	res, err := strconv.ParseInt(id, 10, 64)

	return ID(res), err
}

// MarshalJSON encodes the ID as a JSON string to preserve full int64 precision
func (id ID) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 22)
	buf = append(buf, '"')
	buf = append(buf, id.String()...)
	buf = append(buf, '"')
	return buf, nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid ID format")
	}

	str := string(data[1 : len(data)-1])
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*id = ID(val)
	return nil
}

// for DBs
func (id ID) Value() (driver.Value, error) {
	return int64(id), nil
}

func (id *ID) Scan(src any) error {
	if src == nil {
		*id = 0
		return nil
	}

	var val int64
	switch v := src.(type) {
	case int64:
		val = v
	case []byte:
		var err error
		val, err = strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot scan %T into ID", src)
	}

	*id = ID(val)
	return nil
}

func (id ID) Bytes() []byte {
	return []byte(id.String())
}

func ParseBytes(id []byte) (ID, error) {
	res, err := strconv.ParseInt(string(id), 10, 64)

	return ID(res), err
}

func (id ID) Base2() string {
	return strconv.FormatInt(int64(id), 2)
}

func ParseBase2(id string) (ID, error) {
	res, err := strconv.ParseInt(id, 2, 64)

	return ID(res), err
}

func (id ID) Base32() string {
	return strconv.FormatInt(int64(id), 32)
}

func ParseBase32(id string) (ID, error) {
	val, err := strconv.ParseInt(id, 32, 64)

	return ID(val), err
}

func (id ID) Base36() string {
	return strconv.FormatInt(int64(id), 36)
}

func ParseBase36(id string) (ID, error) {
	res, err := strconv.ParseInt(id, 36, 64)

	return ID(res), err
}

func (id ID) Base64() string {
	return base64.StdEncoding.EncodeToString(id.Bytes())
}

func ParseBase64(id string) (ID, error) {
	val, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return 0, err
	}

	return ParseBytes(val)
}

func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

func ParseIntBytes(id [8]byte) (ID, error) {
	return ID(int64(binary.BigEndian.Uint64(id[:]))), nil
}

func (id ID) Hex() string {
	return strconv.FormatUint(uint64(id), 16)
}

func ParseHex(id string) (ID, error) {
	val, err := strconv.ParseInt(id, 16, 64)
	return ID(val), err
}
