package middlewares

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/context"
	"gopkg.in/ini.v1"
	"reflect"
	"strings"
)

type I18nConfig struct {
	// Default set it if you want a default language
	//
	// Checked: Configuration state, not at runtime
	Default string
	// URLParameter is the name of the url parameter which the language can be indentified
	//
	// Checked: Serving state, runtime
	URLParameter string
	// Languages is a map[string]string which the key is the language i81n and the value is the file location
	//
	// Example of key is: 'en-US'
	// Example of value is: './locales/en-US.ini'
	Languages map[string]string
}

// test file: ../../_examples/miscellaneous/i18n/main_test.go
type I18nMiddleware struct {
	config I18nConfig
}

// ServeHTTP serves the request, the actual middleware's job is here
func (i *I18nMiddleware) ServeHTTP(ctx context.Context) {
	wasByCookie := false

	langKey := ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey()
	language := ctx.Values().GetString(langKey)
	if language == "" {
		// try to get by url parameter
		language = ctx.URLParam(i.config.URLParameter)
		if language == "" {
			// then try to take the lang field from the cookie
			language = ctx.GetCookie(langKey)

			if len(language) > 0 {
				wasByCookie = true
			} else {
				// try to get by the request headers.
				langHeader := ctx.GetHeader("Accept-Language")
				if len(langHeader) > 0 {
					for _, langEntry := range strings.Split(langHeader, ",") {
						lc := strings.Split(langEntry, ";")[0]
						if lc, ok := IsExistSimilar(lc); ok {
							language = lc
							break
						}
					}
				}
			}
		}
		// if it was not taken by the cookie, then set the cookie in order to have it
		if !wasByCookie {
			ctx.SetCookieKV(langKey, language)
		}

		if language == "" {
			language = i.config.Default
		}

		ctx.Values().Set(langKey, language)
	}
	locale := Locale{Lang: language}

	// if unexpected language given, the middleware will  transtlate to the default language, the language key should be
	// also this language instead of the user-given
	if indexLang := locale.Index(); indexLang == -1 {
		locale.Lang = i.config.Default
	}

	translateFuncKey := ctx.Application().ConfigurationReadOnly().GetTranslateFunctionContextKey()
	ctx.Values().Set(translateFuncKey, locale.Tr)
	ctx.Next()
}

// Translate returns the translated word from a context
// the second parameter is the key of the world or line inside the .ini file
// the third parameter is the '%s' of the world or line inside the .ini file
func Translate(ctx context.Context, format string, args ...interface{}) string {
	return ctx.Translate(format, args...)
}

// NewI18nMiddleware returns a new i18n middleware
func NewI18nMiddleware(c I18nConfig) context.Handler {
	if len(c.Languages) == 0 {
		panic("You cannot use this middleware without set the Languages option, please try again and read the _example.")
	}
	i := &I18nMiddleware{config: c}
	firstlanguage := ""
	//load the files
	for k, langFileOrFiles := range c.Languages {
		// remove all spaces.
		langFileOrFiles = strings.Replace(langFileOrFiles, " ", "", -1)
		// note: if only one, then the first element is the "v".
		languages := strings.Split(langFileOrFiles, ",")

		for _, v := range languages { // loop each of the files separated by comma, if any.
			if !strings.HasSuffix(v, ".ini") {
				v += ".ini"
			}
			err := SetMessage(k, v)
			if err != nil && err != ErrLangAlreadyExist {
				panic("Failed to set locale file'" + k + "' Error:" + err.Error())
			}
			if firstlanguage == "" {
				firstlanguage = k
			}
		}
	}
	// if not default language set then set to the first of the i.config.Languages
	if c.Default == "" {
		c.Default = firstlanguage
	}

	SetDefaultLang(i.config.Default)
	return i.ServeHTTP
}

// TranslatedMap returns translated map[string]interface{} from i18n structure.
func TranslatedMap(ctx context.Context, sourceInterface interface{}) map[string]interface{} {
	iType := reflect.TypeOf(sourceInterface).Elem()
	result := make(map[string]interface{})

	for i := 0; i < iType.NumField(); i++ {
		fieldName := reflect.TypeOf(sourceInterface).Elem().Field(i).Name
		fieldValue := reflect.ValueOf(sourceInterface).Elem().Field(i).String()

		result[fieldName] = Translate(ctx, fieldValue)
	}

	return result
}

var (
	// ErrLangAlreadyExist throwed when language is already exists.
	ErrLangAlreadyExist = errors.New("lang already exists")

	locales = &localeStore{store: make(map[string]*locale)}
)

// add support for multi language file per language.
// ini has already implement a  *ini.File#Append
// BUT IT DOESN'T F WORKING, SO:
type localeFiles struct {
	files []*ini.File
}

// Get returns a the value from the "keyName" on the "sectionName"
// by searching all loc.files.
func (loc *localeFiles) GetKey(sectionName, keyName string) (*ini.Key, error) {
	for _, f := range loc.files {
		// returns the first available.
		// section is the same for both files if key exists.
		if sec, serr := f.GetSection(sectionName); serr == nil && sec != nil {
			if k, err := sec.GetKey(keyName); err == nil && k != nil {
				return k, err
			}
		}
	}

	return nil, fmt.Errorf("not found")
}

// Reload reloads and parses all data sources.
func (loc *localeFiles) Reload() error {
	for _, f := range loc.files {
		if err := f.Reload(); err != nil {
			return err
		}
	}
	return nil
}

func (loc *localeFiles) addFile(file *ini.File) error {
	loc.files = append(loc.files, file)
	return loc.Reload()
}

type locale struct {
	id       int
	lang     string
	langDesc string
	message  *localeFiles
}

type localeStore struct {
	langs       []string
	langDescs   []string
	store       map[string]*locale
	defaultLang string
}

// Get target language string
func (d *localeStore) Get(lang, section, format string) (string, bool) {
	if locale, ok := d.store[lang]; ok {
		// println(lang + " language found, let's see keys")
		if key, err := locale.message.GetKey(section, format); err == nil && key != nil {
			// println("value for section= " + section + "and key=" + format + " found")
			return key.Value(), true
		}
	}

	if len(d.defaultLang) > 0 && lang != d.defaultLang {
		// println("use the default lang: " + d.defaultLang)
		return d.Get(d.defaultLang, section, format)
	}

	return "", false
}

func (d *localeStore) Add(lang, langDesc string, source interface{}, others ...interface{}) error {

	file, err := ini.Load(source, others...)
	if err != nil {
		return err
	}
	file.BlockMode = false

	// if already exists add the file on this language.
	lc, ok := d.store[lang]
	if !ok {
		// println("add lang and init message: " + lang)
		// create a new one if doesn't exist.

		lc = new(locale)
		lc.message = new(localeFiles)
		lc.lang = lang
		lc.langDesc = langDesc
		lc.id = len(d.langs)
		// add the first language if not exists.
		d.langs = append(d.langs, lang)
		d.langDescs = append(d.langDescs, langDesc)
		d.store[lang] = lc
	}

	// println("append a file for language: " + lang)

	return lc.message.addFile(file)
}

func (d *localeStore) Reload(langs ...string) (err error) {
	if len(langs) == 0 {
		for _, lc := range d.store {
			if err = lc.message.Reload(); err != nil {
				return err
			}
		}
	} else {
		for _, lang := range langs {
			if lc, ok := d.store[lang]; ok {
				if err = lc.message.Reload(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// SetDefaultLang sets default language which is a indicator that
// when target language is not found, try find in default language again.
func SetDefaultLang(lang string) {
	locales.defaultLang = lang
}

// ReloadLangs reloads locale files.
func ReloadLangs(langs ...string) error {
	return locales.Reload(langs...)
}

// Count returns number of languages that are registered.
func Count() int {
	return len(locales.langs)
}

// ListLangs returns list of all locale languages.
func ListLangs() []string {
	langs := make([]string, len(locales.langs))
	copy(langs, locales.langs)
	return langs
}

func ListLangDescs() []string {
	langDescs := make([]string, len(locales.langDescs))
	copy(langDescs, locales.langDescs)
	return langDescs
}

// IsExist returns true if given language locale exists.
func IsExist(lang string) bool {
	_, ok := locales.store[lang]
	return ok
}

// IsExistSimilar returns true if the language, or something similar
// exists (e.g. en-US maps to en).
// it returns the found name and whether it was able to match something.
// - PATCH by @j-lenoch.
func IsExistSimilar(lang string) (string, bool) {
	_, ok := locales.store[lang]
	if ok {
		return lang, true
	}

	// remove the internationalization element from the IETF code
	code := strings.Split(lang, "-")[0]

	for _, lc := range locales.store {
		if strings.Contains(lc.lang, code) {
			return lc.lang, true
		}
	}

	return "", false
}

// IndexLang returns index of language locale,
// it returns -1 if locale not exists.
func IndexLang(lang string) int {
	if lc, ok := locales.store[lang]; ok {
		return lc.id
	}
	return -1
}

// GetLangByIndex return language by given index.
func GetLangByIndex(index int) string {
	if index < 0 || index >= len(locales.langs) {
		return ""
	}
	return locales.langs[index]
}

func GetDescriptionByIndex(index int) string {
	if index < 0 || index >= len(locales.langDescs) {
		return ""
	}

	return locales.langDescs[index]
}

func GetDescriptionByLang(lang string) string {
	return GetDescriptionByIndex(IndexLang(lang))
}

func SetMessageWithDesc(lang, langDesc string, localeFile interface{}, otherLocaleFiles ...interface{}) error {
	return locales.Add(lang, langDesc, localeFile, otherLocaleFiles...)
}

// SetMessage sets the message file for localization.
func SetMessage(lang string, localeFile interface{}, otherLocaleFiles ...interface{}) error {
	return SetMessageWithDesc(lang, lang, localeFile, otherLocaleFiles...)
}

// Locale represents the information of localization.
type Locale struct {
	Lang string
}

// Tr translates content to target language.
func (l Locale) Tr(format string, args ...interface{}) string {
	return Tr(l.Lang, format, args...)
}

// Index returns lang index of LangStore.
func (l Locale) Index() int {
	return IndexLang(l.Lang)
}

// Tr translates content to target language.
func Tr(lang, format string, args ...interface{}) string {
	var section string

	idx := strings.IndexByte(format, '.')
	if idx > 0 {
		section = format[:idx]
		format = format[idx+1:]
	}

	value, ok := locales.Get(lang, section, format)
	if ok {
		format = value
	}

	if len(args) > 0 {
		params := make([]interface{}, 0, len(args))
		for _, arg := range args {
			if arg == nil {
				continue
			}

			val := reflect.ValueOf(arg)
			if val.Kind() == reflect.Slice {
				for i := 0; i < val.Len(); i++ {
					params = append(params, val.Index(i).Interface())
				}
			} else {
				params = append(params, arg)
			}
		}
		return fmt.Sprintf(format, params...)
	}
	return format
}