<script>
import SideBar from "../components/SideBar.vue";
import DeviceInfo from "../components/DeviceInfo.vue";
import axios from "axios";
export default {
  components: { SideBar, DeviceInfo },
  data() {
    return {
      devices: [],
    };
  },
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      if (token === undefined) {
        return undefined;
      }
      return token;
    },
    async fetchDevices() {
      await axios
        .get("http://localhost:5000/device/", {
          headers: { Authorization: `Bearer ${this.getToken()}` },
        })
        .then((response) => {
          this.devices = response.data.devices;
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
  mounted() {
    this.fetchDevices();
  },
};
</script>

<template>
  <div class="page-container">
    <SideBar />
    <div>
      <h1 v-if="getToken() == undefined">You are logged in as guest</h1>
      <h1 v-if="getToken() != undefined">You are logged in</h1>
      <h1 v-if="this.devices.length == 0">No devices added</h1>
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
.devices-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>
