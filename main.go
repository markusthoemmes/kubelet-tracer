package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
)

const (
	msgContainerKilled = "Killing container with a grace period"
)

type message struct {
	Timestamp float64  `json:"ts"`
	Message   string   `json:"msg"`
	Pod       string   `json:"pod"`
	Pods      []string `json:"pods"`
	Caller    string   `json:"caller"`
}

func main() {
	var (
		pod               string
		stopAfterDeletion bool
	)
	flag.StringVar(&pod, "pod", "", "the pod to analyze the logs for")
	flag.BoolVar(&stopAfterDeletion, "stop-after-deletion", false, "stop log analyzing after seeing a deletion")
	flag.Parse()

	if pod == "" {
		log.Fatalln("No pod provided")
	}
	fmt.Println("Pod: " + pod)

	var msgs []message
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		start := strings.Index(line, "{")
		if start == -1 {
			continue
		}
		line = line[start:]

		if !strings.Contains(line, pod) {
			continue
		}

		var msg message
		json.Unmarshal([]byte(line), &msg)

		if msg.Pod == pod {
			if stopAfterDeletion && msg.Message == msgContainerKilled {
				break
			}

			msgs = append(msgs, msg)
			continue
		}

		if msg.Pods != nil {
			for i := range msg.Pods {
				if msg.Pods[i] == pod {
					msgs = append(msgs, msg)
					break
				}
			}
			continue
		}
	}

	if len(msgs) == 0 {
		log.Fatalln("No messages found for pod " + pod)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Timestamp < msgs[j].Timestamp
	})

	start := msgs[0].Timestamp
	for _, msg := range msgs {
		// Figure out which subsystem the line belongs to.
		subsystem := color.New(color.Bold).SprintFunc()("MISC")
		if strings.HasPrefix(msg.Caller, "volumemanager/") || strings.HasPrefix(msg.Caller, "populator/") || strings.HasPrefix(msg.Caller, "reconciler/") || strings.HasPrefix(msg.Caller, "operationexecutor/") {
			subsystem = color.New(color.Bold, color.FgGreen).SprintFunc()("VOLUME")
		} else if strings.HasPrefix(msg.Caller, "kuberuntime/") || msg.Message == "Generating pod status" {
			subsystem = color.New(color.Bold, color.FgBlue).SprintFunc()("SYNCPOD")
		} else if strings.HasPrefix(msg.Caller, "pleg/") {
			subsystem = color.New(color.Bold, color.FgRed).SprintFunc()("PLEG")
		} else if strings.HasPrefix(msg.Caller, "status/") {
			subsystem = color.New(color.Bold, color.FgHiBlue).SprintFunc()("STATUS")
		} else if strings.HasPrefix(msg.Caller, "kubelet/kubelet_pods") {
			subsystem = color.New(color.Bold, color.FgHiGreen).SprintFunc()("MOUNT")
		}

		fmt.Printf("%d\t%s\t%s\n", int(msg.Timestamp-start), subsystem, truncate(msg.Message, 90))
	}
}

func truncate(msg string, max int) string {
	if len(msg) < max {
		return msg
	}
	return msg[:max-3] + "..."
}
