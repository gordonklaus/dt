package c

import "testing"

func TestCodec(t *testing.T)                  { testCodec(t) }
func TestReadSizeSkipsExtraBits(t *testing.T) { testReadSizeSkipsExtraBits(t) }
