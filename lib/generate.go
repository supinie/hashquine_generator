package lib

import (
	"encoding/hex"
	"sync"
    "os"
    "errors"
)

func gen_collisions(prefix []byte, hashquine_params Hashquine_params) ([][]byte, error) {
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


func Generate(hashquine_params Hashquine_params) ([]byte, error) {
    empty := make([]byte, 0)
    graphic_control_extension, err := hex.DecodeString("21f9040402000000")
    if err != nil {
        return empty, err
    }
    alternatives := make([]Alternatives, 0)     // char_pos, char: coll_pos, coll

    generated_gif := hashquine_params.Background_blocks["header"]
    generated_gif = append(generated_gif, hashquine_params.Background_blocks["lcd"]...)
    generated_gif = append(generated_gif, hashquine_params.Background_blocks["gct"]...)
    
    comment_prefix, err := hex.DecodeString("21fe")
    if err != nil {
        return empty, err
    }
    comment_suffix, err := hex.DecodeString("00")
    if err != nil {
        return empty, err
    }
    generated_gif = append(generated_gif, comment_prefix...)
    generated_gif = append(generated_gif, "Why are you looking here?"...)
    generated_gif = append(generated_gif, comment_suffix...)

    generated_gif = append(generated_gif, graphic_control_extension...)
    generated_gif = append(generated_gif, hashquine_params.Background_blocks["img_descriptor"]...)
    generated_gif = append(generated_gif, hashquine_params.Background_blocks["img_data"]...)
    
    ch := make(chan Collisions)
    wg := &sync.WaitGroup

    for char_index := 0; char_index < 32; char_index++ {
        if hashquine_params.Mask[char_index] != " " {
            continue
        }
        for char := 0; char < 16; char++ {
            wg.Add(1)
            go Gen_collisions()
}
