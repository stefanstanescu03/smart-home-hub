import pandas as pd
from datetime import datetime, timedelta
import numpy as np
from sklearn.ensemble import RandomForestRegressor


def create_features(df, column, lag, steps):

    X_cols = ['hour_sin', 'hour_cos', 'day_sin',
              'day_cos', 'month', 'is_weekend']
    y_cols = []

    df['hour'] = df['timestamp'].dt.hour
    df['day_of_week'] = df['timestamp'].dt.dayofweek
    df['month'] = df['timestamp'].dt.month
    df['is_weekend'] = (df['timestamp'].dt.dayofweek >= 5).astype(int)

    df['hour_sin'] = np.sin(2 * np.pi * df['hour'] / 24)
    df['hour_cos'] = np.cos(2 * np.pi * df['hour'] / 24)

    df['day_sin'] = np.sin(2 * np.pi * df['day_of_week'] / 7)
    df['day_cos'] = np.cos(2 * np.pi * df['day_of_week'] / 7)

    new_cols = {}

    for i in range(1, lag + 1):
        label = "lag_" + str(i)
        new_cols[label] = df[column].shift(i)
        X_cols.append(label)

    for i in range(0, steps):
        label = "step_" + str(i)
        new_cols[label] = df[column].shift(-i)
        y_cols.append(label)

    df = pd.concat([df, pd.DataFrame(new_cols)], axis=1)

    X_forecast = df[X_cols].tail(1).copy()

    df_clean = df.dropna(subset=X_cols + y_cols).copy()

    X = df_clean[X_cols]
    y = df_clean[y_cols]

    return X, y, X_forecast


def forecast(X_train, y_train, X_forecast):

    model = RandomForestRegressor(
        n_estimators=100, max_depth=10, random_state=42)
    model.fit(X_train, y_train)

    return model.predict(X_forecast)


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


def predict(df, lag, steps, param):

    X, y, X_forecast = create_features(df, param, int(lag), int(steps))

    X = X.apply(pd.to_numeric, errors='coerce')
    y = y.apply(pd.to_numeric, errors='coerce')

    X = X.fillna(0)
    y = y.fillna(0)

    X = X.to_numpy(dtype=np.float32)
    y = y.to_numpy(dtype=np.float32)
    X_forecast = X_forecast.to_numpy(dtype=np.float32)

    y_pred = forecast(X, y, X_forecast)

    y_pred = y_pred.flatten()

    return y_pred


def predict_next_values(filename, lag, steps, param):

    df = pd.read_csv(filename)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')

    return predict(df, lag, steps, param)


def predict_next_hours(filename, lag, steps, param):

    df = pd.read_csv(filename)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
    df = df.set_index('timestamp').sort_index()
    hourly_df = df[param].resample('1h').mean()
    hourly_df = hourly_df.reset_index()
    hourly_df['timestamp'] = pd.to_datetime(hourly_df['timestamp'])
    hourly_df = hourly_df.dropna()

    return predict(hourly_df, lag, steps, param)


def predict_next_days(filename, lag, steps, param):

    df = pd.read_csv(filename)
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
    df = df.set_index('timestamp').sort_index()
    day_df = df[param].resample('D').mean().reset_index()
    day_df = day_df.reset_index()
    day_df['timestamp'] = pd.to_datetime(day_df['timestamp'])
    day_df = day_df.dropna()

    return predict(day_df, lag, steps, param)
