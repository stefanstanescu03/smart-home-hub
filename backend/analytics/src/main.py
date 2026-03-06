from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from dotenv import load_dotenv
import os
import engine
import json

app = FastAPI()
load_dotenv("../.env")

origins = [
    "http://localhost:5173",
    "http://127.0.0.1:5173",
    "http://localhost:3000",
]


app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/")
async def root():
    return {"status": "online"}


@app.get("/tail/")
async def get_last_values(filename, num, param):
    csv_path = os.getenv("DATA_LOCATION") + "/" + filename
    values = engine.last_rows(csv_path, num, param)
    return {"values": values}


@app.get("/hour/")
async def get_per_hour(filename, num, param):
    csv_path = os.getenv("DATA_LOCATION") + "/" + filename
    values = engine.last_hours(csv_path, num, param)
    return {"values": values}


@app.get("/day/")
async def get_per_day(filename, num, param):
    csv_path = os.getenv("DATA_LOCATION") + "/" + filename
    values = engine.last_days(csv_path, num, param)
    return {"values": values}


@app.get("/predict/values/")
async def get_per_day(filename, lag, steps, param):
    csv_path = os.getenv("DATA_LOCATION") + "/" + filename
    y_pred = engine.predict_next_values(csv_path, lag, steps, param)
    return {"values": y_pred.tolist()}


@app.get("/predict/hours/")
async def get_per_day(filename, lag, steps, param):
    csv_path = os.getenv("DATA_LOCATION") + "/" + filename
    y_pred = engine.predict_next_hours(csv_path, lag, steps, param)
    return {"values": y_pred.tolist()}


@app.get("/predict/days/")
async def get_per_day(filename, lag, steps, param):
    csv_path = os.getenv("DATA_LOCATION") + "/" + filename
    y_pred = engine.predict_next_days(csv_path, lag, steps, param)
    return {"values": y_pred.tolist()}
