package main

import (
    "bytes"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
)

type Args struct {
    part int
    input string
}

func parse_args() (Args) {
    var part int
    var input string
    flag.IntVar(&part, "part", 1, "part to run")

    flag.Parse()

    if flag.NArg() == 0 {
        fmt.Println("usage: [OPTIONS] input-file")
        flag.Usage()
        os.Exit(1)
    }

    input = flag.Arg(0)

    return Args{part, input}
}

func parse_input(filename string) ([][]int) {
    s, _ := os.ReadFile(filename)
    var elves [][]int

    var elf []int;
    for _, line := range bytes.Split(s, []byte("\n")) {
        if len(line) == 0 {
            elves = append(elves, elf)
            elf = nil
            continue
        }
        calories, err := strconv.Atoi(string(line))
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
        elf = append(elf, calories)
    }

    return elves
}

func calorie_sums(elves [][]int) ([]int) {
    var sums []int
    for _, elf := range elves {
        sum := 0
        for _, calories := range elf {
            sum += calories
        }
        sums = append(sums, sum)
    }
    return sums
}

func max_calories(calories []int) (int) {
    var max int
    for _, v := range calories {
        if v > max {
            max = v
        }
    }
    return max
}

func main() {
    var args = parse_args()

    elves := parse_input(args.input)

    switch args.part {
    case 1:
        calories := calorie_sums(elves)
        max := max_calories(calories)
        fmt.Println(max)
    case 2:
        ;
    default:
        fmt.Fprintf(os.Stderr, "part %d not implemented\n", args.part)
        os.Exit(1)
    }
}
