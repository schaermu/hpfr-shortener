<script lang="ts">
    import { chart } from '../lib/apex-wrapper'

    export let code: string;
    export let location;
    
    let options = {
        chart: {
        type: "line",
        },
        series: [
        {
            name: "hits",
            data: [30, 40, 35, 50, 49, 60, 70, 91, 125],
        },
        ],
        xaxis: {
        categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999],
        },
    };
</script>

<div>
    <div use:chart={options} />
</div>

This site or product includes IP2Location LITE data available from <a href="https://lite.ip2location.com">https://lite.ip2location.com</a>.