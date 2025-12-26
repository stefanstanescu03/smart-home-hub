<script>
import SideBar from "../components/SideBar.vue";
import Card from "../components/Card.vue";
import Table from "../components/Table.vue";
import axios from "axios";
export default {
  components: { SideBar, Card, Table },
  data() {
    return {
      dashboard: "",
      widgets: [],
      available_devices: [],
      widget_type: "",
      device: "",
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
            }))
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
          (widget) => String(widget.DeviceId) === deviceId
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
              <option value="ta">Table</option>
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
  color: #eeeeee;
  background-color: #1a1a1a;
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
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  margin: auto;
  width: 60vw;
  max-width: 500px;
  border: none;
  outline: none;
  box-shadow: none;
  border-radius: 0.3rem;
  background-color: #1a1a1a;
}
</style>
