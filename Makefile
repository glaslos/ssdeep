default: sample
	go test ./...

bench:
	go test -bench=. -run=a^

sample:
	if [ ! -f /tmp/data ]; then \
	head -c 10M < /dev/urandom > /tmp/data; fi
