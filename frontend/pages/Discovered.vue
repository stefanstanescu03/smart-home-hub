<script>
import SideBar from "../components/SideBar.vue";
import axios from "axios";
export default {
  data() {
    return {
      devices: [],
    };
  },
  components: { SideBar },
  methods: {
    async getDiscoveredDevices() {
      await axios
        .get("http://localhost:5000/device/discovered")
        .then((response) => (this.devices = response.data.devices))
        .catch((err) => console.log(err));
    },
  },
  mounted() {
    this.getDiscoveredDevices();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div>
      <h1>List of discovered devices</h1>
      <h1 v-if="this.devices.length == 0">No devices discovered</h1>
      <div>
        <div v-for="device in devices" class="device-container">
          <p>{{ device }}</p>
          <button>Add Device</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
}
.device-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  padding: 1rem;
}
</style>
