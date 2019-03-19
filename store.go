package SimpleStorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cstockton/go-conv"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"sync"
)

type Config struct {
	StorageDir string
}

type Storage struct {
	StorageFile string
	lock        *sync.Mutex
	data        map[string]string
	notify      []chan struct{}
}

func NewSimpleStorage(name string, conf *Config) *Storage {
	s := Storage{}

	// Chose StorageFile
	if envdir := os.Getenv("STORAGE_DIR"); envdir != "" {
		if envdir[len(envdir)-1:] == "/" { // remove trailing /
			envdir = envdir[:len(envdir)-1]
		}
		s.StorageFile = envdir + "/ss_" + name
	} else if conf != nil {
		if e := s.testDir(conf.StorageDir); e != nil {
			fmt.Errorf("Cant create Storage Directory")
			return nil
		}
		s.StorageFile = conf.StorageDir + "/storage"
	} else {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		dir := usr.HomeDir + "/." + name
		if e := s.testDir(dir); e != nil {
			fmt.Errorf("Cant create Storage Directory")
			return nil
		}
		s.StorageFile = dir + "/storage"
	}

	// make sure File exists
	if _, err := os.Stat(s.StorageFile); os.IsNotExist(err) {
		if _, err := os.Create(s.StorageFile); err != nil {
			fmt.Errorf("Cant create Storage Directory")
			return nil
		}
	}

	s.lock = &sync.Mutex{}
	s.read()
	go s.watchFileChange()
	return &s
}

func (this *Storage) GetUpdateChan() chan struct{} {
	res := make(chan struct{})
	this.notify = append(this.notify, res)
	return res
}

func (this *Storage) sendNotify() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, c := range this.notify {
		c <- struct{}{}
	}
}

func (this *Storage) watchFileChange() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()
	//

	if err := watcher.Add(this.StorageFile); err != nil {
		fmt.Println("ERROR ", err)
	}

	for {
		select {
		case <-watcher.Events:
			this.read()
			this.sendNotify()

		case err := <-watcher.Errors:
			fmt.Println("ERROR", err)
		}
	}

}

func (this *Storage) read() {
	data, err := ioutil.ReadFile(this.StorageFile)
	if err != nil {
		this.data = map[string]string{}
		return
	}
	if err = json.Unmarshal(data, &this.data); err != nil {
		this.data = map[string]string{}
		return
	}
}

func (this *Storage) save() {
	this.lock.Lock()
	data, _ := json.Marshal(&this.data)
	err := ioutil.WriteFile(this.StorageFile, data, 644)
	if err != nil {
		fmt.Println("Error saving Data to ", this.StorageFile)
	}
	this.lock.Unlock()
	this.sendNotify()
}

func (this *Storage) testDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// Setters

func (this *Storage) StoreInterface(key string, val interface{}) {
	bytes, _ := json.Marshal(val)
	this.data[key] = string(bytes)
	this.save()
}

func (this *Storage) StoreString(key, value string) {
	this.data[key] = value
	this.save()
}

func (this *Storage) StoreInt(key string, value int) {
	this.data[key] = strconv.Itoa(value)
	this.save()

}

func (this *Storage) StoreUint(key string, value uint) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreInt8(key string, value int8) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreInt32(key string, value int32) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreInt64(key string, value int64) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreuInt8(key string, value uint8) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreUint32(key string, value uint32) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreUint64(key string, value uint64) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreFloat32(key string, value float32) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreFloat64(key string, value float64) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreComplex64(key string, value complex64) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

func (this *Storage) StoreComplex128(key string, value complex128) {
	v, _ := conv.String(value)
	this.data[key] = v
	this.save()
}

// Getters

func (this *Storage) GetInterface(key string, value interface{}) error {
	if val, ok := this.data[key]; ok {
		return json.Unmarshal([]byte(val), value)
	}
	return errors.New("Key not found")
}

func (this *Storage) GetString(key string) (string, error) {
	if val, ok := this.data[key]; ok {
		return val, nil
	}
	return "", errors.New("Key not found")
}

func (this *Storage) GetInt(key string) (int, error) {
	if val, ok := this.data[key]; ok {
		return strconv.Atoi(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetInt8(key string) (int8, error) {
	if val, ok := this.data[key]; ok {
		return conv.Int8(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetInt32(key string) (int32, error) {
	if val, ok := this.data[key]; ok {
		return conv.Int32(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetInt64(key string) (int64, error) {
	if val, ok := this.data[key]; ok {
		return conv.Int64(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetUint(key string) (uint, error) {
	if val, ok := this.data[key]; ok {
		return conv.Uint(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetUint8(key string) (uint8, error) {
	if val, ok := this.data[key]; ok {
		return conv.Uint8(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetUint32(key string) (uint32, error) {
	if val, ok := this.data[key]; ok {
		return conv.Uint32(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetUint64(key string) (uint64, error) {
	if val, ok := this.data[key]; ok {
		return conv.Uint64(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetFloat32(key string) (float32, error) {
	if val, ok := this.data[key]; ok {
		return conv.Float32(val)
	}
	return 0, errors.New("Key not found")
}

func (this *Storage) GetFloat64(key string) (float64, error) {
	if val, ok := this.data[key]; ok {
		return conv.Float64(val)
	}
	return 0, errors.New("Key not found")
}
