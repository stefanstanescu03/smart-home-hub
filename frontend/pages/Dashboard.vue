<script>
import SideBar from "../components/SideBar.vue";
import Card from "../components/Card.vue";
import axios from "axios";
export default {
  components: { SideBar, Card },
  data() {
    return {
      dashboard: "",
      widgets: [],
      available_devices: [],
      widget_type: "",
      device: "",
      ws: null,
    };
  },
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      if (token === undefined) {
        return undefined;
      }
      return token;
    },

    async handleFetchDashboard() {
      try {
        const res = await axios.get(
          `http://localhost:5000/dashboard/${this.$route.params.id}`,
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        );
        this.dashboard = res.data.dashboard;
      } catch (err) {
        console.error(err);
      }
    },
    async handleFetchWidgets() {
      try {
        await this.handleFetchDashboard();

        if (this.dashboard?.ID) {
          const res = await axios.get(
            `http://localhost:5000/widget/${this.dashboard.ID}`,
            {
              headers: { Authorization: `Bearer ${this.getToken()}` },
            }
          );
          this.widgets = res.data.widgets.map((widget) => ({
            ...widget,
            streamData: null,
          }));
        }
      } catch (err) {
        console.error(err);
      }
    },
    async handleGetAvailableDevices() {
      try {
        const res = await axios.get("http://localhost:5000/device/", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.available_devices = res.data.devices;
      } catch (err) {
        console.error(err);
      }
    },
    async handleAddWidget() {
      try {
        await axios.post(
          "http://localhost:5000/widget/create",
          {
            widgettype: this.widget_type,
            deviceId: this.device.ID,
            dashboardId: this.dashboard.ID,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        );
        this.exitDialog();
        this.$router.go(0);
      } catch (err) {
        console.error("Widget creation failed:", err);
      }
    },
    triggerDialog() {
      const dialog = document.querySelector("dialog");
      dialog.show();
    },
    exitDialog() {
      const dialog = document.querySelector("dialog");
      dialog.close();
    },

    setupWebsocket() {
      this.ws = new WebSocket("ws://localhost:5003/streaming");
      this.ws.onopen = () => {
        console.log("Connected to streaming server");
        this.widgets.forEach((widget) => {
          if (widget.DeviceId) {
            this.ws.send(JSON.stringify(`subscribe:${widget.DeviceId}`));
          }
        });
      };
      this.ws.onmessage = (event) => {
        const message = event.data;
        const parts = message.split(",");
        const widget = this.widgets.find(
          (widget) =>
            String(widget.DeviceId) ===
            String(parts[parts.length - 1].split(":")[1])
        );
        if (!widget) return;
        widget.dataStream = message;
      };
      this.ws.onclose = () => {
        console.log("disconnected");
      };
    },
  },
  async mounted() {
    await this.handleFetchWidgets();
    await this.handleGetAvailableDevices();
    this.setupWebsocket();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div class="info-container">
      <h1>{{ this.dashboard.Name }}</h1>
      <div class="widgets-container">
        <div v-for="widget in this.widgets">
          <Card
            v-if="widget.Widgettype == 'ca'"
            :deviceId="widget.DeviceId"
            :streamData="widget.dataStream"
            deviceName="default"
          />
        </div>
        <button class="add-button" @click="triggerDialog">+</button>
      </div>
      <dialog>
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Create a widget</p>
            <button @click="exitDialog" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="widget_type">Type: </label>
            <select name="widget_type" id="widget_type" v-model="widget_type">
              <option value="ca">Card</option>
              <option value="ch">Chart</option>
              <option value="ta">Table</option>
              <option value="al">Alarm</option>
            </select>
          </div>
          <div class="field">
            <label for="device">Device: </label>
            <select name="device" id="device" v-model="device">
              <option v-for="device in this.available_devices" :value="device">
                {{ device.Name }}
              </option>
            </select>
          </div>
          <button @click="handleAddWidget" class="create-button">Add</button>
        </div>
      </dialog>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
}

.info-container {
  width: 50%;
}

.add-button {
  border: none;
  width: 10rem;
  height: 10rem;
  color: #a6a6a6;
  font-weight: bold;
  font-size: 4rem;
  background-color: transparent;
  border: 1px solid #a6a6a6;
  border-style: dashed;
  transition-duration: 300ms;
}

.add-button:hover {
  border: 1px solid #121212;
  border-style: dashed;
  color: #121212;
}

.widgets-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  flex-wrap: wrap;
}

.field {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  gap: 0.5rem;
}

input {
  outline: none;
  box-shadow: none;
  border: none;
  border-bottom: 1px solid #a6a6a6;

  font-size: medium;
  padding: 0.3rem;
  transition-duration: 300ms;
}

input:focus {
  border-bottom: 1px solid #121212;
}

.dialog-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.top-dialog {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.create-button {
  border: none;
  text-decoration: none;
  cursor: pointer;
  background-color: #ff8441;
  color: #121212;
  transition-duration: 300ms;
  padding: 0.5rem;
  font-size: large;
  border-radius: 0.3rem;
}

.create-button:hover {
  background-color: #fe8d50;
}

.cancel-button {
  border: none;
  text-decoration: none;
  cursor: pointer;
  background-color: transparent;
}

dialog {
  border: 1px solid #a6a6a6;
  border-radius: 0.3rem;
}
</style>
