package c

import (
	"fmt"
	"math"
	"runtime/cgo"
	"testing"
	"unsafe"
)

/*
#include "codec.h"
extern dt_error c_write_size_callback(dt_encoder *e, void *user_data);
extern dt_error c_read_size_callback(dt_decoder *d, void *user_data);
extern dt_error c_read(uint8_t *b, uint64_t n, void *user_data);
*/
import "C"

func testCodec(t *testing.T) {
	testValues(t, func(e *C.dt_encoder, x C.bool) C.dt_error { return C.dt_ok }, func(d *C.dt_decoder, x *C.bool) C.dt_error { return C.dt_ok }, false)
	testValues(t, func(e *C.dt_encoder, x C.bool) C.dt_error { return C.dt_write_bool(e, x) }, func(d *C.dt_decoder, x *C.bool) C.dt_error { return C.dt_read_bool(d, x) }, false, true)
	testValues(t, func(e *C.dt_encoder, x C.uint32_t) C.dt_error { return C.dt_write_uint32(e, x) }, func(d *C.dt_decoder, x *C.uint32_t) C.dt_error { return C.dt_read_uint32(d, x) }, 0, 1, 2, 3, math.MaxUint32)
	testValues(t, func(e *C.dt_encoder, x C.uint64_t) C.dt_error { return C.dt_write_uint64(e, x) }, func(d *C.dt_decoder, x *C.uint64_t) C.dt_error { return C.dt_read_uint64(d, x) }, 0, 1, 2, 3, math.MaxUint64)
	testValues(t, func(e *C.dt_encoder, x C.int64_t) C.dt_error { return C.dt_write_int64(e, x) }, func(d *C.dt_decoder, x *C.int64_t) C.dt_error { return C.dt_read_int64(d, x) }, 0, 1, 2, 3, math.MaxInt64, -1, -2, -3, math.MinInt64)
	testValues(t, func(e *C.dt_encoder, x C.uint64_t) C.dt_error { return C.dt_write_var_uint(e, x) }, func(d *C.dt_decoder, x *C.uint64_t) C.dt_error { return C.dt_read_var_uint(d, x) }, 0, 1, 2, 3, math.MaxUint64)
	testValues(t, func(e *C.dt_encoder, x C.int64_t) C.dt_error { return C.dt_write_var_int(e, x) }, func(d *C.dt_decoder, x *C.int64_t) C.dt_error { return C.dt_read_var_int(d, x) }, 0, 1, 2, 3, math.MaxInt64, -1, -2, -3, math.MinInt64)
	testValues(t, func(e *C.dt_encoder, x C.uint64_t) C.dt_error { return C.dt_write_var_uint_4bit(e, x) }, func(d *C.dt_decoder, x *C.uint64_t) C.dt_error { return C.dt_read_var_uint_4bit(d, x) }, 0, 1, 2, 3, math.MaxUint64)
	testValues(t, func(e *C.dt_encoder, x C.float) C.dt_error { return C.dt_write_float32(e, x) }, func(d *C.dt_decoder, x *C.float) C.dt_error { return C.dt_read_float32(d, x) }, 0, 1, 2, 3, math.MaxFloat32, -1, -2, -3, -math.MaxFloat32)
	testValues(t, func(e *C.dt_encoder, x C.double) C.dt_error { return C.dt_write_float64(e, x) }, func(d *C.dt_decoder, x *C.double) C.dt_error { return C.dt_read_float64(d, x) }, 0, 1, 2, 3, math.MaxFloat64, -1, -2, -3, -math.MaxFloat64)
	testValues(t, func(e *C.dt_encoder, x string) C.dt_error {
		return C.dt_write_bytes(e, C.dt_string{
			data: (*C.uint8_t)(unsafe.StringData(x)),
			len:  C.uint64_t(len(x)),
		})
	}, func(d *C.dt_decoder, x *string) C.dt_error {
		var b C.dt_string
		err := C.dt_read_bytes(d, &b)
		*x = unsafe.String((*byte)(b.data), b.len)
		return err
	}, "", "a", "847fqh938", "øˍðº,ßœ≥p«®£ª¢º˝ð-")
}

func testValues[T comparable](t *testing.T, write func(*C.dt_encoder, T) C.dt_error, read func(*C.dt_decoder, *T) C.dt_error, x ...T) {
	t.Helper()
	for offset := range [8]int{} {
		for _, x := range x {
			e := C.dt_new_encoder()
			h := cgo.NewHandle(func(e *C.dt_encoder) C.dt_error {
				for i := offset; i > 0; i-- {
					if err := C.dt_write_bool(e, true); err != C.dt_ok {
						return err
					}
				}
				return write(e, x)
			})
			if err := C.dt_write_size(&e, (*C.dt_write_size_fn)(C.c_write_size_callback), unsafe.Pointer(&h)); err != C.dt_ok {
				t.Fatalf("write_size failed (offset=%d, value=%#v): %v", offset, x, err)
			}
			buf = unsafe.Slice(C.dt_encoder_bytes(&e), C.dt_encoder_len(&e))

			fmt.Printf("buf: %x\n", buf)

			d := C.dt_new_decoder((*C.dt_read_fn)(C.c_read), nil)
			d.n = e.n
			var y T
			h = cgo.NewHandle(func(d *C.dt_decoder) C.dt_error {
				for i := offset; i > 0; i-- {
					var b bool
					if err := C.dt_read_bool(d, (*C.bool)(&b)); err != C.dt_ok {
						return err
					}
					if !b {
						t.Fatalf("expected true from read_bool (offset=%d, value=%#v)", offset, x)
					}
				}
				return read(d, &y)
			})
			if err := C.dt_read_size(&d, (*C.dt_read_size_fn)(C.c_read_size_callback), unsafe.Pointer(&h)); err != C.dt_ok {
				t.Fatalf("read_size failed (offset=%d, value=%#v): %v", offset, x, err)
			}
			if n := d.n - d.j; n != 0 {
				t.Fatalf("%d bits remaining (offset=%d, value=%#v", n, offset, x)
			}
			if x != y {
				t.Fatalf("expected %#v, got %#v (offset %d)", x, y, offset)
			}
			C.dt_delete_encoder(&e)
		}
	}
}

func testReadSizeSkipsExtraBits(t *testing.T) {
	e := C.dt_new_encoder()
	h := cgo.NewHandle(func(e *C.dt_encoder) C.dt_error {
		for i := 0; i < 37; i++ {
			C.dt_write_bool(e, false)
		}
		return C.dt_ok
	})
	if err := C.dt_write_size(&e, (*C.dt_write_size_fn)(C.c_write_size_callback), unsafe.Pointer(&h)); err != C.dt_ok {
		t.Fatalf("write_size failed: %v", err)
	}
	C.dt_write_bool(&e, true)
	buf = unsafe.Slice(C.dt_encoder_bytes(&e), C.dt_encoder_len(&e))

	fmt.Printf("buf: %x\n", buf)

	d := C.dt_new_decoder((*C.dt_read_fn)(C.c_read), nil)
	d.n = e.n
	h = cgo.NewHandle(func(d *C.dt_decoder) C.dt_error {
		for i := 0; i < 11; i++ {
			var b bool
			if err := C.dt_read_bool(d, (*C.bool)(&b)); err != C.dt_ok {
				t.Fatalf("read_bool failed: %v", err)
			}
			if b {
				t.Fatalf("expected false from read_bool")
			}
		}
		return C.dt_ok
	})
	if err := C.dt_read_size(&d, (*C.dt_read_size_fn)(C.c_read_size_callback), unsafe.Pointer(&h)); err != C.dt_ok {
		t.Fatalf("read_size failed: %v", err)
	}
	var b bool
	if err := C.dt_read_bool(&d, (*C.bool)(&b)); err != C.dt_ok {
		t.Fatalf("read_bool failed: %v", err)
	}
	if !b {
		t.Fatalf("expected true from read_bool")
	}
	if n := d.n - d.j; n != 0 {
		t.Fatalf("%d bits remaining", n)
	}
	C.dt_delete_encoder(&e)
}

//export go_write_size_callback
func go_write_size_callback(e *C.dt_encoder, user_data unsafe.Pointer) C.dt_error {
	h := *(*cgo.Handle)(user_data)
	defer h.Delete()
	return h.Value().(func(*C.dt_encoder) C.dt_error)(e)
}

//export go_read_size_callback
func go_read_size_callback(d *C.dt_decoder, user_data unsafe.Pointer) C.dt_error {
	h := *(*cgo.Handle)(user_data)
	defer h.Delete()
	return h.Value().(func(*C.dt_decoder) C.dt_error)(d)
}

var buf []C.uint8_t

//export go_read
func go_read(b *C.uint8_t, n C.uint64_t) C.dt_error {
	dst := unsafe.Slice(b, n)
	if len(buf) < len(dst) {
		return C.dt_error_unexpected_eof
	}
	copy(dst, buf[:n])
	buf = buf[n:]
	return C.dt_ok
}
