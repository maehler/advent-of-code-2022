package main

import (
    "bufio"
    "bytes"
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
        flag.Usage()
        os.Exit(1)
    }

    a.input = flag.Arg(0)
}

type Crate rune

type Stack struct {
    crates []Crate
}

func (s *Stack) pop() Crate {
    r := s.crates[len(s.crates) - 1]
    if len(s.crates) == 1 {
        s.crates = nil
    } else if len(s.crates) > 1 {
        s.crates = s.crates[:len(s.crates) - 1]
    }
    return r
}

func (s *Stack) insert(c Crate) {
    if len(s.crates) == 0 {
        s.crates = append(s.crates, c)
    } else {
        s.crates = append(s.crates[0:1], s.crates...)
        s.crates[0] = c
    }
}

func (s *Stack) push(c Crate) {
    s.crates = append(s.crates, c)
}

func (s Stack) peek() *Crate {
    return &s.crates[len(s.crates) - 1]
}

type Move struct {
    n int
    from int
    to int
}

func do_move(stacks []Stack, move Move) []Stack {
    for i := 0; i < move.n; i++ {
        stacks[move.to].push(stacks[move.from].pop())
    }
    return stacks
}

func parse_input(filename string) ([]Stack, []Move) {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        os.Exit(1)
    }
    reader := bufio.NewReader(file)

    var stacks []Stack = nil
    for {
        // Crates
        line, err := reader.ReadSlice(byte('\n'))
        if err == io.EOF || line[0] == byte('\n') {
            break
        }

        n_stacks := len(line) / 4

        if stacks == nil {
            for i := 0; i < n_stacks; i++ {
                stacks = append(stacks, Stack{})
            }
        }

        n := 0
        for i := 0; i < len(line); i += 4 {
            if line[i] == '[' {
                stacks[n].insert(Crate(line[i + 1]))
            }
            n++
        }
    }

    moves := []Move{}
    for {
        // Moves
        line, err := reader.ReadSlice(byte('\n'))
        if err == io.EOF {
            break
        }

        move := Move{}
        m_pieces := bytes.Split(line[:len(line) - 1], []byte{' '})
        move.n, _ = strconv.Atoi(string(m_pieces[1]))
        move.from, _ = strconv.Atoi(string(m_pieces[3]))
        move.from--
        move.to, _ = strconv.Atoi(string(m_pieces[5]))
        move.to--

        moves = append(moves, move)
    }

    return stacks, moves
}

func main() {
    args := Args{}
    args.parse()

    stacks, moves := parse_input(args.input)

    switch args.part {
    case 1:
        for _, m := range moves {
            stacks = do_move(stacks, m)
        }
        for _, s := range stacks {
            fmt.Printf("%c", *s.peek())
        }
        fmt.Println()
    case 2:
        ;
    default:
        log.Printf("error: not implemented for part %d\n", args.part)
        os.Exit(1)
    }
}
