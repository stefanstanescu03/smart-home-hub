<script>
import axios from "axios";
import SideBar from "../components/SideBar.vue";
export default {
  components: { SideBar },
  data() {
    return {
      username: "",
      password: "",
      failed: false,
    };
  },
  methods: {
    async handleLogin() {
      if (this.password != "" && this.username != "") {
        await axios
          .post("http://localhost:5000/user/login", {
            username: this.username,
            password: this.password,
          })
          .then((response) => {
            this.handleAddCookie(response);
            this.$router.push("/");
          })
          .catch((error) => {
            this.failed = true;
          });
      }
    },
    handleAddCookie(response) {
      const token = response.data.token;
      let d = new Date();
      d.setTime(d.getTime() + 1 * 24 * 60 * 60 * 1000);
      let expires = "expires=" + d.toUTCString();
      document.cookie = "Token=" + token + ";" + expires + ";path=/";
    },
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div class="login-form">
      <div class="field">
        <label for="username">Username</label>
        <input type="text" id="username" name="username" v-model="username" />
      </div>
      <div class="field">
        <label for="password">Password</label>
        <input
          type="password"
          id="password"
          name="password"
          v-model="password"
        />
      </div>
      <button @click="handleLogin">Login</button>
      <p v-if="this.failed == true">Incorrect credentials</p>
    </div>
  </div>
</template>

<style scoped>
.login-form {
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
