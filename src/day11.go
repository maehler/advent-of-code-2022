package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "sort"
    "strconv"
    "strings"
)

type KeepAway struct {
    monkeys []Monkey
}

func (ka *KeepAway) AddMonkey(m Monkey) *Monkey {
    ka.monkeys = append(ka.monkeys, m)
    return &ka.monkeys[len(ka.monkeys) - 1]
}

func (ka *KeepAway) PlayRound(less_worry bool, lcm int) {
    for mi := range ka.monkeys {
        for ii := range ka.monkeys[mi].items {
            ka.monkeys[mi].ModifyWorry(ii)
            if less_worry {
                ka.monkeys[mi].DecreaseWorry(ii)
            }
            ka.monkeys[mi].items[ii] = ka.monkeys[mi].items[ii] % lcm
            to_monkey := ka.monkeys[mi].Test(ii)
            ka.monkeys[to_monkey].Add(ka.monkeys[mi].items[ii])
        }
        ka.monkeys[mi].ClearItems()
    }
}

type Monkey struct {
    items []int
    items_viewed int
    operator string
    operatorvalues []string
    testval int
    test_true_to int
    test_false_to int
}

func (m *Monkey) Add(item int) {
    m.items = append(m.items, item)
}

func (m *Monkey) ClearItems() {
    m.items = nil
}

func (m *Monkey) ModifyWorry(item int) {
    var var1 int
    var var2 int

    if m.operatorvalues[0] == "old" {
        var1 = m.items[item]
    } else {
        var1, _ = strconv.Atoi(m.operatorvalues[0])
    }
    
    if m.operatorvalues[1] == "old" {
        var2 = m.items[item]
    } else {
        var2, _ = strconv.Atoi(m.operatorvalues[1])
    }

    var result int
    switch m.operator {
    case "+":
        result = var1 + var2
    case "-":
        result = var1 - var2
    case "*":
        result = var1 * var2
    case "/":
        result = var1 / var2
    }

    m.items_viewed++

    m.items[item] = result
}

func (m *Monkey) DecreaseWorry(item int) {
    m.items[item] = m.items[item] / 3
}

func (m *Monkey) Test(item int) int {
    if m.items[item] % m.testval == 0 {
        return m.test_true_to
    }
    return m.test_false_to
}

func parse_input(filename string) KeepAway {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s", err)
        os.Exit(1)
    }

    reader := bufio.NewReader(file)

    ka := KeepAway{}

    for {
        line, err := reader.ReadString(byte('\n'))
        if err == io.EOF {
            break
        }

        if strings.HasPrefix(line, "Monkey") {
            m := ka.AddMonkey(Monkey{})

            for {
                line, err := reader.ReadString(byte('\n'))
                if err == io.EOF || len(strings.Trim(line, "\n")) == 0 {
                    break
                }

                line = strings.Trim(line, " \n")

                if strings.HasPrefix(line, "Starting items:") {
                    line = strings.Replace(line, "Starting items: ", "", 1)
                    for _, item := range strings.Split(line, ", ") {
                        i, _ := strconv.Atoi(item)
                        (*m).Add(i)
                    }
                }

                if strings.HasPrefix(line, "Operation:") {
                    line = strings.Replace(line, "Operation: new = ", "", 1)
                    f := strings.Fields(line)
                    (*m).operator = f[1]
                    (*m).operatorvalues = []string{f[0], f[2]}
                }

                if strings.HasPrefix(line, "Test:") {
                    modval, _ := strconv.Atoi(strings.Fields(line)[3])
                    (*m).testval = modval
                    for i := 0; i < 2; i++ {
                        line, _ = reader.ReadString(byte('\n'))
                        line = strings.Trim(line, " ")
                        testto, _ := strconv.Atoi(strings.Fields(line)[5])
                        if strings.HasPrefix(line, "If true:") {
                            (*m).test_true_to = testto
                        } else {
                            (*m).test_false_to = testto
                        }
                    }
                }
            }
        }
    }
    return ka
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

    keepaway := parse_input(input)
    lcm := 1
    for _, m := range keepaway.monkeys {
        lcm *= m.testval
    }

    switch part {
    case 1:
        for i := 0; i < 20; i++ {
            keepaway.PlayRound(true, lcm)
        }
        n_items := make([]int, len(keepaway.monkeys))
        for mi, m := range keepaway.monkeys {
            n_items[mi] = m.items_viewed
        }
        sort.Slice(n_items, func(i, j int) bool { return n_items[j] < n_items[i] })
        fmt.Println(n_items[0] * n_items[1])
    case 2:
        for i := 0; i < 10_000; i++ {
            keepaway.PlayRound(false, lcm)
        }
        n_items := make([]int, len(keepaway.monkeys))
        for mi, m := range keepaway.monkeys {
            n_items[mi] = m.items_viewed
        }
        sort.Slice(n_items, func(i, j int) bool { return n_items[j] < n_items[i] })
        fmt.Println(n_items[0] * n_items[1])
    default:
        log.Printf("error: part %d not implemented\n", part)
        os.Exit(1)
    }
}
