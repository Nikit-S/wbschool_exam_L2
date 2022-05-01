package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"testing"
)

func TestSimple(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/simple/" + strconv.Itoa(i) + ".txt"

		out1, err := exec.Command("./sort", fn).Output()
		if err != nil {
			log.Fatal(err)
		}
		out2, err := exec.Command("sort", fn).Output()
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(out1, out2) {
			t.Fail()
		}
		t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
	}

}

func TestSimpleReverse(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/simple/" + strconv.Itoa(i) + ".txt"

		out1, err := exec.Command("./sort", "-r", fn).Output()
		if err != nil {
			log.Fatal(err)
		}
		out2, err := exec.Command("sort", "-r", fn).Output()
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(out1, out2) {
			t.Fail()
		}
		t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
	}

}

func TestUnique(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/unique/" + strconv.Itoa(i) + ".txt"

		out1, err := exec.Command("./sort", "-u", fn).Output()
		if err != nil {
			t.Log(i, err)
			log.Fatal(err)
		}
		out2, err := exec.Command("sort", "-u", fn).Output()
		if err != nil {
			t.Log(i, err)
			log.Fatal(err)
		}
		if !bytes.Equal(out1, out2) {
			t.Fail()
		}
		t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
	}
}

func TestUniqueReverse(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/unique/" + strconv.Itoa(i) + ".txt"

		out1, err := exec.Command("./sort", "-u", "-r", fn).Output()
		if err != nil {
			t.Log(i, err)
			log.Fatal(err)
		}
		out2, err := exec.Command("sort", "-u", "-r", fn).Output()
		if err != nil {
			t.Log(i, err)
			log.Fatal(err)
		}
		if !bytes.Equal(out1, out2) {
			t.Fail()
		}
		t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
	}
}

func TestSimpleKey(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/simpleKey/" + strconv.Itoa(i) + ".txt"
		for j := range [4]struct{}{} {
			out1, err := exec.Command("./sort", "-k", strconv.Itoa(j+1), fn).Output()
			t.Log("sort", "-k", strconv.Itoa(j+1), fn)
			if err != nil {
				log.Fatal("sort", "-k", strconv.Itoa(j+1), fn, err)
			}
			out2, err := exec.Command("sort", "-k", strconv.Itoa(j+1), fn).Output()
			if err != nil {
				log.Fatal("sort", "-k", strconv.Itoa(j+1), fn, err)
			}
			if !bytes.Equal(out1, out2) {
				t.Fail()
				t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
			}
		}
	}

}

func TestSimpleKeyReverse(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/simpleKey/" + strconv.Itoa(i) + ".txt"
		for j := range [4]struct{}{} {
			out1, err := exec.Command("./sort", "-r", "-k", strconv.Itoa(j+1), fn).Output()
			t.Log("sort", "-r", "-k", strconv.Itoa(j+1), fn)
			if err != nil {
				log.Fatal("sort", "-r", "-k", strconv.Itoa(j+1), fn, err)
			}
			out2, err := exec.Command("sort", "-r", "-k", strconv.Itoa(j+1), fn).Output()
			if err != nil {
				log.Fatal("sort", "-r", "-k", strconv.Itoa(j+1), fn, err)
			}
			if !bytes.Equal(out1, out2) {
				t.Fail()
				t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
			}
		}
	}

}

func TestSimpleKeyUniqueNumReverse(t *testing.T) {
	for i := range [15]struct{}{} {
		fn := "./test_files/uniqueNum/" + strconv.Itoa(i) + ".txt"
		out1, err := exec.Command("./sort", "-u", "-r", fn, fn).Output()
		t.Log("sort", "-u", "-r", fn, fn)
		if err != nil {
			log.Fatal("sort", "-u", "-r", fn, fn, err)
		}
		out2, err := exec.Command("sort", "-u", "-r", fn, fn).Output()
		if err != nil {
			log.Fatal("sort", "-u", "-r", fn, fn, err)
		}
		if !bytes.Equal(out1, out2) {
			t.Fail()
			t.Logf("\n%10s %v\n%10s %v", "Expected:", out1, "Got:", out2)
		}
	}

}
