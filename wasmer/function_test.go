package wasmer

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

func TestInstanceFunctionCall(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (type $sum_t (func (param i32 i32) (result i32)))
			  (func $sum_f (type $sum_t) (param $x i32) (param $y i32) (result i32)
			    local.get $x
			    local.get $y
			    i32.add)
			  (export "sum" (func $sum_f)))
		`),
	)
	defer module.Close()

	assert.NoError(t, err)

	instance, err := NewInstance(module)
	defer instance.Close()

	assert.NoError(t, err)

	sum, err := instance.Exports.GetFunction("sum")

	assert.NoError(t, err)

	result, err := sum(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, result, int32(3))
}

func TestInstanceFunctionCallReturnZeroValue(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (type $test_t (func (param i32 i32)))
			  (func $test_f (type $test_t) (param $x i32) (param $y i32))
			  (export "test" (func $test_f)))
		`),
	)
	defer module.Close()

	assert.NoError(t, err)

	instance, err := NewInstance(module)
	defer instance.Close()

	assert.NoError(t, err)

	test, err := instance.Exports.GetFunction("test")

	assert.NoError(t, err)

	result, err := test(1, 2)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestInstanceFunctionCallReturnMultipleValues(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (type $swap_t (func (param i32 i64) (result i64 i32)))
			  (func $swap_f (type $swap_t) (param $x i32) (param $y i64) (result i64 i32)
			    local.get $y
			    local.get $x)
			  (export "swap" (func $swap_f)))
		`),
	)
	defer module.Close()

	assert.NoError(t, err)

	instance, err := NewInstance(module)
	instance.Close()

	assert.NoError(t, err)

	swap, err := instance.Exports.GetFunction("swap")

	assert.NoError(t, err)

	results, err := swap(41, 42)

	assert.NoError(t, err)
	assert.Equal(t, results, []interface{}{int64(42), int32(41)})
}

func getTestBytes() []byte {
	_, filename, _, _ := runtime.Caller(0)
	modulePath := path.Join(path.Dir(filename), "testdata", "tests.wasm")
	bytes, _ := ioutil.ReadFile(modulePath)

	return bytes
}

func testFunction(t *testing.T) *Instance {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(store, getTestBytes())
	defer module.Close()

	assert.NoError(t, err)

	instance, err := NewInstance(module)
	defer instance.Close()

	assert.NoError(t, err)

	return instance
}

func TestFunctionSum(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("sum")
	result, err := f(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, result, int32(3))
}

func TestFunctionArity0(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("arity_0")
	result, err := f()

	assert.NoError(t, err)
	assert.Equal(t, result, int32(42))
}

func TestFunctionI32I32(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("i32_i32")
	result, err := f(7)

	assert.NoError(t, err)
	assert.Equal(t, result, int32(7))

	result, _ = f(int8(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(uint8(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(int16(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(uint16(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(int32(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(int(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(uint(7))
	assert.Equal(t, result, int32(7))
}

func TestFunctionI64I64(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("i64_i64")
	result, err := f(7)

	assert.NoError(t, err)
	assert.Equal(t, result, int64(7))

	result, _ = f(int8(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(uint8(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int16(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(uint16(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int32(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int64(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(uint(7))
	assert.Equal(t, result, int64(7))
}

func TestFunctionF32F32(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("f32_f32")
	result, err := f(float32(7.42))

	assert.NoError(t, err)
	assert.Equal(t, result, float32(7.42))
}

func TestFunctionF64F64(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("f64_f64")
	result, err := f(7.42)

	assert.NoError(t, err)
	assert.Equal(t, result, float64(7.42))

	result, _ = f(float64(7.42))
	assert.Equal(t, result, float64(7.42))
}

func TestFunctionI32I64F32F64F64(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("i32_i64_f32_f64_f64")
	result, err := f(1, 2, float32(3.4), 5.6)

	assert.NoError(t, err)
	assert.Equal(t, float64(int(result.(float64)*10000000))/10000000, 1+2+3.4+5.6)
}

func TestFunctionBoolCastedtoI32(t *testing.T) {
	instance := testFunction(t)

	f, _ := instance.Exports.GetFunction("bool_casted_to_i32")
	result, err := f()

	assert.NoError(t, err)
	assert.Equal(t, result, int32(1))
}
