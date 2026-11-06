# Ski-Verleih – Aufgabe 1 a) und b)

Lauffähiger Prototyp mit Go-Backend und HTML/JavaScript-Oberfläche. Keine externen Abhängigkeiten, keine Datenbankinstallation. Voraussetzung: Go 1.22 oder neuer.

## Starten

```sh
go run .
```

Im Browser http://127.0.0.1:8080 öffnen. Beenden: Strg+C. Ein anderer Port ist mit `ADDR=127.0.0.1:8081 go run .` möglich, etwa um beide Projekte gleichzeitig zu starten.

```sh
go test ./...
go test -race ./...
go vet ./...
gofmt -w .
```

## Ausprobieren

1. Drei Beispielartikel sind vorhanden; Preis und Details können bearbeitet werden.
2. Eine Person mit fiktivem Namen, Ausweistyp und Ausweisnummer anlegen.
3. Unter Reservierungen Artikel, Person, Beginn und Ende wählen. Die Vorgabe beginnt in der aktuellen Minute.
4. Reservieren und anschließend „Ausgeben“ wählen. Ausgabe ist nur innerhalb des Buchungszeitraums möglich.
5. „Zurückgeben“ beendet die Ausleihe und berechnet die tatsächliche Dauer seit dem reservierten Beginn. Reservierungen lassen sich vor der Ausgabe stornieren.
6. Eine zweite zeitlich überlappende Buchung desselben Artikels wird abgelehnt. Direkt anschließende Reservierungen sind möglich.

## Umsetzung der Aufgabenstellung

- a): Übersicht im Stil des Referenzbildes, Bestand, Personen mit Ausweisdaten, Reservierung/Ausleihe/Rückgabe, benutzerfreundliche Fehlermeldungen.
- b): individuell änderbare Stundensätze je Artikel; `PriceRule` als kleines Interface für alternative Preismodelle. Jede angefangene Stunde zählt. Geldbeträge werden in ganzen Cent gerechnet.
- c) wird separat im Projekt fahrzeug-verleih abgegeben.

## Annahmen und Grenzen

Die Details der Vorlesungsimplementierung lagen nicht vor. Stundenweise Abrechnung und sinnvolle Prototyp-Annahmen wurden vom Auftraggeber bestätigt. Der Tafelhinweis auf ein Verleih-Framework/Plugin-System (d) gehört nicht zu diesen beiden Abgaben.

- Ein Artikel ist genau ein physisches Exemplar; keine Mengenbuchung.
- Preise gelten pro Artikel und angefangener Stunde, ohne Rabatte, Kaution, Steuern oder Tagesdeckel. Maximale Reservierungsdauer: 366 Tage.
- Ein bei Reservierung gespeicherter Stundensatz bleibt für diese Buchung unverändert. Geplante Preise stehen bis zur Rückgabe in der Übersicht.
- Reservierungsintervalle sind halboffen: Ende einer Buchung darf Beginn der nächsten sein. Eine laufende, auch verspätete Ausleihe blockiert die Ausgabe und neue Reservierungen. Bestehende Folgereservierungen bleiben bestehen; nach der Rückgabe ist ihre Ausgabe wieder möglich.
- Kein Ausweis-Upload: nur Nummer und Typ als Text, wie auf der Tafel angegeben.
- Thread-sicherer In-Memory-Store mit Mutex; alle Daten verschwinden beim Neustart. Keine Anmeldung, kein Mehrmandantenbetrieb. Standardmäßig nur lokal erreichbar.
- Personen/Artikel werden hinzugefügt, Artikel auch bearbeitet. Löschen und Bearbeiten von Personen sind nicht Teil dieses Prototyps.
- Keine echten Personen- oder Ausweisdaten verwenden. Die ZIP-Dateien enthalten keine Git-Metadaten.
- Der Hinweis „public Repo / GitHub-Link“ auf der Tafel wird durch die vom Auftraggeber verlangte ZIP-Abgabe ersetzt. Es wurde nichts auf GitHub veröffentlicht.

## Aufbau

`main.go`: Start und eingebettete Webdateien; `internal/domain`: Datenmodelle und Preisregel; `internal/store`: atomare In-Memory-Operationen; `internal/handler`: JSON-API; `web`: Oberfläche.

Die Oberfläche ruft ausschließlich die Go-Endpunkte auf. Benutzertexte werden als Textknoten eingefügt. Der Store wird im Einstieg erzeugt und an Handler übergeben; kein globaler Datenstore.

## API

- `GET /api/state`: aktueller Bestand, Personen, Buchungen.
- `POST /api/items`: Artikel anlegen (`id: 0`) oder ändern (`id` vorhanden); `name`, `details`, `hourlyCents`.
- `POST /api/people`: `name`, `documentType`, `documentNumber`.
- `POST /api/rentals`: `itemId`, `personId`, `start`, `end` (RFC3339-Zeitstempel).
- `POST /api/rentals/{id}/checkout`, `/return`, `/cancel`: Statuswechsel.

Ungültige Eingaben ergeben HTTP 400 mit einem `error`-Text. Tests decken Preisrundung, Grenzen von Zeitintervallen, konkurrierende Reservierungen, Statuswechsel, gespeicherte Tarife und API-Validierung ab.
