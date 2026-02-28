package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	once      sync.Once
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	langUsed  string
)

func T(id string) string {
	return localize(id)
}

func Tf(id string, args ...any) string {
	return fmt.Sprintf(localize(id), args...)
}

func CurrentLanguage() string {
	once.Do(initBundle)
	if langUsed == "" {
		return "pt-BR"
	}
	return langUsed
}

func localize(id string) string {
	once.Do(initBundle)
	if localizer == nil {
		return id
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil || msg == "" {
		return id
	}
	return msg
}

func initBundle() {
	bundle = i18n.NewBundle(language.MustParse("pt-BR"))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	addDefaults()
	lang := resolveLang()
	langUsed = lang
	loadMessages(lang)
	localizer = i18n.NewLocalizer(bundle, lang, "pt-BR")
}

func addDefaults() {
	bundle.AddMessages(language.MustParse("pt-BR"),
		&i18n.Message{ID: "command_name", Other: "reportarbug"},
		&i18n.Message{ID: "modal_title", Other: "Reportar um bug"},
		&i18n.Message{ID: "modal_title_label", Other: "Titulo"},
		&i18n.Message{ID: "modal_title_placeholder", Other: "Resumo curto do bug"},
		&i18n.Message{ID: "modal_desc_label", Other: "Descricao"},
		&i18n.Message{ID: "modal_desc_placeholder", Other: "Descricao detalhada do bug"},
		&i18n.Message{ID: "cooldown_message", Other: "Aguarde %s antes de enviar outro reporte de bug."},
		&i18n.Message{ID: "issue_failed", Other: "Falha ao criar issue no GitHub. Tente novamente mais tarde."},
		&i18n.Message{ID: "issue_created", Other: "Reporte de bug criado: %s"},
		&i18n.Message{ID: "issue_created_simple", Other: "Reporte de bug enviado com sucesso."},
	)
}

func loadMessages(lang string) {
	path := strings.TrimSpace(os.Getenv("I18N_PATH"))
	if path == "" {
		path = "locales"
	}

	tag, err := language.Parse(lang)
	if err != nil {
		return
	}
	base, _ := tag.Base()

	candidates := []string{
		filepath.Join(path, tag.String()+".json"),
		filepath.Join(path, base.String()+".json"),
	}

	for _, file := range candidates {
		if _, err := os.Stat(file); err == nil {
			_, _ = bundle.LoadMessageFile(file)
			return
		}
	}
}

func resolveLang() string {
	lang := normalizeLang(os.Getenv("APP_LANG"))
	if lang != "" {
		return lang
	}
	lang = normalizeLang(os.Getenv("LANG"))
	if lang != "" {
		return lang
	}
	return "pt-BR"
}

func normalizeLang(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.Split(value, ".")[0]
	value = strings.ReplaceAll(value, "_", "-")
	tag, err := language.Parse(value)
	if err != nil {
		return value
	}
	return tag.String()
}
