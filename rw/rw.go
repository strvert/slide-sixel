package rw

import (
    "bufio"
    "strings"
    "os"
)

func ReadUntil(delim byte) (string, error) {
    stdin := bufio.NewReader(os.Stdin)
    var retbuf []byte
    LOOP:
    for {
        ch, err := stdin.ReadByte()
        if err != nil {
            return "none", err
        }
        if ch == delim {
            break LOOP
        }
        retbuf = append(retbuf, ch)
    }
    retstr := string(retbuf)
    return retstr, nil
}

func ScanCommand() ([]string, error) {
    input, err := ReadUntil('\n')
    if err != nil {
        return nil, err
    }
    commands := strings.Split(input, " ")
    return commands, nil
}

