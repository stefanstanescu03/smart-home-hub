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
      <div class="title">
        <h3>{{ deviceName }}</h3>
      </div>

      <button
        @click="$emit('delete')"
        class="icon-button"
        title="Delete stream"
      >
        ✕
      </button>
    </div>

    <div class="table-container">
      <table v-if="valuesArray.length">
        <thead>
          <tr>
            <th v-for="(value, key) in valuesArray[0]" :key="key">
              {{ key }}
            </th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="(row, index) in valuesArray" :key="index">
            <td v-for="(cell, key) in row" :key="key">
              {{ cell }}
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="empty">Waiting for data…</div>
    </div>
  </div>
</template>

<style scoped>
.stream-card {
  background: #1b1b1b;
  border: 1px solid #2a2a2a;
  border-radius: 0.9rem;
  overflow: hidden;
  width: 100%;
  max-width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 0.75rem 1rem;
  border-bottom: 1px solid #2a2a2a;
}

.title h3 {
  margin: 0;
  font-size: large;
  font-weight: 600;
  color: #fff;
}

.subtitle {
  font-size: 0.7rem;
  color: #888;
}

/* Delete button */
.icon-button {
  background: transparent;
  border: 1px solid #333;
  color: #aaa;
  border-radius: 0.4rem;
  cursor: pointer;
  padding: 0.25rem 0.45rem;
  font-size: 0.8rem;

  transition:
    color 150ms ease,
    border-color 150ms ease,
    background 150ms ease;
}

.icon-button:hover {
  color: #e55353;
  border-color: #e55353;
  background: rgba(229, 83, 83, 0.12);
}

.table-container {
  max-height: 260px;
  overflow-y: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: normal;
}

thead {
  position: sticky;
  top: 0;
  z-index: 2;
}

th {
  background: #202020;
  color: #eeeeee;
  font-weight: bold;
  text-transform: uppercase;
  font-size: normal;
  padding: 0.55rem 0.6rem;
  border-bottom: 1px solid #2a2a2a;
}

td {
  padding: 0.45rem 0.6rem;
  border-bottom: 1px solid #262626;
  color: #eeeeee;
}

tbody tr:hover {
  background: #202020;
}

.empty {
  padding: 1rem;
  text-align: center;
  font-size: normal;
  color: #777;
}
</style>
