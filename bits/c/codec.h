#pragma once

#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>

typedef enum {
	dt_ok,
	dt_error_unexpected_eof,
	dt_error_var_int_overflow,
	dt_error_size_exceeds_limit,
	dt_error_out_of_memory,
} dt_error;

int dt_unknown_enum_tag = -1;

typedef struct {
	uint8_t *b;
	uint64_t n, len;
} dt_encoder;

typedef dt_error dt_read_fn(uint8_t *buf, uint64_t len, void *user_data);

typedef struct {
	dt_read_fn *r;
	void       *user_data;
	uint8_t     b[9];
	uint64_t    j, n;
} dt_decoder;

dt_encoder dt_new_encoder();
dt_decoder dt_new_decoder(dt_read_fn r, void *user_data);

void dt_delete_encoder(dt_encoder *e);

uint64_t dt_remaining(dt_decoder *d) { return d->n - d->j; }

dt_error dt_write_bool(dt_encoder *e, bool x);
dt_error dt_read_bool(dt_decoder *d, bool *x);

dt_error dt_write_uint32(dt_encoder *e, uint32_t x);
dt_error dt_read_uint32(dt_decoder *d, uint32_t *x);

dt_error dt_write_uint64(dt_encoder *e, uint64_t x);
dt_error dt_read_uint64(dt_decoder *d, uint64_t *x);

dt_error dt_write_int64(dt_encoder *e, int64_t x);
dt_error dt_read_int64(dt_decoder *d, int64_t *x);

dt_error dt_write_var_uint(dt_encoder *e, uint64_t x);
dt_error dt_read_var_uint(dt_decoder *d, uint64_t *x);

dt_error dt_write_var_int(dt_encoder *e, int64_t x);
dt_error dt_read_var_int(dt_decoder *d, int64_t *x);

dt_error dt_write_var_uint_4bit(dt_encoder *e, uint64_t x);
dt_error dt_read_var_uint_4bit(dt_decoder *d, uint64_t *x);

dt_error dt_write_float32(dt_encoder *e, float x);
dt_error dt_read_float32(dt_decoder *d, float *x);

dt_error dt_write_float64(dt_encoder *e, double x);
dt_error dt_read_float64(dt_decoder *d, double *x);

typedef struct {
	uint64_t len;
	uint8_t *data;
} dt_bytes;

typedef dt_bytes dt_string;

dt_error dt_bytes_set(dt_bytes *b, uint64_t len, uint8_t *data);
void dt_bytes_delete(dt_bytes *b);

dt_error dt_write_bytes(dt_encoder *e, dt_bytes x);
dt_error dt_read_bytes(dt_decoder *d, dt_bytes *x);

typedef dt_error dt_write_size_fn(dt_encoder *e, void *user_data);
typedef dt_error dt_read_size_fn(dt_decoder *e, void *user_data);

dt_error dt_write_size(dt_encoder *e, dt_write_size_fn f, void *user_data);
dt_error dt_read_size(dt_decoder *d, dt_read_size_fn f, void *user_data);

uint8_t *dt_encoder_bytes(dt_encoder *e);
uint64_t dt_encoder_len(dt_encoder *e);
