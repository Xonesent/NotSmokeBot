package progress

import (
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/postgres"
	"NotSmokeBot/pkg/templates"
	"NotSmokeBot/pkg/utilities"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"time"
)

func Notification(ctx context.Context, b *bot.Bot, update *models.Update) {
	currentTime := time.Now()
	timeFlags, timerFlags, timers := []int{11, 16, 20, 21}, []time.Time{}, []*time.Timer{}
	for _, value := range timeFlags {
		timerFlags = append(timerFlags, time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), value, 0, 0, 0, currentTime.Location()))
	}
	for index := range timerFlags {
		if currentTime.After(timerFlags[index]) {
			timerFlags[index] = timerFlags[index].AddDate(0, 0, 1)
		}
		timers = append(timers, time.NewTimer(timerFlags[index].Sub(currentTime)))
	}
	timerFlags[3] = timerFlags[3].AddDate(0, 0, 7)

	go func(ctx context.Context, b *bot.Bot, update *models.Update) {
		for {
			select {
			case <-timers[0].C:
				sendNotification(ctx, b, update)
				timerFlags[0] = timerFlags[0].AddDate(0, 0, 1)
				timers[0].Reset(timerFlags[0].Sub(time.Now()))
			case <-timers[1].C:
				sendNotification(ctx, b, update)
				timerFlags[1] = timerFlags[1].AddDate(0, 0, 1)
				timers[1].Reset(timerFlags[1].Sub(time.Now()))
			case <-timers[2].C:
				sendSurvey(ctx, b, update)
				timerFlags[2] = timerFlags[2].AddDate(0, 0, 1)
				timers[2].Reset(timerFlags[2].Sub(time.Now()))
			case <-timers[3].C:
				sendStatistic(ctx, b, update)
				timerFlags[3] = timerFlags[3].AddDate(0, 0, 7)
				timers[3].Reset(timerFlags[3].Sub(time.Now()))
			}
		}
	}(ctx, b, update)
}

//const image_path = "pkg/images"

func sendNotification(ctx context.Context, b *bot.Bot, update *models.Update) {
	//fileData, errReadFile := os.ReadFile(image_path + "/" + strconv.Itoa(utilities.RandomInRangeInt(1, 4)) + ".png")
	//if errReadFile != nil {
	//	fmt.Printf("error read file, %v\n", errReadFile)
	//	return
	//}
	//
	//params := &bot.SendPhotoParams{
	//	ChatID:  update.Message.Chat.ID,
	//	Photo:   &models.InputFileUpload{Filename: "facebook.png", Data: bytes.NewReader(fileData)},
	//	Caption: "Владимир Путин молодец!",
	//}
	//b.SendPhoto(ctx, params)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Заметано", CallbackData: "quotation_1"},
			},
		},
	}

	num, _ := postgres.GetQuotation(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
	fmt.Println(num, update.CallbackQuery.Message.Message.Chat.ID)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        templates.QuotationsSlice[num],
		ReplyMarkup: kb,
	})
}

func sendSurvey(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Да", CallbackData: "button_1"},
				{Text: "Нет", CallbackData: "button_2"},
			},
		},
	}

	nick, _ := postgres.GetNick(&server.Serve, ctx, postgres.GetNickParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        fmt.Sprintf("Привет, %s, курил сегодня?", nick),
		ReplyMarkup: kb,
	})
}

func sendStatistic(ctx context.Context, b *bot.Bot, update *models.Update) {
	nick, _ := postgres.GetNick(&server.Serve, ctx, postgres.GetNickParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
	money, _ := postgres.GetMoney(&server.Serve, ctx, postgres.GetMoneyParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
	progress, _ := postgres.GetProgress(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
	days := utilities.GetDays(progress)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   fmt.Sprintf("Привет, %s, твой текущий рекорд %d %s за это время ты сэкономил %d", nick, progress, days, money*progress/7),
	})
}
