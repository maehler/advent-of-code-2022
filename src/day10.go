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

type Operation struct {
    name string
    value *int
}

func parse_input(filename string) []Operation {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        os.Exit(1)
    }

    reader := bufio.NewReader(file)

    operations := []Operation{}
    for {
        line, err := reader.ReadString(byte('\n'))
        if err == io.EOF {
            break
        }

        f := strings.Fields(line)

        if f[0] == "noop" {
            operations = append(operations, Operation{f[0], nil})
        } else {
            val, _ := strconv.Atoi(f[1])
            operations = append(operations, Operation{f[0], &val})
        }
    }

    return operations
}

func main() {
    var part int
    var input string
    flag.IntVar(&part, "part", 1, "part to run")
    flag.Parse()

    if flag.NArg() != 1 {
        fmt.Println("usage: [OPTIONS] input-file")
        flag.Usage()
        os.Exit(1)
    }

    input = flag.Arg(0)

    operations := parse_input(input)

    first_check_cycle := 20
    check_interval := 40

    switch part {
    case 1:
        x := 1
        cycle := 0
        strengths := []int{}
        for _, op := range operations {
            if cycle % check_interval == first_check_cycle {
                strengths = append(strengths, x * cycle)
            }
            if op.name == "noop" {
                cycle += 1
            }

            if op.name == "addx" {
                if cycle % check_interval >= first_check_cycle - 2 &&
                        cycle % check_interval < first_check_cycle {
                    cycle_diff := 20 - cycle % check_interval
                    strengths = append(strengths, x * (cycle + cycle_diff))
                }
                cycle += 2
                x += *op.value
            }

            if cycle >= 220 {
                break
            }
        }

        sum := 0
        for _, s := range strengths {
            sum += s
        }
        fmt.Println(sum)
    case 2:
        ;
    default:
        log.Printf("error: part %d not implemented", part)
        os.Exit(1)
    }
}
