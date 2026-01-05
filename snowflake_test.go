package snowflake

import (
	"log"
	"sync"
	"testing"
)

func generateIDs(n int, node int64) {
	sf := NewSnowFlake(node)
	var prev ID

	for i := 0; i < n; i++ {
		val := sf.Generate()
		if prev > val {
			log.Print("ID Decreasing")
		}
		if prev == val {
			log.Print("Duplicate IDs")
		}
		prev = val
	}
}

func TestIncreasing(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(n int) {
			generateIDs(100, int64(n))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestGenerateN(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 1; i < n; i++ {
		if ids[i] != ids[i-1]+1 {
			t.Errorf("Invalid Id %v", ids[i])
		}
	}
}

func TestBase2(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseBase2(id.Base2())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid Base2 Conversion")
		}
	}
}

func TestBase32(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseBase32(id.Base32())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid Base32 Conversion")
		}
	}
}

func TestBase36(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseBase36(id.Base36())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid Base36 Conversion")
		}
	}
}

func TestBase64(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseBase64(id.Base64())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid Base64 Conversion")
		}
	}
}

func TestBytes(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseBytes(id.Bytes())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid Bytes Conversion")
		}
	}
}

func TestHex(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseHex(id.Hex())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid Hex Conversion")
		}
	}
}

func TestIntBytes(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id2, err := ParseIntBytes(id.IntBytes())
		if err != nil {
			t.Error(err)
		}
		if id2 != id {
			t.Error("Invalid IntBytes Conversion")
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id_json, err := id.MarshalJSON()
		if err != nil {
			t.Error(err)
		}

		id2 := ID(0)
		id2.UnmarshalJSON(id_json)

		if id2 != id {
			t.Error("Invalid MarshalJSON Conversion")
		}
	}
}

func TestDBId(t *testing.T) {
	sf := NewSnowFlake(1)

	n := 100
	ids, err := sf.GenerateN(n)
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		id_val, err := id.Value()
		if err != nil {
			t.Error(err)
		}

		id2 := ID(0)
		id2.Scan(id_val)

		if id2 != id {
			t.Error("Invalid DB id Conversion")
		}
	}
}

func BenchmarkThisLib(b *testing.B) {
	b.ResetTimer()
	sf := NewSnowFlake(1)

	for i := 0; i < b.N; i++ {
		sf.Generate()
	}
}

func TestExtractors(t *testing.T) {
	sf := NewSnowFlake(1)
	ids, err := sf.GenerateN(100)
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range ids {
		ts := sf.Timestamp(id)
		nd := sf.Node(id)
		seq := sf.Sequence(id)

		// Reconstruct ID from components
		recon := ID((ts << (int64(sf.nodeBits) + int64(sf.sequenceBits))) |
			(nd << int64(sf.sequenceBits)) | seq)

		if recon != id {
			t.Errorf("Reconstruct mismatch")
		}

		// Verify extracted node matches sf.node (single-node test)
		if nd != sf.node {
			t.Errorf("Wrong node")
		}

	}
}
