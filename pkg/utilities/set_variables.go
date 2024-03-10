package utilities

import (
	"bytes"
	"log"
	"text/template"
)

func ReplaceVariables(text string, variables map[string]interface{}) string {
	//variables := map[string]interface{}{
	//	"Name": bot.EscapeMarkdown(update.Message.From.FirstName),
	//}
	//message := replaceVariables(string(content), variables)

	tmpl, err := template.New("message").Parse(text)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
