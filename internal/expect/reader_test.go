package expect

import (
	"bufio"
	"io"
	"sync"
	"testing"
)

func TestReaderLease1(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() {
		reader.Close()
		writer.Close()
	}()

	//newLease := NewReaderLease(reader)
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//
	var wg1 sync.WaitGroup
	go func() {
		defer wg1.Done()

	}()
}

func TestPip1(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() {
		reader.Close()
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r := bufio.NewReader(reader)
		for {
			b, err := r.ReadByte()
			if err != nil {
				t.Logf("read err %#v", err)
				return
			}
			t.Logf("read: %s", string(b))
		}
	}()
	writer.Write([]byte("abc"))
	writer.Write([]byte("def"))
	writer.Close()
	wg.Wait()
}
