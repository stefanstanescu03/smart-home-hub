<script>
import SideBar from "../components/SideBar.vue";
export default {
  components: { SideBar },
  data() {
    return {
      menuOpen: false,
      available_devices: [],
      automations: [],
      routines: [],
      available_metrics: [],
      new_automation: {
        Name: "",
        Device1: "",
        Device2: "",
        Metric: "",
        Trigger: "",
        Target: "",
        Payload: "",
        ScheduleStart: "",
        ScheduleEnd: "",
        IsScheduled: false,
      },
      new_routine: {
        Name: "",
        Device: "",
        Payload: "",
        StartTime: "",
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
        const res = await fetch("/api/automation/", {
          method: "GET",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });

        if (!res.ok) throw new Error("Failed to fetch automations");

        const data = await res.json();

        this.automations = await Promise.all(
          data.automations.map(async (automation) => ({
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
        const res = await fetch(`/api/automation/delete/${automationId}`, {
          method: "DELETE",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });

        if (!res.ok) throw new Error("Failed to delete automation");

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
    triggerCreateRoutine() {
      const dialog = document.getElementById("create-dialog-routine");
      dialog.show();
    },
    handleCancelCreateRoutine() {
      const dialog = document.getElementById("create-dialog-routine");
      dialog.close();
    },
    async handleGetDeviceName(id) {
      try {
        const res = await fetch(`/api/device/${id}`, {
          method: "GET",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        if (!res.ok) return "default";

        const data = await res.json();
        return data.device.Name;
      } catch (err) {
        return "default";
      }
    },

    async handleGetAvailableDevices() {
      try {
        const res = await fetch("/api/device/", {
          method: "GET",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        if (!res.ok) throw new Error("Failed to fetch available devices");

        const data = await res.json();
        this.available_devices = data.devices;
      } catch (err) {
        console.error(err);
      }
    },

    async handleFetchRoutines() {
      try {
        const res = await fetch("/api/routine/", {
          method: "GET",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        if (!res.ok) throw new Error("Failed to fetch routines");

        const data = await res.json();

        // Using Promise.all to fetch device names concurrently
        this.routines = await Promise.all(
          data.routines.map(async (routine) => ({
            ...routine,
            deviceName: await this.handleGetDeviceName(routine.DeviceId),
          })),
        );
      } catch (err) {
        console.log(err);
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

        const res = await fetch("/api/automation/create", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${this.getToken()}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            Name: this.new_automation.Name,
            Device1Id: this.new_automation.Device1.ID,
            Device2Id: this.new_automation.Device2.ID,
            Payload: this.new_automation.Payload,
            Condition: conditionString,
            ScheduleStart: this.new_automation.ScheduleStart,
            ScheduleEnd: this.new_automation.ScheduleEnd,
          }),
        });

        if (!res.ok) throw new Error("Failed to create automation");
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },

    async handleDeleteRoutine(routineId) {
      try {
        const res = await fetch(`/api/routine/delete/${routineId}`, {
          method: "DELETE",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });

        if (!res.ok) throw new Error("Failed to delete routine");

        this.routines = this.routines.filter(
          (routine) => routine.ID != routineId,
        );
      } catch (err) {
        console.log(err);
      }
    },

    async handleCreateRoutine() {
      try {
        const res = await fetch("/api/routine/create", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${this.getToken()}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            Name: this.new_routine.Name,
            DeviceId: this.new_routine.Device.ID,
            Payload: this.new_routine.Payload,
            StartTime: this.new_routine.StartTime,
          }),
        });

        if (!res.ok) throw new Error("Failed to create routine");
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },

    async fetchMetricsForDevice(id) {
      try {
        const res = await fetch(`/api/device/params/${id}`, {
          method: "GET",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        if (!res.ok) throw new Error("Failed to fetch device parameters");

        const data = await res.json();
        let keys = data.params.split(",");
        keys = keys.filter((key) => key != "timestamp" && key != "");
        keys = keys.map((key) => key.replace(/\[.*?\]/g, ""));
        this.available_metrics = keys;
      } catch (err) {
        console.log(err);
      }
    },
  },
  watch: {
    "new_automation.Device1": async function (newDevice) {
      if (newDevice && newDevice.ID) {
        await this.fetchMetricsForDevice(newDevice.ID);
      }
    },
  },
  async mounted() {
    await this.handleGetAvailableDevices();
    await this.handleFetchAutomations();
    await this.handleFetchRoutines();
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
              <span
                v-if="
                  automation.ScheduleStart != '' && automation.ScheduleEnd != ''
                "
                class="automation-condition"
              >
                {{ automation.ScheduleStart }} - {{ automation.ScheduleEnd }}
              </span>
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

      <h1>My Routines</h1>
      <button class="top-create-button" @click="triggerCreateRoutine">
        Create Routine
      </button>

      <div class="routines-container">
        <div v-for="routine in routines" :key="routine.ID" class="routine-card">
          <div class="routine-timeline">
            <span class="t-start">{{ routine.StartTime }}</span>
            <span class="t-end">{{ routine.EndTime }}</span>
          </div>

          <div class="routine-body">
            <div class="routine-meta">
              <h2 class="routine-name">{{ routine.Name }}</h2>
            </div>

            <div class="routine-target">
              {{ routine.deviceName }}
            </div>
          </div>

          <button class="delete-btn" @click="handleDeleteRoutine(routine.ID)">
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
            <select id="key" v-model="this.new_automation.Metric">
              <option v-if="available_metrics.length === 0" disabled>
                Select a device first
              </option>
              <option
                v-for="metric in available_metrics"
                :key="metric"
                :value="metric"
              >
                {{ metric }}
              </option>
            </select>
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
            <div class="schedule-one-line">
              <input
                type="checkbox"
                v-model="this.new_automation.IsScheduled"
              />
              <label for="schedule">Schedule</label>
            </div>
          </div>

          <div class="field" v-if="this.new_automation.IsScheduled">
            <div class="schedule-one-line">
              <input type="time" v-model="this.new_automation.ScheduleStart" />
              <label for="schedule-times">to</label>
              <input type="time" v-model="this.new_automation.ScheduleEnd" />
            </div>
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

      <dialog id="create-dialog-routine">
        <div class="dialog-header">
          <span>NEW ROUTINE</span>
          <button @click="handleCancelCreateRoutine" class="close-x">
            &times;
          </button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label>Routine Label</label>
            <input type="text" v-model="this.new_routine.Name" />
          </div>

          <div class="field">
            <label for="device2">Action Device</label>
            <select id="device2" v-model="this.new_routine.Device">
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
            <div class="schedule-one-line">
              <label for="schedule-times">Start Time</label>
              <input type="time" v-model="this.new_routine.StartTime" />
            </div>
          </div>

          <div class="field">
            <label for="payload">Command</label>
            <textarea
              id="payload"
              rows="4"
              cols="50"
              v-model="this.new_routine.Payload"
            />
          </div>

          <button @click="handleCreateRoutine" class="add-btn">
            ADD ROUTINE
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
  top: 1rem;
  left: 50%;
  transform: translateX(-50%);

  margin: 0;
  width: 90vw;
  max-width: 380px;
  background-color: #1a1a1a;
  color: #e0e0e0;
  border: 1px solid #444444;
  padding: 0;
  border-radius: 1rem;

  max-height: calc(100vh - 2rem);
  overflow-y: auto;
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
  color: #eeeeee;
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

.schedule-one-line {
  display: flex;
  flex-direction: row;
  gap: 0.5rem;
  align-items: center;
}

.schedule-one-line input[type="checkbox"] {
  appearance: none;
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  border: 1px solid #333;
  border-radius: 4px;
  background-color: transparent;
  display: inline-grid;
  place-content: center;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  position: relative;
}

.schedule-one-line input[type="checkbox"]:checked {
  background-color: #8ac6c9;
  border-color: #8ac6c9;
}

.schedule-one-line input[type="checkbox"]::after {
  content: "";
  width: 10px;
  height: 5px;
  border-left: 2px solid #1a1a1a;
  border-bottom: 2px solid #1a1a1a;
  transform: rotate(-45deg) translate(1px, -1px);
  opacity: 0;
  transition: opacity 0.1s ease-in-out;
}

.schedule-one-line input[type="checkbox"]:checked::after {
  opacity: 1;
}

.routines-container {
  display: flex;
  flex-direction: column;
  padding-top: 2rem;
  gap: 1rem;
  max-width: 600px;
}

.routine-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #1c1c1c;
  border: 1px solid #2a2a2a;
  border-radius: 0.75rem;
  padding: 20px;
}

.routine-timeline {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-right: 12px;
  border-right: 1px solid #222;
  min-width: 50px;
}

.t-start,
.t-end {
  font-family: "Monaco", monospace;
  font-size: 1.5rem;
  font-weight: 700;
  color: #eeeeee;
}

.routine-body {
  flex: 1;
  padding-left: 12px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.routine-meta {
  display: flex;
  align-items: center;
  gap: 6px;
}

.routine-name {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #eee;
  white-space: nowrap;
  overflow: hidden;
}

.delete-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 6px;
  opacity: 0.5;
  transition:
    opacity 0.2s,
    transform 0.2s;
}

.delete-btn:hover {
  opacity: 1;
  transform: scale(1.1);
}
</style>
