package main

import (
    "supinie/hashquine_generator/lib"

    "flag"
    "fmt"
    "strconv"
)
    
type Hashquine struct {
    template_dir         string
    output               string
    hash_img_coordinates [2]int
    mask                 string
    background_blocks    map[string]([]byte)
    chars_img_data       map[uint64]([]byte)
    char_dimension       int
}

func main() {
    var template_dir string
    var output string

    flag.StringVar(&template_dir, "t", "building_blocks", "Specify the location of the directory containing the templates to be used. By default, this will be the building_blocks dir in the root of the repo.")
    flag.StringVar(&output, "o", "hashquine.gif", "Specify the name of the output gif, by default this will be 'hashquine.gif'.")
    flag.Parse()

    background_blocks, err := lib.Read_gif("./building_blocks/background.gif")
    if err != nil {
        fmt.Println(err)
        return
    }
    chars_img_data := make(map[uint64]([]byte))
    for _, character := range "0123456789abcdef" {
        char_block, err := lib.Read_gif("./building_blocks/" + string(character) + ".gif")
        if err != nil {
            fmt.Println(err)
            return
        }
        index, err := strconv.ParseUint(string(character), 16, 0)
        if err != nil {
            fmt.Println(err)
            return
        }
        chars_img_data[index] = char_block["img_data"]
    }
    hashquine := Hashquine{
        template_dir,
        output,
        [2]int{0, -40},
        "1337    deadbeef                                                ",
        background_blocks,
        chars_img_data,
        40,
    }
    fmt.Printf("%v", hashquine)
}
