package lib

type Hashquine struct {
    Template_dir         string
    Output               string
    Hash_img_coordinates [2]int
    Mask                 string
    Background_blocks    map[string]([]byte)
    Chars_img_data       map[uint64]([]byte)
    Char_dimension       int
}

type Alternatives struct {
    Char_pos             [2]int
    Char                 int
    Coll_pos             [2]int
    Coll                 []byte
}
    
