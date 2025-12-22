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
      await axios
        .put(
          `http://localhost:5000/device/update/${id}`,
          {
            name: this.selected_device.name,
            ident: this.selected_device.ident,
            csv_location: this.selected_device.csv_location,
            visibility: this.selected_device.visibility === "public",
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        )
        .then(() => {
          this.handleCancelEdit();
          this.$router.go(0);
        })
        .catch((err) => console.log(err));
    },
    async handleAddDevice() {
      await axios
        .post(
          "http://localhost:5000/device/create",
          {
            name: this.new_device.name,
            ident: this.new_device.ident,
            csv_location: this.new_device.csv_location,
            visibility: this.new_device.visibility === "public",
          },
          { headers: { Authorization: `Bearer ${this.getToken()}` } }
        )
        .then(() => {
          this.handleCancelCreate();
          this.$router.go(0);
        })
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
    <SideBar :menuOpen="this.menuOpen" />
    <div class="info-container">
      <div class="title-container">
        <button class="hamburger" @click="menuOpen = !menuOpen">☰</button>
        <div v-if="menuOpen" class="overlay" @click="menuOpen = false"></div>
        <h1 v-if="getToken() == undefined">You are logged in as guest</h1>
        <h1 v-if="getToken() != undefined">My devices</h1>
      </div>
      <button @click="triggerCreate" class="top-create-button">
        Add a device
      </button>
      <h1 v-if="this.devices.length == 0">No devices added</h1>
      <table class="devices-container">
        <tr v-if="this.devices.length != 0">
          <th>Name</th>
          <th>Ident</th>
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

.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
  color: #eeeeee;
}
.info-container {
  padding: 1rem;
  width: 100%;
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
</style>
