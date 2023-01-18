package reports

import (
	"log"
	"testing"
)

func Test_getCountStringFromOptimization(t *testing.T) {
	k := map[string][]string{
		"a": {"1", "2", "4", "5"},
		"b": {"1", "2", "3", "4", "5"},
		"c": {"1", "5"},
	}

	log.Print(getCountStringFromOptimization(k))
	t.Fail()
}

func Test_nearest10(t *testing.T) {
	log.Print(nearest10(12))
	log.Print(nearest10(22))
	log.Print(nearest10(52))
	log.Print(nearest10(132))
	log.Print(nearest10(222))
	log.Print(nearest10(1322))
	log.Print(nearest10(52332))
	log.Print(nearest10(355254))
	log.Print(nearest10(7323122))

}
