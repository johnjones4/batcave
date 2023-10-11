package worker

import (
	"context"
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
	"github.com/orcaman/writerseeker"
	"github.com/sirupsen/logrus"
)

const (
	Channels   = 1
	SampleRate = 16000
	BitDepth   = 16
)

type VoiceWorker struct {
	workerConcrete
	queue            chan []byte
	buffer           *audio.IntBuffer
	audioChunkBuffer []int16
	stream           *portaudio.Stream
}

func NewVoiceWorker(log logrus.FieldLogger) *VoiceWorker {
	return &VoiceWorker{
		workerConcrete:   newWorkerConcrete(log),
		queue:            make(chan []byte),
		audioChunkBuffer: make([]int16, 64),
	}
}

func (v *VoiceWorker) Setup(errors chan error) error {
	v.workerConcrete.errors = errors
	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	v.stream, err = portaudio.OpenDefaultStream(Channels, 0, SampleRate, len(v.audioChunkBuffer), v.audioChunkBuffer)
	if err != nil {
		return err
	}

	return nil
}

func (v *VoiceWorker) Teardown() error {
	err := portaudio.Terminate()
	if err != nil {
		return err
	}

	err = v.stream.Close()
	if err != nil {
		return err
	}

	return nil
}

func (v *VoiceWorker) Stop() {
	v.log.Debug("Stopping voice")
	err := v.stream.Stop()
	if err != nil {
		v.errors <- err
	}
	v.stop <- true
	<-v.stopped
	v.buffer = nil
	v.log.Debug("Stopped voice")
}

func (v *VoiceWorker) Start(ctx context.Context) {
	v.log.Debug("Starting voice")
	err := v.stream.Start()
	if err != nil {
		v.workerConcrete.errors <- err
		return
	}
	v.buffer = &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: Channels,
			SampleRate:  SampleRate,
		},
		Data:           make([]int, 0),
		SourceBitDepth: BitDepth,
	}
	defer func() {
		v.stopped <- true
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-v.stop:

			v.log.Debugf("Have %d audio frames", len(v.buffer.Data))

			buffer := &writerseeker.WriterSeeker{}
			encoder := wav.NewEncoder(buffer, SampleRate, 16, Channels, 1)
			err = encoder.Write(v.buffer)
			if err != nil {
				v.workerConcrete.errors <- err
				return
			}
			err = encoder.Close()
			if err != nil {
				v.workerConcrete.errors <- err
				return
			}

			bytes, err := io.ReadAll(buffer.Reader())
			if err != nil {
				v.workerConcrete.errors <- err
				return
			}

			v.queue <- bytes

			return
		default:
			if i, _ := v.stream.AvailableToRead(); i > 0 {
				err = v.stream.Read()
				if err != nil {
					v.workerConcrete.errors <- err
					continue
				}
				for _, i := range v.audioChunkBuffer {
					v.buffer.Data = append(v.buffer.Data, int(i))
				}
			}
		}
	}
}

func (v *VoiceWorker) Chan() chan []byte {
	return v.queue
}
