# Prüfprotokoll

Prüfung am 6. November 2026 mit Go 1.25.3 (macOS/ARM64).

- `gofmt`: alle Go-Dateien formatiert.
- `go test -race ./...`: bestanden; Preisrundung, Intervallgrenzen, atomare konkurrierende Reservierungen, Tarifspeicherung, Statuswechsel und API-Eingaben geprüft.
- `go vet ./...`: ohne Befund.
- Ski-Projekt im Browser: Dashboard und Bestandsansicht geladen; fiktive Person angelegt; Artikel reserviert, ausgegeben und zurückgegeben; Endpreis und Freigabe im Dashboard geprüft.
- Beide Projekte verwenden dieselben Webdateien bis auf Titel und Bestandsbeschriftung. Das Fahrzeugprojekt wurde durch eigene Go-Tests geprüft.
- ZIP-Integrität geprüft; Quellcode und Dokumentation enthalten, keine Compilerdateien, Testdaten aus der Browsersitzung oder Git-Metadaten.

Nicht geprüft: Flutter (nicht verwendet), externe Veröffentlichung (nicht erfolgt), produktiver Mehrbenutzerbetrieb.
