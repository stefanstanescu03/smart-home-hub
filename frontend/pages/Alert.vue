<script>
import SideBar from "../components/SideBar.vue";
import axios from "axios";
export default {
  components: { SideBar },
  data() {
    return {
      device: {},
      alerts: [],
      new_alert: {
        subject: "",
        content: "",
        key: "",
        value: "",
        condition: "",
      },
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
    async handleFetchDevice() {
      try {
        const res = await axios.get(`/api/device/${this.$route.params.id}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.device = res.data.device;
      } catch (err) {
        console.log(err);
      }
    },
    async handleCreateAlert() {
      try {
        console.log(this.device);
        let condition = "";
        if (this.new_alert.condition === "above") {
          condition = ">";
        } else if (this.new_alert.condition === "below") {
          condition = "<";
        } else {
          condition = "=";
        }
        const conditionString =
          this.new_alert.key + condition + this.new_alert.value;

        await axios.post(
          "/api/alert/create",
          {
            subject: this.new_alert.subject,
            content: this.new_alert.content,
            condition: conditionString,
            DeviceId: this.device.ID,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        );
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
    async handleFetchAlerts() {
      try {
        const res = await axios.get(`/api/alert/${this.$route.params.id}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.alerts = res.data.alerts;
      } catch (err) {
        console.log(err);
      }
    },
    async handleDeleteAlert(alertId) {
      try {
        await axios.delete(`/api/alert/delete/${alertId}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.alerts = this.alerts.filter((alert) => alert.ID != alertId);
      } catch (err) {
        console.log(err);
      }
    },
    triggerCreateDialog() {
      const dialog = document.getElementById("create-dialog");
      dialog.show();
    },
    handleCancelCreateDialog() {
      const dialog = document.getElementById("create-dialog");
      dialog.close();
    },
  },
  mounted() {
    this.handleFetchDevice();
    this.handleFetchAlerts();
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
        <h1>Alerts for {{ this.device.Name }}</h1>
      </div>
      <button @click="triggerCreateDialog" class="top-create-button">
        Create Alert
      </button>

      <table class="alerts-container">
        <tr v-if="this.alerts.length != 0">
          <th>Subject</th>
          <th>Content</th>
          <th>Condition</th>
          <th>Action</th>
        </tr>
        <tr v-for="alert in this.alerts">
          <td>{{ alert.Subject }}</td>
          <td>{{ alert.Content }}</td>
          <td>{{ alert.Condition }}</td>
          <td>
            <button class="delete-button" @click="handleDeleteAlert(alert.ID)">
              <img src="../public/delete.png" alt="" height="20" width="20" />
            </button>
          </td>
        </tr>
        <h1 v-if="this.alerts.length == 0">No alerts added for this device</h1>
      </table>

      <dialog id="create-dialog">
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Create Alert</p>
            <button @click="handleCancelCreateDialog" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="subject">Subject: </label>
            <input
              type="text"
              id="subject"
              name="subject"
              v-model="this.new_alert.subject"
            />
          </div>
          <div class="field">
            <label for="content">Content: </label>
            <input
              type="text"
              id="content"
              name="content"
              v-model="this.new_alert.content"
            />
          </div>

          <div class="field">
            <label for="key">Key name: </label>
            <input
              type="text"
              id="key"
              name="key"
              v-model="this.new_alert.key"
            />
          </div>

          <div class="field">
            <label for="condition">Condition: </label>
            <select
              name="condition"
              id="condition"
              v-model="this.new_alert.condition"
            >
              <option value="above">above</option>
              <option value="below">below</option>
              <option value="equal">equal</option>
            </select>
          </div>

          <div class="field">
            <label for="value">Value: </label>
            <input
              type="number"
              id="value"
              name="value"
              v-model="this.new_alert.value"
            />
          </div>

          <div class="field"></div>
          <button @click="handleCreateAlert" class="create-button">
            Create
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

.top-create-button {
  border: none;
  background-color: transparent;
  border: 1px solid #eeeeee;
  text-decoration: none;
  cursor: pointer;
  color: #eeeeee;
  padding: 0.5rem;
  border-radius: 0.3rem;
  transition-duration: 300ms;
}

.top-create-button:hover {
  border: 1px solid #a6a6a6;
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

.field {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  gap: 0.5rem;
}

input {
  color: #eeeeee;
  background-color: #252525;
  outline: none;
  box-shadow: none;
  border: none;

  font-size: medium;
  padding: 0.5rem;
  border-radius: 0.3rem;
}
table {
  border-collapse: collapse;
  width: 100%;
}
th {
  border-bottom: 1px solid #c9c9c9;
  padding: 0.3rem;
  text-align: left;
}
td {
  border-bottom: 1px solid #eeeeee;
  padding: 0.3rem;
}
.delete-button {
  border: none;
  text-decoration: none;
  background-color: transparent;
  cursor: pointer;
}
</style>
