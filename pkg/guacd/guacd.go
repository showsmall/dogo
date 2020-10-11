package guacd

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const Delimiter = ';'

type Configuration struct {
	ConnectionID string
	Protocol     string
	Parameters   map[string]string
}

func NewConfiguration() (ret Configuration) {
	ret.Parameters = make(map[string]string)
	return ret
}

func (opt *Configuration) SetParameter(name, value string) {
	opt.Parameters[name] = value
}

func (opt *Configuration) UnSetParameter(name string) {
	delete(opt.Parameters, name)
}

func (opt *Configuration) GetParameter(name string) string {
	return opt.Parameters[name]
}

type Instruction struct {
	Opcode       string
	Args         []string
	ProtocolForm string
}

func NewInstruction(opcode string, args ...string) (ret Instruction) {
	ret.Opcode = opcode
	ret.Args = args
	return ret
}

func (opt *Instruction) String() string {
	if len(opt.ProtocolForm) > 0 {
		return opt.ProtocolForm
	}

	opt.ProtocolForm = fmt.Sprintf("%d.%s", len(opt.Opcode), opt.Opcode)
	for _, value := range opt.Args {
		opt.ProtocolForm += fmt.Sprintf(",%d.%s", len(value), value)
	}
	opt.ProtocolForm += string(Delimiter)
	return opt.ProtocolForm
}

func (opt *Instruction) Parse(content string) Instruction {
	messages := strings.Split(content, ",")

	var args = make([]string, len(messages))
	for i := range messages {
		lm := strings.Split(messages[i], ".")
		args[i] = lm[1]
	}
	return NewInstruction(args[0], args[1:]...)
}

type Tunnel struct {
	rw     *bufio.ReadWriter
	UUID   string
	Config Configuration
}

func NewTunnel(address string, config Configuration) (ret Tunnel) {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("dial error %s", err.Error()))
	}
	ret.rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	ret.Config = config

	selectArg := config.ConnectionID
	if selectArg == "" {
		selectArg = config.Protocol
	}

	ret.WriteInstruction(NewInstruction("select", selectArg))
	args := ret.expect("args")

	// send size
	ret.WriteInstruction(NewInstruction("size", "1024", "768", "96"))

	ret.WriteInstruction(NewInstruction("audio", ""))
	ret.WriteInstruction(NewInstruction("video", ""))
	ret.WriteInstruction(NewInstruction("image", ""))

	parameters := make([]string, len(args.Args))
	for i := range args.Args {
		argName := args.Args[i]
		parameters[i] = config.GetParameter(argName)
	}
	// send connect
	ret.WriteInstruction(NewInstruction("connect", parameters...))

	ready := ret.expect("ready")

	if len(ready.Args) == 0 {
		panic("No connection ID received")
	}

	ret.UUID = ready.Args[0]
	return ret
}

func (opt *Tunnel) WriteInstruction(instruction Instruction) {
	_, err := opt.Write([]byte(instruction.String()))
	if err != nil {
		panic(fmt.Sprintf("write message error %s", err.Error()))
	}
}

func (opt *Tunnel) Write(p []byte) (int, error) {
	nn, err := opt.rw.Write(p)
	if err != nil {
		return nn, err
	}
	err = opt.rw.Flush()
	if err != nil {
		return nn, err
	}
	return nn, nil
}

func (opt *Tunnel) ReadInstruction() (instruction Instruction, err error) {
	msg, err := opt.rw.ReadString(Delimiter)
	if err != nil {
		return instruction, err
	}
	return instruction.Parse(msg), err
}

func (opt *Tunnel) Read() ([]byte, error) {
	return opt.rw.ReadBytes(Delimiter)
}

func (opt *Tunnel) expect(opcode string) (instruction Instruction) {
	instruction, err := opt.ReadInstruction()
	if err != nil {
		panic(fmt.Sprintf("read instruction error %s", err.Error()))
	}

	if opcode != instruction.Opcode {
		panic(fmt.Sprintln(`Expected "%s" instruction but instead received "%s"`, opcode, instruction.Opcode))
	}
	return instruction
}
