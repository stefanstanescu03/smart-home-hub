# Mica descriere

Proiectul este un sistem destinat unificarii si controlului dispozitivelor inteligente. Isi propune sa centralizeze datele colectate de la diversi senzori si pe baza lor sa poata efectua schimburi de date cu alte device-uri + analiza predictiva a datelor, detectie de anomali.
Sistemul este realizat in go (cu framework-ul GIN) si python (cu FastAPI) iar interfata web este realizata cu ajutorul framework-ului Vue.js.

# Comunicare

Pentru a facilita comunicarea dintre server si dispozitive am folosit protocolul MQTT. Ideea de baza este de a folosi un broker care colecteaza informatiile de la clienti, informatii publicate in topicurile aferente, si pe baza topicurilor unde au fost publicate el redirectioneaza datele.

1. Brokerul este implementat cu ajutorul bibliotecii mochi-mqtt
2. Fiecare device cand se conecteaza isi trimite id-ul (aici numit ident de la identificator). Broker-ul verifica daca acesta se afla in baza de date (daca dispozitivul este inregistrat) iar in caz contrar rejecteaza conexiunea. Singurul ident rezervat este SYSTEM_CMD_CLIENT care mereu va fi acceptat (revenim la el mai tarziu)
3. La fiecare publish verificam daca mesajul a fost publicat in topic-ul care trebuie: telemetry/(ident) sau stat/(ident) sau daca device-ul incearca sa dea subscribe la un topic la care are acces: cmd/(ident) care e singurul
4. Cand se publica date la topicul telemetry/ acestea sunt adaugate intr-un fisier csv aferent

## Topicuri

1. telemetry/ - aici dispozitivele trimit informatii, telemetrie
2. stat/ - dispozitivele pot trimite informatii referitoare la starea lor actuala (true sau false l-am gandit mai degraba pentru switch-uri)
3. cmd/ - topicul la care se trimit comenzi catre alte dispozitive

## CMDClient

Client de mqtt care are rolul de intermediar intre utilizator si restul dispozitivelor, implementat cu paho-mqtt. E responsabil cu transmiterea comenzilor catre celelalte dispozitive si mentinerea unei evidente a statusurilor acestora: publica la topicul cmd/ si este subscriber la stat/

Brokerul tine evidenta in memorie a dispozitivelor conectate pentru a nu verifica mereu daca acestea sunt inregistrate in baza de date.

# Transmiterea de date live

Pentru a transmite datele conectate in timp real am apelat la WebSockets. Practic Backend-ul citeste la intervale de timp ultima valoare dintr-un fisier csv si o transmite catre frontend.
In backend, server-ul de streaming gestioneaza conexiunile persistente ale clientilor si distribuie datele catre acestia doar pentru dispozitivele la care acestia au dat subscribe: clientul initial trimite un mesaj de tip subscribe:ID pentru a stii datele de la ce dispozitiv trebuiesc trimise (ID e id-ul dispozitivului din baza de date nu ident). Dupa validare, se retine si locatia datelor device-ului. Dupa care la intervale regulate se transmit valorile dorite sub forma: (param):(val),...,id:(ID). Avem nevoie de id deoarece un client poate fi subscribed la mai multe device-uri cand un dashboard are mai multe widget-uri

# Alerte

Sistemul este capabil sa trimita alerte prin email. Utilizatorul poate crea mai multe alerte pentru diverse dispozitive: poate seta subiectul, content-ul email-ului care urmeaza sa fie trimis, device-ul, parametrul si o conditie. Alertele sunt incarcate in memorie din baza de date si la intervale regulate se verifica daca conditiile de declansare sunt atinse (pentru asta se parseaza conditiile din memorie care sunt de forma "(parametru)>(valoare)"). Se trimit emai-uri doar pentru alertele care permit asta. Intotdeaua alertele vor fi inregistrate in fisierul alerts.log. Pentru a da refresh la buffer-ul de alerte se apeleaza la un semnal declansat in momentul in care se sterge sau se creeaza o alerta noua. Dupa ce s-a trimis o alerta nu se mai pot trimite altele pana cand conditia de declansare devine invalida.

## Cum se trimit emai-urile

Email-urile se trimit prin SMTP (Simple Mail Transfer Protocol). Acesta are nevoie totusi de un server insa il putem folosi foarte usor pe cel de la google.

# Automatizari

Mecanismul este foarte asemanator cu cel al alertelor doar ca de data aceasta in loc sa trimitem un email trimitem o comanda prin intermediul CMDClient.

# Detectia de anomali (Posibil sa-l schimb)

Pentru aceasta am folosit o metoda de tip z-score adaptata pentru date care vin continuu. Deoarece lucrez cu ferestre de date algoritmul este mai degraba o cluster-izare doar ca intr-un singur cluster iar datele cele mai indepartate se considera a fi anomali. In plus centroidul se poate adapta la noile date (datele noi care nu sunt anomali modifica putin centroidul pentru a se adapta schimbarilor lente in timp). Algoritmul este integral implementat in go. Parametrii modelelor sunt retinuti in fisiere json iar metadatele sunt retinute in baza de date. Dupa ce am creat un model il antrenam pe datele deja existente, dupa care el functioneaza printr-un mecanism similar celui de alerte: cand detectam o anomalie trimitem un emai (daca este permis) si scrie alerta in alerts.log.

# Predictia datelor

Se face cu ajutorul algoritmului Random forest. Pentru antrenare se creaza urmatoarele trasaturi din date: ora, ziua saptamanii, luna, daca este weekend, un numar de valori din trecut si un numar de valori din viitor (care constituie si vectorul pe care trebuie sa-l prezicem). Deoarece modelul poate considera ora 23 si 0 ca fiind indepartate, cand de fapt ele sunt doar la o ora distanta, se pot aplica functiile sinus si cosinus pe acestea (cu o frecventa de 1/24) pentru a compensa acest fapt. Acelasi lucru si pentru zile deoarece dupa 7 (duminica) vine 1 (luni) si am avea aceeasi problema. Putem face predictii pe datele viitoare, pe ore si pe zile.

# Ce ar trebui sa contine .env

```sh
.env with structure:
PORT=5000
TELEMETRY_PORT=5001
STREAMING_PORT=5003
HOST=

DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=

SECRET=

LOGS=
DATA_LOCATION=
ANOMALY_MODELS_LOCATION=

SMTP_EMAIL=
EMAIL_PASSWORD=
SMTP=smtp.gmail.com
SMTP_PORT=587

MQTT_CERT_FILE=
MQTT_KEY_FILE=
MQTT_CA_FILE=


```
