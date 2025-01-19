.PHONE: run encode play

DATA_PATH = data

run:
	@go run cmd/main.go

encode:
	@~/Downloads/ffmpeg -f s16le -ar 44100 -ac 1 -i ${DATA_PATH}/output.bin ${DATA_PATH}/output.wav

play:
	@~/Downloads/ffplay -i ${DATA_PATH}/output.wav

play-with-wave:
	@~/Downloads/ffplay -showmode 1 ${DATA_PATH}/output.wav
