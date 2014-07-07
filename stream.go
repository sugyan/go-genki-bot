package main

import (
	"bufio"
	"encoding/json"
)

// Stream type
type Stream struct {
	scanner *bufio.Scanner
}

// NextTweet returns new tweet
func (s *Stream) NextTweet() (tweet *Tweet, err error) {
	for s.scanner.Err() == nil {
		var bytes []byte
		bytes, err = func() ([]byte, error) {
			for {
				if !s.scanner.Scan() {
					return nil, s.scanner.Err()
				}
				bytes := s.scanner.Bytes()
				if len(bytes) > 0 {
					return bytes, nil
				}
			}
		}()
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(bytes, &tweet); err != nil {
			return nil, err
		}
		if tweet.ID > 0 {
			return tweet, nil
		}
	}
	return nil, s.scanner.Err()
}
