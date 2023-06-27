package lib

import (
    "errors"
    "os"
)

func Gen_collisions(prefix []byte) ([][]byte, error) {
    empty := make([][]byte, 0)
    if len(prefix)%64 != 0 {
        return empty, errors.New("Misaligned prefix length")
    }
    tmp_dir, err := os.MkdirTemp("", "test")
    if err != nil {
        return empty, err
    }
    defer os.RemoveAll(tmp_dir)
    return empty, err
}

