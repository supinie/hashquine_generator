package main

import (
    "supinie/hashquine_generator/lib"

    "flag"
    "fmt"
    "strconv"
)
    
func main() {
    var template_dir string
    var output string

    flag.StringVar(&template_dir, "t", "./building_blocks", "Specify the location of the directory containing the templates to be used. By default, this will be the building_blocks dir in the root of the repo.")
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
    hashquine := lib.Hashquine{
        Template_dir: template_dir,
        Output: output,
        Hash_img_coordinates: [2]int{0, -40},
        Mask: "1337    deadbeef                                                ",
        Background_blocks: background_blocks,
        Chars_img_data: chars_img_data,
        Char_dimension: 40,
    }
    fmt.Printf("%v", hashquine)
}
