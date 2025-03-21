package config

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	x, _ := LoadConfig()
	fmt.Println(x.Database.GetDSN())

}
