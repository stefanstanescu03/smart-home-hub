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
        ident: "",
        visibility: "",
        cloud_api: "",
      },
      new_device: {
        ident: "",
        visibility: "",
        cloud_api: "",
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
    async fetchDevices() {
      try {
        if (this.getToken() != undefined) {
          const response = await axios.get("/api/device/user", {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          });
          this.devices = response.data.devices;
        }
      } catch (err) {
        console.log(err);
      }
    },
    async handleDelete(device) {
      const id = device.ID;
      try {
        await axios.delete(`/api/device/delete/${device.ID}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.devices = this.devices.filter((device) => device.ID != id);
      } catch (err) {
        console.log(err);
      }
    },
    async triggerEdit(device) {
      this.selected_device.id = device.ID;
      this.selected_device.name = device.Name;
      this.selected_device.csv_location = device.Csv_location;
      this.selected_device.ident = device.Ident;
      this.selected_device.visibility =
        device.Visibility === true ? "public" : "private";
      this.selected_device.cloud_api = device.Cloud_api;
      const dialog = document.getElementById("edit-dialog");
      dialog.show();
    },
    async triggerCreate() {
      const dialog = document.getElementById("create-dialog");
      dialog.show();
    },
    handleCancelEdit() {
      const dialog = document.getElementById("edit-dialog");
      dialog.close();
    },
    handleCancelCreate() {
      const dialog = document.getElementById("create-dialog");
      dialog.close();
    },
    async handleUpdateDevice() {
      const id = this.selected_device.id;
      console.log(this.selected_device);
      try {
        await axios.put(
          `/api/device/update/${id}`,
          {
            name: this.selected_device.name,
            ident: this.selected_device.ident,
            visibility: this.selected_device.visibility === "public",
            cloud_api: this.selected_device.cloud_api,
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          },
        );
        this.handleCancelEdit();
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
    async handleAddDevice() {
      try {
        await axios.post(
          "/api/device/create",
          {
            name: this.new_device.name,
            ident: this.new_device.ident,
            visibility: this.new_device.visibility === "public",
            cloud_api: this.new_device.cloud_api,
          },
          { headers: { Authorization: `Bearer ${this.getToken()}` } },
        );
        this.handleCancelCreate();
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
  },
  mounted() {
    this.fetchDevices();
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
        <h1 v-if="getToken() == undefined">You are logged in as guest</h1>
        <h1 v-if="getToken() != undefined">My devices</h1>
      </div>
      <button
        @click="triggerCreate"
        class="top-create-button"
        v-if="this.getToken() != undefined"
      >
        Add a device
      </button>
      <h1 v-if="this.devices.length == 0 && this.getToken() != undefined">
        No devices added
      </h1>
      <table class="devices-container">
        <tbody>
          <tr v-if="this.devices.length != 0">
            <th>Name</th>
            <th class="ident">Ident</th>
            <th>Visibility</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
          <DeviceInfo
            v-for="device in this.devices"
            :deviceName="device.Name"
            :ident="device.Ident"
            :visibility="device.Visibility"
            :id="device.ID"
            :should_appear="true"
            @edit="triggerEdit(device)"
            @delete="handleDelete(device)"
          />
        </tbody>
      </table>

      <dialog id="edit-dialog">
        <div class="dialog-header">
          <span>EDIT DEVICE: {{ selected_device.name }}</span>
          <button @click="handleCancelEdit" class="close-x">&times;</button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label for="edit_name">Name</label>
            <input type="text" id="edit_name" v-model="selected_device.name" />
          </div>

          <div class="field">
            <label for="edit_visibility">Visibility</label>
            <select id="edit_visibility" v-model="selected_device.visibility">
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>

          <div class="field">
            <label>Cloud API (Optional)</label>
            <input type="text" v-model="selected_device.cloud_api" />
          </div>

          <button @click="handleUpdateDevice" class="add-btn">
            SAVE CHANGES
          </button>
        </div>
      </dialog>

      <dialog id="create-dialog">
        <div class="dialog-header">
          <span>NEW DEVICE</span>
          <button @click="handleCancelCreate" class="close-x">&times;</button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label>Name</label>
            <input type="text" v-model="new_device.name" />
          </div>

          <div class="field">
            <label>Identificator</label>
            <input type="text" v-model="new_device.ident" />
          </div>

          <div class="field">
            <label>Visibility</label>
            <select v-model="new_device.visibility">
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>

          <div class="field">
            <label>Cloud API (Optional)</label>
            <input type="text" v-model="new_device.cloud_api" />
          </div>

          <button @click="handleAddDevice" class="add-btn">ADD DEVICE</button>
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

.create-button {
  border: none;
  cursor: pointer;
  background-color: #a8dadc;
  color: #121212;
  padding: 0.6rem;
  font-size: 1rem;
  border-radius: 0.4rem;
  transition: background-color 200ms;
}

.create-button:hover {
  background-color: #8ac6c9;
}

.cancel-button {
  border: none;
  background: transparent;
  cursor: pointer;
}

.devices-container {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
  font-family: inherit;
}

.devices-container th {
  text-align: left;
  font-weight: 600;
  padding: 0.6rem 0.75rem;
  border-bottom: 1px solid #606060;
  color: #eeeeee;
}

.devices-container td {
  padding: 0.6rem 0.75rem;
  vertical-align: middle;
  color: #eeeeee;
}

@media (max-width: 900px) {
  .ident {
    display: none;
  }
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
  border-radius: 2px;
}

.add-btn:hover {
  opacity: 0.9;
}
</style>
