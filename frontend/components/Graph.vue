<script>
import { Line } from "vue-chartjs";
import axios from "axios";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  CategoryScale,
} from "chart.js";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  CategoryScale,
);

export default {
  props: ["filename", "param", "num"],
  name: "SimpleChart",
  components: { Line },
  data() {
    return {
      values: [],
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          x: {
            grid: { color: "#222" },
            ticks: { color: "#999", font: { size: 9 } },
          },
          y: {
            grid: { color: "#222" },
            ticks: { color: "#999", font: { size: 9 } },
          },
        },
      },
    };
  },
  computed: {
    chartData() {
      return {
        labels: this.values.map((item) => item.timestamp),
        datasets: [
          {
            label: this.param,
            backgroundColor: "#eeeeee",
            borderColor: "#eeeeee",
            data: this.values.map((item) => item[this.param]),
            tension: 0.3,
            pointRadius: 2,
            borderWidth: 2,
          },
        ],
      };
    },
  },
  methods: {
    // Change this
    async handleFetchValues() {
      try {
        const res = await axios.get(
          `http://localhost:8000/tail/?filename=${this.filename}&num=${this.num}&param=${this.param}`,
        );
        this.values = res.data.values;
      } catch (err) {
        console.log(err);
      }
    },
  },
  async mounted() {
    await this.handleFetchValues();
  },
};
</script>

<template>
  <div class="chart-card">
    <h2 class="chart-title">{{ param }}</h2>
    <div class="chart-container">
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<style scoped>
.chart-card {
  background: #1a1a1a;
  border: 1px solid #333;
  padding: 0.75rem;
  width: 50%;
}

.chart-title {
  color: #888;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.05rem;
  margin: 0 0 0.5rem 0;
  font-family: sans-serif;
}

.chart-container {
  position: relative;

  height: 300px;
}

@media (max-width: 600px) {
  .chart-container {
    height: 150px;
  }
  .chart-card {
    width: 90%;
  }
}
</style>
