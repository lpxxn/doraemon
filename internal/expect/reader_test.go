package expect

import (
	"bufio"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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
		writer.Write([]byte("abc"))
		writer.Write([]byte("def"))
		writer.Close()

	}()

	r := bufio.NewReader(reader)
	for {
		b, err := r.ReadByte()
		if err != nil {
			t.Logf("read err %#v", err)
			return
		}
		t.Logf("read: %s", string(b))
	}
	wg.Wait()
}

func TestPip2(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() {
		reader.Close()
		//writer.Close()
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		newReader, newWriter := io.Pipe()
		defer newReader.Close()
		defer newWriter.Close()

		go func() {
			defer writer.Close()
			_, err := io.Copy(writer, newReader)
			t.Logf("io copy err: %#v", err)
		}()

		//defer newWriter.Close()
		_, err := newWriter.Write([]byte("hello"))
		assert.Nil(t, err)
		_, err = newWriter.Write([]byte("world"))
		assert.Nil(t, err)
		//time.Sleep(5)

		//newWriter.Close()

	}()

	r := make([]byte, 5)
	for {
		b, err := reader.Read(r)
		if err != nil {
			t.Logf("read err %#v", err)
			return
		}
		t.Logf("read: %s, count: %d", string(r), b)
	}
	wg.Wait()
}
