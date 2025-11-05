<script>
import SideBar from "../components/SideBar.vue";
import DeviceInfo from "../components/DeviceInfo.vue";
import axios from "axios";
export default {
  data() {
    return {
      devices: [],
    };
  },
  components: { SideBar, DeviceInfo },
  methods: {
    async getDiscoveredDevices() {
      await axios
        .get("http://localhost:5000/device/public")
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
      <h1>List of public devices</h1>
      <h1 v-if="this.devices.length == 0">No public devices</h1>
      <div v-for="device in this.devices" class="devices-container">
        <DeviceInfo
          :deviceName="device.Name"
          :ip="device.Ip"
          :visibility="device.Visibility"
        />
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
