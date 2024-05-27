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

- `MOZHI_PORT`: Port the webserver listens on (if hosting API)
- `MOZHI_LIBRETRANSLATE_URL`: URL of Libretranslate instance (Example: `MOZHI_LIBRETRANSLATE_URL=https://lt.psf.lt`)
- `MOZHI_DEFAULT_SOURCE_LANG`: Language to default to if no source language is set by user. Defaults to Auto-Detect (or first available language in engines which dont support it)
- `MOZHI_DEFAULT_PREFER_AUTODETECT`: Prefer autodetect if available instead of specified/default source language. Defaults to false
- `MOZHI_DEFAULT_TARGET_LANG`: Language to default to if no target language is set by user. Defaults to English

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
| -------- | ---------- | ----------- | ----- |
| [mozhi.aryak.me](https://mozhi.aryak.me) | No | India | Airtel |
| [translate.bus-hit.me](https://translate.bus-hit.me) | No | Canada | Oracle |
| [nyc1.mz.ggtyler.dev](https://nyc1.mz.ggtyler.dev) | No | USA | Royale Hosting |
| [translate.projectsegfau.lt](https://translate.projectsegfau.lt) | No | Germany / USA / India | Avoro / Racknerd / Airtel |
| [translate.nerdvpn.de](https://translate.nerdvpn.de) | No | Ukraine | vsys.host |
| [mozhi.ducks.party](https://mozhi.ducks.party) | No | Germany | Datalix |
| [mozhi.ducks.party (Tor)](http://42i2bzogwkph3dvoo2bm6srskf7vvabsphw7uzftymbjjlzgfluhnmid.onion) | No | Germany | Datalix |
| [mozhi.pussthecat.org](https://mozhi.pussthecat.org) | No | Germany | Hetzner |

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
