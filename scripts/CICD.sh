#!/bin/bash

# Styling for clarity
GREEN="\033[1;32m"
RED="\033[1;31m"
YELLOW="\033[1;33m"
CYAN="\033[1;36m"
RESET="\033[0m"

read -p "${GREEN} What would you like your commit message to be?\n" message

go fmt ./...
go vet ./...

echo -e "${CYAN} Running tests...${RESET}"

if go test -v ./...; then
  echo -e "${GREEN} All tests passed!${RESET}"
else
  echo -e "${RED} Some tests failed!${RESET}"
  exit 1
fi

git add . && git commit -m message && git push

echo -e "${GREEN} CI pipeline completed successfully!${RESET}"