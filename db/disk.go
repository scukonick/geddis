package db

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

const (
	filename = "geddis.db"
)

const (
	strPrefix = "s"
	arrPrefix = "a"
	mapPrefix = "m"
	ttlPrefix = "t"
)

/*
storage file prefix has the next structure:
s{"stringKey": "stringValue"}
a{"arrKey": ["value1", "value2"] }
m("mapKey": {"a": "b", "c":"d"}}
t("anyKey":"2016-06-06T00:00:00"}

There could be any amount of lines of any type.

*/

type fileStrValue struct {
	Key   string
	Value string
}

type fileArrValue struct {
	Key   string
	Value []string
}

type fileMapValue struct {
	Key   string
	Value map[string]string
}

type fileTTLValue struct {
	Key   string
	Value time.Time
}

func closeC(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("ERR: failed to close: %+v", err)
	}
}

func (s *GeddisStore) runFileStore() {
	if s.storeInterval == 0 {
		return
	}

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		ticker := time.NewTicker(s.storeInterval)
		for {
			select {
			case <-ticker.C:
				s.storeToDiskOnce()
			case <-s.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

// storeToDiskOnce stores everything in storage to disk
// TODO - separate file operations and write operations in order to decrease complexity
func (s *GeddisStore) storeToDiskOnce() {
	log.Print("Storing to disk...")
	defer log.Print("Storing to disk finished")
	tmpFilename := "." + filename
	fullPath := path.Join(s.workDir, tmpFilename)

	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("ERR: failed to open db file: %+v", err)
		return
	}

	bufWriter := bufio.NewWriter(f)

	s.lock.RLock()
	defer s.lock.RUnlock()

	var js []byte

	for key, value := range s.m {
		var data interface{}
		var prefix string
		switch value.(type) {
		case string:
			prefix = strPrefix
			data = &fileStrValue{Key: key, Value: value.(string)}
		case []string:
			prefix = arrPrefix
			data = &fileArrValue{Key: key, Value: value.([]string)}
		case map[string]string:
			prefix = mapPrefix
			data = &fileMapValue{Key: key, Value: value.(map[string]string)}
		}

		js, err = json.Marshal(data)
		if err != nil {
			log.Printf("ERR: failed to marshal data for storing: %+v", err)
			closeC(f)
			return
		}

		_, err = bufWriter.Write([]byte(prefix))
		if err != nil {
			log.Printf("ERR: failed to write prefix: %+v", err)
			closeC(f)
			return
		}

		_, err = bufWriter.Write(js)
		if err != nil {
			log.Printf("ERR: failed to write data: %+v", err)
			closeC(f)
			return
		}

		_, err = fmt.Fprint(bufWriter, "\n")
		if err != nil {
			log.Printf("ERR: failed to write new line: %+v", err)
			closeC(f)
			return
		}
	}

	for _, v := range s.h {
		data := &fileTTLValue{Key: v.key, Value: v.expireAt}
		_, err = bufWriter.Write([]byte(ttlPrefix))
		if err != nil {
			log.Printf("ERR: failed to write prefix: %+v", err)
			closeC(f)
			return
		}

		js, err = json.Marshal(data)
		if err != nil {
			log.Printf("ERR: failed to marshal data for storing: %+v", err)
			closeC(f)
			return
		}

		_, err = bufWriter.Write(js)
		if err != nil {
			log.Printf("ERR: failed to write data: %+v", err)
			closeC(f)
			return
		}

		_, err = fmt.Fprint(bufWriter, "\n")
		if err != nil {
			log.Printf("ERR: failed to write new line: %+v", err)
			closeC(f)
			return
		}
	}

	err = bufWriter.Flush()
	if err != nil {
		log.Printf("ERR: failed to flush db data: %+v", err)
		closeC(f)
		return
	}

	err = f.Sync()
	if err != nil {
		log.Printf("ERR: failed to fsync file: %+v", err)
		closeC(f)
		return
	}

	err = f.Close()
	if err != nil {
		log.Printf("ERR: failed to close db file: %+v", err)
		return
	}

	finalPath := path.Join(s.workDir, filename)

	err = os.Rename(fullPath, finalPath)
	if err != nil {
		log.Printf("ERR: failed to move %s to %s: %v", fullPath, finalPath, err)
	}
}

// loadFromDisk loads database from disk to memory
func (s *GeddisStore) loadFromDisk() {
	if s.storeInterval == 0 {
		return
	}

	log.Printf("Loading from disk...")
	defer log.Printf("Loading from disk finished")

	s.lock.Lock()
	defer s.lock.Unlock()

	fullPath := path.Join(s.workDir, filename)
	f, err := os.Open(fullPath)
	if err != nil {
		log.Printf("ERR: failed to open store")
		return
	}
	defer closeC(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		prefix := string(line[0:1])
		switch prefix {
		case strPrefix:
			err = s.loadStr(line[1:])
		case arrPrefix:
			err = s.loadArr(line[1:])
		case mapPrefix:
			err = s.loadMap(line[1:])
		case ttlPrefix:
			err = s.loadTTL(line[1:])
		default:
			err = errors.New("invalid file")
		}

		if err != nil {
			log.Printf("ERR: failed to read db: %+v", err)
			return
		}
	}
}

func (s *GeddisStore) loadStr(input []byte) error {
	v := &fileStrValue{}
	err := json.Unmarshal(input, v)
	if err != nil {
		return err
	}

	s.m[v.Key] = v.Value
	return nil
}

func (s *GeddisStore) loadArr(input []byte) error {
	v := &fileArrValue{}
	err := json.Unmarshal(input, v)
	if err != nil {
		return err
	}

	s.m[v.Key] = v.Value
	return nil
}

func (s *GeddisStore) loadMap(input []byte) error {
	v := &fileMapValue{}
	err := json.Unmarshal(input, v)
	if err != nil {
		return err
	}

	s.m[v.Key] = v.Value
	return nil
}

func (s *GeddisStore) loadTTL(input []byte) error {
	v := &fileTTLValue{}
	err := json.Unmarshal(input, v)
	if err != nil {
		return err
	}

	heap.Push(&s.h, keyTime{key: v.Key, expireAt: v.Value})
	return nil
}
