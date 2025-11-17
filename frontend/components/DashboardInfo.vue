<script>
export default {
  props: ["name", "visibility", "should_appear", "id"],
};
</script>

<template>
  <tr class="dashboard-container">
    <td>{{ name }}</td>
    <td v-if="visibility == false">private</td>
    <td v-if="visibility == true">public</td>
    <td>
      <div class="action-container" v-if="should_appear == true">
        <button class="delete-button" @click="$emit('edit')">
          <img src="../public/edit.png" alt="" height="25" width="25" />
        </button>
        <button class="delete-button" @click="$emit('delete')">
          <img src="../public/delete.png" alt="" height="20" width="20" />
        </button>
        <button
          class="view-button"
          @click="this.$router.push(`/dashboard/${id}`)"
        >
          View
        </button>
      </div>
    </td>
  </tr>
</template>

<style scoped>
td {
  border-bottom: 1px solid #c9c9c9;
  padding: 0.3rem;
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

.view-button {
  border: none;
  background-color: transparent;
  border: 1px solid #a6a6a6;
  text-decoration: none;
  cursor: pointer;
  color: #121212;
  padding: 0.5rem;
  border-radius: 0.3rem;
  transition-duration: 300ms;
}

.view-button:hover {
  border: 1px solid #121212;
}
</style>
