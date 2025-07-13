package util

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"os"
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

func SaveBase64AudioToWav(base64Data string, filePath string) error {
	if strings.HasPrefix(base64Data, "data:audio/wav;base64,") {
		base64Data = strings.TrimPrefix(base64Data, "data:audio/wav;base64,")
	} else {
		return errors.New("invalid base64 audio format, expected 'data:audio/wav;base64,' prefix")
	}

	audioBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	CreateFile(filePath)
	err = os.WriteFile(filePath, audioBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
