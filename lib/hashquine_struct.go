package lib

type Hashquine_params struct {
    Template_dir         string
    Output               string
    Hash_img_coordinates [2]int
    Mask                 string
    Background_blocks    map[string]([]byte)
    Chars_img_data       [16][]byte
    Char_height, Char_width   int
}

type Alternative_Key struct {
    Char_pos             [2]int
    Char                 int
}

type Alternative_Value struct {
    Coll_pos             int
    Coll                 []byte
}
