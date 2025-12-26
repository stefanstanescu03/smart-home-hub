<script>
import axios from "axios";
export default {
  props: ["deviceName", "ident", "visibility", "should_appear", "id"],
  data() {
    return {
      status: "checking...",
    };
  },
  methods: {
    async checkStatus() {
      try {
        await axios
          .get(`http://localhost:5000/device/ping/${this.ident}`)
          .then((response) => {
            if (response.data.status == true) {
              this.status = "connected";
            } else {
              this.status = "disconnected";
            }
          });
      } catch (err) {
        this.status = "disconnected";
      }
    },
  },
  mounted() {
    this.checkStatus();
  },
};
</script>

<template>
  <tr class="device-container">
    <td>{{ deviceName }}</td>
    <td>{{ ident }}</td>
    <td v-if="visibility == false">private</td>
    <td v-if="visibility == true">public</td>
    <td>
      <span class="disconnected-text" v-if="this.status == 'disconnected'">
        {{ status }}
      </span>
      <span class="connected-text" v-if="this.status == 'connected'">
        {{ status }}
      </span>
    </td>
    <td>
      <div class="action-container" v-if="should_appear == true">
        <button
          class="delete-button"
          @click="this.$router.push(`/alerts/${id}`)"
        >
          <img src="../public/bell.png" alt="" height="20" width="20" />
        </button>
        <button class="delete-button" @click="$emit('edit')">
          <img src="../public/edit.png" alt="" height="25" width="25" />
        </button>
        <button class="delete-button" @click="$emit('delete')">
          <img src="../public/delete.png" alt="" height="20" width="20" />
        </button>
      </div>
    </td>
  </tr>
</template>

<style scoped>
td {
  border-bottom: 1px solid #eeeeee;
  padding: 0.7rem;
}

.delete-button {
  border: none;
  text-decoration: none;
  background-color: transparent;
  cursor: pointer;
}

.action-container {
  display: flex;
  flex-direction: row;
  gap: 0.5rem;
  align-items: center;
}

.disconnected-text {
  display: inline-block;
  font-size: medium;
  padding: 0.3rem 0.6rem;
  border: 1px solid #e43333;
  border-radius: 1rem;
  color: #e43333;
}

.connected-text {
  display: inline-block;
  font-size: medium;
  padding: 0.3rem 0.6rem;
  border: 1px solid #2ba618;
  border-radius: 1rem;
  color: #2ba618;
}
</style>
