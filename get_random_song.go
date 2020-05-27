package main
import (
        "encoding/base64"
        "flag"
        "fmt"
        "math/rand"
        "os"
        "os/signal"
        "strconv"
        "strings"
        "syscall"
        "time"
        "github.com/bwmarrin/discordgo"
        "github.com/asmcos/requests"
        "github.com/tidwall/gjson"
)

var (
        TOKEN string
        CLIENT_ID string
        CLIENT_SECRET string
        SPOTIFY_TOKEN_URL string = "https://accounts.spotify.com/api/token"
        SPOTIFY_API_URL string = "https://api.spotify.com/v1"
)

func init() {
        flag.StringVar(&TOKEN, "t", "", "Bot Token")
        flag.StringVar(&CLIENT_ID, "client_id", "", "Spotify Client ID Code")
        flag.StringVar(&CLIENT_SECRET, "client_secret", "", "Spotify Client Secret Code")
        flag.Parse()

        if TOKEN == "" || (CLIENT_ID == "" && CLIENT_SECRET == "") {
                flag.Usage()
                os.Exit(1)
        }
}

func main() {
        // Create discord session using bot token
        dg, err := discordgo.New("Bot " + TOKEN)
        if err != nil {
                fmt.Println("error creating Discord session,", err)
                return
        }

        // Register the messageCreate function as a callback for events
        dg.AddHandler(message_create)

        // Open a websocket connection to Discord and begin listening
        err = dg.Open()
        if err != nil {
                fmt.Println("error opening connection,", err)
                return
        }

        // Wait until CTRL-C or other term signal is received
        fmt.Println("Bot is now running. Press CTRL-C to exit.")
        sc := make(chan os.Signal, 1)
        signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
        <-sc

        // Close down Discord session.
        dg.Close()
}

func message_create(s *discordgo.Session, m *discordgo.MessageCreate) {
        if m.Author.ID == s.State.User.ID {
                return
        }

        if strings.ToLower(m.Content) == "!playsomething" {
                access_token := get_token(CLIENT_ID, CLIENT_SECRET, SPOTIFY_TOKEN_URL)

                song_dict := request_valid_song(access_token, SPOTIFY_API_URL)
                song_info := song_dict["song_info"]
                external_url := song_dict["external_url"]

                response := "Here's a song: "+song_info+"\n"+external_url
                fmt.Println(response)
                s.ChannelMessageSend(m.ChannelID, response)
        }
}


func get_token(CLIENT_ID string, CLIENT_SECRET string, SPOTIFY_TOKEN_URL string) (string) {
        client_token := CLIENT_ID+":"+CLIENT_SECRET
        client_byte := []byte(client_token)
        b64_client_byte := base64.StdEncoding.EncodeToString(client_byte)

        headers := requests.Header{
              "Authorization":"Basic "+b64_client_byte,
        }

        payload := requests.Datas{
                "grant_type":"client_credentials",
        }
        token_request,_ := requests.Post(SPOTIFY_TOKEN_URL, payload, headers)
        access_token := gjson.Get(token_request.Text(), "access_token")

        return access_token.String()
}

func request_valid_song(access_token string, SPOTIFY_API_URL string) (map[string]string) {
        rand.Seed(time.Now().Unix())
        random_songs_array := [15]string {
                  "%25a%25", "a%25", "%25a",
                  "%25e%25", "e%25", "%25e",
                  "%25i%25", "i%25", "%25i",
                  "%25o%25", "o%25", "%25o",
                  "%25u%25", "u%25", "%25u",
                }
        random_song_choice := random_songs_array[rand.Intn(len(random_songs_array))]
        var song_dict map[string]string
        var max_limit int = 2000
        for {
            fmt.Println("Trying")
            random_offset := strconv.Itoa(rand.Intn(max_limit))
            authorization_header := requests.Header{
                    "Authorization":"Bearer "+access_token,
            }
            request_string := SPOTIFY_API_URL+"/search?query="+random_song_choice+"&offset="+random_offset+"&limit=1&type=track"
            song_requests, err := requests.Get(
                    request_string,
                    authorization_header,
            )

            if err != nil {
              if max_limit > 1000 {
                      max_limit = max_limit - 1000
              } else if max_limit <= 1000 && max_limit > 0 {
                      max_limit = max_limit - 10
              } else {
                      break
              }
            } else {
                    song_info := gjson.Get(song_requests.Text(), "tracks.items.0")
                    fmt.Println(song_info)
                    artist := gjson.Get(song_requests.Text(), "tracks.items.0.artists.0.name").String()
                    song := gjson.Get(song_requests.Text(), "tracks.items.0.name").String()
                    external_url := gjson.Get(song_requests.Text(), "tracks.items.0.album.artists.0.external_urls.spotify").String()
                    song_dict = map[string]string {
                            "song_info": artist+" - "+song,
                            "external_url": external_url,
                    }
                    break
            }

          }
          return song_dict
}
