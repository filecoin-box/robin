package parse

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	c := new(Config)
	err := runtime_viper.Unmarshal(c)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
}
