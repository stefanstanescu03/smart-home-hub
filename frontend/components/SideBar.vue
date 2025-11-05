<script>
export default {
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      if (token === undefined) {
        token = -1;
      }
      return token;
    },
    deleteToken() {
      document.cookie =
        "Token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      this.$router.push("/login");
    },
  },
};
</script>

<template>
  <nav>
    <button @click="this.$router.push('/')">
      <img src="../public/home.png" width="25" height="25" alt="" />
      <span>Home</span>
    </button>
    <button v-if="getToken() != -1" @click="this.$router.push('/account')">
      <img src="../public/user.png" width="25" height="25" alt="" />
      <span>Account</span>
    </button>
    <button v-if="getToken() != -1" @click="this.$router.push('/automations')">
      <img src="../public/robotic-arm.png" width="25" height="25" alt="" />
      <span>Automations</span>
    </button>
    <button @click="this.$router.push('/other')">
      <img src="../public/iot-devices.png" width="25" height="25" alt="" />
      <span>Other Devices</span>
    </button>
    <button v-if="getToken() != -1" @click="this.$router.push('/discovered')">
      <img src="../public/compass.png" width="25" height="25" alt="" />
      <span>Discovered</span>
    </button>
    <button v-if="getToken() == -1" @click="this.$router.push('/login')">
      Login
    </button>
    <button v-if="getToken() == -1" @click="this.$router.push('/signup')">
      Signup
    </button>
    <button v-if="getToken() != -1" @click="deleteToken">Logout</button>
  </nav>
</template>

<style scoped>
nav {
  display: flex;
  flex-direction: column;
  width: 15%;
  background-color: #183b4e;
}
button {
  display: flex;
  align-items: center;
  gap: 1rem;

  border: none;
  text-decoration: none;
  cursor: pointer;
  background-color: #183b4e;
  color: #d7d7d7;
  transition-duration: 300ms;
  padding: 0.5rem;
  text-align: left;
  font-size: large;
}

button:hover {
  background-color: #3e657f;
}

button > img {
  vertical-align: middle;
  filter: invert(87%) sepia(0%) saturate(0%) hue-rotate(154deg) brightness(95%)
    contrast(91%);
  transition: filter 0.3s;
}
</style>
