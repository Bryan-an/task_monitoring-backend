package main_test

import (
	"testing"

	main "github.com/Bryan-an/tasker-backend"
	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	greeting := main.SayHello("Bryan")
	assert.Equal(t, "Hello Bryan!", greeting)

	greeting2 := main.SayHello("Melissa")
	assert.Equal(t, "Hello Melissa!", greeting2)
}
