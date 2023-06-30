package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const COLLISION_DIFF = 123
const COLLISION_LEN = 128

func gen_collision(prefix []byte, tmp_dir string) ([]byte, []byte, error) {
    empty := make([]byte, 0)
    if len(prefix)%64 != 0 {
        err := errors.New("Misaligned prefix length")
        return empty, empty, err
    }
    prefix_file, err := os.CreateTemp(tmp_dir, "prefix")
    if err != nil {
        return empty, empty, err
    }
    coll_a_file, err := os.CreateTemp(tmp_dir, "coll_a*.bit")
    if err != nil {
        return empty, empty, err
    }
    coll_b_file, err := os.CreateTemp(tmp_dir, "coll_b*.bit")
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
    if err != nil {
        return empty, empty, err
    }
    coll_b, err := os.ReadFile(coll_b_file.Name())
    if err != nil {
        return empty, empty, err
    }
    coll_a = coll_a[len(prefix):]
    coll_b = coll_b[len(prefix):]
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
    if len(coll_a) != len(coll_b) && len(coll_a) != COLLISION_LEN {
        goto REGEN
    }
    if coll_a[COLLISION_DIFF] < coll_b[COLLISION_DIFF] {
        return coll_a, coll_b, err
    }
    return coll_b, coll_a, err
}


func test_file_hashes(filename1, filename2 string) (bool, error) {
    file1, err := os.Open(filename1)
    if err != nil {
        return false, err
    }
    defer file1.Close()
    
    file2, err := os.Open(filename2)
    if err != nil {
        return false, err
    }
    defer file2.Close()


    h1 := md5.New()
    if _, err := io.Copy(h1, file1); err != nil {
        return false, err
    }

    h2 := md5.New()
    if _, err := io.Copy(h2, file2); err != nil {
        return false, err
    }

    return bytes.Equal(h1.Sum(nil), h2.Sum(nil)), err 
}


func Generate(hashquine_params Hashquine_params) ([]byte, error) {
    empty := make([]byte, 0)
    graphic_control_extension, err := hex.DecodeString("21f9040402000000")
    if err != nil {
        return empty, err
    }
    alternatives := make(map[Alternative_Key]Alternative_Value)     // char_pos, char: coll_pos, coll

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
            if string([]rune(hashquine_params.Mask)[char_x + (8 * char_y)]) != " " {
                continue
            }
            for char := 0; char < 16; char++ {
                char_img := make([]byte, 0)
                char_img = append(char_img, graphic_control_extension...)
                char_img = append(char_img, descriptor_prefix...)
                char_img = append(char_img, byte(char_x * 40))
                char_img = append(char_img, byte(char_y * 40))
                char_img = append(char_img, byte(hashquine_params.Char_dimension))
                char_img = append(char_img, byte(hashquine_params.Char_dimension))
                char_img = append(char_img, descriptor_suffix...)

                char_img = append(char_img, hashquine_params.Chars_img_data[uint64(char)]...)

                // align with md5 block
                align := 64 - (len(generated_gif) % 64)
                generated_gif = append(generated_gif, byte(align - 1 + COLLISION_DIFF))
                padding, err := hex.DecodeString(strings.Repeat("00", (align - 1)))
                if err != nil {
                    return empty, err
                }
                generated_gif = append(generated_gif, padding...)
                var coll_img, coll_nop []byte
                var coll_p_img, coll_p_nop, pad_len int
                for {
                    fmt.Printf("Generating collision %d\n", char_y + char_x + char)
                    coll_img, coll_nop, err = gen_collision(generated_gif, tmp_dir)
                    if err != nil {
                        return empty, err
                    }
                    offset := COLLISION_LEN - COLLISION_DIFF - 1
                    coll_p_img = int(coll_img[COLLISION_DIFF]) - offset
                    coll_p_nop = int(coll_nop[COLLISION_DIFF]) - offset
                    pad_len = int(coll_p_nop) - int(coll_p_img) - len(char_img) - 4
                    if coll_p_img >= 0 && pad_len >= 0 {
                        break
                    }
                    fmt.Println("Bad collision, retrying")
                }
                char_pos := [2]int{char_x, char_y}
                alternatives[Alternative_Key{char_pos, char}] =  Alternative_Value{len(generated_gif), coll_img}
                generated_gif = append(generated_gif, coll_nop...)

                padding, err = hex.DecodeString(strings.Repeat("00", (coll_p_img + 1)))
                if err != nil {
                    return empty, err
                }
                generated_gif = append(generated_gif, padding...)
                generated_gif = append(generated_gif, char_img...)

                generated_gif = append(generated_gif, comment_prefix...)
                generated_gif = append(generated_gif, byte(pad_len))
                generated_gif = append(generated_gif, comment_suffix...)

                padding, err = hex.DecodeString(strings.Repeat("00", pad_len))
                if err != nil {
                    return empty, err
                }
                generated_gif = append(generated_gif, padding...)
            }
        }
    }
    fmt.Println("bruteforcing hash for fixed digits...")
    current_md5 := md5.Sum(generated_gif)

    for garbage := 0; garbage < (1 << 32); garbage++ {
        fmt.Println("Brute forcing...")
        comment_sub_block := []byte{4, byte(garbage), 0}
        end_comment := []byte{0}
        trailer := []byte{0x3b}

        end := append(append(comment_sub_block, end_comment...), trailer...)


        new_md5 := md5.Sum(append(current_md5[:], end...))

        match := true
        for i, mask_char := range hashquine_params.Mask {
            md5_char := fmt.Sprintf("%02x", new_md5[i])
            if string(mask_char) != " " && string(mask_char) != md5_char {
                match = false
                break
            }
        }

        if match {
            generated_gif = append(generated_gif, end...)
            break
        }
    }
    fmt.Printf("Target hash: %x", md5.Sum(generated_gif))
    hash_slice := md5.Sum(generated_gif)
    for i, char := range hex.EncodeToString(hash_slice[:]) {
        if string(hashquine_params.Mask[i]) != " " {
            continue
        }
        x := i % 4
        y := (i - x)/8
        coll_alternative := alternatives[Alternative_Key{[2]int{x, y}, int(char)}]
        generated_gif = append(generated_gif[:coll_alternative.Coll_pos], coll_alternative.Coll...)
        generated_gif = append(generated_gif, generated_gif[coll_alternative.Coll_pos + len(coll_alternative.Coll):]...)
    }
    return generated_gif, err
}
