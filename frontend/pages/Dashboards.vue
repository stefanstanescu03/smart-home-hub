<script>
import SideBar from "../components/SideBar.vue";
import DashboardInfo from "../components/DashboardInfo.vue";
export default {
  components: { SideBar, DashboardInfo },
  data() {
    return {
      dashboards: [],
      public_dashboards: [],
      dashboard_visibility: "",
      dashboard_name: "",
      selected_dashboard: {},
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
      try {
        const res = await fetch("/api/dashboard/", {
          method: "GET",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });

        if (!res.ok) throw new Error("Failed to fetch dashboards");

        const data = await res.json();
        this.dashboards = data.dashboards;
      } catch (err) {
        console.log(err);
      }
    },

    async handleGetPublicDashboards() {
      try {
        const res = await fetch("/api/dashboard/public", {
          method: "GET",
        });

        if (!res.ok) throw new Error("Failed to fetch public dashboards");

        const data = await res.json();
        this.public_dashboards = data.dashboards;
      } catch (err) {
        console.log(err);
      }
    },

    async handleCreateDashboard() {
      try {
        const res = await fetch("api/dashboard/create", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${this.getToken()}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name: this.dashboard_name,
            visibility: this.dashboard_visibility === "public",
          }),
        });

        if (!res.ok) throw new Error("Dashboard creation failed");

        this.handleCancelCreateDialog();
        this.$router.go(0);
      } catch (err) {
        console.error("Dashboard creation failed:", err);
      }
    },

    async handleDeleteDashboard(dashboard) {
      const id = dashboard.ID;
      try {
        const res = await fetch(`/api/dashboard/delete/${id}`, {
          method: "DELETE",
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });

        if (!res.ok) throw new Error("Failed to delete dashboard");

        this.dashboards = this.dashboards.filter((d) => d.ID != id);
      } catch (err) {
        console.log(err);
      }
    },

    async handleUpdateDashboard() {
      const id = this.selected_dashboard.id;
      try {
        const res = await fetch(`/api/dashboard/update/${id}`, {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${this.getToken()}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name: this.selected_dashboard.name,
            visibility: this.selected_dashboard.visibility === "public",
          }),
        });

        if (!res.ok) throw new Error("Failed to update dashboard");

        this.handleCancelCreateDialog();
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
  },
  mounted() {
    if (this.getToken() != undefined) {
      this.handleGetDashboards();
    }
    this.handleGetPublicDashboards();
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
        <h1>My dashboards</h1>
      </div>
      <h1 v-if="this.dashboards.length == 0 && this.getToken() != undefined">
        No dashboard created
      </h1>
      <button
        @click="triggerCreateDialog"
        class="top-create-button"
        v-if="this.getToken() != undefined"
      >
        Create dashboard
      </button>
      <table class="dashboards-container" v-if="this.dashboards.length != 0">
        <thead>
          <tr>
            <th>Name</th>
            <th>Visibility</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <DashboardInfo
            v-for="dashboard in this.dashboards"
            :name="dashboard.Name"
            :visibility="dashboard.Visibility"
            :should_appear="true"
            :id="dashboard.ID"
            type="private"
            @edit="triggerUpdateDialog(dashboard)"
            @delete="handleDeleteDashboard(dashboard)"
          />
        </tbody>
      </table>

      <h1>Public dashboards</h1>

      <table
        class="dashboards-container"
        v-if="this.public_dashboards.length != 0"
      >
        <thead>
          <tr>
            <th>Name</th>
            <th>Visibility</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <DashboardInfo
            v-for="dashboard in this.public_dashboards"
            :name="dashboard.Name"
            :visibility="dashboard.Visibility"
            :should_appear="true"
            :id="dashboard.ID"
            type="public"
          />
        </tbody>
      </table>

      <dialog id="create-dialog">
        <div class="dialog-header">
          <span>CREATE DASHBOARD</span>
          <button @click="handleCancelCreateDialog" class="close-x">
            &times;
          </button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label for="new_dash_name">Dashboard Name</label>
            <input type="text" id="new_dash_name" v-model="dashboard_name" />
          </div>

          <div class="field">
            <label for="new_dash_visibility">Visibility</label>
            <select id="new_dash_visibility" v-model="dashboard_visibility">
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>

          <button @click="handleCreateDashboard" class="add-btn">
            CREATE DASHBOARD
          </button>
        </div>
      </dialog>

      <dialog id="update-dialog">
        <div class="dialog-header">
          <span>UPDATE DASHBOARD: {{ selected_dashboard.name }}</span>
          <button @click="handleCancelUpdateDialog" class="close-x">
            &times;
          </button>
        </div>

        <div class="dialog-form">
          <div class="field">
            <label for="dash_name">Dashboard Name</label>
            <input
              type="text"
              id="dash_name"
              v-model="selected_dashboard.name"
            />
          </div>

          <div class="field">
            <label for="dash_visibility">Visibility</label>
            <select
              id="dash_visibility"
              v-model="selected_dashboard.visibility"
            >
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </div>

          <button @click="handleUpdateDashboard" class="add-btn">
            UPDATE SETTINGS
          </button>
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
