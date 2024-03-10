package beginning

import (
	"NotSmokeBot/internal/progress"
	"NotSmokeBot/internal/server"
	errlist "NotSmokeBot/pkg/error_list"
	"NotSmokeBot/pkg/postgres"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

const filePathStart = "pkg/templates/start_file.txt"
const filePathBegin = "pkg/templates/progress_file.txt"

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if err := server.Serve.PgDB.TrManager.Do(ctx, func(ctx context.Context) error {
		userId, err := postgres.CreateUser(&server.Serve, ctx, postgres.CreateParams{TgId: int(update.Message.From.ID), LastMes: update.Message.Text})
		_, err = postgres.CreateInfo(&server.Serve, ctx, userId)
		return err
	}); err != nil {
		if err.Error() == errlist.UserExist || err.Error() == errlist.TransactionExist {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "Ты же уже стартовал бота?",
				ParseMode: models.ParseModeMarkdown,
			})
		} else {
			logrus.Error("register %s", err.Error())
		}
		return
	}

	file, err := os.Open(filePathStart)
	if err != nil {
		logrus.Fatalf("opening start file: %s", err.Error())
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		logrus.Fatalf("reading start file: %s", err.Error())
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      string(content),
		ParseMode: models.ParseModeMarkdown,
	})

	done := make(chan struct{})
	timer := time.NewTimer(time.Hour * 24)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-timer.C:
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Ты случаем не забыл о нас?",
				})
				timer.Reset(time.Hour * 24)
			}
		}
	}()

	nickName := getNickname(ctx, b, update)
	err = postgres.UpdateNick(&server.Serve, ctx, postgres.UpdateNickParams{TgId: int(update.Message.From.ID), Nick: nickName})
	if err != nil {
		logrus.Error("nickname ", err.Error())
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Привет, %s", nickName),
	})

	moneyAmount := getMoney(ctx, b, update)
	err = postgres.UpdateMoney(&server.Serve, ctx, postgres.UpdateMoneyParams{TgId: int(update.Message.From.ID), Money: moneyAmount})
	if err != nil {
		logrus.Error("money ", err.Error())
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "ДA!!!", CallbackData: "begin_1"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "ГОТОВЫ БРОСИТЬ КУРИТЬ?",
		ReplyMarkup: kb,
	})

	close(done)
}

func getNickname(ctx context.Context, b *bot.Bot, update *models.Update) string {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Как к тебе обращаться?",
	})
	lastMes, _ := postgres.GetMessage(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.Message.From.ID)})
	for {
		mes, _ := postgres.GetMessage(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.Message.From.ID)})
		if mes != lastMes {
			regex := regexp.MustCompile(`^[^\r\n]{3,20}$`)
			if !regex.MatchString(mes) {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Давай еще одну попытку, ник не должен быть длиннее 20 символов и не должен содержать энтеры",
				})
				lastMes = mes
				continue
			}
			lastMes = mes
			break
		}
	}
	return lastMes
}

func getMoney(ctx context.Context, b *bot.Bot, update *models.Update) int {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Предположи сколько денег ты тратишь в неделю на сиггареты?",
	})
	var rubles float64
	var err error
	lastMes, _ := postgres.GetMessage(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.Message.From.ID)})
	for {
		mes, _ := postgres.GetMessage(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.Message.From.ID)})
		if mes != lastMes {
			rubles, err = strconv.ParseFloat(mes, 64)
			if err != nil || int(rubles) <= 0 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Дружище, скинь чиселку сколько денег на сиггареты тратишь, она обычно еще положительная",
				})
				logrus.Error("money ", err.Error())
				lastMes = mes
				continue
			}
			break
		}
	}

	return int(rubles)
}

func Begin(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	file, err := os.Open(filePathBegin)
	if err != nil {
		logrus.Fatalf("opening start file: %s", err.Error())
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		logrus.Fatalf("reading start file: %s", err.Error())
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		Text:      string(content),
		ParseMode: models.ParseModeMarkdown,
	})

	if update.CallbackQuery.Data == "begin_1" {
		progress.Notification(ctx, b, update)
	}

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
}
