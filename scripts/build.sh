if [ "$NAME" = "" ]; then
	NAME="DiscordChatExporter"
fi

rm -r binaries

mkdir binaries

GOOS=windows go build -o "binaries/${NAME}_windows.exe"
GOOS=darwin go build -o "binaries/${NAME}_mac"
GOOS=linux go build -o "binaries/${NAME}_linux"
