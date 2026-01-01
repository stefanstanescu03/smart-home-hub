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
    async handleChangeNotify(alertId) {
      try {
        const changeAlert = this.alerts.find((alert) => alert.ID === alertId);

        if (!changeAlert) return;

        await axios.put(
          `/api/alert/update/${alertId}`,
          {
            subject: changeAlert.Subject,
            content: changeAlert.Content,
            condition: changeAlert.Condition,
            NotifyEmail: !changeAlert.NotifyEmail,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        );
        changeAlert.NotifyEmail = !changeAlert.NotifyEmail;
      } catch (err) {
        console.error(err);
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
            <div class="action-container">
              <button
                class="delete-button"
                @click="handleDeleteAlert(alert.ID)"
              >
                <img src="../public/delete.png" alt="" height="20" width="20" />
              </button>
              <button
                class="notify-button-not-pressed"
                @click="handleChangeNotify(alert.ID)"
                v-if="alert.NotifyEmail == false"
              >
                Notify
              </button>
              <button
                class="notify-button-pressed"
                @click="handleChangeNotify(alert.ID)"
                v-if="alert.NotifyEmail == true"
              >
                Notify
              </button>
            </div>
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
.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
  color: #eeeeee;
}

.info-container {
  padding: 1.5rem;
  width: 100%;
}

.title-container {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

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

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
  font-family: inherit;
}
tr {
  border-bottom: 1px solid #606060;
}
th {
  text-align: left;
  font-weight: 600;
  padding: 0.6rem 0.75rem;
  color: #eeeeee;
}
td {
  padding: 0.6rem 0.75rem;
  vertical-align: middle;
  color: #eeeeee;
}

.delete-button {
  border: none;
  text-decoration: none;
  background-color: transparent;
  cursor: pointer;
}

.action-container {
  display: flex;
  flex-direction: row;
  gap: 0.5rem;
  align-items: center;
}

.notify-button-not-pressed {
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

.notify-button-pressed {
  border: none;
  background-color: transparent;
  border: 1px solid #c9cc0d;
  text-decoration: none;
  cursor: pointer;
  color: #c9cc0d;
  padding: 0.5rem;
  border-radius: 0.3rem;
  transition-duration: 300ms;
}

input,
select {
  color: #eeeeee;
  background-color: #252525;
  border: none;
  outline: none;
  padding: 0.5rem;
  font-size: 0.95rem;
  border-radius: 0.35rem;
}

select {
  cursor: pointer;
}

dialog {
  position: fixed;
  inset: 0;
  margin: auto;
  width: 90vw;
  max-width: 480px;
  border: none;
  border-radius: 0.5rem;
  background-color: #1a1a1a;
  color: #eeeeee;
}

.dialog-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
}

.top-dialog {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

label {
  font-size: 0.8rem;
  color: #aaaaaa;
}
</style>
