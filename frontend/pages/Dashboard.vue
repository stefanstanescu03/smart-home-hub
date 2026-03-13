<script>
import SideBar from "../components/SideBar.vue";
import Card from "../components/Card.vue";
import Table from "../components/Table.vue";
import axios from "axios";
import WidgetButton from "../components/WidgetButton.vue";
import WidgetSwitch from "../components/WidgetSwitch.vue";
export default {
  components: { SideBar, Card, Table, WidgetButton, WidgetSwitch },
  data() {
    return {
      dashboard: "",
      widgets: [],
      available_devices: [],
      new_widget: {
        widget_type: "",
        device: "",
        label: "",
        payload: "",
        payload2: "",
        channel: "",
      },
      ws: null,
      menuOpen: false,
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
        const res = await axios.get(`/api/dashboard/${this.$route.params.id}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.dashboard = res.data.dashboard;
      } catch (err) {
        console.error(err);
      }
    },
    async handleFetchWidgets() {
      try {
        await this.handleFetchDashboard();

        if (this.dashboard?.ID) {
          const res = await axios.get(`/api/widget/${this.dashboard.ID}`, {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          });
          const widgetsWithDeviceNames = await Promise.all(
            res.data.widgets.map(async (widget) => ({
              ...widget,
              dataStream: null,
              deviceName: await this.handleGetDeviceName(widget.DeviceId),
            })),
          );

          this.widgets = widgetsWithDeviceNames;
        }
      } catch (err) {
        console.error(err);
      }
    },
    async handleDeleteWidget(id) {
      try {
        await axios.delete(`/api/widget/delete/${id}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.widgets = this.widgets.filter((widget) => widget.ID != id);
      } catch (err) {
        console.log(err);
      }
    },
    async handleGetDeviceName(id) {
      try {
        const res = await axios.get(`/api/device/${id}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        return res.data.device.Name;
      } catch (err) {
        return "default";
      }
    },
    async handleGetAvailableDevices() {
      try {
        const res = await axios.get("/api/device/", {
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
          "/api/widget/create",
          {
            widgettype: this.new_widget.widget_type,
            deviceId: this.new_widget.device.ID,
            dashboardId: this.dashboard.ID,
            label: this.new_widget.label,
            payload: this.new_widget.payload,
            payload2: this.new_widget.payload2,
            channel: this.new_widget.channel,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          },
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
      this.ws = new WebSocket(`ws://${location.hostname}:5003/streaming`);
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

        const deviceId = String(parts[parts.length - 1].split(":")[1]);
        const matchingWidgets = this.widgets.filter(
          (widget) => String(widget.DeviceId) === deviceId,
        );

        if (!matchingWidgets.length) return;

        matchingWidgets.forEach((widget) => {
          widget.dataStream = message;
        });

        this.widgets = [...this.widgets];
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
    <SideBar :menuOpen="this.menuOpen" />
    <div class="info-container">
      <div class="title-container">
        <button class="hamburger" @click="menuOpen = !menuOpen">☰</button>
        <div v-if="menuOpen" class="overlay" @click="menuOpen = false"></div>
        <h1>{{ this.dashboard.Name }}</h1>
      </div>
      <div class="widgets-container">
        <div v-for="widget in this.widgets">
          <Card
            v-if="widget.Widgettype == 'ca'"
            :deviceId="widget.DeviceId"
            :streamData="widget.dataStream"
            :deviceName="widget.deviceName"
            @delete="handleDeleteWidget(widget.ID)"
          />
          <Table
            v-if="widget.Widgettype == 'ta'"
            :deviceId="widget.DeviceId"
            :streamData="widget.dataStream"
            :deviceName="widget.deviceName"
            @delete="handleDeleteWidget(widget.ID)"
          />
          <WidgetButton
            v-if="widget.Widgettype == 'btn'"
            :deviceId="widget.DeviceId"
            :deviceName="widget.deviceName"
            :label="widget.Label"
            :payload="widget.Payload"
            :widgetID="widget.ID"
            :channel="widget.Channel"
            @delete="handleDeleteWidget(widget.ID)"
          />
          <WidgetSwitch
            v-if="widget.Widgettype == 'sw'"
            :deviceId="widget.DeviceId"
            :deviceName="widget.deviceName"
            :label="widget.Label"
            :payload="widget.Payload"
            :payload2="widget.Payload2"
            :widgetID="widget.ID"
            :channel="widget.Channel"
            @delete="handleDeleteWidget(widget.ID)"
          />
        </div>
        <button class="add-button" @click="triggerDialog">+</button>
      </div>
      <dialog id="create-widget-dialog">
        <div class="dialog-header">
          <span>CREATE WIDGET</span>
          <button @click="exitDialog" class="close-x">&times;</button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label for="widget_type">Type</label>
            <select id="widget_type" v-model="this.new_widget.widget_type">
              <option value="ca">Card (Metric View)</option>
              <option value="ta">Table (Data Grid)</option>
              <option value="btn">Button (Action)</option>
              <option value="sw">Switch (2 States Action)</option>
            </select>
          </div>

          <div class="field">
            <label for="device">Source Device</label>
            <select id="device" v-model="this.new_widget.device">
              <option
                v-for="dev in available_devices"
                :key="dev.id"
                :value="dev"
              >
                {{ dev.Name }}
              </option>
            </select>
          </div>

          <div
            class="field"
            v-if="
              this.new_widget.widget_type == 'btn' ||
              this.new_widget.widget_type == 'sw'
            "
          >
            <label for="label">Label</label>
            <input type="text" id="label" v-model="this.new_widget.label" />
          </div>

          <div
            class="field"
            v-if="
              this.new_widget.widget_type == 'btn' ||
              this.new_widget.widget_type == 'sw'
            "
          >
            <label for="payload">Payload</label>
            <textarea
              id="payload"
              rows="4"
              cols="50"
              v-model="this.new_widget.payload"
            />
          </div>

          <div class="field" v-if="this.new_widget.widget_type == 'sw'">
            <label for="payload2">Payload 2</label>
            <textarea
              id="payload2"
              rows="4"
              cols="50"
              v-model="this.new_widget.payload2"
            />
          </div>

          <div
            class="field"
            v-if="
              this.new_widget.widget_type == 'btn' ||
              this.new_widget.widget_type == 'sw'
            "
          >
            <label for="payload2">Channel</label>
            <input type="text" id="label" v-model="this.new_widget.channel" />
          </div>

          <button @click="handleAddWidget" class="add-btn">
            ADD TO DASHBOARD
          </button>
        </div>
      </dialog>
    </div>
  </div>
</template>

<style scoped>
.hamburger {
  display: none;
  top: 1rem;
  left: 1rem;
  font-size: 1.8rem;
  background: none;
  border: none;
  color: #eeeeee;
  cursor: pointer;
}

.title-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
}

@media (max-width: 900px) {
  .hamburger {
    display: block;
  }
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.info-container {
  padding: 1rem;
  width: 100%;
}
.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
  color: #eeeeee;
}

.add-button {
  border: none;
  width: 10rem;
  height: 10rem;
  color: #eeeeee;
  font-weight: bold;
  font-size: 4rem;
  background-color: transparent;
  border: 1px solid #eeeeee;
  border-style: dashed;
  transition-duration: 300ms;
}

.add-button:hover {
  border: 1px solid #a6a6a6;
  border-style: dashed;
  color: #a6a6a6;
}

.widgets-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  flex-wrap: wrap;
}

.create-button {
  border: none;
  text-decoration: none;
  cursor: pointer;
  background-color: #a8dadc;
  color: #121212;
  transition-duration: 300ms;
  padding: 0.5rem;
  font-size: large;
  border-radius: 0.3rem;
}

.create-button:hover {
  background-color: #8ac6c9;
}

.cancel-button {
  border: none;
  text-decoration: none;
  cursor: pointer;
  background-color: transparent;
}

dialog {
  position: fixed;
  inset: 0;
  margin: auto;
  width: 90vw;
  max-width: 380px;
  background-color: #1a1a1a;
  color: #e0e0e0;
  border: 1px solid #444444;
  padding: 0;
  border-radius: 1rem;
}

.dialog-header {
  padding: 1rem 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #444444;
}

.dialog-header span {
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 1px;
  color: #b0b0b0;
}

.close-x {
  background: none;
  border: none;
  color: #b0b0b0;
  font-size: 1.5rem;
  cursor: pointer;
}

.dialog-form {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #b0b0b0;
}

input,
textarea,
select {
  background: #161616;
  border: 1px solid #444444;
  padding: 0.6rem;
  color: #e0e0e0;
  border-radius: 2px;
  font-size: 0.9rem;
}

input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: #e0e0e0;
}

.add-btn {
  margin-top: 0.5rem;
  background-color: #8ac6c9;
  color: #1a1a1a;
  border: none;
  padding: 0.75rem;
  font-weight: 700;
  font-size: 0.8rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.add-btn:hover {
  opacity: 0.9;
}
</style>
