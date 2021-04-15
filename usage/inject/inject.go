package inject

// 依赖注入/控制反转 golang实现

import (
	"fmt"
	"reflect"
)

type Injector interface {
	Applicator
	Invoker
	TypeMapper
	SetParent(Injector)
}

// 注入struct
type Applicator interface {
	Apply(interface{}) error
}

// 执行
type Invoker interface {
	Invoke(interface{}) ([]reflect.Value, error)
}

type TypeMapper interface {
	// 注入参数
	Map(interface{}) TypeMapper
	MapTo(val interface{}, ifacePtr interface{}) TypeMapper

	// 获取被注入的参数
	Get(p reflect.Type) reflect.Value
}

type injector struct {
	values map[reflect.Type]reflect.Value
	parent Injector
}

func New() Injector {
	return &injector{
		values: make(map[reflect.Type]reflect.Value),
	}
}

// 对struct的字段进行注入，要求字段必须是可导出的（首字母大写），并且tag设置为`inject`
// val: 底层类型为结构体的指针；
func (inj *injector) Apply(val interface{}) error {
	v := reflect.ValueOf(val)

	for v.Kind() == reflect.Ptr { // 指针转换为指向的变量
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 只对struct进行操作
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		structField := t.Field(i)
		if f.CanSet() && structField.Tag == "inject" {
			ft := f.Type()
			v := inj.Get(ft)
			if !v.IsValid() {
				return fmt.Errorf("Value not found for type %v ", ft)
			}
			f.Set(v)
		}
	}

	return nil
}

// 动态执行函数, 执行前可以通过Map或MapTo来注入参数
// fn: 执行的函数,底层类型必须为func
func (inj *injector) Invoke(fn interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(fn)

	for r, value := range inj.values {
		fmt.Printf("key: %v value: %v \n", r, value)
	}

	// Panic if t is not kind of func
	var in = make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		argType := t.In(i)
		val := inj.Get(argType)
		if !val.IsValid() {
			return nil, fmt.Errorf("Value not found for type %v ", argType)
		}
		in[i] = val
	}

	return reflect.ValueOf(fn).Call(in), nil
}

// 将i2的type与value保存在i.values中
func (inj *injector) Map(i2 interface{}) TypeMapper {
	inj.values[reflect.TypeOf(i2)] = reflect.ValueOf(i2)
	return inj
}

// ifacePtr：接口类型的指针，指定特定的类型当i.values的key；
// 其他内容同i.Map()
func (inj *injector) MapTo(val interface{}, ifacePtr interface{}) TypeMapper {
	inj.values[interfaceOf(ifacePtr)] = reflect.ValueOf(val)
	return inj
}

// 获取type
func (inj *injector) Get(p reflect.Type) reflect.Value {
	val := inj.values[p]
	if !val.IsValid() && inj.parent != nil {
		val = inj.parent.Get(p)
	}
	return val
}

func (inj *injector) SetParent(parent Injector) {
	inj.parent = parent
}

// value: 接口类型的指针;
func interfaceOf(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr { // 用for循环解指针，直到t为非指针类型
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("Called inject.interfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)")
	}

	return t
}
