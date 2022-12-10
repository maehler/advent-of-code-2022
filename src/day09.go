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

func (p Pos) step_towards(other Pos) Pos {
    if p.touching(other) {
        return p
    }

    new_pos := Pos{p.x, p.y}

    x_diff := other.x - p.x
    y_diff := other.y - p.y
    if x_diff > 0 {
        new_pos.x++
    } else if x_diff < 0 {
        new_pos.x--
    }

    if y_diff > 0 {
        new_pos.y++
    } else if y_diff < 0 {
        new_pos.y--
    }

    return new_pos
}

type PlanckRope struct {
    head *Pos
    knots []Pos
    visited map[Pos]bool
}

func NewPlanckRope(n int) PlanckRope {
    v := make(map[Pos]bool)
    v[Pos{0, 0}] = true

    k := make([]Pos, n)
    for i := 0; i < n; i++ {
        k[i] = Pos{0, 0}
    }

    return PlanckRope {
        &k[0],
        k,
        v,
    }
}

func (pr *PlanckRope) move(m Move) {
    for i := 0; i < m.length; i++ {
        switch m.direction {
        case 'U':
            *pr.head = Pos{pr.head.x, pr.head.y - 1}
        case 'R':
            *pr.head = Pos{pr.head.x + 1, pr.head.y}
        case 'D':
            *pr.head = Pos{pr.head.x, pr.head.y + 1}
        case 'L':
            *pr.head = Pos{pr.head.x - 1, pr.head.y}
        }

        for ki := 1; ki < len(pr.knots); ki++ {
            if !pr.knots[ki - 1].touching(pr.knots[ki]) {
                pr.knots[ki] = pr.knots[ki].step_towards(pr.knots[ki -1])
                if ki == len(pr.knots) - 1 {
                    pr.visited[pr.knots[ki]] = true
                }
            }
        }
    }
}

func main() {
    args := Args{}
    args.parse()

    moves := parse_input(args.input)

    switch args.part {
    case 1:
        pr := NewPlanckRope(2)
        for _, move := range moves {
            pr.move(move)
        }
        fmt.Println(len(pr.visited))
    case 2:
        pr := NewPlanckRope(10)
        for _, move := range moves {
            pr.move(move)
        }
        fmt.Println(len(pr.visited))
    default:
        log.Printf("error: part %d not implemented\n", args.part)
    }
}
