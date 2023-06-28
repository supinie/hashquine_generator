package lib

import (
    "encoding/hex"
)

func Generate(hashquine Hashquine) ([]byte, error) {
    empty := make([]byte, 0)
    graphic_control_extension, err := hex.DecodeString("21f9040402000000")
    if err != nil {
        return empty, err
    }
    alternatives := make([]Alternatives, 0)     // char_pos, char: coll_pos, coll

    generated_gif := hashquine.Background_blocks["header"]
    generated_gif = append(generated_gif, hashquine.Background_blocks["lcd"]...)
    generated_gif = append(generated_gif, hashquine.Background_blocks["gct"]...)
    
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
    generated_gif = append(generated_gif, hashquine.Background_blocks["img_descriptor"]...)
    generated_gif = append(generated_gif, hashquine.Background_blocks["img_data"]...)
    

}
