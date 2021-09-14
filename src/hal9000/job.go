package hal9000

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type Job struct {
	ID       string        `json:"id"`
	Interval time.Duration `json:"interval"`
	Name     string        `json:"name"`
}

var ErrorJobNotFound = errors.New("job not found")

type JobState string

const (
	JobStateNormal   = "normal"
	JobStateAbnormal = "abnormal"
)

type JobStatus struct {
	State      JobState  `json:"state"`
	LastUpdate time.Time `json:"lastUpdate"`
}

var jobs []Job

func InitJobs() error {
	bytes, err := ioutil.ReadFile(os.Getenv("JOBS_MANIFEST_PATH"))
	if err != nil {
		return err
	}
	jobs = nil
	err = json.Unmarshal(bytes, &jobs)
	if err != nil {
		return err
	}
	return nil
}

func FindJobById(id string) (Job, error) {
	for _, job := range jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return Job{}, ErrorJobNotFound
}

func ReportJobStatus(id string, status JobStatus) error {
	job, err := FindJobById(id)
	if err != nil {
		return err
	}
	//TODO report status

}
