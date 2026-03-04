import pandas as pd
from datetime import datetime, timedelta


def last_rows(filename, num, param):

    df = pd.read_csv(filename)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
    last_data = df[['timestamp', param]].tail(int(num))

    return last_data.to_dict(orient='records')


def last_hours(filename, num, param):

    df = pd.read_csv(filename)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
    df = df.set_index('timestamp')
    df = df.sort_index()
    latest_time = df.index.max()
    start_time = latest_time - timedelta(hours=int(num))
    df_filtered = df.loc[start_time:latest_time]
    hourly_df = df_filtered[param].resample('1h').mean().reset_index()

    hourly_df['timestamp'] = hourly_df['timestamp'].dt.strftime(
        '%Y-%m-%d %H:00')
    return hourly_df.to_dict(orient='records')


def last_days(filename, num, param):

    df = pd.read_csv(filename)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
    df = df.set_index('timestamp')
    df = df.sort_index()
    latest_time = df.index.max()
    start_time = latest_time - timedelta(days=int(num))
    df_filtered = df.loc[start_time:latest_time]
    hourly_df = df_filtered[param].resample('D').mean().reset_index()

    hourly_df['timestamp'] = hourly_df['timestamp'].dt.strftime(
        '%Y-%m-%d %H:00')
    return hourly_df.to_dict(orient='records')
