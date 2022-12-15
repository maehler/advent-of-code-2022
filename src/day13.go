package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "sort"
    "strings"
)

func parse_input(filename string) []any {
    contents, err := os.ReadFile(filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        os.Exit(1)
    }

    packets := []any{}
    for _, pstring := range strings.Split(strings.TrimSpace(string(contents)), "\n\n") {
        var p1, p2 any
        lists := strings.Split(pstring, "\n")
        json.Unmarshal([]byte(lists[0]), &p1)
        json.Unmarshal([]byte(lists[1]), &p2)

        packets = append(packets, p1, p2)
    }

    return packets
}

func compare(a, b any) int {
    alist, a_ok := a.([]any)
    blist, b_ok := b.([]any)

    if !a_ok && !b_ok {
        // both are ints
        ai, _ := a.(float64)
        bi, _ := b.(float64)
        return int(ai) - int(bi)
    }

    if !a_ok {
        // wrap a in slice
        alist = []any{a}
    }

    if !b_ok {
        // wrap b in slice
        blist = []any{b}
    }

    // now both are lists
    for i := 0; i < len(alist) && i < len(blist); i++ {
        cmp := compare(alist[i], blist[i])
        if cmp < 0 {
            return cmp
        } else if cmp > 0 {
            return 1
        }
    }

    return len(alist) - len(blist)
}

func sort_packets(packets []any) []any {
    less_func := func(i, j int) bool {
        return compare(packets[i], packets[j]) <= 0
    }
    sort.Slice(packets, less_func)
    return packets
}

func main() {
    var part int
    flag.IntVar(&part, "part", 1, "part to run")
    flag.Parse()
    if flag.NArg() != 1 {
        fmt.Println("usage: day13 input-file")
        os.Exit(1)
    }
    
    filename := flag.Arg(0)

    packets := parse_input(filename)

    sorted_indices := []int{}

    switch part {
    case 1:
        sum := 0
        pair := 0
        for i := 0; i < len(packets); i += 2 {
            pair++

            p1 := packets[i]
            p2 := packets[i + 1]

            sorted := compare(p1, p2) <= 0
            if sorted {
                sorted_indices = append(sorted_indices, pair)
                sum += pair
            }
        }

        fmt.Println(sum)
    case 2:
        d1 := []any{[]any{6.}}
        d2 := []any{[]any{2.}}
        packets = append(packets, d1, d2)
        sort_packets(packets)

        res := 1
        for i, packet := range packets {
            if fmt.Sprintf("%v", packet) == "[[2]]" ||
                    fmt.Sprintf("%v", packet) == "[[6]]" {
                res *= i + 1
            }
        }
        fmt.Println(res)
    default:
        log.Printf("error: part %d not implemented", part)
        os.Exit(1)
    }
}
