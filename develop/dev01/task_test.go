package main

import (
	"testing"
)

func TestCorrect(t *testing.T) {
	getTime("0.beevik-ntp.pool.ntp.org")
}

func TestUnCorrect(t *testing.T) {
	getTime("0.beeik-ntp.pool.ntp.org")

}
