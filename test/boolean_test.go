package main

import (
    "testing"
)

func Best1(b *testing.B) {
    for n := 0; n < b.N; n++ {
        a := false
        b := true
        c := true
        z := false

        if (a && b) || (a && c) || (b && c) {
            z = true;
        }

        _ = z
    }
}

func Best2(b *testing.B) {
    for n := 0; n < b.N; n++ {
        a := false
        b := true
        c := true
        z := false

        if (a == b) || (a == c) || (b == c) {
            z = true;
        }

        _ = z
    }
}


func Benchmark1(b *testing.B)  { Best1(b) }
func Benchmark2(b *testing.B)  { Best2(b) }
