package sharedfuncs

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/pubsub"
	pb "github.com/fiveateooate/deployinator/deployproto"
	"github.com/gogo/protobuf/proto"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

// RandString - return a random ascii string
func RandString(n int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(s1.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

// RunCmd - shell out and run something
func RunCmd(cmd string, args []string) ([]byte, error) {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
		return cmdOut, err
	}
	return cmdOut, nil
}

// FileExists - checks if a file exists and returns bool
func FileExists(path string) bool {
	exists := false
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		exists = true
	}
	return exists
}

// ProcessMessage - process a pubsub message
func ProcessMessage(ctx context.Context, msg *pubsub.Message) {
	var message pb.DeployMessage
	err := proto.Unmarshal(msg.Data, &message)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf("Name: %s, Namespace: %s.\n", message.Slug, message.Namespace)
	msg.Ack()
	return
}
