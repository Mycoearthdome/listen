package main

import (
	"fmt"
	"os/exec"
	"sync"
	"time"
)

func runCommand(command string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	fmt.Printf("Command output:\n%s\n", output)
}

func main() {
	var wg sync.WaitGroup

	load := "pactl load-module module-pipe-source source_name=virtual file=irtual.wav format=s16le rate=30000 channels=1"
	setting1 := "pactl set-default-source virtual"
	//the following two to prevent the ticking.
	volume := "pactl set-source-volume alsa_input.pci-0000_00_1b.0.analog-stereo 65535" // 100%
	//volume2 := "pactl set-source-output-volume virtual 43690"                         // 75%

	play := "ffmpeg -re -i legal_acknowledgment.wav -f s16le -ar 30000 -filter:a 'volume=1.0' -ac 1 - > virtual.wav"
	mute := "pactl set-source-mute alsa_input.pci-0000_00_1b.0.analog-stereo toggle" // muted for tick collapse
	setting2 := "pactl set-default-source alsa_input.pci-0000_00_1b.0.analog-stereo"
	unmute := "pactl set-source-mute alsa_input.pci-0000_00_1b.0.analog-stereo toggle" // unmuted regular
	unload := "pactl unload-module module-pipe-source"
	ingress := "parec -d alsa_output.pci-0000_00_1b.0.analog-stereo.monitor --file-format=wav RecordedCall_ingress.wav"
	egress := "parec -d alsa_input.pci-0000_00_1b.0.analog-stereo --file-format=wav RecordedCall_egress.wav"

	// Playing the Acknowledgment...
	wg.Add(4)
	go runCommand(load, &wg)
	time.Sleep(time.Millisecond * 5)
	go runCommand(setting1, &wg)
	time.Sleep(time.Millisecond * 5)
	go runCommand(volume, &wg)
	time.Sleep(time.Millisecond * 5)
	//go runCommand(volume2, &wg)
	//time.Sleep(time.Millisecond * 20)
	go runCommand(play, &wg)
	wg.Wait()
	//fmt.Println("Acknowledgement is done!")

	// Restauring settings...
	wg.Add(4)
	go runCommand(mute, &wg)
	time.Sleep(time.Millisecond * 5)
	go runCommand(setting2, &wg)
	time.Sleep(time.Millisecond * 5)
	go runCommand(unload, &wg)
	time.Sleep(time.Millisecond * 5)
	go runCommand(unmute, &wg)
	time.Sleep(time.Millisecond * 5)
	wg.Wait()
	fmt.Println("Starting to listen...CTRL+C to END.")

	//Recording Call...
	wg.Add(2)
	go runCommand(ingress, &wg)
	go runCommand(egress, &wg)
	wg.Wait()

}
