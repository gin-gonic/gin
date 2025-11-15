package ut

import (
	"strings"

	"github.com/go-playground/locales"
)

// UniversalTranslator holds all locale & translation data
type UniversalTranslator struct {
	translators map[string]Translator
	fallback    Translator
}

// New returns a new UniversalTranslator instance set with
// the fallback locale and locales it should support
func New(fallback locales.Translator, supportedLocales ...locales.Translator) *UniversalTranslator {

	t := &UniversalTranslator{
		translators: make(map[string]Translator),
	}

	for _, v := range supportedLocales {

		trans := newTranslator(v)
		t.translators[strings.ToLower(trans.Locale())] = trans

		if fallback.Locale() == v.Locale() {
			t.fallback = trans
		}
	}

	if t.fallback == nil && fallback != nil {
		t.fallback = newTranslator(fallback)
	}

	return t
}

// FindTranslator trys to find a Translator based on an array of locales
// and returns the first one it can find, otherwise returns the
// fallback translator.
func (t *UniversalTranslator) FindTranslator(locales ...string) (trans Translator, found bool) {

	for _, locale := range locales {

		if trans, found = t.translators[strings.ToLower(locale)]; found {
			return
		}
	}

	return t.fallback, false
}

// GetTranslator returns the specified translator for the given locale,
// or fallback if not found
func (t *UniversalTranslator) GetTranslator(locale string) (trans Translator, found bool) {

	if trans, found = t.translators[strings.ToLower(locale)]; found {
		return
	}

	return t.fallback, false
}

// GetFallback returns the fallback locale
func (t *UniversalTranslator) GetFallback() Translator {
	return t.fallback
}

// AddTranslator adds the supplied translator, if it already exists the override param
// will be checked and if false an error will be returned, otherwise the translator will be
// overridden; if the fallback matches the supplied translator it will be overridden as well
// NOTE: this is normally only used when translator is embedded within a library
func (t *UniversalTranslator) AddTranslator(translator locales.Translator, override bool) error {

	lc := strings.ToLower(translator.Locale())
	_, ok := t.translators[lc]
	if ok && !override {
		return &ErrExistingTranslator{locale: translator.Locale()}
	}

	trans := newTranslator(translator)

	if t.fallback.Locale() == translator.Locale() {

		// because it's optional to have a fallback, I don't impose that limitation
		// don't know why you wouldn't but...
		if !override {
			return &ErrExistingTranslator{locale: translator.Locale()}
		}

		t.fallback = trans
	}

	t.translators[lc] = trans

	return nil
}

// VerifyTranslations runs through all locales and identifies any issues
// eg. missing plural rules for a locale
func (t *UniversalTranslator) VerifyTranslations() (err error) {

	for _, trans := range t.translators {
		err = trans.VerifyTranslations()
		if err != nil {
			return
		}
	}

	return
}
