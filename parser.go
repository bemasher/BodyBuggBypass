package main

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"encoding/json"
	"github.com/bemasher/errhandler"
)

const (
	OUT_FILENAME = "data.json"
	
	SESSION_REGEX = "SESSION-BEGIN.*?_([0-9a-z]+)([A-Z][0-9A-Z]+)"
)

var sessionRegex *regexp.Regexp

type RawSession struct {
	Channel string
	Info string
	Payload []string
}

type Session struct {
	Channel string
	Epoch int64
	Payload interface{}
}

// Produces a string representation of a raw session
func (s RawSession) String() string {
	var payload []string
	for _, p := range s.Payload {
		if len(p) > 16 {
			payload = append(payload, p[:16] + "...")
		} else {
			payload = append(payload, p)
		}
	}
	return fmt.Sprintf("{Channel:%s Info:%s Payload:[%s]}", s.Channel, s.Info, strings.Join(payload, ", "))
}

// Reads raw session data
func (s *RawSession) Read(r *bytes.Buffer) error {
	// Get a line and trim space
	tmp, err := r.ReadString('\n')
	errhandler.Handle("Error reading line: ", err)
	
	tmp = strings.TrimSpace(tmp)
	
	// If the line is a session header
	if sessionRegex.MatchString(tmp) {
		// Parse session header
		matches := sessionRegex.FindStringSubmatch(tmp)
		s.Info = matches[1]
		s.Channel = matches[2]
		
		// Consume following lines until we find another session
		for {
			tmp, err := r.ReadString('\n')
			if err != io.EOF && err != nil {
				return err
			}
			
			// Found another session, unread the line and break loop
			if sessionRegex.MatchString(tmp) {
				*r = *AppendBefore(r, tmp)
				break
			}
			
			// Found a payload line, append it to the payload
			if len(tmp) > 0 {
				s.Payload = append(s.Payload, strings.TrimSpace(tmp))
			}
			
			if err == io.EOF {
				return err
			}
		}
	}
	return nil
}

// Parses timestamp and payload
func (rs RawSession) Parse() (s Session, err error) {
	s.Channel = rs.Channel
	fmt.Sscanf(rs.Info[15:23], "%X", &s.Epoch)
	
	switch rs.Channel {
		case "TIMESTMP": s.Payload, err = Timestamp(rs.Payload)
		case "DIAGNSTC": s.Payload, err = Diagnostic(rs.Payload)
		default: s.Payload, err = Packed(rs.Payload)
	}
	return
}

// Parses 12-bit unaligned integers
func Packed(s []string) (interface{}, error) {
	if len(s) != 1 {
		return nil, fmt.Errorf("Expected only one payload item, got: %d", len(s))
	}
	
	payload := make([]uint16, 0)
	for i := 0; i < len(s[0]); i += 3 {
		n, err := strconv.ParseUint(s[0][i:i + 3], 16, 12)
		if err != nil {
			return nil, err
		}
		
		payload = append(payload, uint16(n))
	}
	
	return payload, nil
}

// Parses channel TIMESTMP
func Timestamp(s []string) (interface{}, error) {
	payload := make([]int64, 0)
	
	for _, i := range s {
		var n int64
		_, err := fmt.Sscanf(i, "%d", &n)
		if err != nil {
			return nil, err
		}
		
		payload = append(payload, n)
	}
	
	return payload, nil
}

// Parses channel DIAGNSTC
type DiagnosticPayload struct {
	Timestamp int64
	I, J int
}

func Diagnostic(s []string) (interface{}, error) {
	payload := make([]DiagnosticPayload, 0)
	
	for _, i := range s {
		var n DiagnosticPayload
		_, err := fmt.Sscanf(i, "%d %d %d", &n.Timestamp, &n.I, &n.J)
		if err != nil {
			return nil, err
		}
		
		payload = append(payload, n)
	}
	
	return payload, nil
}

// Appends the given line to the beginning of the buffer
func AppendBefore(r *bytes.Buffer, line string) (w *bytes.Buffer) {
	w = bytes.NewBuffer(nil)
	w.WriteString(line)
	w.ReadFrom(r)
	return
}

func init() {
	var err error
	sessionRegex, err = regexp.CompilePOSIX(SESSION_REGEX)
	errhandler.Handle("Error compiling regex: ", err)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: parser [INFILE] [[OUTFILE]]")
		os.Exit(1)
	}
	
	logFilename := os.Args[1]
	outFilename := OUT_FILENAME
	
	if len(os.Args) == 3 {
		outFilename = os.Args[2]
	}
	
	// Open the input file
	logFile, err := os.Open(logFilename)
	errhandler.Handle("Error opening data: ", err)
	defer logFile.Close()
	
	// Copy the input file into a buffer
	logBuf := bytes.NewBuffer(nil)
	logBuf.ReadFrom(logFile)
	
	// Create a list for parsed sessions
	sessions := make([]Session, 0)
	done := false
	for !done {
		// Create and read a raw session
		var raw RawSession
		err := raw.Read(logBuf)
		if err != io.EOF {
			errhandler.Handle("Error parsing session: ", err)
		}
		
		// If we hit EOF after the last session, no more data to read
		if err == io.EOF {
			done = true
		}
		
		// Parse the raw session and store it in the sessions list
		session, err := raw.Parse()
		errhandler.Handle("Error parsing session: ", err)
		
		sessions = append(sessions, session)
	}
	
	// Create the output file and encode the session list to json
	outFile, err := os.Create(outFilename)
	errhandler.Handle("Error creating output file: ", err)
	defer outFile.Close()
	
	jsonEncoder := json.NewEncoder(outFile)
	err = jsonEncoder.Encode(sessions)
	errhandler.Handle("Error encoding json output: ", err)
}
