/**
 * @File: bimap
 * @Author: Shuangpeng.Guo
 * @Date: 2021/7/8 3:45 下午
 * 双向map, k-v都唯一, 只适合少量数据使用
 */
package bimap

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var keys []interface{}
var values []interface{}

var lock = new(sync.RWMutex)
var once = new(sync.Once)
var unsafe = false

var ErrKConflict = errors.New("bimap: key conflict")
var ErrVConflict = errors.New("bimap: value conflict")
var ErrKNotExist = errors.New("bimap: key noy exist")
var ErrVNotExist = errors.New("bimap: value noy exist")
var ErrUnsafeLocked = errors.New("bimap: map unsafe locked")

func init() {
	lock = new(sync.RWMutex)
}

// 批量初始化双向map, 只允许执行一次
func Init(_keys interface{}, _values interface{}, _unsafe bool) (err error) {

	once.Do(func() {
		lock.Lock()
		defer func() {
			if r := recover(); nil != r {
				err = errors.New(fmt.Sprintln(r))
			}
			lock.Unlock()
		}()
		keys = []interface{}{}
		values = []interface{}{}
		unsafe = _unsafe

		ik := reflect.ValueOf(_keys)
		iv := reflect.ValueOf(_values)

		if ik.Len() != iv.Len() {
			err = errors.New("k-v length not consistent")
			return
		}

		kMap := make(map[interface{}]interface{})
		vMap := make(map[interface{}]interface{})

		for idx := 0; idx < ik.Len(); idx++ {
			k := ik.Index(idx).Interface()
			if _, ok := kMap[k]; ok {
				err = ErrKConflict
				return
			}
			kMap[k] = struct{}{}
			keys = append(keys, k)
		}
		for idx := 0; idx < iv.Len(); idx++ {
			v := iv.Index(idx).Interface()
			if _, ok := vMap[v]; ok {
				err = ErrVConflict
				return
			}
			vMap[v] = struct{}{}
			values = append(values, v)
		}
		return
	})

	return err
}

// 追加参数 k-v必须都没有冲突
func Add(k interface{}, v interface{}) error {
	if unsafe {
		return ErrUnsafeLocked
	}
	// 顺序遍历
	lock.Lock()
	defer lock.Unlock()
	for _, item := range keys {
		if item == k {
			return ErrKConflict
		}
	}
	for _, item := range values {
		if item == v {
			return ErrVConflict
		}
	}
	keys = append(keys, k)
	values = append(values, v)
	return nil

}

// 根据Key维度添加或更新数据，要求value不能有冲突
func AddOrUpdateByKey(k interface{}, v interface{}) error {
	if unsafe {
		return ErrUnsafeLocked
	}

	lock.Lock()
	defer lock.Unlock()

	for _, item := range values {
		if item == v {
			return ErrVConflict
		}
	}

	for idx, item := range keys {
		if item == k {
			values[idx] = v
			return nil
		}
	}

	keys = append(keys, k)
	values = append(values, v)
	return nil
}

// 根据Value维度添加或更新数据，要求Key不能有冲突
func AddOrUpdateByValue(k interface{}, v interface{}) error {
	if unsafe {
		return ErrUnsafeLocked
	}

	lock.Lock()
	defer lock.Unlock()

	for _, item := range keys {
		if item == k {
			return ErrKConflict
		}
	}

	for idx, item := range values {
		if item == k {
			keys[idx] = k
			return nil
		}
	}

	keys = append(keys, k)
	values = append(values, v)
	return nil
}

// 根据Key删除数据
func DelByKey(k interface{}) error {
	if unsafe {
		return ErrUnsafeLocked
	}

	lock.Lock()
	defer lock.Unlock()

	for idx, item := range keys {
		if item == k {
			keys = append(keys[:idx], keys[idx+1:]...)
			values = append(values[:idx], values[idx+1:]...)
			return nil
		}
	}

	return ErrKNotExist
}

// 根据value删除数据
func DelByValue(v interface{}) error {
	if unsafe {
		return ErrUnsafeLocked
	}

	lock.Lock()
	defer lock.Unlock()

	for idx, item := range values {
		if item == v {
			keys = append(keys[:idx], keys[idx+1:]...)
			values = append(values[:idx], values[idx+1:]...)
			return nil
		}
	}

	return ErrKNotExist
}

// 根据Key获取Value
func GetByKey(k interface{}) (v interface{}, err error) {
	if !unsafe {
		lock.RLock()
		defer lock.RUnlock()
	}

	for idx, item := range keys {
		if item == k {
			return values[idx], nil
		}
	}

	return nil, ErrKNotExist

}

// 根据Value获取Key
func GetByValue(v interface{}) (k interface{}, err error) {
	if !unsafe {
		lock.RLock()
		defer lock.RUnlock()
	}

	for idx, item := range values {
		if item == v {
			return keys[idx], nil
		}
	}

	return nil, ErrVNotExist
}
