package lib

import (
    "fmt"
	"os"
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
		return nil, fmt.Errorf("Invalid GIF file format for %s", filename)
	}
	blocks["header"] = header

	lsd := make([]byte, 7)
	_, err = file.Read(lsd)
	if err != nil {
		return nil, err
	}
	if !bytes.HasSuffix(lsd, []byte{0x80, 0x02, 0x00}) {
		return nil, fmt.Errorf("Invalid logical screen descriptor for %s", filename)
	}
	blocks["lsd"] = lsd

	gct := make([]byte, 6)
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
		return nil, fmt.Errorf("Invalid image descriptor for %s", filename)
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

	trailer := make([]byte, 1)
	_, err = file.Read(trailer)
	if err != nil {
		return nil, err
	}
	if trailer[0] != 0x3b {
		return nil, fmt.Errorf("Invalid GIF trailer for %s", filename)
	}

	blocks["img_data"] = imgData

	return blocks, nil
}

