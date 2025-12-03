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
          if (k != "timestamp" && k != "id") {
            value[k] = v;
          }
        });
        this.valuesArray.push(value);
        console.log(this.valuesArray);
      },
    },
  },
};
</script>

<template>
  <h1>to be made</h1>
</template>
