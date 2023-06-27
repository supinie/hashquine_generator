package lib

import (
    "os"
)

func Read_gif(filename string) (map[string]([]byte), error) {
    empty := make(map[string]([]byte))
    file_data, err := os.ReadFile(filename)
    if err != nil {
        return empty, err
    }
    blocks := make(map[string]([]byte))

    blocks["header"], file_data  = file_data[:6], file_data[6:]
    blocks["lcd"], file_data = file_data[:7], file_data[7:]
    blocks["gct"], file_data = file_data[:16 * 3], file_data[(16 * 3):]
    blocks["img_descriptor"], file_data = file_data[:10], file_data[10:]
    blocks["img_data"] = file_data[:len(file_data) - 1]

    return blocks, err
}
