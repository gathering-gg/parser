package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func watcherSetup() *os.File {
	tmp, err := ioutil.TempFile(os.TempDir(), "gathering-test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	return tmp
}

func watcherTeardown(f *os.File) {
	os.Remove(f.Name())
}

func TestNewWatcher(t *testing.T) {
	f := watcherSetup()
	defer watcherTeardown(f)
	watcher := NewWatcher(f.Name(), 500*time.Millisecond)
	assert.NotNil(t, watcher)
}

func TestWatcherSize(t *testing.T) {
	f := watcherSetup()
	defer watcherTeardown(f)
	_, err := f.Write([]byte("data"))
	assert.Nil(t, err)
	watcher := NewWatcher(f.Name(), 500*time.Millisecond)
	s, err := watcher.size()
	assert.Nil(t, err)
	assert.Equal(t, int64(4), s)
}

func TestWatcherWatchFirstFire(t *testing.T) {
	a := assert.New(t)
	f := watcherSetup()
	defer watcherTeardown(f)
	watcher := NewWatcher(f.Name(), 500*time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				a.True(ok)
				a.NotNil(event)
				a.Equal(int64(0), event.Size)
				done <- true
			}
		}
	}()

	watcher.Start()
	<-done
}

func TestWatcherWatchNoSecondFire(t *testing.T) {
	a := assert.New(t)
	f := watcherSetup()
	defer watcherTeardown(f)
	watcher := NewWatcher(f.Name(), 100*time.Millisecond)
	done := make(chan bool)
	ticks := []int64{}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				a.True(ok)
				a.NotNil(event)
				ticks = append(ticks, event.Size)
			}
		}
	}()
	time.AfterFunc(3*time.Second, func() {
		done <- true
	})
	watcher.Start()
	<-done
	a.Len(ticks, 1)
}

func TestWatcherWatchSecondFire(t *testing.T) {
	a := assert.New(t)
	f := watcherSetup()
	defer watcherTeardown(f)
	watcher := NewWatcher(f.Name(), 100*time.Millisecond)
	done := make(chan bool)
	ticks := []int64{}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				a.True(ok)
				a.NotNil(event)
				ticks = append(ticks, event.Size)
			}
		}
	}()
	time.AfterFunc(3*time.Second, func() {
		done <- true
	})
	time.AfterFunc(1*time.Second, func() {
		f.Write([]byte("data"))
	})
	watcher.Start()
	<-done
	a.Len(ticks, 2)
	a.Equal(int64(4), ticks[1])
}
