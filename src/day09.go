package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
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

func parse_input(filename string) []Move {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        os.Exit(1)
    }

    reader := bufio.NewReader(file)
    moves := []Move{}
    
    for {
        line, err := reader.ReadString(byte('\n'))
        if err == io.EOF {
            break
        }

        move := strings.Fields(line)
        direction := []rune(move[0])
        length, _ := strconv.Atoi(move[1])

        moves = append(moves, Move{direction[0], length})
    }

    return moves
}

type Move struct {
    direction rune
    length int
}

type Pos struct {
    x int
    y int
}

func abs(i int) int {
    if i < 0 {
        return -i
    }
    return i
}

func (p Pos) touching(other Pos) bool {
    diff_x := p.x - other.x
    diff_y := p.y - other.y
    return !(abs(diff_x) > 1 || abs(diff_y) > 1)
}

type PlanckRope struct {
    head Pos
    tail Pos
    visited map[Pos]bool
}

func NewPlanckRope() PlanckRope {
    v := make(map[Pos]bool)
    v[Pos{0, 0}] = true
    return PlanckRope {
        Pos{0, 0},
        Pos{0, 0},
        v,
    }
}

func (pr *PlanckRope) move(m Move) {
    for i := 0; i < m.length; i++ {
        prev_head_pos := pr.head

        switch m.direction {
        case 'U': pr.head = Pos{pr.head.x, pr.head.y - 1}
        case 'R': pr.head = Pos{pr.head.x + 1, pr.head.y}
        case 'D': pr.head = Pos{pr.head.x, pr.head.y + 1}
        case 'L': pr.head = Pos{pr.head.x - 1, pr.head.y}
        }

        if !pr.head.touching(pr.tail) {
            pr.tail = prev_head_pos
            pr.visited[pr.tail] = true
        }
    }
}

func main() {
    args := Args{}
    args.parse()

    moves := parse_input(args.input)
    pr := NewPlanckRope()

    switch args.part {
    case 1: {
        for _, move := range moves {
            pr.move(move)
        }
        fmt.Println(len(pr.visited))
    }
    case 2:
        ;
    default:
        log.Printf("error: part %d not implemented\n", args.part)
    }
}
