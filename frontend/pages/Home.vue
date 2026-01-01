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
        csv_location: "",
        visibility: "",
      },
      new_device: {
        ident: "",
        csv_location: "",
        visibility: "",
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
      try {
        await axios.put(
          `/api/device/update/${id}`,
          {
            name: this.selected_device.name,
            ident: this.selected_device.ident,
            csv_location: this.selected_device.csv_location,
            visibility: this.selected_device.visibility === "public",
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
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
            csv_location: this.new_device.csv_location,
            visibility: this.new_device.visibility === "public",
          },
          { headers: { Authorization: `Bearer ${this.getToken()}` } }
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
      </table>

      <dialog id="edit-dialog">
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Edit {{ this.selected_device.Name }}</p>
            <button @click="handleCancelEdit" class="cancel-button">
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
          <button @click="handleUpdateDevice" class="create-button">
            Modify
          </button>
        </div>
      </dialog>

      <dialog id="create-dialog">
        <div class="dialog-container">
          <div class="top-dialog">
            <p>New device</p>
            <button @click="handleCancelCreate" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="name">Name: </label>
            <input
              type="text"
              id="name"
              name="name"
              v-model="this.new_device.name"
            />
          </div>
          <div class="field">
            <label for="device_csv_location">Data location: </label>
            <input
              type="text"
              id="device_csv_location"
              name="device_csv_location"
              v-model="this.new_device.csv_location"
            />
          </div>
          <div class="field">
            <label for="device_ident">Identificator: </label>
            <input
              type="text"
              id="device_ident"
              name="device_ident"
              v-model="this.new_device.ident"
            />
          </div>
          <div class="field">
            <label for="device_visibility">Visibility: </label>
            <select
              name="device_visibility"
              id="device_visibility"
              v-model="this.new_device.visibility"
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
