<script>
export default {
  props: ["deviceId", "deviceName", "streamData"],
  data() {
    return {
      valuesArray: [],
    };
  },
  watch: {
    streamData: {
      immediate: true,
      deep: true,
      handler(message) {
        if (!message) return;
        const parts = message.split(",");
        let value = {};
        parts.forEach((p) => {
          const [k, v] = p.split(":");
          if (k != "id") {
            value[k] = v;
          }
        });
        this.valuesArray.push(value);
        if (this.valuesArray.length == 100) {
          this.valuesArray = [value];
        }
      },
    },
  },
};
</script>

<template>
  <div class="table-wrapper">
    <h1 class="name-h">{{ deviceName }}</h1>
    <div class="values-container">
      <table>
        <thead>
          <tr>
            <th v-for="(value, key) in valuesArray[0]" :key="key">{{ key }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(value, index) in valuesArray" :key="index">
            <td v-for="(prop, key) in value" :key="key">{{ prop }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.name-h {
  font-size: large;
  margin-bottom: 0.5rem;
}
.table-wrapper {
  width: 100%;
}
.values-container {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #eeeeee;
  border-radius: 0.5rem;
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
td {
  padding: 0.3rem;
}
</style>
