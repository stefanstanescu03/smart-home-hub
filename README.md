# smart-home-hub

Proiect pentru licenta
Momentan in lucru nu avem readme

# TODO

- Finish base CRUD backend
- Try to test it

0/1 - send/receive
1 byte - data
4 bytes - value

Protocol: - Primul byte - clasa (ambiental, de pozitie, electrici, boolean) - Al doilea byte - restul

0001 0001 - temperatura
0001 0010 - umiditate
0001 0011 - presiune
0001 0100 - intensitate luminoasa
0001 0101 - AQI
0001 0110 - CO2
0001 0111 - VOC
0001 1000 - PM2.5
0001 1001 - Nivel UV
0010 0001 - GPS latitudine
0010 0010 - GPS longitudine
0010 0011 - GPS altitudine
0010 0100 - proximitate
0011 0001 - voltaj
0011 0010 - curent
0011 0011 - putere
0011 0100 - nivel baterie
0100 0001 - valoare booleana (usa inchisa/deschisa, bec inchis/deschis)

Restul user defined
