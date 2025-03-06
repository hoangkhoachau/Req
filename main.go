package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func sendRequest(method, url string, header [][2]string,
	data []byte, config map[string]any,
) (*http.Response, error) {
	var err error
	var res *http.Request
	if method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" {
		rd := bytes.NewReader(data)
		res, err = http.NewRequest(method, url, rd)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported method")
	}
	res.Header.Set("Accept", "*/*")
	if _, ok := config["json"]; ok {
		res.Header.Set("Content-Type", "application/json")
		res.Header.Set("Accept", "application/json, */*;q=0.5")
	} else {
		res.Header.Set("Content-Type", http.DetectContentType(data))
	}
	for _, v := range header {
		res.Header.Set(v[0], v[1])
	}
	if _, ok := config["printAll"]; ok {
		printRequest(res)
		println()
	}
	return http.DefaultClient.Do(res)
}

func printRespond(res *http.Response, config map[string]any) {
	defer res.Body.Close()
	if _, ok := config["printAll"]; ok || res.StatusCode/100 != 2 {
		fmt.Printf("%s%s ", Blue, res.Proto)
		switch res.StatusCode / 100 {
		case 2:
			fmt.Print(Green)
		case 3:
			fmt.Print(Blue)
		case 4:
			fmt.Print(Yellow)
		case 5:
			fmt.Print(Red)
		}
		fmt.Printf("%s%s\n", res.Status, Reset)
		for k, v := range res.Header {
			for i := 0; i < len(v); i++ {
				fmt.Printf(Green+"%s"+Reset+": %s\n", k, v[i])
			}
		}
	}
	bod, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if name, ok := config["output"]; ok {
		file, err := os.Create(name.(string))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		_, err = file.Write(bod)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if _, ok := config["headerOnly"]; !ok {
			if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
				prettyJson(bod)
			} else {
				fmt.Println(string(bod))
			}
		}
	}
}

func printRequest(req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Print(Green)
	case "POST":
		fmt.Print(Blue)
	case "PUT":
		fmt.Print(Yellow)
	case "PATCH":
		fmt.Print(Magenta)
	case "DELETE":
		fmt.Print(Red)
	}
	fmt.Printf("%s "+Reset, req.Method)
	fmt.Print(req.URL.RequestURI())
	fmt.Printf(" %s%s%s\n", Blue, req.Proto, Reset)
	for k, v := range req.Header {
		for i := 0; i < len(v); i++ {
			fmt.Printf(Green+"%s"+Reset+": %s\n", k, v[i])
		}
	}
	buf, _ := io.ReadAll(req.Body)
	if strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		prettyJson(buf)
	} else {
		fmt.Println(string(buf))
	}
	req.Body = io.NopCloser(bytes.NewReader(buf))
}

func main() {
	var positionalFlag []string
	config := make(map[string]interface{})
	var data []byte
	stat, err := os.Stdin.Stat()
	if stat.Size() > 0 {
		data, err = io.ReadAll(os.Stdin)
	}
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] == '-' && len(os.Args[i]) > 1 && os.Args[i][1] != '-' {
			switch os.Args[i] {
			case "-h":
				config["headerOnly"] = true
			case "-f":
				config["printAll"] = true
			case "-d":
				if err != nil && i+1 < len(os.Args) {
					data = []byte(os.Args[i+1])
				} else {
					fmt.Println("invalid data")
					return
				}
				i++
			case "-o":
				if i+1 < len(os.Args) {
					config["output"] = os.Args[i+1]
				} else {
					fmt.Println("invalid filename")
					return
				}
				i++
			default:
			}
		} else if os.Args[i][0] == '@' {
			if data == nil {
				data, err = os.ReadFile(os.Args[i][1:])
				if err != nil {
					fmt.Println("Invalid data")
					return
				}
			} else {
				fmt.Println("Invalid data")
				return
			}
		} else {
			positionalFlag = append(positionalFlag, os.Args[i])
		}
	}
	// fmt.Println(positionalFlag)
	url, method := "", ""
	header := [][2]string{}
	// process data flag
	count := len(positionalFlag)
	// TODO: improve this loop
	for ; count > 1 &&
		(strings.Contains(positionalFlag[count-1], ":=") ||
			strings.Contains(positionalFlag[count-1], "==") ||
			strings.Index(positionalFlag[count-1], ":") > 0 ||
			strings.Contains(positionalFlag[count-1], "=")); count-- {
	}
	switch count {
	case 2:
		method = positionalFlag[0]
		fallthrough
	case 1:
		url = positionalFlag[count/2]
	default:
		fmt.Println("Missing url")
		return
	}

	// process method and url
	method = strings.ToUpper(method)
	if len(url) > 0 && url[0] == ':' {
		if len(url) > 1 && ('0' > url[1] || url[1] > '9') {
			url = fmt.Sprintf(":80%s", url[1:])
		}
		url = "localhost" + url
	}
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	// process data
	jsonObj := map[string]interface{}{}
	err = json.Unmarshal(data, &jsonObj)
	if err == nil {
		config["json"] = true
	}
	for i := count; i < len(positionalFlag); i++ {
		if strings.Contains(positionalFlag[i], "==") {
			split := strings.SplitAfterN(positionalFlag[i], "==", 2)
			if len(split[1]) > 0 && split[1][0] == '@' {
				content, err := os.ReadFile(split[1][1:])
				if err != nil {
					fmt.Printf("Can't read %s\n", split[1][1:])
					return
				}
				split[1] = string(content[:len(content)-1])
			}
			if !strings.Contains(url, "?") {
				url += "?"
			} else {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", split[0][:len(split[0])-2], split[1])
		} else if strings.Contains(positionalFlag[i], ":=") {
			split := strings.SplitAfterN(positionalFlag[i], ":=", 2)
			subObj := map[string]interface{}{}
			if len(split[1]) > 0 && split[1][0] == '@' {
				content, err := os.ReadFile(split[1][1:])
				if err != nil {
					fmt.Printf("Can't read %s\n", split[1][1:])
					return
				}
				split[1] = string(content)
			}
			err := json.Unmarshal([]byte(split[1]), &subObj)
			if err != nil {
				fmt.Println(err)
				return
			}
			jsonObj[split[0][:len(split[0])-2]] = subObj
		} else if strings.Contains(positionalFlag[i], "=") {
			split := strings.SplitAfterN(positionalFlag[i], "=", 2)
			if len(split[1]) > 0 && split[1][0] == '@' {
				content, err := os.ReadFile(split[1][1:])
				if err != nil {
					fmt.Printf("Can't read %s\n", split[1][1:])
					return
				}
				split[1] = string(content)
			}

			// jsonObj[split[0][:len(split[0])-1]] = split[1]

			for i := 0; i < len(split[0])-1; i++ {
				break
			}

			jsonObj[split[0][:len(split[0])-1]] = split[1]
		} else if strings.Contains(positionalFlag[i], ":") {
			split := strings.SplitAfterN(positionalFlag[i], ":", 2)
			if len(split[1]) > 0 && split[1][0] == '@' {
				content, err := os.ReadFile(split[1][1:])
				if err != nil {
					fmt.Printf("Can't read %s\n", split[1][1:])
					return
				}
				split[1] = string(content)
			}
			header = append(header, [2]string{split[0][:len(split[0])-1], split[1]})
			// header[split[0][:len(split[0])-1]] = split[1]
		}
	}
	if len(jsonObj) != 0 {
		data, err = json.Marshal(jsonObj)
		if err != nil {
			fmt.Println(err)
			return
		}
		config["json"] = true
		if method == "" {
			method = "POST"
		}
	}
	if method == "" {
		method = "GET"
	}

	resp, err := sendRequest(method, url, header, data, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	printRespond(resp, config)
}

func prettyJson(buf []byte) {
	dec := json.NewDecoder(bytes.NewReader(buf))
	bufwt := bufio.NewWriter(os.Stdout)
	defer bufwt.Flush()
	key := false
	oldNewLine := false
	currentDelim := Stack[string]{}
	for {
		newLine, comma := true, false
		token, err := dec.Token()
		var printString string
		if err != nil {
			return
		}
		switch token {
		case json.Delim('{'), json.Delim('['):
			key = token == json.Delim('{')
			printString = fmt.Sprint(token)
			currentDelim.Push(printString)
		case json.Delim('}'), json.Delim(']'):
			currentDelim.Pop()
			printString = fmt.Sprint(token)
			comma = dec.More()
		default:
			if key {
				printString = fmt.Sprintf("\"%s\": ", token)
				newLine = false
			} else {
				printString = colorize(token)
				comma = dec.More()
				val, _ := currentDelim.Top()
				key = val == "[" && dec.More()
			}
			key = !key
		}
		if oldNewLine {
			spaces := currentDelim.Len()
			if printString == "[" || printString == "{" {
				spaces--
			}
			for range spaces {
				fmt.Fprint(bufwt, "  ")
			}
		}
		fmt.Fprint(bufwt, printString)
		if comma {
			fmt.Fprint(bufwt, ",")
		}
		if newLine {
			fmt.Fprint(bufwt, "\n")
			bufwt.Flush()
		}
		oldNewLine = newLine
	}
}

func colorize(tok json.Token) string {
	switch tok := tok.(type) {
	case float64:
		return fmt.Sprintf(Yellow+"%v"+Reset, tok)
	case bool:
		return fmt.Sprintf(Blue+"%t"+Reset, tok)
	case string:
		return fmt.Sprintf(Green+"\"%s\""+Reset, tok)
	case nil:
		return Gray + "null" + Reset
	}
	return ""
}

type Stack[T any] struct {
	buf []T
}

func (s *Stack[T]) Push(x T) {
	s.buf = append(s.buf, x)
}

func (s *Stack[T]) Pop() {
	if len(s.buf) > 0 {
		s.buf = s.buf[:len(s.buf)-1]
	}
}

func (s *Stack[T]) Top() (T, bool) {
	var x T
	if len(s.buf) > 0 {
		return s.buf[len(s.buf)-1], true
	}
	return x, false
}

func (s *Stack[T]) Len() int {
	return len(s.buf)
}
