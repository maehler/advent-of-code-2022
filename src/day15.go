package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strings"
    "strconv"
)

func Abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

type Pos struct {
    x, y int
}

func (p Pos) Manhattan(other Pos) int {
    x_diff := Abs(p.x - other.x)
    y_diff := Abs(p.y - other.y)
    return x_diff + y_diff
}

type Sensor struct {
    pos Pos
    closest_beacon *Beacon
}

func (s Sensor) Radius() int {
    return s.pos.Manhattan(s.closest_beacon.pos)
}

type Beacon struct {
    pos Pos
}

func parse_input(filename string) (map[Pos]Sensor, map[Pos]*Beacon) {
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        os.Exit(1)
    }

    reader := bufio.NewReader(file)

    sensors := map[Pos]Sensor{}
    beacons := map[Pos]*Beacon{}

    for {
        line, err := reader.ReadString(byte('\n'))
        if err == io.EOF {
            break
        }

        line = strings.TrimSpace(line)
        components := strings.Fields(line)

        sensor_x, _ := strconv.Atoi(strings.Split(strings.Trim(components[2], ","), "=")[1])
        sensor_y, _ := strconv.Atoi(strings.Split(strings.Trim(components[3], ":"), "=")[1])

        beacon_x, _ := strconv.Atoi(strings.Split(strings.Trim(components[8], ","), "=")[1])
        beacon_y, _ := strconv.Atoi(strings.Split(components[9], "=")[1])

        beacon_pos := Pos{beacon_x, beacon_y}
        _, ok := beacons[beacon_pos]
        if !ok {
            beacons[beacon_pos] = &Beacon{beacon_pos}
        }

        sensor_pos := Pos{sensor_x, sensor_y}
        sensors[sensor_pos] = Sensor{sensor_pos, beacons[beacon_pos]}
    }

    return sensors, beacons
}

func main() {
    var part int
    flag.IntVar(&part, "part", 1, "part to run")
    flag.Parse()

    input := flag.Arg(0)

    sensors, beacons := parse_input(input)

    switch part {
    case 1:
        min_x := 0
        max_x := 0
        for _, s := range sensors {
            if s.pos.x + s.Radius() > max_x {
                max_x = s.pos.x + s.Radius()
            }
            if s.pos.x - s.Radius() < min_x {
                min_x = s.pos.x - s.Radius()
            }
        }

        n_spots := 0
        for i := min_x; i <= max_x; i++ { 
            p := Pos{i, 2_000_000}
            within_radius := false
            _, b_ok := beacons[p]
            for _, s := range sensors {
                if p.Manhattan(s.pos) <= s.Radius() && !b_ok {
                    within_radius = true
                    break
                }
            }
            if within_radius {
                n_spots++
            }
        }
        fmt.Println(n_spots)

    case 2:
        ;
    default:
        log.Printf("error: part %d not implemented\n", part)
        os.Exit(1)
    }
}
