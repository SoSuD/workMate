package faketask

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Task struct {
	ID         uuid.UUID
	Name       string
	CreatedAt  time.Time
	FinishedAt time.Time
	Status     string
	Result     string
}

type Result struct {
	Value string
	Err   error
}

func LongTask(resCh chan<- Result) {
	minimal := 180
	maximal := 300
	duration := time.Second * time.Duration(rand.Intn(maximal-minimal)+minimal)
	time.Sleep(duration)
	resCh <- Result{Value: "done", Err: nil}
}
