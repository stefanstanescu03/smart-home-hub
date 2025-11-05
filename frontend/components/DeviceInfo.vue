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
  <tr class="device-container">
    <td>{{ deviceName }}</td>
    <td>{{ ip }}</td>
    <td v-if="visibility == false">private</td>
    <td v-if="visibility == true">public</td>
    <td>{{ status }}</td>
  </tr>
</template>

<style scoped>
td {
  border-bottom: 1px solid #c9c9c9;
  padding: 0.3rem;
}
</style>
