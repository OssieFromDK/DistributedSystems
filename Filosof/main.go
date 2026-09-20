package main

import (
    "fmt"
    "time"
)

type SpisepindChV2 struct {
    id      int
    request chan chan bool
}

func SpisepindV2(id int) SpisepindChV2 {
    return SpisepindChV2{
        id:      id,
        request: make(chan chan bool),
    }
}

func RunSpisepind(pind SpisepindChV2) {
    for {
        replyChan := <-pind.request
        replyChan <- true
    }
}

func FilosofV2(navn string, venstreID, højreID int, venstrePind, højrePind SpisepindChV2) {
    for i := 0; i < 3; i++ {
        fmt.Println("Philosopher", navn, "tænker")
        time.Sleep(300 * time.Millisecond)

        if venstreID < højreID {
            tagPind(venstrePind)
            fmt.Println("Philosopher", navn, "tog venstre spisepind", venstreID)

            tagPind(højrePind)
            fmt.Println("Philosopher", navn, "tog højre spisepind", højreID)
        } else {
            tagPind(højrePind)
            fmt.Println("Philosopher", navn, "tog højre spisepind", højreID)

            tagPind(venstrePind)
            fmt.Println("Philosopher", navn, "tog venstre spisepind", venstreID)
        }

        fmt.Println("Philosopher", navn, "spiser (", i+1, "/3 )")
        time.Sleep(300 * time.Millisecond)

        fmt.Println("Philosopher", navn, "lagde begge pinde tilbage")
    }

    fmt.Println("Philosopher", navn, "er færdig med at spise")
}

func tagPind(pind SpisepindChV2) {
    reply := make(chan bool)
    pind.request <- reply
    <-reply
}

func main() {

    pind1 := SpisepindV2(0)
    pind2 := SpisepindV2(1)
    pind3 := SpisepindV2(2)
    pind4 := SpisepindV2(3)
    pind5 := SpisepindV2(4)

    go RunSpisepind(pind1)
    go RunSpisepind(pind2)
    go RunSpisepind(pind3)
    go RunSpisepind(pind4)
    go RunSpisepind(pind5)

    go FilosofV2("F1", 0, 1, pind1, pind2)
    go FilosofV2("F2", 1, 2, pind2, pind3)
    go FilosofV2("F3", 2, 3, pind3, pind4)
    go FilosofV2("F4", 3, 4, pind4, pind5)
    go FilosofV2("F5", 4, 0, pind5, pind1)

    time.Sleep(9 * time.Second)
}