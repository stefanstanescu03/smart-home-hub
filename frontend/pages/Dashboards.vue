<script>
import SideBar from "../components/SideBar.vue";
import DashboardInfo from "../components/DashboardInfo.vue";
import axios from "axios";
export default {
  components: { SideBar, DashboardInfo },
  data() {
    return {
      dashboards: [],
      dashboard_visibility: "",
      dashboard_name: "",
      selected_dashboard: {},
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
    triggerCreateDialog() {
      const dialog = document.getElementById("create-dialog");
      dialog.show();
    },
    handleCancelCreateDialog() {
      const dialog = document.getElementById("create-dialog");
      dialog.close();
    },
    triggerUpdateDialog(dashboard) {
      this.selected_dashboard.id = dashboard.ID;
      this.selected_dashboard.name = dashboard.Name;
      this.selected_dashboard.visibility =
        dashboard.Visibility === true ? "public" : "private";
      const dialog = document.getElementById("update-dialog");
      dialog.show();
    },
    handleCancelUpdateDialog() {
      const dialog = document.getElementById("update-dialog");
      dialog.close();
    },
    async handleGetDashboards() {
      await axios
        .get("http://localhost:5000/dashboard/", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        })
        .then((res) => {
          this.dashboards = res.data.dashboards;
        })
        .catch((err) => console.log(err));
    },
    async handleCreateDashboard() {
      try {
        await axios.post(
          "http://localhost:5000/dashboard/create",
          {
            name: this.dashboard_name,
            visibility: this.dashboard_visibility === "public",
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        );

        this.handleCancelCreateDialog();
        this.$router.go(0);
      } catch (err) {
        console.error("Dashboard creation failed:", err);
      }
    },
    async handleDeleteDashboard(dashboard) {
      const id = dashboard.ID;
      await axios
        .delete(`http://localhost:5000/dashboard/delete/${id}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        })
        .then(() => {
          this.dashboards = this.dashboards.filter((d) => d.ID != id);
        })
        .catch((err) => console.log(err));
    },
    async handleUpdateDashboard() {
      const id = this.selected_dashboard.id;
      await axios
        .put(
          `http://localhost:5000/dashboard/update/${id}`,
          {
            name: this.selected_dashboard.name,
            visibility: this.selected_dashboard.visibility === "public",
          },
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          }
        )
        .then(() => {
          this.handleCancelCreateDialog();
          this.$router.go(0);
        })
        .catch((err) => console.log(err));
    },
  },
  mounted() {
    this.handleGetDashboards();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div class="info-container">
      <h1>Available dashboards</h1>
      <h1 v-if="this.dashboards.length == 0">No dashboard created</h1>
      <button @click="triggerCreateDialog" class="top-create-button">
        Create dashboard
      </button>
      <table class="dashboards-container" v-if="this.dashboards.length != 0">
        <tr>
          <th>Name</th>
          <th>Visibility</th>
          <th>Action</th>
        </tr>
        <DashboardInfo
          v-for="dashboard in this.dashboards"
          :name="dashboard.Name"
          :visibility="dashboard.Visibility"
          :should_appear="true"
          :id="dashboard.ID"
          @edit="triggerUpdateDialog(dashboard)"
          @delete="handleDeleteDashboard(dashboard)"
        />
      </table>

      <dialog id="create-dialog">
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Create dashboard</p>
            <button @click="handleCancelCreateDialog" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="name">Name: </label>
            <input type="text" id="name" name="name" v-model="dashboard_name" />
          </div>
          <div class="field">
            <label for="device_visibility">Visibility: </label>
            <select
              name="device_visibility"
              id="device_visibility"
              v-model="dashboard_visibility"
            >
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>
          <button @click="handleCreateDashboard" class="create-button">
            Create
          </button>
        </div>
      </dialog>

      <dialog id="update-dialog">
        <div class="dialog-container">
          <div class="top-dialog">
            <p>Update dashboard</p>
            <button @click="handleCancelUpdateDialog" class="cancel-button">
              <img src="../public/close.png" height="20" width="20" alt="" />
            </button>
          </div>
          <div class="field">
            <label for="name">Name: </label>
            <input
              type="text"
              id="name"
              name="name"
              v-model="this.selected_dashboard.name"
            />
          </div>
          <div class="field">
            <label for="device_visibility">Visibility: </label>
            <select
              name="device_visibility"
              id="device_visibility"
              v-model="this.selected_dashboard.visibility"
            >
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>
          <button @click="handleUpdateDashboard" class="create-button">
            Update
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
  width: 50%;
}
table {
  border-collapse: collapse;
  width: 100%;
}
th {
  border-bottom: 1px solid #eeeeee;
  padding: 0.3rem;
  text-align: left;
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
