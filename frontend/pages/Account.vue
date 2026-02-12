<script>
import axios from "axios";
import SideBar from "../components/SideBar.vue";

export default {
  components: { SideBar },
  data() {
    return {
      username: "",
      password: "",
      email: "",
      failed: false,
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
    async getAccount() {
      try {
        const response = await axios.get("/api/user/info", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        });
        this.username = response.data.user.Username;
        this.email = response.data.user.Email;
      } catch (err) {
        console.log(err);
      }
    },
    async handleUpdate() {
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
      }
    },
  },
  mounted() {
    this.getAccount();
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
        <h1>Edit your account</h1>
      </div>
      <div class="signup-form">
        <div class="field">
          <label for="username">Username</label>
          <input type="text" id="username" name="username" v-model="username" />
        </div>

        <div class="field">
          <label for="email">Email</label>
          <input type="email" id="email" name="email" v-model="email" />
        </div>

        <div class="field">
          <label for="password">Current password</label>
          <input
            type="password"
            id="password"
            name="password"
            v-model="password"
          />
        </div>
        <button class="update-button" @click="handleUpdate">Update</button>
        <p class="incorrect" v-if="failed == true">Incorrect password</p>
        <p v-if="changed == true">Account updated</p>
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

.signup-form {
  margin-top: 2rem;
  padding: 1.5rem;
  width: 100%;
  max-width: 380px;
  background-color: #1a1a1a;
  border-radius: 1rem;

  display: flex;
  flex-direction: column;
  gap: 1.2rem;
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

.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
  color: #eeeeee;
}

input {
  background: #161616;
  border: 1px solid #444444;
  padding: 0.6rem;
  color: #e0e0e0;
  border-radius: 2px;
  font-size: 0.9rem;
}

input:focus {
  outline: none;
  border-color: #e0e0e0;
}

.update-button {
  margin-top: 0.5rem;
  background-color: #8ac6c9;
  color: #1a1a1a;
  border: none;
  padding: 0.75rem;
  font-weight: 700;
  font-size: 1rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.update-button:hover {
  opacity: 0.9;
}

@media (max-width: 900px) {
  .signup-form {
    margin-left: auto;
    margin-right: auto;
    width: 100%;
    max-width: 320px;
  }
}
</style>
