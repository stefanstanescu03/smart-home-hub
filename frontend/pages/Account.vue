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
      await axios
        .get("http://localhost:5000/user/info", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        })
        .then((response) => {
          this.username = response.data.user.Username;
          this.email = response.data.user.Email;
        })
        .catch((err) => console.log(err));
    },
    async handleUpdate() {
      await axios
        .put(
          "http://localhost:5000/user/update",
          {
            username: this.username,
            email: this.email,
            password: this.password,
          },
          { headers: { Authorization: `Bearer ${this.getToken()}` } }
        )
        .then((response) => {
          this.changed = true;
          this.failed = false;
        })
        .catch((err) => {
          this.changed = false;
          this.failed = true;
        });
    },
  },
  mounted() {
    this.getAccount();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
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
      <button @click="handleUpdate">Update</button>
      <p v-if="failed == true">Password incorrect</p>
      <p v-if="changed == true">Account updated</p>
    </div>
  </div>
</template>

<style scoped>
.signup-form {
  display: flex;
  flex-direction: column;
  width: 20rem;
  gap: 1rem;
  padding-top: 2rem;
}
.field {
  display: flex;
  flex-direction: column;
}
.page-container {
  display: flex;
  flex-direction: row;
  gap: 25%;
  height: 100%;
}

label {
  font-size: large;
}

input {
  outline: none;
  box-shadow: none;
  border: 1px solid #a6a6a6;

  font-size: medium;
  padding: 0.3rem;
  border-radius: 0.3rem;
}

button {
  border: none;
  text-decoration: none;
  cursor: pointer;
  background-color: #ff8441;
  color: #121212;
  transition-duration: 300ms;
  padding: 0.5rem;
  font-size: large;
  border-radius: 0.3rem;
}

button:hover {
  background-color: #fe8d50;
}
</style>
