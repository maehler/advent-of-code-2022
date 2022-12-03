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

var move_score = map[string]int{
    "rock": 1,
    "paper": 2,
    "scissors": 3,
}

type Match struct {
    me string
    other string
}

func (m Match) rigged_result() int {
    score := 0
    switch m.me {
    case "rock": // X: lose
        score += 0
    case "paper": // Y: draw
        score += 3
    case "scissors": // Z: win
        score += 6
    default:
        // unreachable
        log.Println("error: you can't be here")
        os.Exit(1)
    }

    switch m.me {
    case "rock":
        // lose
        switch m.other {
        case "rock":
            score += move_score["scissors"]
        case "paper":
            score += move_score["rock"]
        case "scissors":
            score += move_score["paper"]
        }
    case "paper":
        // draw
        score += move_score[m.other]
    case "scissors":
        // win
        switch m.other {
        case "rock":
            score += move_score["paper"]
        case "paper":
            score += move_score["scissors"]
        case "scissors":
            score += move_score["rock"]
        }
    default:
        // unreachable
        log.Println("error: %s doesn't exist in this game", m.other)
        os.Exit(1)
    }

    return score
}

func (m Match) result() int {
    score := move_score[m.me]
    if m.me == m.other {
        // draw
        score += 3
    }

    // win
    if m.me == "rock" && m.other == "scissors" {
        return score + 6
    }

    if m.me == "paper" && m.other == "rock" {
        return score + 6
    }

    if m.me == "scissors" && m.other == "paper" {
        return score + 6
    }

    // lose
    return score
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
        score := 0
        for _, move := range moves {
            score += move.rigged_result()
        }
        fmt.Println(score)
    }
    default: {
        log.Println("not implemented for part %d", part)
        os.Exit(1)
    }
}}
