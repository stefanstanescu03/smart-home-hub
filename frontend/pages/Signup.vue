<script>
import axios from "axios";
import SideBar from "../components/SideBar.vue";

export default {
  components: { SideBar },
  data() {
    return {
      username: "",
      password: "",
      repeated_password: "",
      email: "",
    };
  },
  methods: {
    async handleSignup() {
      if (this.password != "" && this.username != "" && this.email != "") {
        await axios
          .post("http://localhost:5000/user/signup", {
            username: this.username,
            email: this.email,
            password: this.password,
          })
          .then((response) => {
            this.$router.push("/login");
          })
          .catch((error) => {
            console.log(error);
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
        <label for="password">Password</label>
        <input
          type="password"
          id="password"
          name="password"
          v-model="password"
        />
      </div>

      <div class="field">
        <label for="password">Repeat password</label>
        <input
          type="password"
          id="repeated_password"
          name="repeated_password"
          v-model="repeated_password"
        />
      </div>
      <button @click="handleSignup">Signup</button>
    </div>
  </div>
</template>

<style scoped>
.signup-form {
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
