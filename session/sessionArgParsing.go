package argparsing

import (
    "os"
    "bufio"
    "io"
    "strconv"
    "fmt"
)

const DefaultPoolSize int = 256
const DefaultTimeout int = 30000

//reads a file into a []string split at '\n'
func readFile(path string) ([]string, error) {
    ioFile, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer ioFile.Close()

    reader := bufio.NewReader(ioFile)
    file := []string{}
    for {
        line, err := reader.ReadString(byte('\n'))
        if len(line) > 1 && line[len(line)-1] == '\n' {
            line = line[:len(line)-1] //removing the trailing '\n'
        }
        file = append(file, line)
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }
    }
    return file, nil
}

//struct containing settings args for current session
type SessionArgs struct {
    BaseUrl string
    Words []string
    UserAgents []string
    PoolSize int
    Timeout int
}

//parses os args for program flags
func (s *SessionArgs) parseFlags(flagArgs []string) error {
    for i := 0; i < len(flagArgs); i++ {
        arg := flagArgs[i]
        var err error
        switch arg {
        case "-ua": //User-Agent flag
            if i >= len(flagArgs)-1 {
                return fmt.Errorf("error: no arg provided for flag: '%s'", arg)
            }
            s.UserAgents, err = readFile(flagArgs[i+1])
            if err != nil {
                return fmt.Errorf("error reading User-Agents file: %s", err.Error())
            }
            i += 1
        case "-p": //pool-size flag
            if i >= len(flagArgs)-1 {
                return fmt.Errorf("error: no arg provided for flag: '%s'", arg)
            }
            s.PoolSize, err = strconv.Atoi(flagArgs[i+1])
            if err != nil {
                return fmt.Errorf("error parsing pool-size arg: %s", err.Error())
            }
            i += 1
        case "-t": //timeout flag
            if i >= len(flagArgs)-1 {
                return fmt.Errorf("error: no arg provided for flag: '%s'", arg)
            }
            s.Timeout, err = strconv.Atoi(flagArgs[i+1])
            if err != nil {
                return fmt.Errorf("error parsing timeout arg: %s", err.Error())
            }
            i += 1
        default:
            return fmt.Errorf("error: invalid flag: '%s'", arg)
        }
    }
    return nil
}

//parses os args for program arguments/flags
func (s *SessionArgs) ParseArgs(args []string) error {    
    if len(args) < 3 {
        return fmt.Errorf("error: insufficient args: min: 2 (URL, path-to-words-file), got: %d", len(args)-1)
    }
    //url
    s.BaseUrl = args[1]
    //words
    var err error
    s.Words, err = readFile(args[2])
    if err != nil {
        return fmt.Errorf("error reading words file: %s", err.Error())
    }
    //remaining args (user-agents, timeout, pool-size)
    if len(args) > 3 {
        err := s.parseFlags(args[3:])
        if err != nil {
            return err
        }
    }
    return nil
}
