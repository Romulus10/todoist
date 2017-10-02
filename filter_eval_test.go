package main

import (
	"testing"
	"time"

	"github.com/sachaos/todoist/lib"
	"github.com/stretchr/testify/assert"
)

func testFilterEval(t *testing.T, f string, item todoist.Item, expect bool) {
	actual, _ := Eval(Filter(f), item)
	assert.Equal(t, expect, actual, "they should be equal")
}

func TestEval(t *testing.T) {
	testFilterEval(t, "", todoist.Item{}, true)
}

func TestPriorityEval(t *testing.T) {
	testFilterEval(t, "p1", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p2", todoist.Item{Priority: 1}, false)
}

func TestBoolInfixOpExp(t *testing.T) {
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 1}, true)
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 2}, true)
	testFilterEval(t, "p1 | p2", todoist.Item{Priority: 3}, false)

	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 1}, false)
	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 2}, false)
	testFilterEval(t, "p1 & p2", todoist.Item{Priority: 3}, false)
}

func TestDueOnEval(t *testing.T) {
	timeNow := time.Date(2017, time.October, 2, 1, 0, 0, 0, time.Local) // JST: Mon 2 Oct 2017 00:00:00
	setNow(timeNow)

	testFilterEval(t, "today", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 15:00:00 +0000"}, true)  // JST: Mon 2 Oct 2017 00:00:00
	testFilterEval(t, "today", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 14:59:59 +0000"}, true)  // JST: Mon 2 Oct 2017 23:59:59
	testFilterEval(t, "today", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 15:00:00 +0000"}, false) // JST: Tue 3 Oct 2017 00:00:00

	testFilterEval(t, "yesterday", todoist.Item{DueDateUtc: "Sun 1 Oct 2017 14:59:59 +0000"}, true)   // JST: Sun 1 Oct 2017 23:59:59
	testFilterEval(t, "yesterday", todoist.Item{DueDateUtc: "Sat 30 Sep 2017 15:00:00 +0000"}, true)  // JST: Sun 1 Oct 2017 00:00:00
	testFilterEval(t, "yesterday", todoist.Item{DueDateUtc: "Sat 30 Sep 2017 14:59:59 +0000"}, false) // JST: Sat 30 Sept 2017 23:59:59
	testFilterEval(t, "tomorrow", todoist.Item{DueDateUtc: "Mon 2 Oct 2017 15:00:00 +0000"}, true)    // JST: Tue 3 Oct 2017 00:00:00
	testFilterEval(t, "tomorrow", todoist.Item{DueDateUtc: "Tue 3 Oct 2017 14:59:59 +0000"}, true)    // JST: Tue 3 Oct 2017 23:59:59
	testFilterEval(t, "tomorrow", todoist.Item{DueDateUtc: "Tue 3 Oct 2017 15:00:00 +0000"}, false)   // JST: Wed 4 Oct 2017 00:00:00
}
