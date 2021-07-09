/**
 * @File: bimap_test.go
 * @Author: Shuangpeng.Guo
 * @Date: 2021/7/8 5:45 下午
 */
package bimap

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		_keys   interface{}
		_values interface{}
		_unsafe bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"T1", args{[]string{"张三", "李四", "王五"}, []string{"zhangsan", "lisi", "wangwu"}, false}, false},
		{"T2", args{[]string{"张三", "李四", "王五"}, []string{"zhangsan", "lisi", "wangwu", "zhaoliu"}, false}, true},
		{"T3", args{[]string{"张三", "李四", "王五"}, []string{"zhangsan", "lisi", "zhangsan"}, false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(tt.args._keys, tt.args._values, tt.args._unsafe); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		k interface{}
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"T1", args{"张三", "ZhangSan"}, false},
		{"T2", args{"李四", "LiSi"}, false},
		{"T3", args{"张三", "zhangSan2"}, true},
		{"T4", args{"王五", "LiSi"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Add(tt.args.k, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetByKey(t *testing.T) {
	_ = Init([]string{"张三", "李四", "王五"}, []string{"zhangsan", "lisi", "wangwu"}, true)

	type args struct {
		k interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantV   interface{}
		wantErr bool
	}{
		{"T1", args{"张三"}, "zhangsan", false},
		{"T2", args{"赵六"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, err := GetByKey(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("GetByKey() gotV = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestGetByValue(t *testing.T) {
	_ = Init([]string{"张三", "李四", "王五"}, []string{"zhangsan", "lisi", "wangwu"}, true)

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantK   interface{}
		wantErr bool
	}{
		{"T1", args{"zhangsan"}, "张三", false},
		{"T2", args{"zhaoliu"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotK, err := GetByValue(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("GetByValue() gotK = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}

func TestAddOrUpdateByKey(t *testing.T) {
	type args struct {
		k interface{}
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddOrUpdateByKey(tt.args.k, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("AddOrUpdateByKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddOrUpdateByValue(t *testing.T) {
	type args struct {
		k interface{}
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddOrUpdateByValue(tt.args.k, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("AddOrUpdateByValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelByKey(t *testing.T) {
	type args struct {
		k interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DelByKey(tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("DelByKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelByValue(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DelByValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DelByValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
