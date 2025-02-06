.PHONE: run encode play

DATA_PATH = data

run:
	@go run cmd/main.go

encode:
	@~/Downloads/ffmpeg -f s16le -ar 44100 -ac 1 -i ${DATA_PATH}/output.bin ${DATA_PATH}/output.wav

# - -f s16le : Spécifie le format d'entrée des données audio brutes :
#   - s16le signifie signed 16-bit little-endian, ce qui correspond au format PCM 16 bits.
# - -ac 1 : Définit le nombre de canaux à 1 (mono).

play:
	@~/Downloads/ffplay -i ${DATA_PATH}/output.wav

play-with-wave:
	@~/Downloads/ffplay -showmode 1 ${DATA_PATH}/output.wav
