package lib

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
	if !bytes.HasSuffix(lcd, []byte{0x80, 0x02, 0x00}) {
        fmt.Printf("%x\n", lcd)
		return nil, fmt.Errorf("Invalid logical screen descriptor for %s", filename)
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
	// if imgDescriptor[0] != 0x2c || imgDescriptor[9] != 0 {
	// 	return nil, fmt.Errorf("Invalid image descriptor")
	// }
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
        fmt.Printf("%x\n", subBlockSize)
		imgData = append(imgData, subBlockSize[0])
		if subBlockSize[0] == 0 {
			break
		}
		subBlockData := make([]byte, subBlockSize[0])
		_, err = file.Read(subBlockData)
		if err != nil {
			return nil, err
		}
        fmt.Printf("%x\n", subBlockData)
		imgData = append(imgData, subBlockData...)
	}
	blocks["img_data"] = imgData

	trailer := make([]byte, 2)
	_, err = file.Read(trailer)
	if err != nil {
		return nil, err
	}
    fmt.Printf("%x\n", trailer)
	if trailer[1] != 0x3b {
		return nil, fmt.Errorf("Invalid GIF trailer")
	}

	return blocks, nil
}
