<script>
import axios from "axios";
import SideBar from "../components/SideBar.vue";
import Graph from "../components/Graph.vue";
export default {
  components: { SideBar, Graph },
  data() {
    return {
      menuOpen: false,
      device: {},
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
    async fetchParams() {
      try {
        if (this.getToken() != undefined) {
          const res = await axios.get(
            `/api/device/params/${this.$route.params.id}`,
            {
              headers: { Authorization: `Bearer ${this.getToken()}` },
            },
          );
          return res.data.params;
        }
      } catch (err) {
        console.log(err);
      }
    },
    async handleFetchDevice() {
      try {
        if (this.getToken() != undefined) {
          const response = await axios.get(
            `/api/device/${this.$route.params.id}`,
            {
              headers: { Authorization: `Bearer ${this.getToken()}` },
            },
          );
          this.device = response.data.device;
          const path = this.device.Csv_location;
          this.device.filename = path.split("/").pop().split("\\").pop();

          const paramsStr = await this.fetchParams();
          let paramsArr = paramsStr.split(",");
          paramsArr = paramsArr.filter(
            (param) => param != "" && param != "timestamp",
          );
          this.device.params = paramsArr;
        }
      } catch (err) {
        console.log(err);
      }
    },
  },
  async mounted() {
    await this.handleFetchDevice();
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
        <h1>Analyze data collected by: {{ this.device.Name }}</h1>
      </div>
      <div class="graphs-container">
        <Graph
          v-if="device.filename"
          v-for="param in this.device.params"
          :filename="this.device.filename"
          :param="param"
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
  color: #eeeeee;
}

.graphs-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.info-container {
  padding: 1.5rem;
  width: 100%;
}

.title-container {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.hamburger {
  display: none;
  font-size: 1.8rem;
  background: none;
  border: none;
  color: #eeeeee;
  cursor: pointer;
}

@media (max-width: 900px) {
  .hamburger {
    display: block;
  }
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1000;
}
</style>
