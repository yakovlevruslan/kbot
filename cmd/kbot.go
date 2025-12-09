package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v4"
)

var (
	TELE_Token  = os.Getenv("TELE_TOKEN")  // Токен @BotFather
	OWM_API_Key = os.Getenv("OWM_API_KEY") // Токен openweathermap.org

)

const BaseURL = "https://api.openweathermap.org/data/2.5/weather"

// WeatherResponse API
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// kbotCmd команда kbot, яка запускає бота
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Запускає Telegram-бота для отримання погоди",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started. Connecting to Telegram...\n", appVersion)

		if TELE_Token == "" {
			log.Fatal("Помилка: Не встановлена змінна TELE_TOKEN")
		}
		if OWM_API_Key == "" {
			log.Fatal("Помилка: Не встановлена змінна OWM_API_KEY")
		}

		kbot, err := telebot.NewBot(telebot.Settings{
			Token:  TELE_Token,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatalf("Check TELE_Token env. %s", err)
		}

		// Обробник команди /start
		kbot.Handle("/start", func(m telebot.Context) error {
			return m.Send(fmt.Sprintf("Привіт! Я Wbot %s. Щоб дізнатись погоду, відправ /w [місто]. Наприклад: /w Kyiv", appVersion))
		})

		// Обробник команди /w (weather)
		kbot.Handle("/w", func(m telebot.Context) error {
			parts := strings.SplitN(m.Text(), " ", 2)
			if len(parts) < 2 {
				return m.Send("Будьласка вкажіть місто. Наприклад: /w Kyiv")
			}
			city := parts[1]

			weatherData, err := getCurrentWeather(city, OWM_API_Key)
			if err != nil {
				log.Printf("Помилка отримання погоди для %s: %v", city, err)
				return m.Send(fmt.Sprintf("Не вдалося отримати погоду для міста \"%s\". Це вірна назва вміста?", city))
			}

			// Формуємо відповідь
			responseMsg := fmt.Sprintf(
				"Погода в мітс *%s*:\n\n"+
					"🌡 Температура: `%.1f °C`\n"+
					"🧍 Відчувається як: `%.1f °C`\n"+
					"💧 Вологість: `%d%%`\n"+
					"📝 Опис: *%s*",
				weatherData.Name,
				weatherData.Main.Temp,
				weatherData.Main.FeelsLike,
				weatherData.Main.Humidity,
				weatherData.Weather[0].Description,
			)

			return m.Send(responseMsg, telebot.ModeMarkdown)
		})

		log.Println("Bot is running...")
		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}

// getCurrentWeather запит до API OpenWeatherMap
func getCurrentWeather(city, apiKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=ua", BaseURL, city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Відповідь з помилками від OWM API (наприклад, City not found)
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Помилка API запросу, статус: %d, відповідь: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}
