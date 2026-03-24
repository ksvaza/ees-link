# ees-link

## Projekta apraksts

**ees-link** (E-ES link) ir Energoefektivitātes sacensību (E-ES) oficiālā digitālā platforma - pilnas kaudzes tīmekļa lietotne, kas nodrošina centralizētu piekļuvi visai sacensību informācijai gan dalībniekiem, gan organizatoriem, gan skatītājiem.

Platforma piedāvā publisku sākumlapu ar sacensību nolikumu, jaunumiem, komandu sarakstiem un iepriekšējo gadu arhīvu. Sacensību laikā rezultāti tiek attēloti reāllaikā pa posmiem (kvalifikācija, efektivitāte, tehniskais dizains, autobūves zināšanas, fināls). Papildus sacensībām platforma ļauj jebkurai komandai, kas izmanto organizatoru barošanas un uzskaites sistēmu, pieslēgt savu transportlīdzekli datu nolasīšanai un enerģijas patēriņa analīzei arī ārpus sacensībām.

## Komanda

| Vārds Uzvārds | Loma |
|---|---|
| Adrians Dacko | Aizmugurgals |
| Kārlis Svaža | Aizmugurgals |
| Lauris Začests | Priekšgals |

## Tehnoloģiju kaudze

### Aizmugurgals
- **Go** — galvenā programmēšanas valoda
- **httprouter** — HTTP maršrutēšana
- **logrus** — loģēšana
- **viper** — vides mainīgo pārvaldīšana
- **PostgreSQL** — datu glabāšana (gan parastās, gan hronoloģiskās tabulas)
- **Eclipse Mosquitto (MQTT)** — divvirzienu komunikācija ar sacensību elektroautomašīnām (dati JSON formātā)
- **WebSocket** — reāllaika rezultātu pārraide bez lapas pārlādes
- **Goroutines** — asinhrona, daudzpavedienu apstrāde veiktspējas uzlabošanai

### Priekšgals
- **React** — lietotāja saskarnes satvars
- **TypeScript** — programmēšanas valoda
- **HTML / CSS** — struktūra un stils
- **Vite** — priekšgala būvēšanas rīks

### Infrastruktūra
- **Docker** (Alpine Linux bāze) — aizmugurgala un priekšgala apkopošana vienā konteinerī
- **SWAG** — Nginx, Certbot un Fail2Ban apvienots Docker rīks
- **Nginx** — apgrieztais starpnieks (HTTP → HTTPS novirzīšana)
- **Let's Encrypt + Certbot** — HTTPS sertifikāti
- **Apache Ant** — būvēšanas automatizācija

### Izstrādes rīki
- **Git / GitHub** — versiju kontrole
- **VS Code** — koda redaktors
- **Swagger** — API galapunktu dokumentēšana
- **Jira** — uzdevumu pārvaldība
- **Autentifikators** — lietotāju autentifikācija (paroles šifrētas ar SHA-512)

## Nepieciešamā programmatūra

- Apache Ant
- Docker Desktop ar WSL 2.0
- Go 1.23.4 vai jaunāka
- VS Code ar Code Runner un Go paplašinājumiem (lokālai palaišanai)

## Sistēmas arhitektūra

```
                          ┌─────────────────────────────────────────┐
                          │           DOCKER (Alpine Linux)         │
                          │                                         │
  [React priekšgals] ◄──► │  HTTPS / Nginx ◄──► REST API            │
                          │               ◄──► WebSocket            │
                          │                                         │
                          │  HTTP aizmugurgals (Go)                 │
                          │    ├── Autentifikators (SHA-512)        │
                          │    ├── DB apstrādātājs ──► PostgreSQL   │
                          │    └── MQTT apstrādātājs                │
                          └────────────────┬────────────────────────┘
                                           │ MQTT (Eclipse Mosquitto)
                                    ┌──────┴──────┐
                              [Mašīna 1] [Mašīna 2] [Mašīna 3]
```

Priekšgals komunicē ar aizmugurgalu caur REST API (parasti) un WebSocket (reāllaika dati). MQTT nodrošina divvirzienu savienojumu ar sacensību elektroautomašīnām.

## Vides mainīgo iestatīšana pirms būvēšanas

Ja projekts tiks būvēts un palaists Docker konteinerā, tad vides mainīgie tiek iestatīti `dockercompose/docker-compose.yml` failā.

Laižot projektu lokāli, nepieciešams saknē izveidot .env failu ar vides mainīgajiem: LOGFILE un CLEARLOG.

Piemērs `.env` failam ar noklusējuma vērtībām:

```
LOGFILE=./logs/log.txt
CLEARLOG=true
```

## Palaišanas instrukcija

### Lokāla palaišana (VS Code)
1. Atveriet projektu VS Code.
2. Atveriet `main.go` failu.
3. Noklikšķiniet uz "Run Code" ikonas, kas parādās virs `main` funkcijas, vai izmantojiet īsinājumtaustiņu `Ctrl+Alt+N` (Windows/Linux) vai `Cmd+Option+N` (Mac).
4. Varētu būt nepieciešams apstiprināt Windows ugunsmūra dialoglodziņu un atļaut piekļuvi tīklam.

### Dokerizēta palaišana
1. Atveriet termināli projekta saknē.
2. Pārliecinieties, ka Docker Desktop ir palaists un darbojas.
3. Izpildiet komandu `ant win`, lai būvētu Docker konteineri. Pirms būvēšanas automātiski tiks palaisti testi, taču tos var manuāli palaist arī ar `ant test`.
4. Pēc būvēšanas izpildiet komandu `docker compose up`, lai palaistu konteineri. Pirmajā palaišanas reizē Docker lejuplādēs nepieciešamos attēlus un izveidos konteineri, kas var aizņemt kādu laiku.

## Esošā funkcionalitāte

Pašlaik ees-link ir Go aizmugurgals ar pamata HTTP serveri. Palaižot projektu, tas:

1. Nolasa vides konfigurāciju no `.env` faila
2. Inicializē loģēšanas sistēmu ar norādīto log failu
3. Uzsāk HTTP serveri uz porta **1884**
4. Servē statiskos failus no `public/` direktorijas
5. Piedāvā šādus API galapunktus:
   - `GET /api/test` — testa galapunkts, atgriež HTML atbildi
   - `POST /api/submit-form` — formas datu pieņemšana, atgriež JSON
6. Pieraksta visas operācijas log failā ar laika zīmogiem un līmeņiem (`info`, `warning`, `error`)

Lokāli laižot, priekšgalam var piekļūt ar `http://localhost:1884`.

Dokerizēti laižot, priekšgalam var piekļūt ar `http://localhost:6767`

## Sagaidāmā darbība

Palaižot projektu, vajadzētu redzēt log failā ierakstus par servera startu un API izsaukumiem. Piekļūstot `http://localhost:1884/api/test`, vajadzētu redzēt HTML atbildi, un mājaslapā izpildot formu, vajadzētu redzēt JSON atbildi ar nosūtītajiem datiem. Visi šie notikumi tiks ierakstīti log failā ar atbilstošiem līmeņiem.
