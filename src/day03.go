package main

import (
    "bytes"
    "flag"
    "fmt"
    "log"
    "os"
)

type Args struct {
    part int
    input string
}

func (a *Args) parse() {
    flag.IntVar(&a.part, "part", 1, "the part to run")
    flag.Parse()

    if flag.NArg() != 1 {
        log.Println("usage: [OPTIONS] input-file")
        os.Exit(1)
    }

    a.input = flag.Arg(0)
}

type Knapsack struct {
    c1 string
    c2 string
}

func (k Knapsack) common_items() []rune {
    var common []rune
    for _, i1 := range []rune(k.c1) {
        for _, i2 := range []rune(k.c2) {
            if i1 == i2 {
                if !contains(common, i1) {
                    common = append(common, i1)
                }
            }
        }
    }
    return common
}

func contains(runes []rune, x rune) bool {
    for _, r := range runes {
        if r == x {
            return true
        }
    }
    return false
}

func distinct(runes []rune) []rune {
    rune_map := make(map[rune] int)

    for _, r := range runes {
        rune_map[r]++
    }

    distinct_runes := make([]rune, len(rune_map))
    
    i := 0
    for k := range rune_map {
        distinct_runes[i] = k
        i++
    }

    return distinct_runes
}

func get_badge_item(knapsacks []Knapsack) rune {
    items := make(map[rune]int)
    for _, k := range knapsacks {
        for _, r := range distinct([]rune(k.c1)) {
            items[r]++
        }
        for _, r := range distinct([]rune(k.c2)) {
            items[r]++
        }
    }

    for r, c := range items {
        if c == 3 {
            return r
        }
    }

    log.Println("error: no bade item found")
    os.Exit(1)

    return 0
}

func get_priority_map() map[rune]int {
    priority_map := make(map[rune]int)
    for i, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
        priority_map[c] = i + 1
    }
    return priority_map
}

func parse_input(filename string) []Knapsack {
    b, err := os.ReadFile(filename)
    if err != nil {
        log.Println("error: %s", err)
        os.Exit(1)
    }

    var knapsacks []Knapsack
    for _, line := range bytes.Split(b, []byte("\n")) {
        if len(line) == 0 {
            continue
        }
        c1 := line[:len(line)/2]
        c2 := line[len(line)/2:]
        knapsacks = append(knapsacks, Knapsack{string(c1), string(c2)})
    }

    return knapsacks
}

func main() {
    args := Args{}
    args.parse()

    knapsacks := parse_input(args.input)

    priority_map := get_priority_map()

    switch args.part {
    case 1:
        sum := 0
        for _, k := range knapsacks {
            common := k.common_items()
            priority := 0
            for _, c := range common {
                priority += priority_map[c]
            }
            sum += priority
        }
        fmt.Println(sum)
    case 2:
        sum := 0
        for i := 0; i < len(knapsacks); i += 3 {
            badge_item := get_badge_item(knapsacks[i:i+3])
            sum += priority_map[badge_item]
        }
        fmt.Println(sum)
    default:
        log.Println("error: no part %d", args.part)
    }
}
