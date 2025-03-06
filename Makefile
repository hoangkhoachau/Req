.PHONY: clean
req: main.go
	go build
clean:
	rm ./req
