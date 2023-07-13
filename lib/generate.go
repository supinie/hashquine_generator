package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"strconv"
)

const COLLISION_DIFF = 123
const COLLISION_LEN = 128

func gen_collision(prefix []byte, tmp_dir string) ([]byte, []byte, error) {
    if len(prefix)%64 != 0 {
        fmt.Println(len(prefix))
        err := errors.New("Misaligned prefix length")
        return nil, nil, err
    }
    prefix_file, err := os.CreateTemp(tmp_dir, "prefix")
    if err != nil {
        return nil, nil, err
    }
    coll_a_file, err := os.CreateTemp(tmp_dir, "coll_a*.bit")
    if err != nil {
        return nil, nil, err
    }
    coll_b_file, err := os.CreateTemp(tmp_dir, "coll_b*.bit")
    if err != nil {
        return nil, nil, err
    }

    defer os.Remove(prefix_file.Name())
    defer os.Remove(coll_a_file.Name())
    defer os.Remove(coll_b_file.Name())

    if _, err := prefix_file.Write(prefix); err != nil {
        return nil, nil, err
    }
    fastcoll_args := []string{"-p", prefix_file.Name(), "-o", coll_a_file.Name(), coll_b_file.Name()}
REGEN:
    err = exec.Command("fastcoll", fastcoll_args...).Run()
    if err != nil {
        return nil, nil, err
    }
    coll_a, err := os.ReadFile(coll_a_file.Name())
    if err != nil {
        return nil, nil, err
    }
    coll_b, err := os.ReadFile(coll_b_file.Name())
    if err != nil {
        return nil, nil, err
    }
    coll_a = coll_a[len(prefix):]
    coll_b = coll_b[len(prefix):]
    if bytes.Equal(coll_a, coll_b) {
        goto REGEN
    }
    equality, err := test_file_hashes(coll_a_file.Name(), coll_b_file.Name())
    if err != nil {
        return nil, nil, err
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


func Generate(hp *Hashquine_params) ([]byte, error) {
    graphic_control_extension := []byte{0x21, 0xf9, 0x04, 0x04, 0x02, 0x00, 0x00, 0x00}

    alternatives := make(map[Alternative_Key]Alternative_Value)     // [char_x, char_y], char: coll_pos, coll

    generated_gif := append(hp.Background_blocks["header"], hp.Background_blocks["lsd"]...)
    generated_gif = append(generated_gif, hp.Background_blocks["gct"]...)
    
    comment := []byte("If you see this, please don't use md5\n")
    generated_gif = append(generated_gif, []byte{0x21, 0xfe}...)
    generated_gif = append(generated_gif, byte(len(comment)))
    generated_gif = append(generated_gif, comment...)
    generated_gif = append(generated_gif, 0x00)

    generated_gif = append(generated_gif, graphic_control_extension...)
    generated_gif = append(generated_gif, hp.Background_blocks["img_descriptor"]...)
    generated_gif = append(generated_gif, hp.Background_blocks["img_data"]...)

    generated_gif = append(generated_gif, []byte{0x21, 0xfe}...)
    
    tmp_dir, err := os.MkdirTemp("", "tmp_dir")
    if err != nil {
        return nil, err
    }
    defer os.RemoveAll(tmp_dir)

    top, left := hp.Hash_img_coordinates[0], hp.Hash_img_coordinates[1]
    for char_y := 0; char_y < 4; char_y++ {
        top += hp.Char_height
        left = -40
        for char_x := 0; char_x < 8; char_x++ {
            left += hp.Char_width
            if string([]rune(hp.Mask)[char_x + (8 * char_y)]) != " " {
                continue
            }
            for char := 0; char < 16; char++ {
                char_img := append(graphic_control_extension, []byte{0x2c}...)

                left_bytes := make([]byte, 2)
                top_bytes := make([]byte, 2)
                binary.LittleEndian.PutUint16(left_bytes, uint16(left))
                binary.LittleEndian.PutUint16(top_bytes, uint16(top))
                char_img = append(char_img, left_bytes...)
                char_img = append(char_img, top_bytes...)
                char_img = append(char_img, byte(hp.Char_width), byte(hp.Char_height), 0x00)
                
                char_img = append(char_img, hp.Chars_img_data[char]...)

                // align with md5 block
                align := 64 - (len(generated_gif) % 64)
                aligned_bytes := make([]byte, align - 1)
                generated_gif = append(generated_gif, byte(align - 1 + COLLISION_DIFF))
                generated_gif = append(generated_gif, aligned_bytes...)

                var coll_img, coll_nop []byte
                var coll_p_img, pad_len int
                for {
                    fmt.Printf("Generating collision for position (%d, %d), character %d\n", char_x,  char_y, char)
                    coll_img, coll_nop, err = gen_collision(generated_gif, tmp_dir)
                    if err != nil {
                        return nil, err
                    }
                    offset := COLLISION_LEN - COLLISION_DIFF - 1
                    coll_p_img = int(coll_img[COLLISION_DIFF]) - offset
                    coll_p_nop := int(coll_nop[COLLISION_DIFF]) - offset
                    pad_len = coll_p_nop - coll_p_img - len(char_img) - 4
                    if pad_len < 0 {
                        pad_len += 64
                    }
                    if coll_p_img >= 0 && pad_len >= 0 {
                        break
                    }
                    fmt.Println("Bad collision, retrying")
                }
                char_pos := [2]int{char_x, char_y}
                alternatives[Alternative_Key{char_pos, char}] =  Alternative_Value{len(generated_gif), coll_img}
                generated_gif = append(generated_gif, coll_nop...)
                generated_gif = append(generated_gif, make([]byte, coll_p_img)...)
                generated_gif = append(generated_gif, 0x00)

                generated_gif = append(generated_gif, char_img...)

                generated_gif = append(generated_gif, []byte{0x21, 0xfe}...)
                generated_gif = append(generated_gif, byte(pad_len))
                generated_gif = append(generated_gif, make([]byte, pad_len)...)
            }
        }
    }
    current_md5 := md5.Sum(generated_gif)

    fmt.Println("Brute forcing...")
    for garbage := 0; garbage < (1 << 32); garbage++ {
        end := []byte{0x04, byte(garbage), 0x00, 0x3b}

        new_md5 := md5.Sum(append(current_md5[:], end...))

        match := true
        new_md5_iterable := []rune(fmt.Sprintf("%x", new_md5))
        for i, mask_char := range hp.Mask {
            md5_char := new_md5_iterable[i]
            if string(mask_char) != " " && string(mask_char) != string(md5_char) {
                match = false
                break
            }
        }

        if match {
            fmt.Println("Found bruteforce")
            generated_gif = append(generated_gif, end...)
            break
        }
    }
    if len(generated_gif) == 0 { 
        err := errors.New("Did not find GIF matching md5 mask")
        return nil, err
    }

    fmt.Printf("Target hash: %x", md5.Sum(generated_gif))
    hash_rune_slice := []rune(fmt.Sprintf("%x", md5.Sum(generated_gif)))
    mask_rune_slice := []rune(hp.Mask)
    for i, char := range hash_rune_slice {
        if string(mask_rune_slice[i]) != " " {
            continue
        }
        x := i % 8
        y := i / 8
        new_char, err := strconv.ParseInt(string(char), 16, 0)
        if err != nil {
            return nil, err
        }
        coll_alternative := alternatives[Alternative_Key{[2]int{x, y}, int(new_char)}]
        fmt.Printf("Substituting char %v in position (%v, %v)\n", new_char, x, y)
        generated_gif = append(generated_gif[:coll_alternative.Coll_pos], append(coll_alternative.Coll, generated_gif[coll_alternative.Coll_pos + len(coll_alternative.Coll):]...)...)
    }
    return generated_gif, err
}
