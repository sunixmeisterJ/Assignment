# Projektregeln

Prototyp: einfach und verständlich halten. Go-Standardbibliothek, keine zusätzlichen Abhängigkeiten ohne Rückfrage. Domain-Logik in Struct-Methoden, thread-sicherer Store, getrennte HTTP-Handler und HTML/JS-Oberfläche. Keine echten Ausweisdaten oder Secrets speichern.

Start: `go run .`; Tests: `go test ./...`; Nebenläufigkeit: `go test -race ./...`; Format: `gofmt -w .`; Prüfung: `go vet ./...`.

Geld in Cent; angefangene Stunden aufrunden; bei Buchung Tarifsatz einfrieren. Verfügbarkeit und Buchung atomar unter dem Store-Mutex prüfen. Frontend kommuniziert nur über `/api/` mit dem Backend. Keine generierten Dateien von Hand ändern. Commit-Nachrichten kurz und imperativ.
