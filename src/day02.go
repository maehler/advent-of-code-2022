package main

import (
    "bytes"
    "flag"
    "fmt"
    "log"
    "os"
)

func parse_args() (int, string) {
    var part int
    var input string

    flag.IntVar(&part, "part", 1, "part to run")

    flag.Parse()

    if flag.NArg() != 1 {
        log.Printf("usage: [OPTIONS] input-file")
        flag.Usage()
        os.Exit(1)
    }

    input = flag.Arg(0)

    return part, input
}

var decoding = map[string]string{
    "A": "rock",
    "X": "rock",
    "B": "paper",
    "Y": "paper",
    "C": "scissors",
    "Z": "scissors",
}

var win_score = map[string]int{
    "rock": 1,
    "paper": 2,
    "scissors": 3,
}

type Match struct {
    me string
    other string
}

func (m Match) result() int {
    if m.me == m.other {
        switch m.me {
        case "rock":
            return 3 + 1
        case "paper":
            return 3 + 2
        case "scissors":
            return 3 + 3
        default:
            return 0
        }
    }

    if m.me == "rock" {
        if m.other == "scissors" {
            return 6 + 1
        } else {
            return 1
        }
    }

    if m.me == "paper" {
        if m.other == "rock" {
            return 6 + 2
        } else {
            return 2
        }
    }

    if m.me == "scissors" {
        if m.other == "paper" {
            return 6 + 3
        } else {
            return 3
        }
    }

    return 0
}

func parse_input(filename string) []Match {
    s, err := os.ReadFile(filename)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    var moves []Match
    for _, line := range bytes.Split(s, []byte("\n")) {
        if len(line) == 0 {
            continue;
        }
        line := bytes.Split(line, []byte(" "))
        moves = append(moves, Match{
            decoding[string(line[1])],
            decoding[string(line[0])],
        })
    }

    return moves
}

func main() {
    part, input := parse_args()

    moves := parse_input(input)

    switch part {
    case 1: {
        score := 0
        for _, move := range moves {
            score += move.result()
        }
        fmt.Println(score)
    }
    case 2: {
        ;
    }
    default: {
        log.Println("not implemented for part %d", part)
        os.Exit(1)
    }
}}
