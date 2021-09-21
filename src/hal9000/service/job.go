package service

import (
	"encoding/json"
	"fmt"
	"hal9000/types"
	"hal9000/util"
	"io/ioutil"
	"os"
	"time"
)

type JobConcrete struct {
	ID       string        `json:"id"`
	Interval time.Duration `json:"interval"`
	Name     string        `json:"name"`
}

func (j JobConcrete) GetID() string {
	return j.ID
}

func (j JobConcrete) GetInterval() time.Duration {
	return j.Interval
}

func (j JobConcrete) GetName() string {
	return j.Name
}

type jobProviderConcrete struct {
	jobs []types.Job
}

func InitJobProvider(runtime *types.Runtime) (types.JobProvider, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("JOBS_MANIFEST_PATH"))
	if err != nil {
		return nil, err
	}
	var jobsConcrete []JobConcrete
	err = json.Unmarshal(bytes, &jobsConcrete)
	if err != nil {
		return nil, err
	}
	jobs := make([]types.Job, len(jobsConcrete))
	for i, j := range jobsConcrete {
		jobs[i] = j
	}
	jp := jobProviderConcrete{jobs}

	go (func() {
		for {
			for _, j := range jp.jobs {
				var status types.JobStatus
				key := keyForJob(j)
				err := (*(*runtime).KVStore()).GetInterface(key, &status)
				if err != nil {
					(*(*runtime).Logger()).LogError(err)
				} else {
					expectedUpdate := status.LastUpdate.Add(j.GetInterval())
					if expectedUpdate.Before(time.Now()) {
						err = jp.ReportJobStatus(runtime, &j, &types.JobStatusInfo{
							State:       types.JobStateLate,
							Description: fmt.Sprintf("Expected update at %s", expectedUpdate.Format(time.RFC1123)),
						})
						if err != nil {
							(*(*runtime).Logger()).LogError(err)
						}
					}
				}
			}
			time.Sleep(time.Hour)
		}
	})()

	return &jp, nil
}

func (jp *jobProviderConcrete) FindJobById(id string) (*types.Job, error) {
	for _, job := range jp.jobs {
		if job.GetID() == id {
			return &job, nil
		}
	}
	return nil, util.ErrorJobNotFound
}

func (jp *jobProviderConcrete) ReportJobStatus(runtime *types.Runtime, job *types.Job, info *types.JobStatusInfo) error {
	status := types.JobStatus{
		Info:       *info,
		LastUpdate: time.Now(),
	}

	key := keyForJob(*job)
	err := (*(*runtime).KVStore()).SetInterface(key, status, time.Time{})
	if err != nil {
		return err
	}

	m := types.ResponseMessage{Text: fmt.Sprintf("Job status update: %s: %s / %s", (*job).GetName(), status.Info.State, status.Info.Description)}
	(*(*runtime).AlertQueue()).Enqueue(m)

	return nil
}

func keyForJob(job types.Job) string {
	return fmt.Sprintf("job_%s", job.GetID())
}
