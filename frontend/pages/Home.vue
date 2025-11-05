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
    <div class="info-container">
      <h1 v-if="getToken() == undefined">You are logged in as guest</h1>
      <h1 v-if="getToken() != undefined">Your devices</h1>
      <h1 v-if="this.devices.length == 0">No devices added</h1>
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
