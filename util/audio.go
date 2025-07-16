package util

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	sampleRate    = 22100
	bitsPerSample = 16
	numChannels   = 1
	wavHeaderSize = 44
)

func writeWavHeader(f *os.File, dataSize int64) error {
	byteRate := sampleRate * numChannels * bitsPerSample / 8
	blockAlign := numChannels * bitsPerSample / 8
	chunkSize := 36 + dataSize

	f.Seek(0, 0)
	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, uint32(chunkSize))
	f.Write([]byte("WAVE"))
	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1)) // PCM
	binary.Write(f, binary.LittleEndian, uint16(numChannels))
	binary.Write(f, binary.LittleEndian, uint32(sampleRate))
	binary.Write(f, binary.LittleEndian, uint32(byteRate))
	binary.Write(f, binary.LittleEndian, uint16(blockAlign))
	binary.Write(f, binary.LittleEndian, uint16(bitsPerSample))
	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, uint32(dataSize))
	return nil
}

// 流式写入函数
func StreamWritePCMBase64ToWav(filePath, base64PCM string, isFinal bool) error {
	var dataSize int64
	var err error

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.Size() == 0 {
		err = writeWavHeader(f, 0)
		if err != nil {
			return err
		}
	}

	if base64PCM != "" {
		pcmBytes, err := base64.StdEncoding.DecodeString(base64PCM)
		if err != nil {
			return err
		}
		f.Seek(0, io.SeekEnd)
		_, err = f.Write(pcmBytes)
		if err != nil {
			return err
		}
	}

	if isFinal {
		fi, err = f.Stat()
		if err != nil {
			return err
		}
		dataSize = fi.Size() - wavHeaderSize
		err = writeWavHeader(f, dataSize)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveBase64WebmToMp3(base64Data string, mp3FilePath string) error {
	if strings.HasPrefix(base64Data, "data:audio/webm;base64,") {
		base64Data = strings.TrimPrefix(base64Data, "data:audio/webm;base64,")
	} else {
		return errors.New("invalid base64 audio format, expected 'data:audio/webm;base64,' prefix")
	}

	audioBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	tmpWebmFile, err := os.CreateTemp("", "*.webm")
	if err != nil {
		return err
	}
	tmpWebmFileName := tmpWebmFile.Name()
	defer func() {
		tmpWebmFile.Close()
		os.Remove(tmpWebmFileName)
	}()

	if _, err := tmpWebmFile.Write(audioBytes); err != nil {
		return err
	}
	tmpWebmFile.Close()

	cmd := exec.Command("ffmpeg", "-y", "-i", tmpWebmFileName, "-vn", "-acodec", "libmp3lame", mp3FilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("ffmpeg error: " + err.Error() + "\n" + string(output))
	}

	return nil
}

func WriteWAVHeader(buf *bytes.Buffer, dataLen int, sampleRate int, bitsPerSample int, channels int) {
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, uint32(36+dataLen))
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16)) // PCM
	binary.Write(buf, binary.LittleEndian, uint16(1))  // PCM format
	binary.Write(buf, binary.LittleEndian, uint16(channels))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	byteRate := sampleRate * channels * bitsPerSample / 8
	binary.Write(buf, binary.LittleEndian, uint32(byteRate))
	blockAlign := channels * bitsPerSample / 8
	binary.Write(buf, binary.LittleEndian, uint16(blockAlign))
	binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample))
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, uint32(dataLen))
}

func ChunksToWavBase64(chunks []string) (string, error) {
	var audioData bytes.Buffer
	for _, chunk := range chunks {
		decoded, err := base64.StdEncoding.DecodeString(chunk)
		if err != nil {
			return "", fmt.Errorf("base64 decode chunk error: %v", err)
		}
		audioData.Write(decoded)
	}
	raw := audioData.Bytes()
	dataLen := len(raw)

	var wav bytes.Buffer
	WriteWAVHeader(&wav, dataLen, 24000, 16, 1)
	wav.Write(raw)

	wavBase64 := base64.StdEncoding.EncodeToString(wav.Bytes())
	dataURL := "data:audio/wav;base64," + wavBase64
	return dataURL, nil
}
