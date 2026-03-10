package json

import (
	"reflect"
	"testing"
)

func TestMarshalToString(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		if result := MarshalToString("hello"); result != `"hello"` {
			t.Errorf(`MarshalToString("hello") = %s, want "hello"`, result)
		}
		if result := MarshalToString(123); result != `123` {
			t.Errorf("MarshalToString(123) = %s, want 123", result)
		}
		if result := MarshalToString(nil); result != `null` {
			t.Errorf("MarshalToString(nil) = %s, want null", result)
		}
	})

	t.Run("slice", func(t *testing.T) {
		result := MarshalToString([]int{1, 2, 3})
		if result != `[1,2,3]` {
			t.Errorf("MarshalToString([]int{1,2,3}) = %s, want [1,2,3]", result)
		}
	})
}

func TestCopyByJson(t *testing.T) {
	t.Run("copy map", func(t *testing.T) {
		src := map[string]int{"a": 1, "b": 2}
		var dst map[string]int
		err := CopyByJson(&dst, src)
		if err != nil {
			t.Fatalf("CopyByJson failed: %v", err)
		}
		if !reflect.DeepEqual(dst, src) {
			t.Errorf("dst = %v, want %v", dst, src)
		}
		dst["c"] = 3
		if _, exists := src["c"]; exists {
			t.Error("src should not be affected by dst modification")
		}
	})

	t.Run("copy struct", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		src := Person{Name: "Alice", Age: 30}
		var dst Person
		err := CopyByJson(&dst, src)
		if err != nil {
			t.Fatalf("CopyByJson failed: %v", err)
		}
		if dst != src {
			t.Errorf("dst = %v, want %v", dst, src)
		}
	})

	t.Run("nil dst", func(t *testing.T) {
		err := CopyByJson(nil, map[string]int{"a": 1})
		if err == nil {
			t.Error("expected error for nil dst")
		}
	})

	t.Run("nil src", func(t *testing.T) {
		var dst map[string]int
		err := CopyByJson(&dst, nil)
		if err == nil {
			t.Error("expected error for nil src")
		}
	})
}

func TestMapToStruct(t *testing.T) {
	type Config struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	t.Run("valid conversion", func(t *testing.T) {
		src := map[string]interface{}{"host": "localhost", "port": 8080}
		var dst Config
		err := MapToStruct(&dst, src)
		if err != nil {
			t.Fatalf("MapToStruct failed: %v", err)
		}
		if dst.Host != "localhost" || dst.Port != 8080 {
			t.Errorf("dst = %+v, want {Host:localhost, Port:8080}", dst)
		}
	})
}

func TestCopy(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	t.Run("copy struct", func(t *testing.T) {
		src := User{Name: "Bob", Age: 25}
		var dst User
		err := Copy(&dst, src)
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}
		if dst != src {
			t.Errorf("dst = %v, want %v", dst, src)
		}
	})
}

func TestMapCopy(t *testing.T) {
	src := map[string]interface{}{"a": 1, "b": "hello"}
	dst := MapCopy(src)

	if !reflect.DeepEqual(dst, src) {
		t.Errorf("dst = %v, want %v", dst, src)
	}

	dst["c"] = 3
	if _, exists := src["c"]; exists {
		t.Error("src should not be affected by dst modification (shallow copy issue)")
	}

	// Note: MapCopy is shallow, nested maps will share references
	srcNested := map[string]interface{}{"a": map[string]int{"x": 1}}
	dstNested := MapCopy(srcNested)
	dstNested["a"].(map[string]int)["x"] = 999
	if srcNested["a"].(map[string]int)["x"] != 999 {
		t.Log("MapCopy is shallow copy - nested maps share references")
	}
}

func TestDeepCopyMap(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		src := map[string]interface{}{"a": 1, "b": "hello"}
		dst := DeepCopyMap(src)

		if !reflect.DeepEqual(dst, src) {
			t.Errorf("dst = %v, want %v", dst, src)
		}

		dst["c"] = 3
		if _, exists := src["c"]; exists {
			t.Error("src should not be affected by dst modification")
		}
	})

	t.Run("nested map", func(t *testing.T) {
		src := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{"c": 2},
		}
		dst := DeepCopyMap(src)

		dst["d"] = 3
		dstB := dst["b"].(map[string]interface{})
		dstB["c"] = 4

		// src should remain unchanged
		if src["a"] != 1 {
			t.Errorf("src[a] = %v, want 1", src["a"])
		}
		if src["b"].(map[string]interface{})["c"] != 2 {
			t.Errorf("src[b][c] = %v, want 2", src["b"].(map[string]interface{})["c"])
		}
	})

	t.Run("with slice", func(t *testing.T) {
		src := map[string]interface{}{
			"a": []int{1, 2, 3},
		}
		dst := DeepCopyMap(src)

		dst["a"].([]int)[0] = 999
		if src["a"].([]int)[0] != 1 {
			t.Error("slice in src should not be modified")
		}
	})

	t.Run("with pointer", func(t *testing.T) {
		val := 42
		src := map[string]interface{}{
			"ptr": &val,
		}
		dst := DeepCopyMap(src)

		// Verify the pointer was copied
		dstPtr := dst["ptr"].(*int)
		*dstPtr = 999
		if *src["ptr"].(*int) != 42 {
			t.Error("pointer in src should not be modified")
		}
	})

	t.Run("with struct", func(t *testing.T) {
		type Point struct {
			X, Y int
		}
		src := map[string]interface{}{
			"point": Point{X: 1, Y: 2},
		}
		dst := DeepCopyMap(src)

		dst["point"] = Point{X: 999, Y: 999}
		if src["point"].(Point).X != 1 {
			t.Error("struct in src should not be modified")
		}
	})
}

func TestDeepCopySlice(t *testing.T) {
	src := []int{1, 2, 3}
	dst := DeepCopySlice(src).([]int)

	if !reflect.DeepEqual(dst, src) {
		t.Errorf("dst = %v, want %v", dst, src)
	}

	dst[0] = 999
	if src[0] != 1 {
		t.Error("src should not be affected by dst modification")
	}
}

func TestDeepCopyPointer(t *testing.T) {
	t.Run("int pointer", func(t *testing.T) {
		val := 42
		src := &val
		dst := DeepCopyPointer(src).(*int)

		if *dst != *src {
			t.Errorf("*dst = %d, want %d", *dst, *src)
		}

		*dst = 999
		if *src != 42 {
			t.Error("src should not be affected by dst modification")
		}
	})
}

func TestDeepCopyStruct(t *testing.T) {
	type Point struct {
		X, Y int
	}

	src := Point{X: 1, Y: 2}
	dst := DeepCopyStruct(src).(Point)

	if dst != src {
		t.Errorf("dst = %v, want %v", dst, src)
	}

	dst.X = 999
	if src.X != 1 {
		t.Error("src should not be affected by dst modification")
	}
}

func TestDeepCopyGob(t *testing.T) {
	t.Run("copy map", func(t *testing.T) {
		src := map[string]int{"a": 1, "b": 2}
		var dst map[string]int
		err := DeepCopyGob(src, &dst)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		if !reflect.DeepEqual(dst, src) {
			t.Errorf("dst = %v, want %v", dst, src)
		}
	})

	t.Run("copy struct", func(t *testing.T) {
		type Data struct {
			Value int
		}
		src := Data{Value: 42}
		var dst Data
		err := DeepCopyGob(src, &dst)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		if dst != src {
			t.Errorf("dst = %v, want %v", dst, src)
		}
	})

	t.Run("encode error", func(t *testing.T) {
		// Channel cannot be gob encoded
		type Bad struct {
			Ch chan int
		}
		src := Bad{Ch: make(chan int)}
		var dst Bad
		err := DeepCopyGob(src, &dst)
		if err == nil {
			t.Error("expected error for channel field")
		}
	})
}