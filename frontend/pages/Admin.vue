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
    async handleMakeAdmin(accountId) {
      try {
        await axios.put(
          `/api/user/toggle-role/${accountId}`,
          {},
          {
            headers: { Authorization: `Bearer ${this.getToken()}` },
          },
        );
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
    async handleDeleteAccount(accountId) {
      try {
        await axios.delete(`/api/user/delete/${accountId}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.$router.go(0);
      } catch (err) {
        console.log(err);
      }
    },
    async handleDeleteDevice(deviceId) {
      try {
        await axios.delete(`/api/device/delete/${deviceId}`, {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.devices = this.devices.filter((device) => device.ID != deviceId);
        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;
      } catch (err) {
        console.log(err);
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
        <table>
          <tbody>
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
                  <button
                    class="delete-button"
                    @click="this.handleDeleteAccount(account.ID)"
                  >
                    <img
                      src="../public/delete.png"
                      alt=""
                      height="20"
                      width="20"
                    />
                  </button>
                  <button
                    class="action-button"
                    @click="this.handleMakeAdmin(account.ID)"
                    v-if="account.Role === 'user'"
                  >
                    Make admin
                  </button>
                  <button
                    class="action-button"
                    @click="this.handleMakeAdmin(account.ID)"
                    v-if="account.Role === 'admin'"
                  >
                    Make user
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="table-container">
        <h2>Added devices</h2>
        <table>
          <tbody>
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
                  <button
                    class="delete-button"
                    @click="this.handleDeleteDevice(device.ID)"
                  >
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
          </tbody>
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
