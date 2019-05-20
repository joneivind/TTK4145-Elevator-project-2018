# TTK 4145 - Elevator project 2018
Heisprosjekt i faget TTK4145 Sanntidsprogrammering 2018 @ NTNU

### Når en heis mottar en lokal ordre
- Cab order taes av den lokale heisen, backup lagres på disk
- Hall order beregner ordrekost og ordren sendes til heisen med lavest kost
- Hvis en ordre feiler lokalt, redistribueres ordren via ny kostberegning
- Hvis en heis mister nettverk, vil de andre heisene distribuere ordrene til heisen

### Når en heis mottar en nettverksordre
- Hver heis lytter etter ordre som er assignet til seg selv
- Ordren sjekkes for duplikater og legges til i ordrelisten

Alle heiser blir oppdatert med status og ordreliste til de andre heisene

### Kostfunksjon
Ordrekostnad beregnes ut i fra status (ledig, opptatt, feil, etc.), retning, antall ordre og avstand.

### Knapper
| Knapp | Funksjonalitet |
| ------ | ------ |
| Hall 1-4 | Stopper på valgt etasje, fungerer ikke hvis heisen er alene i nettverket |
| Cab 1-4 | Stopper på valgt etasje, gjenopptar ved oppstart hvis ikke utført |
| Stopp | Stopper heisen i neste etasje |
| Obstruksjon | Aktiverer 'Easter egg' i terminalvinduet |

### For å kjøre programmet
Skriv følgende i terminalvinduet:
```sh
$ ElevatorServer  
```
```sh
$ chmod +x Nyan-elevator  
```
```sh
$ ./Nyan-elevator 
```
Programmet er skrevet i Go. Nettverksmodul og heisdriver er skrevet av [klasbo](https://github.com/klasbo)
