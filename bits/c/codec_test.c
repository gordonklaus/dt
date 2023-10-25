#include "codec.h"
#include "_cgo_export.h"

dt_error c_write_size_callback(dt_encoder *e, void *user_data) {
	return go_write_size_callback(e, user_data);
}

dt_error c_read_size_callback(dt_decoder *d, void *user_data) {
	return go_read_size_callback(d, user_data);
}

dt_error c_read(uint8_t *b, uint64_t n, void *user_data) {
	return go_read(b, n);
}
