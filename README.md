## Discord Bot to Play Random Spotify Songs
This simple Discord app runs in the background and runs a random Spotify song every time a user types `!playsomething`.

### Build
This assumes you have a working Go environment setup and the following dependencies installed on your system:
- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [Requests](https://github.com/asmcos/requests)
- [gJSON](https://github.com/tidwall/gjson)

To install, download the file and run:  
```sh
go build get_random_song.go
```

### Usage
You need a developer Spotify account for this to work. You also need your Discord Bot Token.

```
./get_random_song --help
Usage of ./get_random_song.exe:
  -client_id string
        Spotify Client ID Code
  -client_secret string
        Spotify Client Secret Code
  -t string
        Bot Token
```
