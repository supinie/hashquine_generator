package lib

import (
)

func Generate(hashquine Hashquine) ([]byte, error) {
    empty := make([]byte, 0)
    graphic_control_extension := []byte("21f9040402000000")
    alternatives := make([]Alternatives, 0)     // char_pos, char: coll_pos, coll

    generated_gif := hashquine.Background_blocks["header"]
    generated_gif = append(generated_gif, hashquine.Background_blocks["lcd"]...)
    
}
