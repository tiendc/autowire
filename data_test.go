package autowire

import (
	"context"
	"errors"
)

var (
	errTest1 = errors.New("errTest1")
)

//
// SERVICES
//

type ServiceBase interface {
	InitArgs() []any
}

type Service1 interface {
	ServiceBase
}

type Service2 interface {
	ServiceBase
}

type Service3 interface {
	ServiceBase
}

type Service4 interface {
	ServiceBase
}

type Service5 interface {
	ServiceBase
}

type serviceBase struct {
	initArgs []any
}

func (s serviceBase) InitArgs() []any {
	return s.initArgs
}

type service1 struct {
	serviceBase
}

type service2 struct {
	serviceBase
}

type service3 struct {
	serviceBase
}

type service4 struct {
	serviceBase
}

type service5 struct {
	serviceBase
}

// CREATE FUNCTIONS - Service 1

func NewSrv1_OK() Service1 {
	return &service1{}
}

func NewSrv1_OK_With_Need_Srv2_Srv3(s2 Service2, s3 Service3) Service1 {
	return &service1{serviceBase{initArgs: []any{s2, s3}}}
}

func NewSrv1_OK_With_Need_Srv2_Srv3_Struct1(s2 Service2, s3 Service3, st1 *Struct1_OK) Service1 {
	return &service1{serviceBase{initArgs: []any{s2, s3, st1}}}
}

func NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice(s2 Service2, s3 Service3, ints []int) Service1 {
	return &service1{serviceBase{initArgs: []any{s2, s3, ints}}}
}

func NewSrv1_OK_With_Nil_Err() (Service1, error) {
	return &service1{}, nil
}

func NewSrv1_OK_With_Need_Ctx(ctx context.Context) Service1 {
	return &service1{serviceBase{initArgs: []any{ctx}}}
}

func NewSrv1_Fail_With_Err() (Service1, error) {
	return nil, errTest1
}

func NewSrv1_Fail_With_0_Ret_Value() {
}

func NewSrv1_Fail_With_3_Ret_Values() (Service1, int, error) {
	return nil, 123, errTest1
}

func NewSrv1_Fail_With_Non_Err_At_Last() (Service1, int) {
	return service1{}, 123
}

func NewSrv1_Fail_With_Dup_Arg_Type(i1 int, s2 Service2, i2 int) (Service1, error) {
	return service1{}, nil
}

func NewSrv1_Fail_With_Variadic(s2 Service2, ii ...int) Service1 {
	return service1{}
}

func NewSrv1_Fail_Need_Srv1(s1 Service1) (Service1, error) {
	return service1{}, nil
}

// CREATE FUNCTIONS - Service 2

func NewSrv2_OK() Service2 {
	return &service2{}
}

func NewSrv2_OK_With_Need_Srv4_Srv5(s4 Service4, s5 Service5) (Service2, error) {
	return &service2{serviceBase{initArgs: []any{s4, s5}}}, nil
}

// CREATE FUNCTIONS - Service 3

func NewSrv3_OK() Service3 {
	return &service3{}
}

// CREATE FUNCTIONS - Service 4

func NewSrv4_OK() Service4 {
	return &service4{}
}

func NewSrv4_OK_With_Need_Srv1(s1 Service1) (Service4, error) {
	return &service4{serviceBase{initArgs: []any{s1}}}, nil
}

// CREATE FUNCTIONS - Service 5

func NewSrv5_OK() Service5 {
	return &service5{}
}

//
// STRUCTS
//

type Struct1_OK struct {
	Int   int
	Str   string
	Slice []int
	Map   map[string]int64
	Map2  map[string]string
}

type Struct2_Dup_Field_Type struct {
	Struct1_OK
	AnotherSlice []int
}

type Struct3_Empty struct {
}

type Struct4_OK_Dup_Field_Type_Unexported struct {
	Struct1_OK
	unexportedSlice []int
}

type Struct5_OK struct {
	IntP *int
}

type Struct6_OK_Nested_Anonymous struct {
	Int     int
	Nested1 struct {
		Int16   int16
		Str     string
		Nested2 struct {
			Slice []int
			Map   map[int]int
		}
	}
}

type Struct7_OK_Nested struct {
	Int     int
	Nested1 Nested1
}

type Nested1 struct {
	Int16   int16
	Str     string
	Nested2 *Nested2
}

type Nested2 struct {
	Slice []int
	Map   map[int]int
}

var (
	struct1_OK = Struct1_OK{
		Int:   123,
		Str:   "abc",
		Slice: []int{1, 2, 3},
		Map:   map[string]int64{"abc": 123},
		Map2:  map[string]string{"abc": "123"},
	}
	struct2_Dup_Field_Type               = Struct2_Dup_Field_Type{}
	struct3_Empty                        = Struct3_Empty{}
	struct4_OK_Dup_Field_Type_Unexported = Struct4_OK_Dup_Field_Type_Unexported{
		Struct1_OK: struct1_OK,
	}
	struct5_OK                  = Struct5_OK{}
	struct6_OK_Nested_Anonymous = Struct6_OK_Nested_Anonymous{
		Int: 123,
	}
	struct7_OK_Nested = Struct7_OK_Nested{
		Int: 123,
		Nested1: Nested1{
			Int16: 16,
			Str:   "abc",
			Nested2: &Nested2{
				Slice: []int{1, 2, 3},
				Map:   map[int]int{1: 1, 2: 2, 3: 3},
			},
		},
	}
)
