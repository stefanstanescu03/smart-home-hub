import paho.mqtt.client as mqtt
import threading
import time
import random
import json

BROKER = "localhost"
PORT = 5001

devices = [
    {"ident": "AQ923", "pub": "telemetry/AQ923", "sub": None},
    {"ident": "POWER002", "pub": "telemetry/POWER002", "sub": None},
    {"ident": "PIR001", "pub": "telemetry/PIR001", "sub": None},
    {"ident": "RELAY001", "pub": "stat/RELAY001", "sub": "cmd/RELAY001"},
    {"ident": "RELAY002", "pub": "stat/RELAY002", "sub": "cmd/RELAY002"}
]


def on_connect(client, userdata, flags, rc, properties=None):
    if rc == 0:
        print(f"Connected: {userdata['ident']}")
        if userdata.get("topic_sub"):
            client.subscribe(userdata["topic_sub"])
            print(f"Subscribed to {userdata['topic_sub']}")
    else:
        print(f"Failed to connect {userdata['ident']}, code {rc}")


def on_message(client, userdata, msg):
    payload = msg.payload.decode()
    print(f"\n[{userdata['ident']}] Received Command on {msg.topic}: {payload}")
    if (userdata['ident'] == 'RELAY001'):
        if payload == "ON_RELAY1":
            client.publish('stat/RELAY001/stat_relay1', 'true')
        if payload == "OFF_RELAY1":
            client.publish('stat/RELAY001/stat_relay1', 'false')
        if payload == "ON_RELAY2":
            client.publish('stat/RELAY001/stat_relay2', 'true')
        if payload == "OFF_RELAY2":
            client.publish('stat/RELAY001/stat_relay2', 'false')
    if (userdata['ident'] == 'RELAY002'):
        if payload == "ON_RELAY":
            client.publish('stat/RELAY002', 'true')
        if payload == "OFF_RELAY":
            client.publish('stat/RELAY002', 'false')


def create_device(ident, topic_pub=None, topic_sub=None, interval=5):

    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id=ident,
        userdata={"ident": ident, "topic_sub": topic_sub}
    )
    client.on_connect = on_connect
    client.on_message = on_message

    client.connect(BROKER, PORT, 60)
    client.loop_start()

    try:
        while True:
            if topic_pub:
                if ident == "AQ923":
                    temp = random.uniform(15.0, 25.0)
                    humidity = random.randint(40, 60)
                    co2 = random.randint(400, 1000)
                    data = [
                        {"n": "temperature", "u": "Cel", "v": temp},
                        {"n": "humidity",    "u": "%RH", "v": humidity},
                        {"n": "co2",         "u": "ppm", "v": co2}
                    ]
                    payload = json.dumps(data)
                    client.publish(topic_pub, payload)
                if ident == "POWER002":
                    current_w = random.uniform(
                        150.0, 160.0) if random.random() > 0.2 else 5.0
                    voltage = random.uniform(228.0, 232.0)
                    amps = current_w / voltage
                    data = [
                        {"n": "power", "u": "W", "v": current_w},
                        {"n": "voltage", "u": "V", "v": voltage},
                        {"n": "current", "u": "A", "v": amps}
                    ]
                    payload = json.dumps(data)
                    client.publish(topic_pub, payload)
                if ident == "PIR001":
                    is_somebody = random.choice([0, 1])
                    data = [{"n": "status", "v": is_somebody}]
                    payload = json.dumps(data)
                    client.publish(topic_pub, payload)

            time.sleep(interval)
    except KeyboardInterrupt:
        client.disconnect()


if __name__ == "__main__":

    threads = []
    for d in devices:
        t = threading.Thread(target=create_device, args=(
            d["ident"], d["pub"], d["sub"], 5))
        t.daemon = True
        threads.append(t)
        t.start()

    print("--- Virtual Home Devices Running ---")
    while True:
        time.sleep(1)
