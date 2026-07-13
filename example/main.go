package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	snowflake "github.com/Shudharshan07/idflake/v2"
)

type order struct {
	ID       snowflake.ID `json:"id"`
	Customer string       `json:"customer"`
}

func main() {
	sf := snowflake.NewSnowFlake(
		7,
		snowflake.WithNodeBits(10),
		snowflake.WithSequenceBits(12),
	)

	id := sf.Generate()
	fmt.Println("single id:", id)
	fmt.Println("int64:", id.Int64())
	fmt.Println("node:", sf.Node(id))
	fmt.Println("timestamp:", sf.Timestamp(id), "ms since idflake epoch")
	fmt.Println("sequence:", sf.Sequence(id))

	ids, err := sf.GenerateN(5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("batch:", ids)

	encoded := id.Base36()
	decoded, err := snowflake.ParseBase36(encoded)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("base36:", encoded)
	fmt.Println("base36 round trip:", decoded == id)

	payload, err := json.Marshal(order{
		ID:       id,
		Customer: "ada",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("json:", string(payload))

	var stored snowflake.ID
	dbValue, err := id.Value()
	if err != nil {
		log.Fatal(err)
	}
	if err := stored.Scan(dbValue); err != nil {
		log.Fatal(err)
	}
	fmt.Println("database round trip:", stored == id)

	createdAt := time.UnixMilli(sf.Timestamp(id) + 1288834974657)
	fmt.Println("created at:", createdAt.UTC().Format(time.RFC3339Nano))
}
