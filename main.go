package main

import (
	"fmt"
	"voice_friend/inputs"
)

func main(){
    fmt.Println("Hello World")

    MicrophoneChan := make(chan inputs.ChanIn)
    inp := inputs.NewMicrophoneInput(MicrophoneChan)
    inp.Run()
}
