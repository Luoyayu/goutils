//+build darwin

package main

import (
    "golang.org/x/crypto/ssh/terminal"
    "log"
    "os"
    "syscall"
)

func read1() {
    fd := int(os.Stdin.Fd())
    if !terminal.IsTerminal(fd) {
        panic("stdin is not a terminal; skipping!")
    }

    st, err := terminal.GetState(fd)
    if err != nil {
        log.Fatalf("failed to get terminal state from GetState: %s", err)
    }

    defer terminal.Restore(fd, st)
    raw, err := terminal.MakeRaw(fd)
    if err != nil {
        log.Fatalf("failed to get terminal state from MakeRaw: %s", err)
    }

    if *st != *raw {
        log.Fatalf("states do not match; was %v, expected %v", raw, st)
    }
    for {
        p := make([]byte, 4)

        _, err = syscall.Read(fd, p)
        if err != nil {
            panic(err)
        } else {
            log.Println(p[0])
        }
    }
}

func main() {
    read1()
}
