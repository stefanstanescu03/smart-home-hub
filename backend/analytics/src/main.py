from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from dotenv import load_dotenv
import os
import engine

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
