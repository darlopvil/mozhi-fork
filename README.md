<div align="center">
  <img src="public/assets/mozhi.png" width="128" height="128" alt="Mozhi logo">
  <h1>Mozhi — personal fork</h1>
</div>

Fork of [aryak/mozhi](https://codeberg.org/aryak/mozhi), a private alternative frontend for many
translation engines, and of its engine library
[aryak/libmozhi](https://codeberg.org/aryak/libmozhi) → [darlopvil/libmozhi-fork](https://github.com/darlopvil/libmozhi-fork).

I maintain it for my own instance: I don't publish instances or send changes upstream.
Issues are welcome at [darlopvil/mozhi-fork](https://github.com/darlopvil/mozhi-fork/issues).

License: AGPL-3.0, same as upstream ([LICENSE](LICENSE)).

---

## Changes from upstream

- Reverso and LibreTranslate removed: Reverso is IP-blocked on servers and LibreTranslate's public instances shut down.
- DeepL uses the official API when an API key is set (the free deeplx otherwise).
- MyMemory uses POST requests plus an optional email (more quota, longer texts).
- `MOZHI_<ID>_ENABLED` works for every engine and everywhere: web UI, API, All mode and TTS. All mode only uses the enabled engines, in a stable order.
- Text length is capped (`MOZHI_MAX_TEXT_LENGTH`, 5000 characters by default).
- Audio is only generated when you press play, and is sent with its real type (MP3 or WAV).

## Fixes

- One failing or hanging engine crashed the whole instance or stalled All mode. Every engine now has its own timeout, panics are recovered, and All mode answers within 50 s.
- Texts containing `&` or `+` broke Google (empty translation), its dictionary and Yandex's transliteration.
- Google TTS fixed.
- Yandex and MyMemory crashed on error responses, and MyMemory's daily quota warning was shown as the translation.
- DuckDuckGo crashed when its token request failed; it now retries on a fresh connection during its blocking windows.
- Gemini hung when its model was saturated; it now uses the native API and falls back across models.
- Free LLM models retired by their providers broke their engines until the configuration was edited by hand.
- The OpenRouter engines (one free pool shared by all its users) and Mistral's free tier (2 requests a minute at most) answered 429 just when they were needed. They were replaced by engines with a quota of their own, or with a public web endpoint that needs no key.
- ModernMT's web translator stopped answering ("your account status is suspended"): ModernMT turns into Lara, which is paid, and closes at the end of 2026, so the engine was removed.
- The chosen language was lost when switching engines, because each engine names some languages its own way (`es` / `es-ES`, `zh-CN` / `zh-Hans`, `nb` / `no`).
- TTS links weren't escaped, and `/switchlanguages` answered with a 301 that browsers cached.
- Only one engine could be disabled at a time (chained `_ENABLED` checks).
- Yandex showed its error page ("Request-URI Too Large") as the transliteration of long texts, and swapped the source and target transliterations.
- TexTra and Groq voices were offered for languages they don't speak, so the other side of a translation never played; their WAV audio showed a duration of tens of hours.
- When TexTra was down (e.g. a maintenance), its voice went silent and /api/health only said "no access token": now Google reads instead, and the error says why (under maintenance, the token endpoint's status, or TexTra's reason).

## Features

- New engines: Gemini, TexTra (NICT, Japanese), Groq, Qwen (through Groq), Cohere, Microsoft Translator (with romanization of non-Latin scripts) and Langbly, plus Papago and Volcengine through the public endpoints of their own web translators.
- Model preference lists: when a provider retires a model, Gemini, Groq, Qwen and Cohere move on to the next one without touching the configuration (Gemini also when a model is saturated).
- Dictionary and alternative words: senses, examples and paraphrases from the LLM engines, alternatives from MyMemory, TexTra's system dictionary for Japanese pairs, and Microsoft Translator's dictionary with examples for pairs with English.
- Text-to-speech with Microsoft (neural voices for most languages, through Azure Speech with a key, and through Edge's read-aloud service, with no key, as its backup or without a key), TexTra (Japanese, English, Chinese, Korean) and Groq (English, Arabic), besides Google. Each side of a translation is read in its own language: by the engine's own voice when it has that language, otherwise by Microsoft, or by Google.
- Long texts are read aloud in streamed chunks: audio starts as soon as the first chunk is ready.
- TexTra's small request limit (about 3 requests per 2 minutes for the whole account) is shared sensibly: translation first, the dictionary only when a slot is left, and Google reads the audio when TexTra has no slot.
- Japanese dictionary: a short Japanese text (a word or a phrase, in the source or in a single engine's translation) gets a panel with Jisho (readings, English senses, JLPT level), jpdb (senses with examples, frequency, conjugated forms), MOJi (pitch accent, romaji, Chinese senses) and Weblio (a Japanese monolingual dictionary), whatever the engine. Also as JSON at `/api/jadict?word=`.
- Favourite languages pinned at the top of the selectors (`MOZHI_PINNED_LANGS`).
- Service icons in the engine selector, sorted alphabetically.
- `/api/health`: the state of every engine, taken from real requests. `?engine=<id>` (plus `&kind=tts` for speech; the dictionaries are `jisho`, `jpdb`, `moji` and `weblio`, and Edge's voices `edge`) answers 503 while that engine is failing, for uptime monitors. An engine that is failing, unused since start or unused for 12 hours is re-checked first with one short request (at most every 2 minutes), so a broken engine shows up even when nobody uses it, and turns green again as soon as it recovers. TexTra is only re-checked when it has slots to spare, and error messages are shown without keys, account ids or IP addresses.

## Configuration

Everything is set through environment variables in a `.env` file that is **not** versioned (use
placeholders, never commit real keys):

```bash
# DeepL (optional; without a key it falls back to the free deeplx)
MOZHI_DEEPL_API_KEY=your_key_here            # deepl.com/pro-api

# Gemini (aistudio.google.com, free, no card)
MOZHI_GEMINI_API_KEY=your_key_here
# MOZHI_GEMINI_MODEL=gemini-3.5-flash,gemini-flash-lite-latest,gemini-flash-latest,gemini-3.1-flash-lite,gemini-3.5-flash-lite  # optional: preference list

# TexTra / NICT (mt-auto-minhon-mlt.ucri.jgn-x.jp, free registration)
MOZHI_TEXTRA_API_KEY=your_key_here
MOZHI_TEXTRA_API_SECRET=your_secret_here
MOZHI_TEXTRA_LOGIN_ID=your_login_id

# Groq (console.groq.com, free)
MOZHI_GROQ_API_KEY=your_key_here
# MOZHI_GROQ_MODEL=openai/gpt-oss-20b,openai/gpt-oss-120b  # optional: preference list, first active one is used
# MOZHI_GROQ_TTS_VOICE=diana                 # optional

# Qwen, through Groq: same key as Groq, rate-limited separately
# MOZHI_QWEN_MODEL=qwen/qwen3.8-27b          # optional: preference list

# Cohere (dashboard.cohere.com, free trial key: 1,000 calls a month)
MOZHI_COHERE_API_KEY=your_key_here
# MOZHI_COHERE_MODEL=north-small-translate-1.0,command-a-translate-08-2025  # optional: translation models
# MOZHI_COHERE_DICT_MODEL=command-a-plus-05-2026,command-a-03-2025          # optional: dictionary models

# Microsoft Translator (an Azure Translator resource; free tier F0: 2 million characters a month)
MOZHI_MICROSOFT_API_KEY=your_key_here
# MOZHI_MICROSOFT_REGION=westeurope          # only for regional resources, not global ones
# MOZHI_MICROSOFT_SPEECH_KEY=your_key_here    # text-to-speech: an Azure Speech resource (free tier F0: 500,000 characters a month)
# MOZHI_MICROSOFT_SPEECH_REGION=westeurope    # the Speech resource's region
# MOZHI_MICROSOFT_TTS_VOICES=ja=ja-JP-NanamiNeural,es=es-ES-ElviraNeural  # optional: voice per language
# MOZHI_EDGE_ENABLED=false                 # Edge's read-aloud voices (no key): Azure Speech's backup, on by default

# Langbly (langbly.com; 500,000 characters a month free, set a spending limit to never go over it)
MOZHI_LANGBLY_API_KEY=your_key_here

# MyMemory (optional: an email raises the daily quota)
# MOZHI_MYMEMORY_EMAIL=you@example.com

# Japanese dictionaries (no keys): on by default
# MOZHI_JISHO_ENABLED=false                # also MOZHI_JPDB_, MOZHI_MOJI_ and MOZHI_WEBLIO_ENABLED

# General
# MOZHI_PORT=3000                            # MOZHI_HOST: interface to listen on (all by default)
# MOZHI_DEFAULT_ENGINE=google
# MOZHI_DEFAULT_SOURCE_LANG=auto             # MOZHI_DEFAULT_PREFER_AUTODETECT=true to prefer "Detect language"
# MOZHI_DEFAULT_TARGET_LANG=en
# MOZHI_PINNED_LANGS=ja,es,en                # languages pinned at the top of the selectors
# MOZHI_MAX_TEXT_LENGTH=5000                 # characters; 0 = no limit
```

Any engine can be disabled with `MOZHI_<ID>_ENABLED=false` (web UI, API, All mode and TTS at once).
IDs: `google deepl duckduckgo mymemory yandex gemini textra groq qwen cohere microsoft langbly papago volcengine`.
Japanese dictionaries: `jisho jpdb moji weblio`.

The API is documented with Swagger at `/api/swagger/index.html` on the instance.

## Local build (self-hosted)

This fork is built locally, without a registry. Expected layout:

    ~/mozhi/      → this repo (code + docker-compose.yml + .env, not versioned)
    ~/libmozhi/   → libmozhi fork (the engines), wired through a `replace` in go.mod

To rebuild after a change and recreate the container:

    ./rebuild.sh

Requires `docker-compose.yml` to use `image: mozhi-fork:local`.

## Credits (upstream)

Mozhi (Tamil *moḻi*, "language") is the work of [Arya](https://aryak.me) and contributors:

- [Arya](https://aryak.me): creator
- [Midou36o](https://midou.dev): made the logo
- [py_](https://github.com/supercolbat): Design files
- [Missuo](https://github.com/missuo): creating gDeepLX that does the hard part of making DeepL work
- [translatepy](https://github.com/Animenosekai/translate): giving me the format of request for yandex engine
- [SimplyTranslate](https://codeberg.org/simpleweb/simplytranslate): Inspiration and base code for the webui
- [Rimgo](https://codeberg.org/rimgo/rimgo): Code for embedding html in binary
