NAME=nextime

build:
	go build -trimpath -ldflags "-s -w" .

build-all: $(NAME)_darwin_amd64 $(NAME)_darwin_arm64 $(NAME)_linux_amd64 $(NAME)_windows_amd64 $(NAME)_openbsd_amd64

$(NAME)_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .
$(NAME)_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o $@ .
$(NAME)_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .
$(NAME)_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .
$(NAME)_openbsd_amd64:
	GOOS=openbsd GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o $@ .

clean:
	rm -f $(NAME) $(NAME)_*

.PHONY: build build-all clean
