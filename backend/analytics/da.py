import pandas as pd

df = pd.read_csv(
    "/home/stefan/Documents/Projects/licenta/data/Plant_1_Weather_Sensor_Data.csv")
df = df[["AMBIENT_TEMPERATURE", "DATE_TIME"]].copy()

df['timestamp_obj'] = pd.to_datetime(df['DATE_TIME'])

# 2. Convert to Unix Seconds (int)
# We divide by 10**9 because .astype(int) returns nanoseconds
df['timestamp'] = df['timestamp_obj'].astype('int64') // 1_000_000

df[["AMBIENT_TEMPERATURE", "timestamp"]].to_csv(
    'weather_sensor.csv', index=False)
