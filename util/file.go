package util

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const (
	sampleRate    = 22100
	bitsPerSample = 16
	numChannels   = 1
	wavHeaderSize = 44
)

func FileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFolder(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		logrus.Errorf("Create folder %s error: %v", path, err)
	}
	return errors.WithStack(err)
}

func CreateFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if err := CreateFolder(basePath); err != nil {
		return nil, errors.WithStack(err)
	}
	return os.Create(path)
}

func JsonToFile(dst string, data interface{}) error {
	str, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logrus.Errorf("JsonToFile error: %v", err)
		return errors.WithStack(err)
	}
	err = os.WriteFile(dst, str, 0700)
	if err != nil {
		logrus.Errorf("WriteFile error: %v", err)
		return errors.WithStack(err)
	}
	return nil
}

func SaveUploadFile(dst string, header *multipart.FileHeader, file *multipart.File) error {
	out, err := CreateFile(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer out.Close()
	_, err = io.Copy(out, *file)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

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
