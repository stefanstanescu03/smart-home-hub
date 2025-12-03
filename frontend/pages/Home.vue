<script>
import SideBar from "../components/SideBar.vue";
import DeviceInfo from "../components/DeviceInfo.vue";
import axios from "axios";
export default {
  components: { SideBar, DeviceInfo },
  data() {
    return {
      devices: [],
      selected_device: {
        id: "",
        ip: "",
        csv_location: "",
        visibility: "",
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
    async fetchDevices() {
      await axios
        .get("http://localhost:5000/device/", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        })
        .then((response) => {
          this.devices = response.data.devices;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async handleDelete(device) {
      const id = device.ID;
      await axios
        .delete(`http://localhost:5000/device/delete/${device.ID}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        })
        .then(
          () =>
            (this.devices = this.devices.filter((device) => device.ID != id))
        )
        .catch((err) => console.log(err));
    },
    async triggerEdit(device) {
      console.log(device);
      this.selected_device.id = device.ID;
      this.selected_device.name = device.Name;
      this.selected_device.csv_location = device.Csv_location;
      this.selected_device.ip = device.Ip;
      this.selected_device.visibility =
        device.Visibility === true ? "public" : "private";
      const dialog = document.querySelector("dialog");
      dialog.show();
    },
    handleCancel() {
      const dialog = document.querySelector("dialog");
      dialog.close();
    },
    async handleUpdateDevice() {
      const id = this.selected_device.id;
      await axios
        .put(
          `http://localhost:5000/device/update/${id}`,
          {
            name: this.selected_device.name,
            ip: this.selected_device.ip,
            csv_location: this.selected_device.csv_location,
            visibility: this.selected_device.visibility === "public",
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        )
        .then(this.handleCancel())
        .catch((err) => console.log(err));
    },
  },
  mounted() {
    this.fetchDevices();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div class="info-container">
      <h1 v-if="getToken() == undefined">You are logged in as guest</h1>
      <h1 v-if="getToken() != undefined">My devices</h1>
      <h1 v-if="this.devices.length == 0">No devices added</h1>
      <table class="devices-container">
        <tr v-if="this.devices.length != 0">
          <th>Name</th>
          <th>Ip</th>
          <th>Visibility</th>
          <th>Status</th>
          <th>Action</th>
        </tr>
        <DeviceInfo
          v-for="device in this.devices"
          :deviceName="device.Name"
          :ip="device.Ip"
          :visibility="device.Visibility"
          :should_appear="true"
          @edit="triggerEdit(device)"
          @delete="handleDelete(device)"
        />
      </table>

      <dialog>
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Edit device with ip: {{ this.selected_device.ip }}</p>
            <button @click="handleCancel" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="name">Name: </label>
            <input
              type="text"
              id="name"
              name="name"
              v-model="this.selected_device.name"
            />
          </div>
          <div class="field">
            <label for="device_csv_location">Data directory: </label>
            <input
              type="text"
              id="device_csv_location"
              name="device_csv_location"
              v-model="this.selected_device.csv_location"
            />
          </div>
          <div class="field">
            <label for="device_visibility">Visibility: </label>
            <select
              name="device_visibility"
              id="device_visibility"
              v-model="this.selected_device.visibility"
            >
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>
          <button @click="handleUpdateDevice" class="create-button">Add</button>
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
  width: 50%;
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
  width: 60vw;
  max-width: 500px;
  border: none;
  outline: none;
  box-shadow: none;
  border-radius: 0.3rem;
  background-color: #1a1a1a;
}
</style>
