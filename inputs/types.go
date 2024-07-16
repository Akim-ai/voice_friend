package inputs

type OptsFunc func(OptsInput)

type OptsInput interface{}

type ChanIn struct{
    flag int8
    message string
}

type WavHeader struct {
	RiffHeader       [4]byte
	FileSize         uint32
	WaveHeader       [4]byte
	FmtHeader        [4]byte
	FmtChunkSize     uint32
	AudioFormat      uint16
	NumChannels      uint16
	SampleRate       uint32
	ByteRate         uint32
	BlockAlign       uint16
	BitsPerSample    uint16
	DataHeader       [4]byte
	DataBytes        uint32
}
