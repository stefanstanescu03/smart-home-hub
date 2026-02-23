# smart-home-hub

Proiect pentru licenta
Momentan in lucru nu avem readme

# Run

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
LOGS=/home/stefan/Documents/Projects/licenta/logs
SMTP_EMAIL=st.stefan24@gmail.com
EMAIL_PASSWORD=
SMTP=smtp.gmail.com
SMTP_PORT=587
```

# Milestones

1. Cloud
2. Automations
3. Forecasting

# References

1. Gin - https://gin-gonic.com/en/docs/
2. GORM - https://gorm.io/docs/
3. RESTful API tutorial - https://go.dev/doc/tutorial/web-service-gin
4. JWT - https://pkg.go.dev/github.com/golang-jwt/jwt/v5
5. MQTT Broker - https://github.com/mochi-mqtt/server
6. MQTT Client - https://github.com/eclipse-paho/paho.mqtt.golang
7. Websockets - https://pkg.go.dev/github.com/gorilla/websocket
8. Anomaly detection:
   - https://engineering.yale.edu/application/files/3817/3714/7395/tr222.pdf (Calcul std online)
   - https://vldb.org/pvldb/vol14/p141-tran.pdf (Idee de baza)
9. Forecasting:
   - https://www.sciencedirect.com/science/article/pii/S1877050920310267?ref=pdf_download&fr=RR-2&rr=9d1f898af91bc9d1
