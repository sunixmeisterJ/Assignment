# Spezifikation – Verleihprototyp

## Ziel
Ein lokales Verleihsystem verwaltet einzelne Artikel, Personen und zeitgebundene Buchungen. Die Oberfläche orientiert sich am vorgegebenen Dashboard. Das Fahrzeugprojekt verwendet eigenständige Beispieldaten und dieselben fachlichen Regeln.

## Fachliche Regeln und Abnahme
- Ein Artikel hat eine Bezeichnung, Details und einen positiven Stundensatz in Cent. Neue Artikel erscheinen sofort in Bestand und Buchungsformular.
- Eine Person hat Name, Ausweistyp und Ausweisnummer. Derselbe Ausweis darf nicht zweimal registriert werden.
- Eine Reservierung referenziert vorhandene Artikel und Personen. Ende liegt nach Beginn; maximal 366 Tage. Beginn darf wegen minutengenauer Eingabe höchstens eine Minute zurückliegen.
- Reservierungen desselben Artikels dürfen sich nicht überlappen. Direkt aufeinanderfolgende Intervalle sind zulässig. Auch gleichzeitige Anfragen dürfen nur eine überlappende Buchung erzeugen.
- Zustände: reserviert → ausgeliehen → zurückgegeben; alternativ reserviert → storniert. Andere Übergänge werden abgewiesen.
- Ausgabe ist nur im reservierten Zeitraum möglich. Ein noch ausgeliehener Artikel blockiert jede neue Reservierung und weitere Ausgabe.
- Je angefangener Stunde wird der gespeicherte Stundensatz berechnet. Bei Rückgabe wird ab reserviertem Beginn bis zum tatsächlichen Rückgabezeitpunkt abgerechnet. Preisänderungen beeinflussen bestehende Buchungen nicht.
- Das Dashboard zeigt aktuell verfügbare und ausgeliehene Artikel, Gesamtbestand und Personenanzahl.
- Neustart setzt den Zustand auf drei Beispielartikel und keine Personen/Buchungen zurück.

## Technische Grenzen
Go-Standardbibliothek; Structs und Methoden; kleines Interface für Preisregeln; Store mit Mutex; getrennte Domain-, Store- und HTTP-Schichten. HTML/JavaScript greift ausschließlich über die HTTP-API auf Daten zu. Keine persistenten Daten und keine Anmeldung. Ein Plugin-Framework aus Teil d) ist nicht enthalten.

## Änderbarkeit
Tarife werden über PriceRule erweitert, Stammdaten in store.New angepasst und Eingaben zentral im Store validiert. README beschreibt Start, API und Annahmen; CLAUDE.md hält die Entwicklungsregeln fest.
