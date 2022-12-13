package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "sort"
)

type Pos [2]int

type PosQueue struct {
    q []PosNode
}

func (pq PosQueue) Len() int {
    return len(pq.q)
}

func (pq *PosQueue) Enqueue(p PosNode) {
    pq.q = append(pq.q, p)
}

func (pq *PosQueue) Dequeue() PosNode {
    front := pq.q[0]
    pq.q = pq.q[1:]
    return front
}

type PosNode struct {
    height int
    pos Pos
    parent *PosNode
}

func (p PosNode) Len() (int, map[Pos]bool) {
    path := map[Pos]bool{}
    length := 0
    current_node := &p
    for {
        if current_node == nil {
            break
        }

        path[current_node.pos] = true
        current_node = current_node.parent
        length++
    }
    return length - 1, path
}

func AbsDiff(a, b int) int {
    diff := a - b
    if diff < 0 {
        return -diff
    }
    return diff
}

func IntAbs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func bfs(mat [][]int, pos Pos, end Pos) *PosNode {
    // matrix dimensions
    width := len(mat[0])
    height := len(mat)

    // qeueue to store positions to visit
    q := PosQueue{}
    q.Enqueue(PosNode{mat[pos[1]][pos[0]], pos, nil})

    // visited positions
    visited := map[Pos]bool{
        pos: true,
    }

    for {
        if q.Len() == 0 {
            break
        }

        current_node := q.Dequeue()

        if current_node.pos == end {
            return &current_node
        }

        // for each neighbour
        for dx := -1; dx <= 1; dx++ {
            for dy := -1; dy <= 1; dy++ {
                if dx == dy || ((dx == 0) == (dy == 0)) {
                    // no diagonals
                    continue
                }
                
                if current_node.pos[0] + dx < 0 || current_node.pos[1] + dy < 0 ||
                        current_node.pos[0] + dx > width - 1 || current_node.pos[1] + dy > height - 1 {
                    // position out of bounds
                    continue
                }

                neighbour := [2]int{current_node.pos[0] + dx, current_node.pos[1] + dy}

                // too high
                if mat[neighbour[1]][neighbour[0]] - current_node.height > 1 {
                    continue
                }

                if !visited[neighbour] {
                    visited[neighbour] = true
                    q.Enqueue(PosNode{mat[neighbour[1]][neighbour[0]], neighbour, &current_node})
                }
            }
        }
    }

    return nil
}

func parse_input(filename string) ([][]int, Pos, Pos) {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        os.Exit(1)
    }

    reader := bufio.NewReader(file)

    mat := [][]int{}
    current_row := []int{}
    var offset int
    col, row := 0, 0

    var start [2]int
    var end [2]int

    for {
        r, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }

        if r == '\n' {
            mat = append(mat, current_row)
            current_row = nil
            row++
            col = 0
            continue
        }

        var intval int
        switch r {
        case 'S':
            intval = int('a')
            offset = intval
            start = Pos{col, row}
        case 'E':
            intval = int('z') + 1
            end = Pos{col, row}
        default:
            intval = int(r)
        }

        current_row = append(current_row, intval)
        col++
    }
    
    // subtract offset
    for y, row := range mat {
        for x := range row {
            mat[y][x] -= offset
        }
    }

    return mat, start, end
}

func main() {
    var part int
    var print_path bool
    flag.IntVar(&part, "part", 1, "part to run")
    flag.BoolVar(&print_path, "print", false, "print the final path")
    flag.Parse()

    if flag.NArg() != 1 {
        fmt.Println("usage: [OPTIONS] input-file")
        flag.Usage()
        os.Exit(1)
    }

    input := flag.Arg(0)

    mat, start, end := parse_input(input)

    switch part {
    case 1:
        end_node := bfs(mat, start, end)
        if end_node != nil {
            pathlen, path := end_node.Len()

            if print_path {
                for i, row := range mat {
                    for j := range row {
                        if path[[2]int{j, i}] {
                            fmt.Printf("#")
                        } else {
                            fmt.Printf(".")
                        }
                    }
                    fmt.Println()
                }
            }

            fmt.Printf("Path length from %v to %v: %d\n", start, end, pathlen)
        } else {
            fmt.Println("No path found...")
        }
    case 2:
        pathlens := []int{}
        for i, row := range mat {
            for j := range row {
                if mat[i][j] == 0 {
                    end_node := bfs(mat, [2]int{j, i}, end)
                    if end_node != nil {
                        pathlen, _ := end_node.Len()
                        pathlens = append(pathlens, pathlen)
                    }
                }
            }
        }
        sort.Ints(pathlens)
        fmt.Println(pathlens[0])
    case 3:
        log.Printf("error: part %d not implemented\n", part)
        os.Exit(1)
    }
}
