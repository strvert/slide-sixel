package termutil

import (
    termios "github.com/k0kubun/go-termios"
)

type Termutil struct {
    defterm termios.Termios
}

func (term *Termutil) Init() {
    if err := term.defterm.GetAttr(termios.Stdin); err != nil {
        panic(err)
    }
}

func (term *Termutil) SetCanon() {
    var canonTerm termios.Termios

    if err := canonTerm.GetAttr(termios.Stdin); err != nil {
        panic(err)
    }

    canonTerm.LFlag ^= termios.ICANON
    if err := canonTerm.SetAttr(termios.Stdin, termios.TCSANOW); err != nil {
        panic(err)
    }
}

func (term *Termutil) SetUncanon() {
    if err := term.defterm.SetAttr(termios.Stdin, termios.TCSANOW); err != nil {
        panic(err)
    }
}
