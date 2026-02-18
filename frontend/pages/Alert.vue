<script>
import SideBar from "../components/SideBar.vue";
import axios from "axios";
export default {
  components: { SideBar },
  data() {
    return {
      device: {},
      alerts: [],
      models: [],
      new_alert: {
        subject: "",
        content: "",
        key: "",
        value: "",
        condition: "",
      },
      new_model: {
        param: "",
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
          },
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
          },
        );
        changeAlert.NotifyEmail = !changeAlert.NotifyEmail;
      } catch (err) {
        console.error(err);
      }
    },
    async handleCreateModel() {
      try {
        console.log(this.device);

        await axios.post(
          "/api/model/anomaly/create",
          {
            param: this.new_model.param,
            DeviceId: this.device.ID,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          },
        );
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
    async handleFetchModels() {
      try {
        const res = await axios.get(
          `/api/model/anomaly/device/${this.$route.params.id}`,
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          },
        );
        this.models = res.data.models;
        console.log(this.models);
      } catch (err) {
        console.log(err);
      }
    },
    async handleChangeNotifyAD(modelId) {
      try {
        const changeModel = this.models.find((model) => model.ID === modelId);

        if (!changeModel) return;

        await axios.put(
          `/api/model/anomaly/update/${modelId}`,
          {
            NotifyEmail: !changeModel.NotifyEmail,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          },
        );
        changeModel.NotifyEmail = !changeModel.NotifyEmail;
      } catch (err) {
        console.error(err);
      }
    },

    async handleDeleteModel(modelId) {
      try {
        await axios.delete(`/api/model/anomaly/${modelId}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.models = this.models.filter((model) => model.ID != modelId);
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
    triggerCreateDialogAD() {
      const dialog = document.getElementById("create-dialog-ad");
      dialog.show();
    },
    handleCancelCreateDialogAD() {
      const dialog = document.getElementById("create-dialog-ad");
      dialog.close();
    },
  },
  mounted() {
    this.handleFetchDevice();
    this.handleFetchAlerts();
    this.handleFetchModels();
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
      <div class="buttons-container">
        <button @click="triggerCreateDialog" class="top-create-button">
          Create Alert
        </button>
        <button @click="triggerCreateDialogAD" class="top-create-button">
          Create Anomaly Detector
        </button>
      </div>

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

      <div class="models-container">
        <div v-for="model in models" :key="model.ID" class="model-card">
          <div class="model-main">
            <div class="model-param">{{ model.Param }}</div>
            <div class="model-meta">Anomaly Detection Model</div>
          </div>

          <div class="model-actions">
            <button
              class="model-notify inactive"
              v-if="!model.NotifyEmail"
              @click="handleChangeNotifyAD(model.ID)"
            >
              Notify
            </button>

            <button
              class="model-notify active"
              v-else
              @click="handleChangeNotifyAD(model.ID)"
            >
              Notify
            </button>

            <button class="model-delete" @click="handleDeleteModel(model.ID)">
              <img src="../public/delete.png" alt="" />
            </button>
          </div>
        </div>
      </div>

      <dialog id="create-dialog">
        <div class="dialog-header">
          <span>CREATE ALERT</span>
          <button @click="handleCancelCreateDialog" class="close-x">
            &times;
          </button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label for="subject">Subject</label>
            <input type="text" id="subject" v-model="new_alert.subject" />
          </div>

          <div class="field">
            <label for="content">Notification Content</label>
            <input type="text" id="content" v-model="new_alert.content" />
          </div>

          <div class="field">
            <label for="key">Metric / Key Name</label>
            <input type="text" id="key" v-model="new_alert.key" />
          </div>

          <div class="form-row">
            <div class="field">
              <label for="condition">Condition</label>
              <select id="condition" v-model="new_alert.condition">
                <option value="above">Above</option>
                <option value="below">Below</option>
                <option value="equal">Equal</option>
              </select>
            </div>

            <div class="field">
              <label for="value">Threshold Value</label>
              <input type="number" id="value" v-model="new_alert.value" />
            </div>
          </div>

          <button @click="handleCreateAlert" class="add-btn">
            CREATE ALERT
          </button>
        </div>
      </dialog>

      <dialog id="create-dialog-ad">
        <div class="dialog-header">
          <span>CREATE ANOMALY MODEL</span>
          <button @click="handleCancelCreateDialogAD" class="close-x">
            &times;
          </button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label for="key">Key name</label>
            <input type="text" id="key" v-model="new_model.param" />
          </div>

          <div class="info-box">
            <span class="info-icon">!</span>
            <p>
              Model requires baseline data. Ensure the sensor has collected data
              for
              <strong>at least 24 hours</strong> before initialization.
            </p>
          </div>

          <button @click="handleCreateModel" class="add-btn">
            CREATE MODEL
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
select {
  background: #161616;
  border: 1px solid #444444;
  padding: 0.6rem;
  color: #e0e0e0;
  border-radius: 2px;
  font-size: 0.9rem;
}

input:focus,
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

.buttons-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
}

.models-container {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 1rem;
  max-width: 720px;
}

.model-card {
  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 0.9rem 1.1rem;
  border-radius: 0.75rem;

  background: #1c1c1c;
  border: 1px solid #2a2a2a;
  width: 100%;
}

.model-main {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
  min-width: 0;
}

.model-param {
  font-size: 1.05rem;
  font-weight: 600;
  color: #eeeeee;
  word-break: break-word;
}

.model-meta {
  font-size: 0.75rem;
  color: #9a9a9a;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.model-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-left: 1rem;
}

.model-notify {
  padding: 0.4rem 0.7rem;
  font-size: 0.75rem;
  border-radius: 0.4rem;
  cursor: pointer;
  border: 1px solid transparent;
  background: transparent;
  transition: all 200ms ease;
}

.model-notify.inactive {
  border-color: #555;
  color: #cccccc;
}

.model-notify.active {
  border-color: #c9cc0d;
  color: #c9cc0d;
}

.model-delete {
  border: none;
  background: transparent;
  cursor: pointer;
}

.model-delete img {
  width: 18px;
  height: 18px;
}
</style>
