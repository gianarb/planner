package planner

import (
	"context"
	"strings"
	"testing"
)

const Seq = "seq"

type FakePlan struct {
	P       []Procedure
	Counter int
}

func (p *FakePlan) Create(ctx context.Context) ([]Procedure, error) {
	if p.Counter > 0 {
		return nil, nil
	}
	p.Counter++
	return p.P, nil
}

func (p *FakePlan) Name() string {
	return "fake"
}

type Two struct {
	Seq     chan string
	Counter int
}

func (o *Two) Name() string {
	return "two"
}

func (o *Two) Do(ctx context.Context) ([]Procedure, error) {
	if o.Counter > 0 {
		return nil, nil
	}
	o.Seq <- "two"
	o.Counter++
	return []Procedure{
		&One{
			String: "from-two",
			Seq:    o.Seq,
		},
	}, nil
}

type One struct {
	Seq    chan string
	String string
}

func (o *One) Name() string {
	return "one"
}

func (o *One) Do(ctx context.Context) ([]Procedure, error) {
	o.Seq <- "one"
	return nil, nil
}

func TestExecutionSingleStep(t *testing.T) {
	order := []string{}
	seq := make(chan string, 1)

	ctx := context.Background()
	p := &FakePlan{
		P:       []Procedure{},
		Counter: 0,
	}
	p.P = append(p.P, &One{
		Seq: seq,
	})
	s := NewScheduler()
	s.Execute(ctx, p)
	close(seq)
	for ss := range seq {
		order = append(order, ss)
	}
	if strings.Join(order, " ") != "one" {
		t.Errorf("expected \"one\". Got \"%s\".", strings.Join(order, " "))
	}
}

func TestExecutionStepTwoThatReturnStepOne(t *testing.T) {
	order := []string{}
	seq := make(chan string, 2)

	ctx := context.Background()
	p := &FakePlan{
		P:       []Procedure{},
		Counter: 0,
	}
	p.P = append(p.P, &Two{
		Seq: seq,
	})
	s := NewScheduler()
	s.Execute(ctx, p)
	close(seq)
	for ss := range seq {
		order = append(order, ss)
	}
	if strings.Join(order, " ") != "two one" {
		t.Errorf("expected \"two one\". Got \"%s\".", strings.Join(order, " "))
	}
}

func TestExecutionStepTwoAndStapOne(t *testing.T) {
	order := []string{}
	seq := make(chan string, 3)

	ctx := context.Background()
	p := &FakePlan{
		P:       []Procedure{},
		Counter: 0,
	}
	p.P = append(p.P, &Two{
		Seq: seq,
	})
	p.P = append(p.P, &One{
		Seq:    seq,
		String: "plan",
	})
	s := NewScheduler()
	s.Execute(ctx, p)
	close(seq)
	for ss := range seq {
		order = append(order, ss)
	}
	if strings.Join(order, " ") != "two one one" {
		t.Errorf("expected \"two one one\". Got \"%s\".", strings.Join(order, " "))
	}
}
