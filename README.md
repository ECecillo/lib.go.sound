# lib.go.sound

## TODO

- [] Utiliser la feature `embed` de Golang pour créer un test qui va checker si ce que l'on a écris dans un buffer
  - Corrspond à ce que l'on a sauvegardé dans notre fichier de référence (retrouver la vidéo qui présente ce machin)
- Have fun with other stuff

## Commands

All commands for this project can be found in Makefile.

### Encode binary data

```sh
~/Downloads/ffmpeg -f s16le -ar 44100 -ac 1 -i output.bin output.wav
```

Explications des options :

- -f s16le : Spécifie le format d'entrée des données audio brutes :
  - s16le signifie signed 16-bit little-endian, ce qui correspond au format PCM 16 bits.
- -ar 44100 : Définit la fréquence d'échantillonnage (sample rate) à 44100 Hz.
- -ac 1 : Définit le nombre de canaux à 1 (mono).
- -i output.bin : Spécifie le fichier d'entrée brut (output.bin).
- output.wav : Spécifie le fichier de sortie encodé au format WAV.
