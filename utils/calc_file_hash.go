package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func CalcFileHash(r io.ReadSeeker) string {
	buf := make([]byte, 4*1024*1024) // 4MB Cut
	h := sha256.New()

	for {
		bytesRead, err := r.Read(buf)
		if bytesRead > 0 {
			h.Write(buf[:bytesRead])
		}
		if err != nil {
			if err != io.EOF {
				return ""
			}
			break
		}
	}

	_, _ = r.Seek(0, io.SeekStart) // Reset read seeker

	return hex.EncodeToString(h.Sum(nil))
}
