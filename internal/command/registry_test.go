package command

import "testing"

type stubCommand struct {
	name string
}

func (c stubCommand) Name() string {
	return c.name
}

func (stubCommand) Description() string {
	return "test command"
}

func (stubCommand) Execute() Result {
	return Result{}
}

func TestRegistryRegisterAndFind(t *testing.T) {
	registry := NewRegistry()
	expected := stubCommand{name: "test"}

	registry.Register(expected)

	actual, found := registry.Find("test")
	if !found {
		t.Fatal("Find(test) did not find registered command")
	}
	if actual.Name() != expected.Name() {
		t.Errorf("Find(test).Name() = %q, want %q", actual.Name(), expected.Name())
	}
}

func TestRegistryCommandsAreSorted(t *testing.T) {
	registry := NewRegistry()
	registry.Register(stubCommand{name: "zebra"})
	registry.Register(stubCommand{name: "alpha"})

	commands := registry.Commands()

	if len(commands) != 2 {
		t.Fatalf("len(Commands()) = %d, want 2", len(commands))
	}
	if commands[0].Name() != "alpha" || commands[1].Name() != "zebra" {
		t.Errorf(
			"Commands() names = [%q, %q], want [alpha, zebra]",
			commands[0].Name(),
			commands[1].Name(),
		)
	}
}
