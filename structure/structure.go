package structure

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// 模型存储 start
var models = newModel()

func newModel() *ModelManager {
	return &ModelManager{
		models: make(map[string]*Intermediate),
		lock:   new(sync.RWMutex),
	}
}

type ModelManager struct {
	models map[string]*Intermediate
	lock   *sync.RWMutex
}

func (this *ModelManager) Get(name string) *Intermediate {
	this.lock.RLock()
	inter, ok := this.models[name]
	if !ok {
		inter = nil
	}
	this.lock.RUnlock()
	return inter
}

func (this *ModelManager) Set(name string, inter *Intermediate) bool {
	this.lock.Lock()
	r := true
	if val, ok := this.models[name]; !ok {
		this.models[name] = inter
	} else if val != inter {
		this.models[name] = inter
	} else {
		r = false
	}
	this.lock.Unlock()
	return r
}

// 模型存储 end

type Intermediate struct {
	Pk     string                   // 查询依赖值(主键名)
	PkZero interface{}              // 主键零值
	KindKV []KindKV                 // 入参struct Index 对应出参Name
	Type   reflect.Type             // 类型
	Slaves map[string]*Intermediate // 从节点值
}

type KindKV struct {
	Name  string
	Index int
}

const (
	tagName = "gsp"
	tagPk   = "pk"
)

/******************************************************************************
 **函数名称: ReflectStructs
 **功    能: 平铺数据结构化
 **输入参数:
 **    in: 输入数据, 格式要求：ptr -> slice -> ptr -> struct
 **    out: 输出数据, 格式要求，ptr -> slice -> ptr -> struct -> ptr -> struct
 **输出参数:
 **    err: panic信息
 **返    回:
 **实现描述:
 **    通过反射读取输入数据的idx.name.tag等信息,通过tag中制定的规则将对应Value set
 **    到输出结构中.
 **注意事项:
 **    本方法目前未做更多的格式适配, 请严格按照demo中入参格式设置，任何约定外操作都
 **    可能会引起非预期panic, 另当主键为当前类型0值时会被忽略
 **作    者: # gsp # 2020-06-22 14:27:03 #
 ******************************************************************************/
func Structure(in, out interface{}) (err error) {
	defer func() {
		if r := recover(); nil != r {
			err = errors.New(fmt.Sprintln(r))
		}
	}()

	// 反射句柄
	iValue := reflect.ValueOf(in)
	iType := iValue.Type()
	oValue := reflect.ValueOf(out)
	oType := oValue.Type()

	// 严格规范格式
	if iType.Kind() != reflect.Ptr ||
		iType.Elem().Kind() != reflect.Slice ||
		iType.Elem().Elem().Kind() != reflect.Ptr ||
		iType.Elem().Elem().Elem().Kind() != reflect.Struct ||
		oType.Kind() != reflect.Ptr ||
		oType.Elem().Kind() != reflect.Slice ||
		oType.Elem().Elem().Kind() != reflect.Ptr ||
		oType.Elem().Elem().Elem().Kind() != reflect.Struct ||
		!oValue.Elem().CanSet() {
		panic(errors.New("<Reflect>: type error"))
	}

	modelName := iType.Elem().Elem().Elem().PkgPath() + "." + iType.Elem().Elem().Elem().Name()

	master := models.Get(modelName)
	if nil == master {
		master = &Intermediate{
			Type:   oType.Elem(),
			Slaves: make(map[string]*Intermediate),
		}

		// 从入参tag中获取字段名
		inSTyp := iType.Elem().Elem().Elem()
		for idx := 0; idx < inSTyp.NumField(); idx++ {
			filed, isPK := getTag(inSTyp.Field(idx).Tag.Get(tagName))
			buildIntermediate(filed, master, isPK, inSTyp, idx)
		}
		models.Set(modelName, master)
	}

	structureHandle(master, iValue.Elem(), oValue.Elem())

	return nil
}

func buildIntermediate(filed []string, master *Intermediate,
	isPk bool, inSType reflect.Type, idx int) {
	l := len(filed)
	if l == 0 {
		return
	}
	if l == 1 {
		if isPk {
			master.Pk = inSType.Field(idx).Name
			master.PkZero = reflect.New(inSType.Field(idx).Type).Elem().Interface()
		}
		master.KindKV = append(master.KindKV, KindKV{Name: filed[0], Index: idx})
		return
	}

	if l > 1 {
		if _, ok := master.Slaves[filed[0]]; !ok {
			slaveTyp, ok := master.Type.Elem().Elem().FieldByName(filed[0])
			if !ok {
				panic(errors.New("<Reflect>: slave key not found"))
			}
			master.Slaves[filed[0]] = &Intermediate{
				Slaves: make(map[string]*Intermediate),
				Type:   slaveTyp.Type,
			}
		}
		buildIntermediate(filed[1:], master.Slaves[filed[0]], isPk, inSType, idx)
	}
}

// 解析Tag
// tag格式 "一级从表1___一级从表2__field, pk"
func getTag(str string) (filed []string, isPk bool) {
	str = strings.Replace(str, " ", "", -1)
	_arg1 := strings.Split(str, ",")
	if len(_arg1) >= 2 && _arg1[1] == tagPk {
		isPk = true
	}

	return strings.Split(_arg1[0], "__"), isPk
}

// 递归处理函数
func structureHandle(master *Intermediate, raw, out reflect.Value) {

	res := make(map[interface{}]reflect.Value)             // 结果集（无序, 主key为pk对应值）
	var resSort []interface{}                              // 维护结果集有序关系
	raws := make(map[interface{}]map[string]reflect.Value) // 子节点入参

	// 遍历入参, 取参数值set到出参
	for idx := 0; idx < raw.Len(); idx++ {
		pkId := raw.Index(idx).Elem().FieldByName(master.Pk).Interface()
		if pkId == master.PkZero {
			// TODO 如果主键为该类型零值则忽略
			continue
		}
		if _, ok := res[pkId]; !ok {
			v := reflect.New(master.Type.Elem().Elem())
			// 设置master值
			for _, kindKV := range master.KindKV {
				v.Elem().FieldByName(
					kindKV.Name).Set(raw.Index(idx).Elem().Field(kindKV.Index))
			}

			// 设置slave key & slave 原始数据
			raws[pkId] = make(map[string]reflect.Value)
			for sKey, slave := range master.Slaves {
				v.Elem().FieldByName(sKey).Set(reflect.New(slave.Type).Elem())
				raws[pkId][sKey] = reflect.New(raw.Type()).Elem()
			}

			res[pkId] = v
			resSort = append(resSort, pkId)
		}
		// 设置从表原始值
		for sKey, _ := range master.Slaves {
			raws[pkId][sKey] = reflect.Append(raws[pkId][sKey], raw.Index(idx))
		}
	}

	// 递归处理从节点
	for _, pkId := range resSort {
		for sKey, sValue := range master.Slaves {
			structureHandle(
				sValue, raws[pkId][sKey], res[pkId].Elem().FieldByName(sKey))
		}
	}

	resSlicePtr := reflect.New(master.Type)
	resSlice := resSlicePtr.Elem()
	for _, pkId := range resSort {
		resSlice = reflect.Append(resSlice, res[pkId])
	}

	out.Set(resSlice)
}
