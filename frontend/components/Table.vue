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
          if (k !== "id") {
            value[k] = v;
          }
        });

        this.valuesArray.push(value);

        if (this.valuesArray.length > 100) {
          this.valuesArray.shift();
        }
      },
    },
  },
};
</script>

<template>
  <div class="stream-card">
    <div class="card-header">
      <span class="header-title">{{ deviceName }}</span>
      <button
        @click="$emit('delete')"
        class="minimal-delete"
        title="Remove Stream"
      >
        &times;
      </button>
    </div>

    <div class="table-scroll-area">
      <table v-if="valuesArray.length">
        <thead>
          <tr>
            <th v-for="(value, key) in valuesArray[0]" :key="key">
              {{ key }}
            </th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="(row, index) in valuesArray.slice().reverse()"
            :key="index"
          >
            <td v-for="(cell, key) in row" :key="key">
              {{ cell }}
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="waiting-state">Waiting for data...</div>
    </div>
  </div>
</template>

<style scoped>
.stream-card {
  background-color: #1a1a1a;
  border: 1px solid #333;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #333;
}

.header-title {
  font-size: 1rem;
  font-weight: 600;
  color: #eeeeee;
}

.minimal-delete {
  background: none;
  border: none;
  color: #666;
  font-size: 1.4rem;
  cursor: pointer;
  line-height: 1;
}

.minimal-delete:hover {
  color: #ff4444;
}

.table-scroll-area {
  max-height: 300px;
  overflow-y: auto;
  background-color: #1a1a1a;
}

table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

thead {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: #252525;
}

th {
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: #888;
  border-bottom: 1px solid #333;
}

td {
  padding: 0.7rem 1rem;
  font-size: 0.95rem;
  color: #eeeeee;
  border-bottom: 1px solid #252525;
  font-family: monospace;
}

tbody tr:hover {
  background-color: #222222;
}

.waiting-state {
  padding: 3rem 1rem;
  text-align: center;
  color: #555;
  font-size: 0.9rem;
}

.table-scroll-area::-webkit-scrollbar {
  width: 6px;
}
.table-scroll-area::-webkit-scrollbar-track {
  background: #1a1a1a;
}
.table-scroll-area::-webkit-scrollbar-thumb {
  background: #333;
}
</style>
