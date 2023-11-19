#include "codec.h"

// type Value interface {
// 	Write(*dt_encoder)
// 	Read(*dt_decoder) error
// }

// Write(w io.Writer, v Value) error {
// 	e := Newdt_encoder()
// 	v.Write(e)
// 	_, err := w.Write(e->Bytes())
// 	return err;
// }

// Read(r io.Reader, v Value) error {
// 	return v.Read(Newdt_decoder(r))
// }

dt_encoder dt_new_encoder() {
	dt_encoder e = {
		.b   = calloc(8, 1),
		.len = 8,
	};
	return e;
}

dt_decoder dt_new_decoder(dt_read_fn r, void *user_data) {
	dt_decoder d = {
		.r         = r,
		.user_data = user_data,
		.n         = UINT64_MAX,
	};
	return d;
}

void dt_delete_encoder(dt_encoder *e) {
	free(e->b);
	e->b = NULL;
	e->n = 0;
	e->len = 0;
}

static uint64_t dt_remaining(dt_decoder *d) { return d->n - d->j; }

#define ei (e->n % 8)
#define ex ((uint64_t*)&e->b[e->n/8])

#define di (d->j % 8)
#define dx (*(uint64_t*)&d->b[0])

static dt_error dt_grow(dt_encoder *e, uint64_t n) {
	uint8_t *b = realloc(e->b, e->len + n);
	if (b == NULL) {
		return dt_error_out_of_memory;
	}
	memset(&b[e->len], 0, n);
	e->b = b;
	e->len += n;
	return dt_ok;
}

static dt_error dt_read(dt_decoder *d, uint64_t n) {
	uint64_t start = (di + 7) / 8;
	return d->r(&d->b[start], n, d->user_data);
}

dt_error dt_write_bool(dt_encoder *e, bool x) {
	if (x) {
		*ex |= 1 << ei;
	}
	e->n++;
	if (ei == 0) {
		dt_error err = dt_grow(e, 1);
		if (err != dt_ok) {
			return err;
		}
	}
	return dt_ok;
}

dt_error dt_read_bool(dt_decoder *d, bool *x) {
	if (dt_remaining(d) < 1) {
		return dt_error_unexpected_eof;
	}
	if (di == 0) {
		dt_error err = dt_read(d, 1);
		if (err != dt_ok) {
			return err;
		}
	}
	*x = ((dx>>di)&1) != 0;
	d->j++;
	return dt_ok;
}

dt_error dt_write_uint32(dt_encoder *e, uint32_t x) {
	*ex |= ((uint64_t)x) << ei;
	e->n += 32;
	dt_error err = dt_grow(e, 4);
	if (err != dt_ok) {
		return err;
	}
	return dt_ok;
}

dt_error dt_read_uint32(dt_decoder *d, uint32_t *x) {
	if (dt_remaining(d) < 32) {
		return dt_error_unexpected_eof;
	}
	dt_error err = dt_read(d, 4);
	if (err != dt_ok) {
		return err;
	}
	*x = dx >> di;
	d->j += 32;
	d->b[0] = d->b[4];
	return dt_ok;
}

dt_error dt_write_uint64(dt_encoder *e, uint64_t x) {
	*ex |= x << ei;
	e->n += 64;
	dt_error err = dt_grow(e, 8);
	if (err != dt_ok) {
		return err;
	}
	*ex |= x >> (64 - ei);
	return dt_ok;
}

dt_error dt_read_uint64(dt_decoder *d, uint64_t *x) {
	if (dt_remaining(d) < 64) {
		return dt_error_unexpected_eof;
	}
	dt_error err = dt_read(d, 8);
	if (err != dt_ok) {
		return err;
	}
	*x = dx >> di;
	d->j += 64;
	d->b[0] = d->b[8];
	*x |= dx << (64 - di);
	return dt_ok;
}

dt_error dt_write_int64(dt_encoder *e, int64_t x) { return dt_write_uint64(e, x); }
dt_error dt_read_int64(dt_decoder *d, int64_t *x) {
	uint64_t i = 0;
	dt_error err = dt_read_uint64(d, &i);
	*x = i;
	return err;
}

dt_error dt_write_var_uint(dt_encoder *e, uint64_t x) {
	while (true) {
		uint64_t y = x & ((1<<7) - 1);
		x >>= 7;
		if (x != 0) {
			y |= 1 << 7;
		}
		*ex |= y << ei;
		e->n += 8;
		dt_error err = dt_grow(e, 1);
		if (err != dt_ok) {
			return err;
		}
		if (x == 0) {
			break;
		}
	}
	return dt_ok;
}

dt_error dt_read_var_uint(dt_decoder *d, uint64_t *x) {
	for (int shift = 0; shift < 64; shift += 7) {
		if (dt_remaining(d) < 8) {
			return dt_error_unexpected_eof;
		}
		dt_error err = dt_read(d, 1);
		if (err != dt_ok) {
			return err;
		}

		uint64_t y = dx >> di;
		d->j += 8;
		d->b[0] = d->b[1];
		*x |= (y & ((1<<7) - 1)) << shift;
		if ((y&(1<<7)) == 0) {
			return dt_ok;
		}
	}
	return dt_error_var_int_overflow;
}

static uint64_t zigzag(int64_t x) { return (x >> 63) ^ (x << 1); }
static int64_t zagzig(uint64_t x) { return (x >> 1) ^ -(x & 1); }

dt_error dt_write_var_int(dt_encoder *e, int64_t x) { return dt_write_var_uint(e, zigzag(x)); }

dt_error dt_read_var_int(dt_decoder *d, int64_t *x) {
	uint64_t u = 0;
	dt_error err = dt_read_var_uint(d, &u);
	*x = zagzig(u);
	return err;
}

dt_error dt_write_var_uint_4bit(dt_encoder *e, uint64_t x) {
	while (true) {
		uint64_t y = x & ((1<<3) - 1);
		x >>= 3;
		if (x != 0) {
			y |= 1 << 3;
		}
		*ex |= y << ei;
		e->n += 4;
		if (ei < 4) {
			dt_error err = dt_grow(e, 1);
			if (err != dt_ok) {
				return err;
			}
		}
		if (x == 0) {
			break;
		}
	}
	return dt_ok;
}

dt_error dt_read_var_uint_4bit(dt_decoder *d, uint64_t *x) {
	for (int shift = 0; shift < 64; shift += 3) {
		if (dt_remaining(d) < 4) {
			return dt_error_unexpected_eof;
		}
		if ((d->j+3)%8 < 4) {
			dt_error err = dt_read(d, 1);
			if (err != dt_ok) {
				return err;
			}
		}

		uint64_t y = dx >> di;
		d->j += 4;
		if (di < 4) {
			d->b[0] = d->b[1];
		}
		*x |= (y & ((1<<3) - 1)) << shift;
		if ((y&(1<<3)) == 0) {
			return dt_ok;
		}
	}
	return dt_error_var_int_overflow;
}

dt_error dt_write_float32(dt_encoder *e, float x) { return dt_write_uint32(e, *(uint32_t*)&x); }
dt_error dt_read_float32(dt_decoder *d, float *x) { return dt_read_uint32(d, (uint32_t*)x); }

dt_error dt_write_float64(dt_encoder *e, double x) { return dt_write_uint64(e, *(uint64_t*)&x); }
dt_error dt_read_float64(dt_decoder *d, double *x) { return dt_read_uint64(d, (uint64_t*)x); }

dt_error dt_bytes_set(dt_bytes *b, uint64_t len, uint8_t *data) {
	dt_bytes_delete(b);
	b->data = calloc(len, 1);
	if (b->data == NULL && len > 0) {
		return dt_error_out_of_memory;
	}
	b->len = len;
	memcpy(b->data, data, len);
	return dt_ok;
}

void dt_bytes_delete(dt_bytes *b) {
	free(b->data);
	b->data = NULL;
	b->len = 0;
}

dt_error dt_write_bytes(dt_encoder *e, dt_bytes x) {
	dt_error err = dt_write_var_uint(e, x.len);
	if (err != dt_ok) {
		return err;
	}
	err = dt_grow(e, x.len);
	if (err != dt_ok) {
		return err;
	}
	while (x.len > 7) {
		*ex |= *(uint64_t*)x.data << ei;
		e->n += 56;
		x.data += 7;
		x.len -= 7;
	}
	while (x.len > 0) {
		*ex |= *x.data << ei;
		e->n += 8;
		x.data++;
		x.len--;
	}
	return dt_ok;
}

dt_error dt_read_bytes(dt_decoder *d, dt_bytes *x) {
	uint64_t len = 0;
	dt_error err = dt_read_var_uint(d, &len);
	if (err != dt_ok) {
		return err;
	}

	if (dt_remaining(d) < 8*len) {
		return dt_error_unexpected_eof;
	}

	if (x->len != len) {
		dt_bytes_delete(x);
		x->data = calloc(len, 1);
		if (x->data == NULL && len > 0) {
			return dt_error_out_of_memory;
		}
		x->len = len;
	}

	uint8_t *data = x->data;
	while (len >= 7) {
		dt_error err = dt_read(d, 7);
		if (err != dt_ok) {
			return err;
		}
		*(uint64_t*)data = (dx >> di) & (((uint64_t)1<<56) - 1);
		d->j += 56;
		d->b[0] = d->b[7];
		data += 7;
		len -= 7;
	}
	while (len > 0) {
		dt_error err = dt_read(d, 1);
		if (err != dt_ok) {
			return err;
		}
		data[0] = dx >> di;
		d->j += 8;
		d->b[0] = d->b[1];
		data++;
		len--;
	}

	return dt_ok;
}

dt_error dt_write_size(dt_encoder *e, dt_write_size_fn f, void *user_data) {
	// dt_write_size takes advantage of the fact that the size of the payload (a var_uint) occupies a whole number of bytes to avoid having to bit shift the payload, which could be expensive for large payloads.
	// The fact that the payload is written immediately with the right bit offset is also nice because it makes it possible for nested objects to do byte alignment, which would be good for large byte arrays.

	dt_encoder e2 = dt_new_encoder();
	e2.n = ei;
	dt_error err = f(&e2, user_data);
	if (err != dt_ok) {
		dt_delete_encoder(&e2);
		return err;
	}

	uint64_t size = e2.n - ei;
	err = dt_write_var_uint(e, size);
	if (err != dt_ok) {
		dt_delete_encoder(&e2);
		return err;
	}

	*ex |= e2.b[0];
	uint64_t i = e->n/8 + 1;
	uint64_t oldlen = dt_encoder_len(e);
	e->n += size;
	uint64_t newlen = dt_encoder_len(e);
	err = dt_grow(e, newlen-oldlen);
	if (err != dt_ok) {
		dt_delete_encoder(&e2);
		return err;
	}
	memcpy(e->b+i, e2.b+1, e2.len-1);

	dt_delete_encoder(&e2);
	return dt_ok;
}

dt_error dt_read_size(dt_decoder *d, dt_read_size_fn f, void *user_data) {
	uint64_t size = 0;
	dt_error err = dt_read_var_uint(d, &size);
	if (err != dt_ok) {
		return err;
	}
	if (dt_remaining(d) < size) {
		return dt_error_size_exceeds_limit;
	}

	uint64_t n = d->n;
	d->n = d->j + size;

	if (f != NULL) {
		err = f(d, user_data);
		if (err != dt_ok) {
			return err;
		}
	}

	uint64_t k = (d->n-1)/8 - (d->j-1)/8;
	if (k > 0) {
		uint8_t buf[k];
		memset(buf, 0, k);
		dt_error err = d->r(buf, k, d->user_data);
		if (err != dt_ok) {
			return err;
		}
		d->b[0] = buf[k-1];
	}
	d->j = d->n;
	d->n = n;

	return dt_ok;
}

uint8_t *dt_encoder_bytes(dt_encoder *e) { return e->b; }
uint64_t dt_encoder_len(dt_encoder *e) { return (e->n+7)/8; }
