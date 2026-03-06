package internal

type Command struct {
	Args    []string
	InFile  string
	OutFile string
}

type Pipeline struct {
	Commands []Command
}
