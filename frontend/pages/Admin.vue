<script>
import axios from "axios";
import SideBar from "../components/SideBar.vue";

export default {
  components: { SideBar },
  data() {
    return {
      accounts: [],
      devices: [],
      changed: false,
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
    async getAccounts() {
      try {
        const response = await axios.get("/api/user/all", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.accounts = response.data.accounts;
      } catch (err) {
        console.log(err);
      }
    },
    async getDevices() {
      try {
        const response = await axios.get("/api/device/all", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.devices = response.data.devices;
      } catch (err) {
        console.log(err);
      }
    },
    async handleUpdate() {
      if (this.username != "" && this.email != "" && this.password != "") {
        try {
          await axios.put(
            "api/user/update",
            {
              username: this.username,
              email: this.email,
              password: this.password,
            },
            { headers: { Authorization: `Bearer ${this.getToken()}` } },
          );
          this.changed = true;
          this.failed = false;
        } catch (err) {
          this.changed = false;
          this.failed = true;
          this.errorMessage = err.response.data.error;
        }
      } else {
        this.changed = false;
        this.failed = true;
        this.errorMessage = "All fields must be completed";
      }
    },
  },
  async mounted() {
    await this.getAccounts();
    await this.getDevices();
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
        <h1>Admin page</h1>
      </div>
      <div class="table-container">
        <h2>Created accounts</h2>
        <table class="alerts-container">
          <tr>
            <th>Username</th>
            <th>Email</th>
            <th>Role</th>
            <th>Action</th>
          </tr>
          <tr v-for="account in accounts">
            <td>{{ account.Username }}</td>
            <td>{{ account.Email }}</td>
            <td>{{ account.Role }}</td>
            <td>
              <div class="action-container">
                <button class="delete-button">
                  <img
                    src="../public/delete.png"
                    alt=""
                    height="20"
                    width="20"
                  />
                </button>
                <button class="action-button">Make admin</button>
              </div>
            </td>
          </tr>
        </table>
      </div>
      <div class="table-container">
        <h2>Added devices</h2>
        <table class="alerts-container">
          <tr>
            <th>Name</th>
            <th>Ident</th>
            <th>Visibility</th>
            <th>Owner</th>
            <th>Action</th>
          </tr>
          <tr v-for="device in devices">
            <td>{{ device.Name }}</td>
            <td>{{ device.Ident }}</td>
            <td>{{ device.Visibility ? "public" : "private" }}</td>
            <td>{{ device.Owner }}</td>
            <td>
              <div class="action-container">
                <button class="delete-button">
                  <img src="../public/bell.png" alt="" height="20" width="20" />
                </button>
                <button class="delete-button">
                  <img src="../public/edit.png" alt="" height="25" width="25" />
                </button>
                <button class="delete-button">
                  <img
                    src="../public/delete.png"
                    alt=""
                    height="20"
                    width="20"
                  />
                </button>
              </div>
            </td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
h1 {
  color: #eeeeee;
}
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

.incorrect {
  color: #e43333;
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
td {
  padding: 0.6rem 0.75rem;
  vertical-align: middle;
  color: #eeeeee;
}

.delete-button {
  border: none;
  text-decoration: none;
  background-color: transparent;
  cursor: pointer;
}

.action-container {
  display: flex;
  flex-direction: row;
  gap: 0.5rem;
  align-items: center;
}

.action-button {
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

.action-button:hover {
  border: 1px solid #a6a6a6;
}
</style>
