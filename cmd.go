package buckis

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const commandSize = 512

type stateChangingCommand uint32

const (
	SAVE = iota
	DONT_SAVE
)

// TODO:
// - a way to set a struct obj to hashes in hset
// - write test for all data structures
// - continue more implementation for full-text, graph, like auto-sugesstion, finc nearest neighbors
// - look for means to support blockchain and list without the use of generics
// - further lookup of blockchain

const (
	SET stateChangingCommand = iota
	INCRBY

	HSET
	HINCRBY

	SADD
	SREM
	SMOVE

	ZADD
	ZRANGESTORE
	ZINCRBY
	ZREM

	RADD

	RPUSH
	LPUSH
	LPOP
	RPOP

	BCADD

	FTCREATE
)

type command struct {
	instruction stateChangingCommand
	args        []any
}

func (instruction stateChangingCommand) toInt() uint16 {
	return uint16(instruction)
}

func newCommand(instruction stateChangingCommand, args ...any) *command {
	return &command{instruction, args}
}

func (d *dict) Load() {
	info, _ := d.aof.Stat()

	var offset int64

	for offset = 0; offset < info.Size(); offset += commandSize {
		commandBytes := make([]byte, commandSize)

		_, err := d.aof.ReadAt(commandBytes, offset)
		if err != nil {
			log.Fatal(err)
		}

		cmd := command{}

		var pos uint16

		instruction := binary.BigEndian.Uint16(commandBytes[pos:])
		cmd.instruction = stateChangingCommand(instruction)
		pos += 2

		lengthOfArgs := binary.BigEndian.Uint16(commandBytes[pos:])
		pos += 2

		var args []any

		for i := uint16(0); i < lengthOfArgs; i++ {
			lengthOfArg := binary.BigEndian.Uint16(commandBytes[pos:])
			pos += 2

			var arg any
			err := json.Unmarshal(commandBytes[pos:pos+lengthOfArg], &arg)
			if err != nil {
				log.Fatal(err)
			}

			args = append(args, arg)

			pos += lengthOfArg
		}

		cmd.args = args

		fmt.Println(cmd)

		d.commandLoadQueue <- cmd
	}
}

func (d *dict) backgroundLoad() {
	for {
		select {
		case cmd := <-d.commandLoadQueue:
			d.runCommand(cmd)
		}
	}
}

func (d *dict) listenForCommands() {
	for cmd := range d.commandChan {
		//	fmt.Println(cmd, "commands")

		// to persist on disk we would need
		// 1) number of arguments
		// 2) length of each argument
		// 3) total offset
		// 4) command

		lengthOfArgs := uint16(len(cmd.args))

		var pos uint16

		commandBytes := make([]byte, commandSize)

		// stateChangingCommand
		binary.BigEndian.PutUint16(commandBytes[pos:], cmd.instruction.toInt())
		pos += 2

		// number of arguments
		binary.BigEndian.PutUint16(commandBytes[pos:], lengthOfArgs)
		pos += 2

		for _, arg := range cmd.args {
			b, err := json.Marshal(arg)

			if err != nil {
				log.Fatal(err)
			}

			lengthOfArg := uint16(len(b))

			binary.BigEndian.PutUint16(commandBytes[pos:], lengthOfArg)
			pos += 2

			copy(commandBytes[pos:], b)

			pos += lengthOfArg
		}

		stat, err := d.aof.Stat()
		if err != nil {
			log.Fatal(err)
		}

		d.aof.WriteAt(commandBytes, stat.Size())

		stat, err = d.aof.Stat()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(stat.Size())
	}
}

func (d *dict) Save() {
	pwd, _ := os.Getwd()

	args := os.Args

	fmt.Println(pwd, args)

	pid := os.Getpid()
	ppid := os.Getppid()
	log.Printf("pid: %d, ppid: %d, args: %s", pid, ppid, os.Args)

}
