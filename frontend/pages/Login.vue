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
}
.field {
  display: flex;
  flex-direction: column;
}
.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
}
</style>
