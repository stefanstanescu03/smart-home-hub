<script>
import axios from "axios";
export default {
  props: ["deviceName", "ip", "visibility"],
  data() {
    return {
      status: "checking...",
    };
  },
  methods: {
    async checkStatus() {
      await axios
        .get(`http://localhost:5000/device/ping/${this.ip}`)
        .then((response) => {
          if (response.data.status == true) {
            this.status = "connected";
          } else {
            this.status = "disconnected";
          }
        })
        .catch((error) => (this.status = "disconnected"));
    },
  },
  mounted() {
    this.checkStatus();
  },
};
</script>

<template>
  <div class="device-container">
    <p>{{ deviceName }}</p>
    <p>{{ ip }}</p>
    <p v-if="visibility == false">private</p>
    <p v-if="visibility == true">public</p>
    <p>{{ status }}</p>
    <button>View</button>
  </div>
</template>

<style scoped>
.device-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
}
</style>
