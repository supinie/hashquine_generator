package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sync"
)

const COLLISION_DIFF = 123
const COLLISION_LEN = 128

func gen_collisions(prefix []byte, tmp_dir string) ([]byte, []byte, error) {
    empty := make([]byte, 0)
    if len(prefix)%64 != 0 {
        err := errors.New("Misaligned prefix length")
        return empty, empty, err
    }
    prefix_file, err := os.CreateTemp(tmp_dir, "prefix")
    if err != nil {
        return empty, empty, err
    }
    coll_a_file, err := os.CreateTemp(tmp_dir, "coll_a")
    if err != nil {
        return empty, empty, err
    }
    coll_b_file, err := os.CreateTemp(temp_dir, "coll_b")
    if err != nil {
        return empty, empty, err
    }

    defer os.Remove(prefix_file.Name())
    defer os.Remove(coll_a_file.Name())
    defer os.Remove(coll_b_file.Name())

    if _, err := prefix_file.Write(prefix); err != nil {
        return empty, empty, err
    }
    fastcoll_args := []string{"-p", prefix_file.Name(), "-o", coll_a_file.Name(), coll_b_file.Name()}
REGEN:
    err = exec.Command("fastcoll", fastcoll_args...).Run()
    if err != nil {
        return empty, empty, err
    }
    coll_a, err := os.ReadFile(coll_a_file.Name())
    coll_b, err := os.ReadFile(coll_b_file.Name())
    if bytes.Equal(coll_a, coll_b) {
        goto REGEN
    }
    equality, err := test_file_hashes(coll_a_file.Name(), coll_b_file.Name())
    if err != nil {
        return empty, empty, err
    }
    if !equality {
        goto REGEN
    }
}


func test_file_hashes(filename1, filename2 string) (bool, err) {
    file1, err := os.Open(filename1)
    if err != nil {
        return md5.New(), err
    }
    defer file1.Close()
    
    file2, err := os.Open(filename2)
    if err != nil {
        return md5.New(), err
    }
    defer file2.Close()


    h1 := md5.New()
    if _, err := io.Copy(h1, file1); err != nil {
        return md5.New(), err
    }

    h2 := md5.New()
    if _, err := io.Copy(h2, file2); err != nil {
        return md5.New(), err
    }

    return bytes.Equal(h1.Sum(nil), h2.Sum(nil)), err 
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
    
    tmp_dir, err := os.MkdirTemp("", "tmp_dir")
    if err != nil {
        return empty, err
    }
    defer os.RemoveAll(tmp_dir)

    descriptor_prefix, err := hex.DecodeString("2c")
    if err != nil {
        return empty, err
    }
    descriptor_suffix, err := hex.DecodeString("00")
    if err != nil {
        return empty, err
    }

    for char_y := 0; char_y < 4; char_y++ {
        for char_x := 0; char_x < 8; char_x++ {
            if hashquine_params.Mask[char_index] != " " {
                continue
            }
            for char := 0; char < 16; char++ {
                char_img := make([]byte, 0)
                char_img = append(char_img, graphic_control_extension...)
                char_img = append(char_img, descriptor_prefix...)
                char_img = append(char_img, byte[char_x * 40])
                char_img = append(char_img, byte[char_y * 40])
                char_img = append(char_img, byte[hashquine_params.Char_dimension])
                char_img = append(char_img, byte[hashquine_params.Char_dimension])
                char_img = append(char_img, descriptor_prefix...)

                char_img = append(char_img, hashquine_params.Chars_img_data[char]...)

                // align with md5 block
                align := 64 - (len(generated_gif) % 64)
                generated_gif = append(generated_gif, byte[align - 1 + COLLISION_DIFF])
                padding, err := hex.DecodeString("00" * (align - 1))
                if err != nil {
                    return empty, err
                }
                generated_gif = append(generated_gif, padding...)
                coll_img, coll_nop, err := gen_collision(generated_gif, tmp_dir)
                if err != nil {
                    return empty, err
                }

}
