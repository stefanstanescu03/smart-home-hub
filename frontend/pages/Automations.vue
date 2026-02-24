<script>
import axios from "axios";
import SideBar from "../components/SideBar.vue";
export default {
  components: { SideBar },
  data() {
    return {
      menuOpen: false,
      available_devices: [],
      automations: [],
      new_automation: {
        Name: "",
        Device1: "",
        Device2: "",
        Metric: "",
        Trigger: "",
        Target: "",
        Payload: "",
      },
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
    async handleFetchAutomations() {
      try {
        const res = await axios.get("/api/automation/", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.automations = await Promise.all(
          res.data.automations.map(async (automation) => ({
            ...automation,
            device1Name: await this.handleGetDeviceName(automation.Device1Id),
            device2Name: await this.handleGetDeviceName(automation.Device2Id),
          })),
        );
      } catch (err) {
        console.log(err);
      }
    },
    async handleDeleteAutomation(automationId) {
      try {
        await axios.delete(`/api/automation/delete/${automationId}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.automations = this.automations.filter(
          (automation) => automation.ID != automationId,
        );
      } catch (err) {
        console.log(err);
      }
    },
    triggerCreate() {
      const dialog = document.getElementById("create-dialog");
      dialog.show();
    },
    handleCancelCreate() {
      const dialog = document.getElementById("create-dialog");
      dialog.close();
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
    async handleCreateAutomation() {
      try {
        let condition = "";
        if (this.new_automation.Trigger === "above") {
          condition = ">";
        } else if (this.new_automation.Trigger === "below") {
          condition = "<";
        } else {
          condition = "=";
        }
        const conditionString =
          this.new_automation.Metric + condition + this.new_automation.Target;

        await axios.post(
          "/api/automation/create",
          {
            Name: this.new_automation.Name,
            Device1Id: this.new_automation.Device1.ID,
            Device2Id: this.new_automation.Device2.ID,
            Payload: this.new_automation.Payload,
            Condition: conditionString,
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
  },
  async mounted() {
    await this.handleGetAvailableDevices();
    await this.handleFetchAutomations();
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
        <h1>My automations</h1>
      </div>
      <button class="top-create-button" @click="triggerCreate">
        Create automation
      </button>

      <div class="automations-container">
        <div
          v-for="automation in automations"
          :key="automation.ID"
          class="automation-card"
        >
          <div class="automation-content">
            <div class="automation-header">
              <h2 class="automation-name">{{ automation.Name }}</h2>
              <span class="automation-condition">{{
                automation.Condition
              }}</span>
            </div>

            <div class="devices-row">
              <div class="device-badge">{{ automation.device1Name }}</div>
              <div class="device-connector">→</div>
              <div class="device-badge">{{ automation.device2Name }}</div>
            </div>
          </div>

          <button
            class="delete-btn"
            @click="handleDeleteAutomation(automation.ID)"
          >
            <img src="../public/delete.png" alt="Delete" />
          </button>
        </div>
      </div>

      <dialog id="create-dialog">
        <div class="dialog-header">
          <span>NEW AUTOMATION</span>
          <button @click="handleCancelCreate" class="close-x">&times;</button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label>Automation Label</label>
            <input type="text" v-model="this.new_automation.Name" />
          </div>

          <div class="field">
            <label for="device1">Trigger Device</label>
            <select id="device1" v-model="this.new_automation.Device1">
              <option
                v-for="dev in available_devices"
                :key="dev.id"
                :value="dev"
              >
                {{ dev.Name }}
              </option>
            </select>
          </div>

          <div class="field">
            <label for="device2">Action Device</label>
            <select id="device2" v-model="this.new_automation.Device2">
              <option
                v-for="dev in available_devices"
                :key="dev.id"
                :value="dev"
              >
                {{ dev.Name }}
              </option>
            </select>
          </div>

          <div class="field">
            <label for="key">Metric / Key Name</label>
            <input type="text" id="key" v-model="this.new_automation.Metric" />
          </div>

          <div class="field">
            <label for="condition">Trigger Condition</label>
            <select id="condition" v-model="this.new_automation.Trigger">
              <option value="above">Above</option>
              <option value="below">Below</option>
              <option value="equal">Equal</option>
            </select>
          </div>

          <div class="field">
            <label for="value">Target Value</label>
            <input
              type="number"
              id="value"
              v-model="this.new_automation.Target"
            />
          </div>

          <div class="field">
            <label for="payload">Command</label>
            <textarea
              id="payload"
              rows="4"
              cols="50"
              v-model="this.new_automation.Payload"
            />
          </div>

          <button @click="handleCreateAutomation" class="add-btn">
            ADD AUTOMATION
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
  background: rgba(0, 0, 0, 0.6);
  z-index: 1000;
}

.top-create-button {
  margin-bottom: 1rem;
  border: 1px solid #eeeeee;
  background: transparent;
  color: #eeeeee;
  padding: 0.5rem 0.75rem;
  border-radius: 0.4rem;
  cursor: pointer;
  transition: all 200ms ease;
}

.top-create-button:hover {
  color: #8c8c8c;
  border-color: #8c8c8c;
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

.automations-container {
  display: flex;
  flex-direction: column;
  padding-top: 2rem;
  gap: 1rem;
  max-width: 600px;
}

.automation-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #1c1c1c;
  border: 1px solid #2a2a2a;
  border-radius: 0.75rem;
  padding: 20px;
}

.automation-content {
  flex-grow: 1;
}

.automation-header {
  display: flex;
  align-items: baseline;
  gap: 1rem;
  margin-bottom: 1rem;
}

.automation-name {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #eeeeee;
}

.automation-condition {
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #eeeeee;
  background: #2a2a2a;
  padding: 0.5rem;
  border-radius: 0.4rem;
}

.devices-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.device-badge {
  background: #2a2a2a;
  color: #eeeeee;
  padding: 0.5rem;
  border-radius: 0.4rem;
  font-size: 1rem;
  font-weight: 500;
}

.device-connector {
  color: #cbd5e0;
  font-weight: bold;
}

.delete-btn {
  background: none;
  border: none;
  padding: 8px;
  cursor: pointer;
  border-radius: 8px;
  display: flex;
  align-items: center;
}

.delete-btn img {
  width: 20px;
  height: 20px;
}
</style>
