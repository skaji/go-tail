package tail

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"time"
)

const lineFeed = '\n'

type Config struct {
	Logger     Logger
	Path       string
	ReadLength int
	Tick       time.Duration
}

func NewConfig(path string) *Config {
	return &Config{
		Logger:     log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		Path:       path,
		ReadLength: 2 * 1024,
		Tick:       200 * time.Millisecond,
	}
}

func Run(ctx context.Context, config *Config) <-chan string {
	ch := make(chan string)
	f := &file{
		buffer: make([]byte, 0, 64*1024),
		file:   nil,
		logger: config.Logger,
		path:   config.Path,
		tick:   config.Tick,
		tmp:    make([]byte, config.ReadLength),
	}
	go f.run(ctx, ch)
	return ch
}

type file struct {
	buffer []byte
	file   *os.File
	logger Logger
	path   string
	tick   time.Duration
	tmp    []byte
}

func (f *file) run(ctx context.Context, ch chan<- string) {
	defer func() {
		if f.file != nil {
			f.file.Close()
		}
	}()
	defer close(ch)
	ticker := time.NewTicker(f.tick)
	defer ticker.Stop()
	var line []byte
	var err error
TOPLOOP:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for {
				if line == nil {
					line, err = f.nextLine()
					if err != nil {
						if err != io.EOF {
							f.logger.Println(err)
						}
						continue TOPLOOP
					}
				}
				select {
				case <-ctx.Done():
					return
				case ch <- string(line):
					line = nil
				}
			}
		}
	}
}

func (f *file) read() error {
	if f.file == nil {
		file, err := os.Open(f.path)
		if err != nil {
			return err
		}
		f.file = file
	}
	n, err := f.file.Read(f.tmp)
	if err != nil {
		return err
	}
	f.buffer = append(f.buffer, f.tmp[:n]...)
	return nil
}

func (f *file) nextLine() ([]byte, error) {
	for {
		if i := bytes.IndexByte(f.buffer, lineFeed); i != -1 {
			b := f.buffer[:i+1]
			f.buffer = f.buffer[i+1:]
			return b, nil
		} else if err := f.read(); err != nil {
			return nil, err
		}
	}
}
