package redis

import (
	"github.com/redis-go/redcon"
)

// Command flags. Please check the command table defined in the redis.c file
// for more information about the meaning of every flag.
const (
	CMD_WRITE             CmdFlag = iota + 1 /* "w" flag */
	CMD_READONLY                             /* "r" flag */
	CMD_DENYOOM                              /* "m" flag */
	CMD_MODULE                               /* Command exported by module. */
	CMD_ADMIN                                /* "a" flag */
	CMD_PUBSUB                               /* "p" flag */
	CMD_NOSCRIPT                             /* "s" flag */
	CMD_RANDOM                               /* "R" flag */
	CMD_SORT_FOR_SCRIPT                      /* "S" flag */
	CMD_LOADING                              /* "l" flag */
	CMD_STALE                                /* "t" flag */
	CMD_SKIP_MONITOR                         /* "M" flag */
	CMD_ASKING                               /* "k" flag */
	CMD_FAST                                 /* "F" flag */
	CMD_MODULE_GETKEYS                       /* Use the modules getkeys interface. */
	CMD_MODULE_NO_CLUSTER                    /* Deny on Redis Cluster. */
)

// A command can be registered.
type Command struct {
	// The command name.
	name string

	// Handler
	handler CommandHandler

	// Command flags
	flags map[CmdFlag]struct{} // Use map as a set data structure
}

func NewCommand(name string, handler CommandHandler, flags ...CmdFlag) *Command {
	mFlags := make(map[CmdFlag]struct{}, len(flags))
	for _, f := range flags {
		mFlags[f] = struct{}{}
	}

	return &Command{
		name:    name,
		handler: handler,
		flags:   mFlags,
	}
}

// Command flag type.
type CmdFlag uint

// Commands map
type Commands map[string]*Command

// The CommandHandler is triggered when the received
// command equals a registered command.
//
// However the CommandHandler is executed by the Handler,
// so if you implement an own Handler make sure the CommandHandler is called.
type CommandHandler func(c *Client, cmd redcon.Command)

// Is called when a request is received,
// after Accept and if the command is not registered.
//
// However UnknownCommand is executed by the Handler,
// so if you implement an own Handler make sure to include UnknownCommand.
type UnknownCommand func(c *Client, cmd redcon.Command)

// Gets registered commands name.
func (cmd *Command) Name() string {
	return cmd.name
}

// RegisterCommands adds commands to the redis instance.
// If a cmd already exists the handler is overridden.
func (r *Redis) RegisterCommands(cmds []*Command) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	for _, cmd := range cmds {
		r.registerCommand(cmd)
	}
}

// RegisterCommand adds a command to the redis instance.
// If cmd already exists the handler is overridden.
func (r *Redis) RegisterCommand(cmd *Command) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.registerCommand(cmd)
}
func (r *Redis) registerCommand(cmd *Command) {
	r.getCommands()[cmd.Name()] = cmd
}

// UnregisterCommand removes a command.
func (r *Redis) UnregisterCommand(name string) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	delete(r.commands, name)
}

// Command returns the registered command or nil if not exists.
func (r *Redis) Command(name string) *Command {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.command(name)
}

func (r *Redis) command(name string) *Command {
	return r.commands[name]
}

// Commands returns the commands map.
func (r *Redis) Commands() Commands {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.getCommands()
}

func (r *Redis) getCommands() Commands {
	return r.commands
}

// CommandExists checks if one or more commands are registered.
func (r *Redis) CommandExists(cmds ...string) bool {
	regCmds := r.Commands()

	// TODO does this make the performance better because it does not create a loop every time?
	if len(cmds) == 1 {
		_, ex := regCmds[cmds[0]]
		return ex
	}

	for _, cmd := range cmds {
		if _, ex := regCmds[cmd]; !ex {
			return false
		}
	}
	return true
}

// FlushCommands removes all commands.
func (r *Redis) FlushCommands() {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.commands = make(Commands)
}

// CommandHandlerFn returns the CommandHandler of cmd.
func (r *Redis) CommandHandlerFn(name string) CommandHandler {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.command(name).handler
}

// UnknownCommandFn returns the UnknownCommand function.
func (r *Redis) UnknownCommandFn() UnknownCommand {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.unknownCommand
}
