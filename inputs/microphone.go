package inputs

import (
	"github.com/gordonklaus/portaudio"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"strings"
)


type Microphone struct{
    opts optsMicrophone
    chanIn chan ChanIn   
    // OutputChan Chan
}

type optsMicrophone struct{
    
}

func defaultMicrophoneOpts() optsMicrophone{
    return optsMicrophone{}
}

func NewMicrophoneInput(chanIn chan ChanIn, opts ...OptsFunc) *Microphone{
    o :=  defaultMicrophoneOpts()

    for _, fn := range opts {
        fn(&o)
    }
    mic := Microphone{
        opts: o,
        chanIn: chanIn,
    }
    fmt.Println("Returned a mic")
    return &mic
}

func (m *Microphone) Run() {
	fmt.Println("Recording. Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fileName := "test"
	if !strings.HasSuffix(fileName, ".wav") {
		fileName += ".wav"
	}
	f, err := os.Create(fileName)
	chk(err)

	writeWavHeader(f)

	nSamples := 0

	defer func() {
		updateWavSizes(f, nSamples)
		chk(f.Close())
	}()

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	for {
		chk(stream.Read())
		chk(binary.Write(f, binary.LittleEndian, in))
		nSamples += len(in)

		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())
}

func writeWavHeader(w *os.File) {
	// Write the WAV RIFF chunk descriptor
	w.Write([]byte("RIFF"))
	// Placeholder for total file size
	binary.Write(w, binary.LittleEndian, int32(0))
	// WAVE format
	w.Write([]byte("WAVE"))

	// fmt subchunk
	w.Write([]byte("fmt "))
	binary.Write(w, binary.LittleEndian, int32(16)) // Subchunk1Size
	binary.Write(w, binary.LittleEndian, int16(1))  // AudioFormat
	binary.Write(w, binary.LittleEndian, int16(1))  // NumChannels
	binary.Write(w, binary.LittleEndian, int32(44100)) // SampleRate
	binary.Write(w, binary.LittleEndian, int32(44100*1*4)) // ByteRate
	binary.Write(w, binary.LittleEndian, int16(1*4)) // BlockAlign
	binary.Write(w, binary.LittleEndian, int16(32))  // BitsPerSample

	// data subchunk
	w.Write([]byte("data"))
	binary.Write(w, binary.LittleEndian, int32(0)) // Placeholder for data size
}

func updateWavSizes(w *os.File, nSamples int) {
	totalBytes := 4 + 8 + 16 + 8 + 4*nSamples
	// Update the total file size
	w.Seek(4, 0)
	binary.Write(w, binary.LittleEndian, uint32(totalBytes-8))
	// Update the data chunk size
	w.Seek(40, 0)
	binary.Write(w, binary.LittleEndian, uint32(4*nSamples))
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

