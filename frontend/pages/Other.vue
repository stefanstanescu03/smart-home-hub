<script>
import SideBar from "../components/SideBar.vue";
import DeviceInfo from "../components/DeviceInfo.vue";
import axios from "axios";
export default {
  data() {
    return {
      devices: [],
      menuOpen: false,
    };
  },
  components: { SideBar, DeviceInfo },
  methods: {
    async getDiscoveredDevices() {
      try {
        const response = await axios.get("/api/device/public");
        this.devices = response.data.devices;
      } catch (err) {
        console.log(err);
      }
    },
  },
  mounted() {
    this.getDiscoveredDevices();
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
        <h1>List of public devices</h1>
      </div>
      <h1 v-if="this.devices.length == 0">No public devices</h1>
      <table class="devices-container">
        <tr v-if="this.devices.length != 0">
          <th>Name</th>
          <th class="ident">Ident</th>
          <th>Visibility</th>
          <th>Status</th>
        </tr>
        <DeviceInfo
          v-for="device in this.devices"
          :deviceName="device.Name"
          :ident="device.Ident"
          :visibility="device.Visibility"
          :should_appear="false"
        />
      </table>
    </div>
  </div>
</template>

<style scoped>
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

.ident {
  display: block;
}

@media (max-width: 900px) {
  .hamburger {
    display: block;
  }

  .ident {
    display: none;
  }
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.info-container {
  padding: 1rem;
  width: 100%;
}

.page-container {
  display: flex;
  flex-direction: row;
  gap: 2rem;
  height: 100%;
  color: #eeeeee;
}

table {
  border-collapse: collapse;
  width: 100%;
}
th {
  border-bottom: 1px solid #eeeeee;
  padding: 0.3rem;
  text-align: left;
}
</style>
