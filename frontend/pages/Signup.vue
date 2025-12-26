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
      menuOpen: false,
    };
  },
  methods: {
    async handleSignup() {
      if (
        this.password != "" &&
        this.username != "" &&
        this.email != "" &&
        this.repeated_password == this.password
      ) {
        try {
          await axios.post("/api/user/signup", {
            username: this.username,
            email: this.email,
            password: this.password,
          });
          this.$router.push("/login");
        } catch (err) {
          console.log(err);
        }
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
    <SideBar :menuOpen="this.menuOpen" />
    <div class="info-container">
      <div class="title-container">
        <button class="hamburger" @click="menuOpen = !menuOpen">☰</button>
        <div v-if="menuOpen" class="overlay" @click="menuOpen = false"></div>
        <h1>Create an account</h1>
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
        <button class="signup-button" @click="handleSignup">Signup</button>
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
  display: flex;
  flex-direction: column;
  width: 20rem;
  gap: 1rem;
  padding-top: 2rem;
  color: #eeeeee;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
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
  color: #eeeeee;
  background-color: #252525;
  outline: none;
  box-shadow: none;
  border: 1px solid #eeeeee;

  font-size: medium;
  padding: 0.5rem;
  border-radius: 0.3rem;
}
.signup-button {
  color: #eeeeee;
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

.signup-button:hover {
  background-color: #8ac6c9;
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
