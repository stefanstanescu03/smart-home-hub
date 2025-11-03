<script>
export default {
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      console.log(token);
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
    <RouterLink to="/">Home</RouterLink>
    <RouterLink v-if="getToken() == -1" to="/login">Login</RouterLink>
    <RouterLink v-if="getToken() == -1" to="/signup">Signup</RouterLink>
    <button v-if="getToken() != -1" @click="deleteToken">Logout</button>
  </nav>
</template>

<style>
nav {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  border-right: 1px solid black;
  padding: 1rem;
}
</style>
