> **Fork notice:** versión modificada de [aryak/mozhi](https://codeberg.org/aryak/mozhi),
> mantenida por [darlopvil](https://github.com/darlopvil). Cambios y roadmap en los
> [issues](https://github.com/darlopvil/mozhi-fork/issues). Licencia AGPLv3, igual que upstream.

---

## Qué añade este fork

Traductores, diccionario y voz sobre el Mozhi original:

**Motores nuevos**
- **DeepL oficial** — usa la API oficial cuando hay API key (si no, cae al deeplx gratuito).
- **Gemini** (LLM) — traducción de alta calidad con contexto y matiz.
- **TexTra (NICT)** — motor japonés de un organismo público japonés, gratuito. Su límite es de unas 3 peticiones cada 2 minutos para toda la cuenta (traducción, diccionario y TTS juntos): Mozhi lleva la cuenta, da prioridad a la traducción, pide el diccionario solo si sobra hueco y, si no queda hueco para el audio, lo lee Google.
- **Groq** (LLM) — modelos abiertos (gpt-oss) a gran velocidad.
- **OpenRouter** — tres motores separados por familia de modelos: ERNIE (Baidu), Gemma (Google) y Nemotron (NVIDIA).
  Los modelos gratis vigentes de cada familia se descubren solos: si OpenRouter jubila uno, el motor
  pasa al siguiente sin tocar la configuración.
- **Mistral** (LLM) — modelos europeos.

**Diccionario / palabras alternativas** (senses, ejemplos, sinónimos, antónimos)
- En **Gemini**, **Groq**, **OpenRouter** y **Mistral**: para una palabra da acepciones con
  ejemplos; para una frase corta, paráfrasis; para texto largo se omite.
- En **MyMemory**: traducciones alternativas de su memoria.
- En **TexTra**: diccionario de sistema para pares con japonés (en↔ja, ja↔ko, ja↔zh, etc.).

**Texto a voz (TTS)**
- **TexTra** (japonés, inglés, chino, coreano) y **Groq** (inglés, árabe), además del Google ya existente.
- Textos largos en streaming: se parten por frases según el límite de cada motor (Google 200 caracteres, TexTra 120, Groq 600) y el audio empieza a sonar en cuanto está el primer trozo.

**Correcciones y limpieza**
- TTS de Google arreglado.
- MyMemory con peticiones POST + email opcional (más cuota y textos largos).
- Motores muertos fuera: Reverso (bloqueo de IP en servidores) y LibreTranslate (instancias públicas caídas).
- Bug del `_ENABLED` encadenado (solo desactivaba un motor) corregido.
- Mejoras de UI: iconos de cada servicio en el selector y lista ordenada alfabéticamente.

## Configuración de los motores nuevos

Cada motor se activa poniendo su API key en el `.env` (que **no** se versiona). Nombres y dónde
obtener cada clave (usa placeholders, nunca subas claves reales):

```bash
# DeepL (opcional; sin key cae al deeplx gratuito)
MOZHI_DEEPL_API_KEY=your_key_here            # deepl.com/pro-api

# Gemini (aistudio.google.com, gratis sin tarjeta)
MOZHI_GEMINI_API_KEY=your_key_here
# MOZHI_GEMINI_MODEL=gemini-3.5-flash,gemini-flash-lite-latest,gemini-flash-latest  # opcional: lista de preferencia

# TexTra / NICT (mt-auto-minhon-mlt.ucri.jgn-x.jp, registro gratuito)
MOZHI_TEXTRA_API_KEY=your_key_here
MOZHI_TEXTRA_API_SECRET=your_secret_here
MOZHI_TEXTRA_LOGIN_ID=your_login_id

# Groq (console.groq.com, gratis)
MOZHI_GROQ_API_KEY=your_key_here
# MOZHI_GROQ_MODEL=openai/gpt-oss-20b,openai/gpt-oss-120b  # opcional: lista de preferencia, se usa el primero activo
# MOZHI_GROQ_TTS_VOICE=diana                 # opcional

# OpenRouter (openrouter.ai) — una sola key para los tres motores
MOZHI_OPENROUTER_API_KEY=your_key_here
# MOZHI_OPENROUTER_GPTOSS_FAMILY=baidu/ernie       # opcionales: familia de cada motor
# MOZHI_OPENROUTER_GEMMA_FAMILY=google/gemma
# MOZHI_OPENROUTER_NEMOTRON_FAMILY=nvidia/nemotron
# MOZHI_OPENROUTER_GPTOSS_MODEL / _GEMMA_MODEL / _NEMOTRON_MODEL  # opcionales: modelo preferido, mientras exista

# Mistral (console.mistral.ai)
MOZHI_MISTRAL_API_KEY=your_key_here
# MOZHI_MISTRAL_MODEL=mistral-small-latest   # opcional
```

También se puede desactivar cualquier motor con `MOZHI_<NOMBRE>_ENABLED=false` (web, API y modo All), y limitar la longitud del texto con `MOZHI_MAX_TEXT_LENGTH` (por defecto 5000 caracteres; 0 = sin límite).

Idiomas favoritos: con `MOZHI_PINNED_LANGS=ja,es,en` esos idiomas salen los primeros en los dos desplegables (cada motor con su propio código: `es` es `es-ES` en MyMemory).

## Build local (self-hosted)

Este fork se construye localmente, sin registry. Estructura esperada:

    ~/mozhi/      → este repo (código + docker-compose.yml + .env, no versionados)
    ~/libmozhi/   → fork de libmozhi (el motor); enganchado vía `replace` en go.mod

Para reconstruir tras un cambio y recrear el contenedor:

    ./rebuild.sh

Requiere que `docker-compose.yml` use `image: mozhi-fork:local`.
Las variables de entorno (API keys de DeepL, Gemini, TexTra, etc.) van en `.env`, que no se versiona.

---

<div align="center">
  <img src="public/assets/mozhi.png" width="192" height="192" alt="Mozhi logo">
  <h1>Mozhi</h1>
  <a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
    <img alt="License: AGPLv3" src="https://shields.io/badge/License-AGPL%20v3-blue.svg">
  </a>
  <a href="https://matrix.to/#/#mozhi:frei.chat">
  	<img alt="Matrix" src="https://img.shields.io/badge/matrix-000000?style=for-the-badge&logo=Matrix&logoColor=white">
  </a>
  <h3>Mozhi (spelt moḻi) is an alternative-frontend for many translation engines.</h3>
</div>

It was initially made as a maintained fork/rewrite of [simplytranslate](https://codeberg.org/SimpleWeb/SimplyTranslate-Web), but has grown to have a lot more features as well!

I'm initially focusing on the api and engines, but eventually Mozhi will have a functioning CLI and webapp.

## Supported Engines:

- Google
- Reverso
- DeepL
- LibreTranslate
- Yandex
- MyMemory
- DuckDuckGo ( 1-1 with Bing Translate )

## Projects that use Mozhi

- [select2translate](https://codeberg.org/aryak/select2translate) - Translate text from your selection clipboard using Mozhi
- [Crow Translate](https://invent.kde.org/office/crow-translate) - KDE Project written in C++ / Qt that allows you to translate and speak text using Mozhi

## Where is the engine code?

The engine code has recently been split from the main codebase. Please check [aryak/libmozhi](https://codeberg.org/aryak/libmozhi) for it.

## Installing

You can either use [docker](https://codeberg.org/aryak/mozhi/src/branch/master/compose.yml) or the build artifacts from [CI jobs on git.projectsegfau.lt](https://git.projectsegfau.lt/arya/mozhi/actions).

## Building

```
GOPRIVATE=codeberg.org/aryak/libmozhi # Get latest commit since proxy server is a bit slow
go mod download
go run github.com/swaggo/swag/cmd/swag@latest init --parseDependency
go build -o mozhi
```

## API Docs

Mozhi makes use of swagger (using the fiber middleware) to manage the documentation of the API.

You can find it in /api/swagger of any instance ([example](https://mozhi.aryak.me/api/swagger/index.html)).

## Why does Reverso not work?

Reverso sometimes blocks IPs of servers hosting mozhi, and since it doesn't have IPv6, an IP Rotator won't be viable. For more information, check out [#27](https://codeberg.org/aryak/mozhi/issues/27)

## Configuration

Features of Mozhi can be customized and toggled on/off using Environment Variables.

- `MOZHI_HOST`: Host address the webserver listens on (if hosting API). Defaults to listening on all interfaces
- `MOZHI_PORT`: Port the webserver listens on (if hosting API). Defaults to `3000`
- `MOZHI_LIBRETRANSLATE_URL`: URL of Libretranslate instance (Example: `MOZHI_LIBRETRANSLATE_URL=https://lt.psf.lt`)
- `MOZHI_DEFAULT_SOURCE_LANG`: Language to default to if no source language is set by user. Defaults to Auto-Detect (or first available language in engines which dont support it)
- `MOZHI_DEFAULT_PREFER_AUTODETECT`: Prefer autodetect if available instead of specified/default source language. Defaults to false
- `MOZHI_DEFAULT_TARGET_LANG`: Language to default to if no target language is set by user. Defaults to English
- `MOZHI_DEFAULT_ENGINE`: Engine to default to if no engine is set by user. Defaults to `google`

These envvars turn off/on engines. By default all of them are enabled.

- `MOZHI_GOOGLE_ENABLED`
- `MOZHI_REVERSO_ENABLED`
- `MOZHI_DEEPL_ENABLED`
- `MOZHI_LIBRETRANSLATE_ENABLED`
- `MOZHI_YANDEX_ENABLED`
- `MOZHI_MYMEMORY_ENABLED`
- `MOZHI_DUCKDUCKGO_ENABLED`

## Instances

| Link | Cloudflare | Country | ISP |
| --- | --- | --- | --- |
| [mozhi.aryak.me](https://mozhi.aryak.me) | No | India | Airtel |
| [translate.projectsegfau.lt](https://translate.projectsegfau.lt) | No | Germany / USA / India | Avoro / Racknerd / Airtel |
| [translate.nerdvpn.de](https://translate.nerdvpn.de) | No | Ukraine | vsys.host |
| [mozhi.ducks.party](https://mozhi.ducks.party) | No | Germany | Datalix |
| [mozhi.pussthecat.org](https://mozhi.pussthecat.org) | No | Germany | Hetzner |
| [mozhi.adminforge.de](https://mozhi.adminforge.de) | No | Germany | Hetzner |
| [translate.privacyredirect.com](https://translate.privacyredirect.com) | No | Finland | Private WebHost |
| [mozhi.canine.tools](https://mozhi.canine.tools) | No | USA | RoyaleHosting |
| [mzh.dc09.xyz](https://mzh.dc09.xyz) | No | Russia | Beget |
| [mozhi.frontendfriendly.xyz (Tor)](http://mozhi.wsuno6lnjdcsiok5mrxvl6e2bdex7nhsqqav6ux7tkwrqiqnulejfbyd.onion) | No | USA | Hetzner |
| [mozhi.ducks.party (Tor)](http://42i2bzogwkph3dvoo2bm6srskf7vvabsphw7uzftymbjjlzgfluhnmid.onion) | No | Germany | Datalix |
| [mozhi.r4fo.com](https://mozhi.r4fo.com) | No | Netherlands | Oracle |
| [mozhi.r4fo.com (Tor)](http://mozhi.r4focoma7gu2zdwwcjjad47ysxt634lg73sxmdbkdozanwqslho5ohyd.onion) | No | Netherlands | Oracle |
| [mozhi.bloat.cat](https://mozhi.bloat.cat) | No | Germany | Datalix |
| [mozhi.catsarch.com](https://mozhi.catsarch.com) | No | USA / Germany | Netcup |
| [mozhi.catsarch.com (Tor)](http://mozhi.catsarchywsyuss6jdxlypsw5dc7owd5u5tr6bujxb7o6xw2hipqehyd.onion) | No | USA / Germany | Netcup |
| [mozhi.catsarch.com (I2P)](http://b5jb6gilzl43u5js4d7jtcqmsk3xdjfbiowudij5yyhpm5bub3kq.b32.i2p) | No | USA / Germany | Netcup |

## Features

- An all mode where the responses of all supported engines will be shown.
- Autodetect which will show the language that was detected
- Text-To-Speech for multiple engines
- A good API (subjective :P)
- All the stuff you expect from a translation utility :)

## Etymology

Mozhi is the word in Tamil for language. Simple as that :P

## Credits

- [Arya](https://aryak.me): creator
- [Midou36o](https://midou.dev): made the logo
- [py_](https://github.com/supercolbat): Design files
- [Missuo](https://github.com/missuo): creating gDeepLX that does the hard part of making DeepL work
- [translatepy](https://github.com/Animenosekai/translate): giving me the format of request for yandex engine
- [SimplyTranslate](https://codeberg.org/simpleweb/simplytranslate): Inspiration and base code for the webui
- [Rimgo](https://codeberg.org/rimgo/rimgo): Code for embedding html in binary
- [Bnyro](https://me.chatoyer.de): Parallelization of all engines