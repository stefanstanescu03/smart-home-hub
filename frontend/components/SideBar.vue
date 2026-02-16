<script>
import axios from "axios";

export default {
  props: ["menuOpen"],
  data() {
    return {
      admin: false,
    };
  },
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      return token ?? -1;
    },
    deleteToken() {
      document.cookie =
        "Token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      this.$router.push("/login");
      // props.menuOpen = false;
    },
    navigate(path) {
      this.$router.push(path);
      // props.menuOpen = false;
    },
    async isAdmin() {
      if (this.getToken() == -1) {
        this.admin = false;
      }
      try {
        const response = await axios.get("/api/user/info", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        const role = response.data.user.Role;
        if (role !== "admin") {
          this.admin = false;
        } else {
          this.admin = true;
        }
      } catch (err) {
        console.log(err);
        this.admin = false;
      }
    },
  },
  async mounted() {
    await this.isAdmin();
  },
};
</script>

<template>
  <nav :class="{ open: menuOpen }">
    <button @click="navigate('/')">
      <img src="../public/home.png" width="25" height="25" />
      <span>Home</span>
    </button>

    <button v-if="getToken() != -1" @click="navigate('/account')">
      <img src="../public/user.png" width="25" height="25" />
      <span>Account</span>
    </button>

    <button @click="navigate('/dashboards')">
      <img src="../public/dashboard.png" width="25" height="25" />
      <span>Dashboards</span>
    </button>

    <button v-if="getToken() != -1" @click="navigate('/automations')">
      <img src="../public/robotic-arm.png" width="25" height="25" />
      <span>Automations</span>
    </button>

    <button v-if="getToken() != -1" @click="navigate('/other')">
      <img src="../public/iot-devices.png" width="25" height="25" />
      <span>Other Devices</span>
    </button>

    <button v-if="admin == true" @click="navigate('/admin')">
      <img src="../public/security.png" width="25" height="25" />
      <span>Admin</span>
    </button>

    <button v-if="getToken() == -1" @click="navigate('/login')">Login</button>
    <button v-if="getToken() == -1" @click="navigate('/signup')">Signup</button>
    <button v-if="getToken() != -1" @click="deleteToken">Logout</button>
  </nav>
</template>

<style scoped>
.hamburger {
  display: none;
  top: 1rem;
  left: 1rem;
  z-index: 1001;
  font-size: 1.8rem;
  background: none;
  border: none;
  color: white;
  cursor: pointer;
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

nav {
  display: flex;
  flex-direction: column;
  width: 15%;
  min-width: 220px;
  background-color: #1a1a1a;
  height: 100vh;
}

button {
  display: flex;
  align-items: center;
  gap: 1rem;
  border: none;
  background-color: #1a1a1a;
  color: #eeeeee;
  cursor: pointer;
  padding: 0.7rem;
  font-size: large;
}

button:hover {
  background-color: rgb(40, 41, 43);
}

button > img {
  filter: invert(87%);
}

@media (max-width: 900px) {
  .hamburger {
    display: block;
  }

  nav {
    position: fixed;
    left: -100%;
    top: 0;
    width: 70%;
    max-width: 300px;
    transition: left 0.3s ease;
    z-index: 1001;
  }

  nav.open {
    left: 0;
  }
}
</style>
