package main

import (
	"context"
	"os"

	"github.com/IxDay/templ-exp/templates"
)

func main() {
	component := templates.Hello("John")
	component.Render(context.Background(), os.Stdout)
}
