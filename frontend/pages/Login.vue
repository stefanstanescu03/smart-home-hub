<script>
import SideBar from "../components/SideBar.vue";
export default {
  components: { SideBar },
  data() {
    return {
      username: "",
      password: "",
      failed: false,
      errorMessage: "",
      menuOpen: false,
    };
  },
  methods: {
    async handleLogin() {
      if (this.password != "" && this.username != "") {
        try {
          const response = await fetch("/api/user/login", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              username: this.username,
              password: this.password,
            }),
          });

          const data = await response.json();

          if (!response.ok) {
            throw new Error(data.error || "Login failed");
          }

          this.handleAddCookie(data);
          this.$router.push("/");
        } catch (err) {
          this.errorMessage = err.message;
          this.failed = true;
        }
      } else {
        this.errorMessage = "All fields must be completed";
        this.failed = true;
      }
    },
    handleAddCookie(data) {
      const token = data.token;

      if (!token) {
        console.error("No token found in response data");
        return;
      }

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
      <div class="login-form">
        <span>Enter in your account</span>
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
        <button class="login-button" @click="handleLogin">Login</button>
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

.info-container {
  padding: 1.5rem;
  width: 100%;
  margin-left: 15%;
}

@media (max-width: 900px) {
  .hamburger {
    display: block;
  }
  .info-container {
    margin-left: 0;
    padding: 1rem;
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

.login-form {
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

input,
select {
  background: #161616;
  border: 1px solid #444444;
  padding: 0.6rem;
  color: #e0e0e0;
  border-radius: 2px;
  font-size: 0.9rem;
}

input:focus,
select:focus {
  outline: none;
  border-color: #e0e0e0;
}

.login-button {
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

.login-button:hover {
  opacity: 0.9;
}

span {
  font-size: 1rem;
  font-weight: 700;
  letter-spacing: 1px;
  color: #e0e0e0;
}

@media (max-width: 900px) {
  .login-form {
    margin-left: auto;
    margin-right: auto;
    width: 100%;
    max-width: 320px;
  }
}
</style>
