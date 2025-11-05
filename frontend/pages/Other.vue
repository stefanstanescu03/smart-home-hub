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
    <div class="info-container">
      <h1>List of public devices</h1>
      <h1 v-if="this.devices.length == 0">No public devices</h1>
      <table class="devices-container">
        <tr v-if="this.devices.length != 0">
          <th>Name</th>
          <th>Ip</th>
          <th>Visibility</th>
          <th>Status</th>
        </tr>
        <DeviceInfo
          v-for="device in this.devices"
          :deviceName="device.Name"
          :ip="device.Ip"
          :visibility="device.Visibility"
        />
      </table>
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
.info-container {
  width: 50%;
}
.info-container {
  width: 50%;
}
table {
  border-collapse: collapse;
  width: 100%;
}
th {
  border-bottom: 1px solid #c9c9c9;
  padding: 0.3rem;
  text-align: left;
}
</style>
