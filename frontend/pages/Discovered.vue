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
      menuOpen: false,
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
    <SideBar :menuOpen="this.menuOpen" />
    <div class="info-container">
      <div class="title-container">
        <button class="hamburger" @click="menuOpen = !menuOpen">☰</button>
        <div v-if="menuOpen" class="overlay" @click="menuOpen = false"></div>
        <h1>List of discovered devices</h1>
      </div>
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
  color: #eeeeee;
  background-color: #252525;
  outline: none;
  box-shadow: none;
  border: none;

  font-size: medium;
  padding: 0.5rem;
  border-radius: 0.3rem;
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
