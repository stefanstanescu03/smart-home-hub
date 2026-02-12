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

Anomaly Detection Pipeline: table in DB with model data (user, device, param, model_location, notify_email). After that everything works
like alerts: fetch current models, feed them with new data. If we create a new one we train it

# Tasks

1. Make the forms look a bit better
2. Add validations for forms
