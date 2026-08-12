> **Fork notice:** versión modificada de [aryak/mozhi](https://codeberg.org/aryak/mozhi),
> mantenida por [darlopvil](https://github.com/darlopvil). Cambios y roadmap en los
> [issues](https://github.com/darlopvil/mozhi-fork/issues). Licencia AGPLv3, igual que upstream.


## Build local (self-hosted)

Este fork se construye localmente, sin registry. Estructura esperada:

    ~/mozhi/      → este repo (código + docker-compose.yml + .env, no versionados)
    ~/libmozhi/   → fork de libmozhi (el motor); enganchado vía `replace` en go.mod

Para reconstruir tras un cambio y recrear el contenedor:

    ./rebuild.sh

Requiere que `docker-compose.yml` use `image: mozhi-fork:local`.
Las variables de entorno (API keys de DeepL, Gemini, TexTra, etc.) van en `.env`, que no se versiona.

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
