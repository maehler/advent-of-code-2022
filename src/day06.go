package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
)

type Args struct {
    part int
    input string
}

func (a *Args) parse() {
    flag.IntVar(&a.part, "part", 1, "part to run")
    flag.Parse()

    if flag.NArg() != 1 {
        fmt.Println("usage: [OPTIONS] input-filename")
        flag.Usage()
        os.Exit(1)
    }

    a.input = flag.Arg(0)
}

type Queue struct {
    size int
    items []rune
}

func NewQueue(size int) Queue {
    return Queue{size, make([]rune, size)}
}

func (q *Queue) push(v rune) {
    q.items = append(q.items, v)
    if len(q.items) > q.size {
        q.items = q.items[1:]
    }
}

func (q Queue) all_unique() bool {
    c := make(map[rune]int)

    for _, v := range q.items {
        c[v]++
    }

    for _, n := range c {
        if n > 1 {
            return false
        }
    }

    return true
}

func parse_input(filename string, length int) int {
    file, err := os.Open(filename)
    if err != nil {
        log.Println("error: %s", err)
        os.Exit(1)
    }
    defer file.Close()
    reader := bufio.NewReader(file)

    buf := NewQueue(length)
    n := 1
    for {
        r, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }
        
        if r == '\n' {
            continue
        }

        buf.push(r)
        if n >= length && buf.all_unique() {
            return n
        }
        n++
    }
    return 1
}

func main() {
    args := Args{}
    args.parse()


    switch args.part {
    case 1:
        fmt.Println(parse_input(args.input, 4))
    case 2:
        fmt.Println(parse_input(args.input, 14))
    default:
        log.Printf("error: part %d not implemented", args.part)
        os.Exit(1)
    }
}
