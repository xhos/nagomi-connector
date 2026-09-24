package runner

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"nagomi-connector/internal/api"
	"nagomi-connector/internal/domain"
	"nagomi-connector/internal/provider"
)

type testJobs struct {
	cursor    *time.Time
	status    string
	completed bool
}

func (*testJobs) ListSyncJobs(context.Context) ([]api.SyncJob, error) {
	return []api.SyncJob{{ID: 1}}, nil
}
func (j *testJobs) CompleteSyncJob(_ context.Context, _ int64, cursor *time.Time, status *string) error {
	j.cursor = cursor
	j.status = *status
	j.completed = true
	return nil
}

type testProvider struct{ err error }

func (testProvider) Name() string { return "test" }
func (p testProvider) Poll(context.Context) ([]domain.Transaction, error) {
	return []domain.Transaction{{}}, p.err
}

type testSink struct{ err error }

func (s testSink) CreateTransactions(context.Context, string, []domain.Transaction) error {
	return s.err
}

func TestSyncOutcome(t *testing.T) {
	for _, stage := range []string{"success", "factory", "poll", "sink"} {
		t.Run(stage, func(t *testing.T) {
			failure := errors.New("failed")
			jobs := &testJobs{}
			sink := testSink{}
			if stage == "sink" {
				sink.err = failure
			}
			r := New(jobs, sink, func(api.SyncJob) (provider.Provider, error) {
				if stage == "factory" {
					return nil, failure
				}
				p := testProvider{}
				if stage == "poll" {
					p.err = failure
				}
				return p, nil
			}, log.New(io.Discard))
			r.RunOnce(context.Background())
			if !jobs.completed {
				t.Fatal("outcome not recorded")
			}
			if stage == "success" {
				if jobs.cursor == nil || jobs.status != "active" {
					t.Fatal("success must advance cursor and recover status")
				}
			} else if jobs.cursor != nil || jobs.status != "broken" {
				t.Fatal("failure must preserve cursor and report broken")
			}
		})
	}
}
