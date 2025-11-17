<script>
import SideBar from "../components/SideBar.vue";
import axios from "axios";
export default {
  data() {
    return {
      devices: [],
      selected_ip: "",
      device_name: "",
      device_csv_location: "",
      device_visibility: "",
    };
  },
  components: { SideBar },
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      if (token === undefined) {
        return undefined;
      }
      return token;
    },
    async getDiscoveredDevices() {
      await axios
        .get("http://localhost:5000/device/discovered")
        .then((response) => (this.devices = response.data.devices))
        .catch((err) => console.log(err));
    },
    triggerDialog(ip) {
      this.selected_ip = ip;
      const dialog = document.querySelector("dialog");
      dialog.show();
    },
    async handleAddDevice() {
      await axios
        .post(
          "http://localhost:5000/device/create",
          {
            name: this.device_name,
            ip: this.selected_ip,
            csv_location: this.device_csv_location,
            visibility: this.device_visibility === "public",
          },
          { headers: { Authorization: `Bearer ${this.getToken()}` } }
        )
        .then(() => {
          const dialog = document.querySelector("dialog");
          dialog.close();
        })
        .catch((err) => console.log(err));
    },
    handleCancel() {
      const dialog = document.querySelector("dialog");
      dialog.close();
    },
  },
  mounted() {
    this.getDiscoveredDevices();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div>
      <h1>List of discovered devices</h1>
      <h1 v-if="this.devices.length == 0">No devices discovered</h1>
      <div class="discovered-list">
        <div v-for="device in devices" class="device-container">
          <p>Device ip: {{ device }}</p>
          <button @click="triggerDialog(device)" class="add-button">
            <img
              src="../public/plus-symbol-button.png"
              height="25"
              width="25"
              alt=""
            />
          </button>
        </div>
      </div>
      <dialog>
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Adding device with ip: {{ selected_ip }}</p>
            <button @click="handleCancel" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="name">Name: </label>
            <input type="text" id="name" name="name" v-model="device_name" />
          </div>
          <div class="field">
            <label for="device_csv_location">Data directory: </label>
            <input
              type="text"
              id="device_csv_location"
              name="device_csv_location"
              v-model="device_csv_location"
            />
          </div>
          <div class="field">
            <label for="device_visibility">Visibility: </label>
            <select
              name="device_visibility"
              id="device_visibility"
              v-model="device_visibility"
            >
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>
          <button @click="handleAddDevice" class="create-button">Add</button>
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
.discovered-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.device-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  border: 1px solid rgb(176, 176, 176);
  border-radius: 0.5rem;
  justify-content: space-between;
  width: 50%;
  padding-left: 0.5rem;
  padding-right: 0.5rem;
}
.add-button {
  background-color: transparent;
  border: none;
  cursor: pointer;
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
