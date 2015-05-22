package storage

import (
    "errors"
)

type Object interface{}
type Storage map[string]Object

type safeMap chan command
type command struct {
    command string
    key string
    object Object
    result chan result
}
type result struct {
    value interface{}
    err error
}

var version = "0.1"
var storage Storage
var sf safeMap

func init() {
    storage = make(Storage)
    sf = make(safeMap, 100)

    go run(sf)
}

func run(safe safeMap) {
    for {
        com := <- safe
        switch com.command {
            case "get": 
                value, err := get(com.key)
                com.result <- result{value, err}
            case "set":
                value, err := set(com.key, com.object)
                com.result <- result{value, err}
            case "del":
                value, err := del(com.key)
                com.result <- result{value, err}
            case "count":
                value, err := count()
                com.result <- result{value, err}
        }
    }
}


func Set(key string, val Object) (bool, error) {
    result := send("set", key, val)
    return result.value.(bool), result.err
}
func Get(key string) (Object, bool) {
    result := send("get", key, nil)
    if result.err != nil {
        return nil, false
    }
    value, ok := result.value.(Object)
    if !ok {
        return nil, false
    }
    return value, true
}

func Del(key string) (bool, error){
    result := send("del", key, nil)
    return result.value.(bool), result.err
}
func Count() (int, error) {
    result := send("count", "", nil)
    return result.value.(int), result.err
}

func send(com, key string, val Object) result{
    sr := make(chan result)
    sf <- command{com, key, val, sr}
    return <- sr
}

func set(key string, value Object) (bool, error) {
    if len(storage) > 1000 {
        for key, _ := range storage {
            del(key)
            continue
        }
    }
    storage[key] = value
    return true, nil
}

func get(key string) (val Object, err error) {
    val, ok := storage[key]
    if !ok {
        return nil, errors.New("No Key in Storage")
    }
    return val, nil
}

func del(key string) (bool, error) {
    delete(storage, key)
    return true, nil
}

func count() (int, error) {
    return len(storage), nil
}