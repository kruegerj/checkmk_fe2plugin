# checkmk_fe2plugin
FE2 Monitoring Plugin für checkmk


Dieses muss auf den FE2 Server unter %programdata%/checkmk/agent/local kopiert werden und in der config.yaml muss ein Authorization Token gesetzt werden.

```yaml
hostname: 127.0.0.1
port: 83
protocol: http
token: <TOKEN aus dem Monitoring Plugin>
```
Es kann entweder jeder Endpunkt separat oder alle akutell (Dezember 2023) dokumentierten Endpunkte überwacht werden.
die _all Datei enthält alle Endpunkte aus der Dokumentation 
https://alamos-support.atlassian.net/wiki/spaces/documentation/pages/1683226637/Monitoring-Schnittstelle. Ansonsten sollten die Namen entsprechende Rückschlüsse zulassen.
