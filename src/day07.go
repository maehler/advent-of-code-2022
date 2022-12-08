package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "sort"
    "strings"
    "strconv"
)

type FileInterface interface {
    getname() string
    getsize() int
    isdir() bool
}

type Tree struct {
    root *DirNode
}

type DirNode struct {
    name string
    parent *DirNode
    dirs []*DirNode
    files []FileNode
    // children []*FileInterface
}

func (d DirNode) isdir() bool {
    return true
}

func (d DirNode) getname() string {
    return d.name
}

func (d DirNode) getsize() int {
    sum := 0
    for _, fi := range d.dirs {
        sum += (*fi).getsize()
    }
    for _, fi := range d.files {
        sum += fi.getsize()
    }
    return sum
}

func (d DirNode) contains(fname string) FileInterface {
    for _, fi := range d.dirs {
        if (*fi).getname() == fname { 
            return fi
        }
    }
    for _, fi := range d.files {
        if fi.getname() == fname {
            return fi
        }
    }
    return nil
}

func (d *DirNode) addchild(f FileInterface) {
    if f.isdir() {
        fd, _ := f.(DirNode)
        d.dirs = append(d.dirs, &fd)
    } else {
        fd, _ := f.(FileNode)
        d.files = append(d.files, fd)
    }
}

type FileNode struct {
    name string
    size int
}

func (f FileNode) isdir() bool {
    return false
}

func (f FileNode) getsize() int {
    return f.size
}

func (f FileNode) getname() string {
    return f.name
}

type Args struct {
    part int
    input string
}

func (a *Args) parse() {
    flag.IntVar(&a.part, "part", 1, "part to run")
    flag.Parse()

    if flag.NArg() != 1 {
        fmt.Println("usage: [OPTIONS] input-file")
        flag.Usage()
        os.Exit(1)
    }

    a.input = flag.Arg(0)
}


type Command string

const (
    cd Command = "cd"
    ls         = "ls"
)

type CommandLine struct {
    exec Command
    args []string
}

func parse_command(line string) CommandLine {
    args := strings.Fields(line)

    c := CommandLine{Command(args[0]), args[1:]}

    return c
}

func parse_ls(line string) FileInterface {
    f := strings.Fields(line)
    if f[0] == "dir" {
        return DirNode{f[1], nil, nil, nil}
    }

    size, _ := strconv.Atoi(f[0])

    return FileNode{f[1], size}
}

func parse_input(filename string) *Tree {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s", err)
        os.Exit(1)
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    
    tree := &Tree{}
    var current_dir *DirNode

    for {
        r, _, err := reader.ReadRune()
        if err == io.EOF {
            break
        }

        if r == '$' {
            // command
            line, err := reader.ReadString(byte('\n'))
            if err == io.EOF {
                break
            }
            command := parse_command(line)
            if command.exec == cd {
                // change dir and add to tree
                if current_dir == nil {
                    current_dir = &DirNode{command.args[0], nil, nil, nil}
                    tree.root = current_dir
                } else if command.args[0] == ".." {
                    // go to parent
                    current_dir = current_dir.parent
                } else {
                    edir := current_dir.contains(command.args[0])
                    if edir != nil {
                        edir, ok := edir.(*DirNode)
                        if !ok {
                            log.Println("error: could not get *DirNode from FileInterface")
                            os.Exit(1)
                        }
                        // go to already created dir
                        current_dir = edir
                    } else {
                        // create new directory
                        d := &DirNode{command.args[0], current_dir, nil, nil}
                        // add to current directory children
                        current_dir.addchild(d)
                        current_dir = d
                    }
                }
            } else if command.exec == ls {
                // parse directory contents
                if len(command.args) != 0 {
                    log.Println("error: arguments to ls not supported")
                    os.Exit(1)
                }
                for {
                    line, err := reader.ReadString(byte('\n'))
                    if err == io.EOF {
                        break
                    }
                    
                    f := parse_ls(line)
                    if f.isdir() {
                        dirnode := f.(DirNode)
                        dirnode.parent = current_dir
                        current_dir.addchild(dirnode)
                    } else {
                        current_dir.addchild(f)
                    }

                    next_byte, err := reader.Peek(1)
                    if err == io.EOF || next_byte[0] == byte('$') {
                        break
                    }
                }
            }
        }
    }

    return tree
}

func (d *DirNode) print(indent string) {
    fmt.Printf("%s%s: %d/\n", indent, d.getname(), d.getsize())
    indent = fmt.Sprintf("%s  ", indent)
    for _, dir := range d.dirs {
        dir.print(indent)
    }
    for _, file := range d.files {
        fmt.Printf("%s%s: %d\n", indent, file.getname(), file.getsize())
    }
}

func (d *DirNode) dirsizes() []int {
    sizes := []int{}
    sizes = append(sizes, d.getsize())
    for _, dir := range d.dirs {
        sizes = append(sizes, dir.dirsizes()...)
    }
    return sizes
}

func main() {
    args := Args{}
    args.parse()

    tree := parse_input(args.input)

    tree.root.print("")

    switch args.part {
    case 1:
        sum := 0
        sizes := tree.root.dirsizes()
        for _, s := range sizes {
            if s < 100000 {
                sum += s
            }
        }
        fmt.Println(sum)
    case 2:
        total_size := 70_000_000
        free_space := total_size - tree.root.getsize()
        sizes := tree.root.dirsizes()
        sort.Ints(sizes)
        for _, s := range sizes {
            if s + free_space >= 30_000_000 {
                fmt.Println(s)
                break
            }
        }
    default:
        log.Printf("error: no implementation for part %d", args.part)
        os.Exit(1)
    }
}
