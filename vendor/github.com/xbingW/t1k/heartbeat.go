package t1k

import (
	"io"

	"github.com/xbingW/t1k/t1k"
)

func DoHeartbeat(s io.ReadWriter) error {
	h := t1k.MakeHeader(t1k.MASK_FIRST|t1k.MASK_LAST, 0)
	_, err := s.Write(h.Serialize())
	if err != nil {
		return err
	}
	_, err = readDetectionResult(s)
	return err
}
