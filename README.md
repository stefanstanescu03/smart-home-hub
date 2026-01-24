# smart-home-hub

Proiect pentru licenta
Momentan in lucru nu avem readme

# Run

```sh
.env with structure:
PORT=5000
TELEMETRY_PORT=5001
HOST=0.0.0.0
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=127.0.0.1
DB_PORT=3306
LOGS=
SECRET=
```

Anomaly Detection Pipeline: table in DB with model data (user, device, param, model_location, notify_email). After that everything works
like alerts: fetch current models, feed them with new data. If we create a new one we train it

# Tasks

1. Create database table [done]
2. Create the CRUD for models [done]
3. Figure out how to save models data and how to load them [done]
4. See how to handle user wanting to create new model [done]
5. Create the checking loop like you did with alerts
