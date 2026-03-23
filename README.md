# ees-link

## Projekta apraksts


## Nepieciešamā programmatūra

Apache Ant

Docker Desktop ar WSL 2.0

Go 1.23.4 vai jaunāka

VS Code ar Code Runner un Go paplašinājumiem (lokālai palaišanai) 

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

## Sagaidāmā darbība

Palaižot projektu, ...

