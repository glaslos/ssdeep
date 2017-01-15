bench:
	go test -bench=. -run=a^

sample:
	head -c 10M < /dev/urandom > /tmp/data
