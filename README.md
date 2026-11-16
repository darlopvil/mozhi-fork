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

## What this fork adds

**New engines**
- **Official DeepL** — uses the official API when an API key is set (otherwise falls back to the free deeplx).
- **Gemini** (LLM) — through the native API, with a model preference list: a saturated or retired
  model falls through to the next one.
- **TexTra (NICT)** — Japanese engine from a Japanese public research institute, free. It allows
  about 3 requests per 2 minutes for the whole account (translation, dictionary and TTS together):
  Mozhi keeps count, gives translation priority, only requests the dictionary when a slot is left,
  and lets Google read the audio when TexTra has no slot.
- **Groq** (LLM) — open models (gpt-oss) at high speed, with a model preference list.
- **OpenRouter** — three engines, one per model family: ERNIE (Baidu), Gemma (Google) and Nemotron
  (NVIDIA). The current free models of each family are discovered automatically: when OpenRouter
  retires one, the engine moves to the next without touching the configuration.
- **Mistral** (LLM) — European models.

**Dictionary / alternative words** (senses, examples, synonyms, antonyms)
- **Gemini**, **Groq**, **OpenRouter** and **Mistral**: senses with examples for a single word,
  paraphrases for a short phrase, skipped for long text.
- **MyMemory**: alternative translations from its memory.
- **TexTra**: system dictionary for language pairs with Japanese (en↔ja, ja↔ko, ja↔zh, etc.).

**Text-to-speech (TTS)**
- **TexTra** (Japanese, English, Chinese, Korean) and **Groq** (English, Arabic), besides Google.
- Long texts are streamed: split at sentence boundaries within each engine's limit (Google 200
  characters, TexTra 120, Groq 600), and audio starts playing as soon as the first chunk is ready.
- Audio is only generated when you press play (`preload="none"`) and is sent with its real type (MP3 or WAV).

**Reliability**
- A failing or hanging engine no longer crashes the instance or holds the others back: every engine
  has its own timeout, panics are recovered, and All mode answers within 50 s.
- DuckDuckGo retries on a fresh connection during its blocking windows.
- Yandex and MyMemory no longer show their errors (or MyMemory's quota warning) as the translation.
- `_ENABLED` is honoured in the web UI, the API, All mode and TTS, and text length is capped.

**Fixes**
- Texts containing `&` or `+` broke Google (empty translation), its dictionary and Yandex's transliteration.
- Unescaped TTS links; `/switchlanguages` uses 303 instead of a 301 that browsers cached.
- Google TTS fixed; MyMemory uses POST requests plus an optional email (more quota, longer texts).
- Dead engines removed: Reverso (IP-blocked on servers) and LibreTranslate (public instances shut down).

**Interface**
- Service icons in the engine selector, sorted alphabetically.
- The chosen language survives engine switches, even when engines name it differently
  (`es` / `es-ES`, `zh-CN` / `zh-Hans`, `nb` / `no`).
- Favourite languages pinned at the top of the selectors.

## Configuration

Everything is set through environment variables in a `.env` file that is **not** versioned (use
placeholders, never commit real keys):

```bash
# DeepL (optional; without a key it falls back to the free deeplx)
MOZHI_DEEPL_API_KEY=your_key_here            # deepl.com/pro-api

# Gemini (aistudio.google.com, free, no card)
MOZHI_GEMINI_API_KEY=your_key_here
# MOZHI_GEMINI_MODEL=gemini-3.5-flash,gemini-flash-lite-latest,gemini-flash-latest  # optional: preference list

# TexTra / NICT (mt-auto-minhon-mlt.ucri.jgn-x.jp, free registration)
MOZHI_TEXTRA_API_KEY=your_key_here
MOZHI_TEXTRA_API_SECRET=your_secret_here
MOZHI_TEXTRA_LOGIN_ID=your_login_id

# Groq (console.groq.com, free)
MOZHI_GROQ_API_KEY=your_key_here
# MOZHI_GROQ_MODEL=openai/gpt-oss-20b,openai/gpt-oss-120b  # optional: preference list, first active one is used
# MOZHI_GROQ_TTS_VOICE=diana                 # optional

# OpenRouter (openrouter.ai) — one key for the three engines
MOZHI_OPENROUTER_API_KEY=your_key_here
# MOZHI_OPENROUTER_GPTOSS_FAMILY=baidu/ernie       # optional: model family of each engine
# MOZHI_OPENROUTER_GEMMA_FAMILY=google/gemma
# MOZHI_OPENROUTER_NEMOTRON_FAMILY=nvidia/nemotron
# MOZHI_OPENROUTER_GPTOSS_MODEL / _GEMMA_MODEL / _NEMOTRON_MODEL  # optional: preferred model, while it exists

# Mistral (console.mistral.ai)
MOZHI_MISTRAL_API_KEY=your_key_here
# MOZHI_MISTRAL_MODEL=mistral-small-latest   # optional

# MyMemory (optional: an email raises the daily quota)
# MOZHI_MYMEMORY_EMAIL=you@example.com

# General
# MOZHI_PORT=3000                            # MOZHI_HOST: interface to listen on (all by default)
# MOZHI_DEFAULT_ENGINE=google
# MOZHI_DEFAULT_SOURCE_LANG=auto             # MOZHI_DEFAULT_PREFER_AUTODETECT=true to prefer "Detect language"
# MOZHI_DEFAULT_TARGET_LANG=en
# MOZHI_PINNED_LANGS=ja,es,en                # languages pinned at the top of the selectors
# MOZHI_MAX_TEXT_LENGTH=5000                 # characters; 0 = no limit
```

Any engine can be disabled with `MOZHI_<ID>_ENABLED=false` (web UI, API, All mode and TTS at once).
IDs: `google deepl duckduckgo mymemory yandex gemini textra groq gptoss gemma nemotron mistral`.

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
