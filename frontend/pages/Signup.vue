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
      errorMessage: "",
      failed: false,
    };
  },
  methods: {
    async handleSignup() {
      if (this.password != "" && this.username != "" && this.email != "") {
        if (this.password != this.repeated_password) {
          this.errorMessage = "Passwords must match";
          this.failed = true;
        } else {
          try {
            await axios.post("/api/user/signup", {
              username: this.username,
              email: this.email,
              password: this.password,
            });
            this.$router.push("/login");
          } catch (err) {
            this.errorMessage = err.response.data.error;
            this.failed = true;
          }
        }
      } else {
        this.errorMessage = "All fields must be completed";
        this.failed = true;
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
      </div>
      <div class="signup-form">
        <span>Create an account</span>

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

        <p style="line-height: 1.6; color: #b0b0b0; font-size: 0.85rem">
          • Minimum 8 characters<br />
          • Uppercase & lowercase letters<br />
          • At least one number<br />
          • At least one special character<br />
        </p>

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
        <p class="incorrect" v-if="this.failed == true">
          {{ this.errorMessage }}
        </p>
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
.page-container {
  display: flex;
  flex-direction: row;
  gap: 25%;
  height: 100%;
}

label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #b0b0b0;
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
.signup-button {
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

.signup-button:hover {
  opacity: 0.9;
}

span {
  font-size: 1rem;
  font-weight: 700;
  letter-spacing: 1px;
  color: #e0e0e0;
}

.incorrect {
  color: #e43333;
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
