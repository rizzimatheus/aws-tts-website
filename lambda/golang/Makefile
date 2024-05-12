
NEW_POST=newpost
NEW_POST_DYNAMODB=newpost-amplify
CONVERT_AUDIO=convertaudio
GET_POSTS=getposts

build: build_newpost build_newpost_amplify build_convertaudio build_getpost
	@echo "Building binaries..."

## build_newpost: builds the newpost binary as a linux executable
build_newpost:
	@echo "Building newpost binary..."
	cd ./${NEW_POST} && \
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/bootstrap .
	cd ./${NEW_POST}/bin && zip ${NEW_POST} bootstrap
	@echo "Done!"

## build_newpost_amplify: builds the newpost amplify binary as a linux executable
build_newpost_amplify:
	@echo "Building newpost_amplify binary..."
	cd ./${NEW_POST_DYNAMODB} && \
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/bootstrap .
	cd ./${NEW_POST_DYNAMODB}/bin && zip ${NEW_POST_DYNAMODB} bootstrap
	@echo "Done!"

## build_convertaudio: builds the convertaudio binary as a linux executable
build_convertaudio:
	@echo "Building convertaudio binary..."
	cd ./${CONVERT_AUDIO} && \
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/bootstrap .
	cd ./${CONVERT_AUDIO}/bin && zip ${CONVERT_AUDIO} bootstrap
	@echo "Done!"

## build_getpost: builds the getpost binary as a linux executable
build_getpost:
	@echo "Building getpost binary..."
	cd ./${GET_POSTS} && \
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/bootstrap .
	cd ./${GET_POSTS}/bin && zip ${GET_POSTS} bootstrap
	@echo "Done!"

## go_work_init: initializes and writes a new go.work file in the current directory, 
## creating a new workspace at the current directory.
go_work_init:
	go work init ./util ./newpost ./newpost-amplify ./getposts ./convertaudio

.PHONY: build build_newpost build_newpost_amplify build_convertaudio build_getpost go_work
