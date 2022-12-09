package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
)

type Args struct {
    part int
    input string
}

func (a *Args) parse() {
    flag.IntVar(&a.part, "part", 1, "part to run")
    flag.Parse()

    if flag.NArg() != 1 {
        fmt.Println("usage: [OPTIONS] input-file")
        os.Exit(1)
    }

    a.input = flag.Arg(0)
}

func parse_input(filename string) [][]int {
    file, err := os.Open(filename)
    if err != nil {
        log.Println("error: %s", err)
        os.Exit(1)
    }

    reader := bufio.NewReader(file)

    forest := [][]int{}
    row := []int{}

    for {
        r, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }

        if r == '\n' {
            forest = append(forest, row)
            row = nil
            continue
        }

        h, _ := strconv.Atoi(string(r))
        row = append(row, h)
    }

    return forest
}

type Forest struct {
    width int
    height int
    trees [][]int
    visibility [][]bool
}

func (f *Forest) check_visibility() {
    // West
    for y, row := range f.trees {
        max := -1
        for x, t := range row {
            if t > max {
                max = t
                f.visibility[y][x] = true
            }
        }
    }

    // North
    for x := 0; x < f.width; x++ {
        max := -1
        for y := 0; y < f.height; y++ {
            if f.trees[y][x] > max {
                max = f.trees[y][x]
                f.visibility[y][x] = true
            }
        }
    }

    // East
    for y := 0; y < f.height; y++ {
        max := -1
        for x := f.width - 1; x >= 0; x-- {
            if f.trees[y][x] > max {
                max = f.trees[y][x]
                f.visibility[y][x] = true
            }
        }
    }

    // South
    for x := 0; x < f.width; x++ {
        max := -1
        for y := f.height - 1; y >= 0; y-- {
            if f.trees[y][x] > max {
                max = f.trees[y][x]
                f.visibility[y][x] = true
            }
        }
    }
}

func (f Forest) count_visible() int {
    sum := 0
    for _, row := range f.visibility {
        for _, t := range row {
            if t {
                sum++
            }
        }
    }
    return sum
}

func (f Forest) print() {
    for _, row := range f.visibility {
        for _, t := range row {
            if t {
                fmt.Printf("1 ")
            } else {
                fmt.Printf("0 ")
            }
        }
        fmt.Println()
    }
}

func NewForest(trees [][]int) Forest {
    height := len(trees)
    width := len(trees[0])
    visibility := make([][]bool, height)
    for y := range visibility {
        visibility[y] = make([]bool, width)
    }
    return Forest {
        width,
        height,
        trees,
        visibility,
    }
}

func main() {
    args := Args{}
    args.parse()

    forest := NewForest(parse_input(args.input))
    forest.check_visibility()

    switch args.part {
    case 1:
        fmt.Println(forest.count_visible())
    case 2:
        ;
    default:
        log.Printf("error: part %d not implemented", args.part)
        os.Exit(1)
    }
}
