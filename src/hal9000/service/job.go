package service

import (
	"encoding/json"
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

func InitJobProvider() (types.JobProvider, error) {
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
	return jobProviderConcrete{jobs}, nil
}

func (jp jobProviderConcrete) StartAbnormalJobWatchLoop(runtime types.Runtime) error {
	return nil
	//TODO
}

func (jp jobProviderConcrete) FindJobById(id string) (types.Job, error) {
	for _, job := range jp.jobs {
		if job.GetID() == id {
			return job, nil
		}
	}
	return nil, util.ErrorJobNotFound
}

func (jp jobProviderConcrete) UpdateJobStatus(job types.Job, status types.JobStatus) error {
	return nil //TODO
}

func (jp jobProviderConcrete) SendJobAlert(job types.Job, status types.JobStatus) error {
	return nil //TODO
}

func (jp jobProviderConcrete) ReportJobStatus(id string, status types.JobStatus) error {
	job, err := jp.FindJobById(id)
	if err != nil {
		return err
	}
	err = jp.UpdateJobStatus(job, status)
	if err != nil {
		return err
	}
	err = jp.SendJobAlert(job, status)
	if err != nil {
		return err
	}
	return nil
}
