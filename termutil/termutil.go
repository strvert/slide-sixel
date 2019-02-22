package termutil

import (
    "fmt"
    "strings"
    "strconv"

    termios "github.com/k0kubun/go-termios"

    "../rw"
)

type Termutil struct {
    defterm termios.Termios
    curterm termios.Termios
}

func (term *Termutil) Init() error {
    if err := term.defterm.GetAttr(termios.Stdin); err != nil {
        return err
    }
    term.curterm = term.defterm
    return nil
}

func (term *Termutil) SetUncanon() error {
    term.curterm.LFlag &^= termios.ICANON
    if err := term.curterm.SetAttr(termios.Stdin, termios.TCSANOW); err != nil {
        return err
    }
    return nil
}

func (term *Termutil) LoadBefore() error {
    if err := term.defterm.SetAttr(termios.Stdin, termios.TCSANOW); err != nil {
        return err
    }
    return nil
}

func (term *Termutil) SetEcho(state bool) error {
    if err := term.curterm.GetAttr(termios.Stdin); err != nil {
        return err
    }

    if state {
        term.curterm.LFlag &^= termios.ECHO
    } else {
        term.curterm.LFlag ^= termios.ECHO
    }

    if err := term.curterm.SetAttr(termios.Stdin, termios.TCSANOW); err != nil {
        return err
    }
    return nil
}



type CtrlSeqs struct {}

func (ctr *CtrlSeqs) GetWindowSize() (uint, uint, error) {
    var term Termutil
    term.Init()
    term.SetUncanon() // 非カノニカルモードに
    term.SetEcho(true)

    fmt.Print("\x1b[14;;t")
    input, err := rw.ReadUntil('t')
    if err != nil {
        return 0, 0, err
    }

    input = input[6:]
    strsize := strings.Split(input, ";")
    height, err := strconv.Atoi(strsize[0])
    if err != nil {
        return 0, 0, err
    }
    width, err := strconv.Atoi(strsize[1])
    if err != nil {
        return 0, 0, err
    }

    term.LoadBefore()
    return uint(height), uint(width), nil
}

func (ctr *CtrlSeqs) ClearScreen() {
    fmt.Print("\x1b[2J")
    fmt.Print("\x1b[H")
}

func (ctr *CtrlSeqs) ToggleFullScreen() {
    fmt.Print("\x1b[10;2;t")
}
