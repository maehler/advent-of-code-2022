package main

import (
    "bytes"
    "flag"
    "fmt"
    "log"
    "os"
    "sort"
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

func int_slice_sum(s []int) (int) {
    var sum int;
    for _, v := range s {
        sum += v
    }
    return sum
}

func calorie_sums(elves [][]int) ([]int) {
    var sums []int
    for _, elf := range elves {
        sum := int_slice_sum(elf)
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

func top3_calories(calories []int) ([]int) {
    top3 := calories
    sort.Ints(top3)
    return top3[len(top3)-3:]
}

func main() {
    var args = parse_args()

    elves := parse_input(args.input)
    calories := calorie_sums(elves)

    switch args.part {
    case 1:
        max := max_calories(calories)
        fmt.Println(max)
    case 2:
        top3 := top3_calories(calories)
        top3_sum := int_slice_sum(top3)
        fmt.Println(top3_sum)
    default:
        fmt.Fprintf(os.Stderr, "part %d not implemented\n", args.part)
        os.Exit(1)
    }
}
