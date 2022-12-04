package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
)

type Args struct {
    part int
    input string
}

func (a *Args) parse() {
    flag.IntVar(&a.part, "part", 1, "part to run")
    flag.Parse()

    if flag.NArg() != 1 {
        log.Println("usage: [OPTIONS] input-filename")
        os.Exit(1)
    }

    a.input = flag.Arg(0)
}

type Range struct {
    start int
    end int
}

func (r Range) overlaps(other Range) bool {
    return other.start >= r.start && other.start <= r.end ||
        other.end >= r.start && other.end <= r.end
}

func (r Range) contains(other Range) bool {
    return other.start >= r.start && other.end <= r.end
}

func parse_input(filename string) [][]Range {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error opening file: %s\n", err)
        os.Exit(1)
    }
    reader := bufio.NewReader(file)

    ranges := [][]Range{}
    current_range := Range{}
    current_pair := []Range{}
    buf := []rune{}
    for {
        r, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }
        if r == '-' {
            // start found
            current_range.start, _ = strconv.Atoi(string(buf))
            buf = nil
            continue
        }

        if r == ',' || r == '\n' {
            // end coordinate found
            current_range.end, _ = strconv.Atoi(string(buf))
            buf = nil
            
            current_pair = append(current_pair, current_range)
            current_range = Range{}

            if r == '\n' {
                ranges = append(ranges, current_pair)
                current_pair = nil
            }
            continue
        }

        buf = append(buf, r)
    }

    file.Close()

    return ranges
}

func main() {
    args := Args{}
    args.parse()

    ranges := parse_input(args.input)

    switch args.part {
    case 1:
        sum := 0
        for _, pair := range ranges {
            if pair[0].contains(pair[1]) || pair[1].contains(pair[0]) {
                sum++
            }
        }
        fmt.Println(sum)
    case 2:
        ;
    default:
        log.Printf("error: no such part: %d", args.part)
        os.Exit(1)
    }
}
