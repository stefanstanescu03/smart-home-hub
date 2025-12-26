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
    <div class="title-container">
      <h1 class="name-h">{{ deviceName }}</h1>
      <button @click="$emit('delete')" class="delete-button">Delete</button>
    </div>
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
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border: 1px solid #eeeeee;
  border-radius: 0.5rem;
  padding-top: 0.5rem;
}
.values-container {
  max-height: 200px;
  overflow-y: auto;
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

.title-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  align-items: center;
  padding: 0.5rem;
  border-bottom: 1px solid #eeeeee;
}

.delete-button {
  border: none;
  background-color: transparent;
  border: 1px solid #eeeeee;
  text-decoration: none;
  cursor: pointer;
  color: #eeeeee;
  padding: 0.5rem;
  border-radius: 0.3rem;
  transition-duration: 300ms;
}

.delete-button:hover {
  border: 1px solid #e43333;
  color: #e43333;
}
</style>
