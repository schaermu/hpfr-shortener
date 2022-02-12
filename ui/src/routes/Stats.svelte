<script lang="ts">
  import ApexCharts from 'apexcharts';
  import { onMount } from 'svelte';
  import ApiClient from '../lib/api';

  window['ApexCharts'] = ApexCharts;

  export let code: string;

  let chartDom;

  const client = new ApiClient();
  let options = {
    chart: {
      type: 'area',
      stacked: false,
      height: 350,
      zoom: {
        type: 'x',
        enabled: true,
        autoScaleYaxis: true,
      },
      toolbar: {
        autoSelected: 'zoom',
      },
    },
    title: {
      text: 'Access Statistics',
      align: 'left',
    },
    fill: {
      type: 'gradient',
      gradient: {
        shadeIntensity: 1,
        inverseColors: false,
        opacityFrom: 0.5,
        opacityTo: 0,
        stops: [0, 90, 100],
      },
    },
    xaxis: {
      type: 'datetime',
    },
    tooltip: {
      shared: false,
      y: {
        formatter: function (val) {
          return val.toFixed(0);
        },
      },
    },
    series: [],
    noData: {
      text: 'Loading...',
    },
  };

  onMount(() => {
    let chart = new ApexCharts(chartDom, options);
    chart.render();

    client.getStatistics(code.replace('+', '')).then((stats) => {
      chart.updateSeries(
        [
          {
            name: 'Hits',
            data: stats.hitTimeData,
          },
        ],
        true
      );
    });
  });
</script>

<div>
  <div bind:this={chartDom} id="chart" />
</div>

This site or product includes IP2Location LITE data available from<a
  href="https://lite.ip2location.com">https://lite.ip2location.com</a
>.
