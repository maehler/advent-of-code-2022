package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    "strconv"
)

type Pos struct {
    x int
    y int
}

type Line struct {
    points []Pos
}

func (l Line) MaxX() int {
    max := 0
    for _, p := range l.points {
        if p.x > max {
            max = p.x
        }
    }
    return max
}

func (l Line) MaxY() int {
    max := 0
    for _, p := range l.points {
        if p.y > max {
            max = p.y
        }
    }
    return max
}

type Cave struct {
    board [][]rune
    start Pos
    n_grains int
}

func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func NewCave(lines []Line) Cave {
    max_x, max_y := 0, 0
    for _, l := range lines {
        lx, ly := l.MaxX(), l.MaxY()
        if lx > max_x {
            max_x = lx
        }
        if ly > max_y {
            max_y = ly
        }
    }

    max_x += 2
    max_y++

    board := make([][]rune, max_y + 1)
    for y := 0; y <= max_y; y++ {
        board[y] = make([]rune, max_x + 1)
        for x := 0; x <= max_x; x++ {
            if x == 0 || x == max_x || y == max_y {
                board[y][x] = '~'
            } else {
                board[y][x] = '.'
            }
        }
    }
    
    for _, line := range lines {
        for i := 0; i < len(line.points) - 1; i++ {
            x1 := line.points[i].x
            x2 := line.points[i + 1].x
            y1 := line.points[i].y
            y2 := line.points[i + 1].y
            for x := Min(x1, x2); x <= Max(x1, x2); x++ {
                for y := Min(y1, y2); y <= Max(y1, y2); y++ {
                    board[y][x] = '#'
                }
            }
        }
    }

    return Cave{board, Pos{500, 0}, 0}
}

func (c Cave) Print() {
    for _, row := range c.board {
        for _, x := range row {
            fmt.Printf("%c", x)
        }
        fmt.Printf("\n")
    }
}

func (c Cave) IsFree(p Pos) bool {
    return c.board[p.y][p.x] == '.' || c.board[p.y][p.x] == '~'
}

func (c Cave) Oob(p Pos) bool {
    if p.x < 0 || p.x >= len(c.board[0]) {
        return true
    }

    return p.y >= len(c.board)
}

func (c Cave) NextPos(p Pos) Pos {
    down := Pos{p.x, p.y + 1}
    down_left := Pos{down.x - 1, down.y}
    down_right := Pos{down.x + 1, down.y}

    if c.IsFree(down) {
        return down
    }

    if c.IsFree(down_left) {
        return down_left
    }

    if c.IsFree(down_right) {
        return down_right
    }

    return p
}

func (c *Cave) Set(p Pos, r rune) {
    c.board[p.y][p.x] = r
}

func (c Cave) Get(p Pos) rune {
    return c.board[p.y][p.x]
}

func (c *Cave) Add() bool {
    prev_pos := c.start
    for {
        current_pos := c.NextPos(prev_pos)
        if c.Get(current_pos) == '~' {
            // fell off the edge
            return false
        }
        if current_pos == prev_pos {
            // settled
            c.n_grains++
            c.Set(current_pos, 'o')
            return true
        }
        prev_pos = current_pos
    }
}

func parse_input(filename string) Cave {
    contents, _ := os.ReadFile(filename)

    lines := []Line{}

    for _, line := range strings.Split(strings.TrimSpace(string(contents)), "\n") {
        points := []Pos{}
        for _, l := range strings.Split(line, " -> ") {
            str_pos := strings.Split(l, ",")
            x, _ := strconv.Atoi(str_pos[0])
            y, _ := strconv.Atoi(str_pos[1])
            points = append(points, Pos{x, y})
        }
        lines = append(lines, Line{points})
    }

    return NewCave(lines)
}

func main() {
    var part int
    flag.IntVar(&part, "part", 1, "part to run")
    flag.Parse()

    input := flag.Arg(0)

    cave := parse_input(input)

    for {
        if !cave.Add() {
            break
        }
    }

    fmt.Println(cave.n_grains)
}
