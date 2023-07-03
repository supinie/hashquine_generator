package lib

// import (
// 	"fmt"
// 	"os"
// )

// func Read_gif(filename string) (map[string]([]byte), error) {
//     empty := make(map[string]([]byte))
//     file_data, err := os.ReadFile(filename)
//     if err != nil {
//         return empty, err
//     }
//     blocks := make(map[string]([]byte))

//     blocks["header"], file_data  = file_data[:6], file_data[6:]
//     blocks["lcd"], file_data = file_data[:7], file_data[7:]
//     blocks["gct"], file_data = file_data[:16 * 3], file_data[(16 * 3):]
//     blocks["img_descriptor"], file_data = file_data[:10], file_data[10:]
//     blocks["img_data"], file_data = file_data[:1], file_data[1:]
//     for i := 0; i < 50; i++ {
//         new_block := make([]byte, 0)
//         new_block, file_data = file_data[:1], file_data[1:]
//         fmt.Printf("%x\n", new_block)
//     }
//     return blocks, err
// }

import (
	"os"
    "fmt"
    "bytes"
)

func Read_gif(filename string) (map[string][]byte, error) {
	blocks := make(map[string][]byte)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	header := make([]byte, 6)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}
	if string(header) != "GIF87a" && string(header) != "GIF89a" {
		return nil, fmt.Errorf("Invalid GIF file format")
	}
	blocks["header"] = header

	lcd := make([]byte, 7)
	_, err = file.Read(lcd)
	if err != nil {
		return nil, err
	}
    fmt.Printf("%x\n", lcd)
	if !bytes.HasSuffix(lcd, []byte{0x80, 0x02, 0x00}) {
		return nil, fmt.Errorf("Invalid logical screen descriptor")
	}
	blocks["lcd"] = lcd

	gct := make([]byte, 16*3)
	_, err = file.Read(gct)
	if err != nil {
		return nil, err
	}
	blocks["gct"] = gct

	imgDescriptor := make([]byte, 10)
	_, err = file.Read(imgDescriptor)
	if err != nil {
		return nil, err
	}
	if imgDescriptor[0] != 0x2c || imgDescriptor[9] != 0 {
		return nil, fmt.Errorf("Invalid image descriptor")
	}
	blocks["img_descriptor"] = imgDescriptor

	imgData := make([]byte, 1)
	_, err = file.Read(imgData)
	if err != nil {
		return nil, err
	}
	for {
		subBlockSize := make([]byte, 1)
		_, err = file.Read(subBlockSize)
		if err != nil {
			return nil, err
		}
		imgData = append(imgData, subBlockSize...)
		if subBlockSize[0] == 0 {
			break
		}
		subBlockData := make([]byte, subBlockSize[0])
		_, err = file.Read(subBlockData)
		if err != nil {
			return nil, err
		}
		imgData = append(imgData, subBlockData...)
	}
	blocks["img_data"] = imgData

	trailer := make([]byte, 1)
	_, err = file.Read(trailer)
	if err != nil {
		return nil, err
	}
	if trailer[0] != 0x3b {
		return nil, fmt.Errorf("Invalid GIF trailer")
	}

	return blocks, nil
}
