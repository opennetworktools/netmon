<template>
    <div class="chart-container">
        <svg class="line-chart" :viewBox="viewBox">
        </svg>
    </div>
</template>

<script>
import * as d3 from "d3"
import packetsData from '../assets/packets.json'

export default {
    name: 'RealTimeChart',
    props: {
        width: {
            default: 1200,
            type: Number,
        },
        height: {
            default: 400,
            type: Number,
        }
    },
    data() {
        return {
            padding: 20,
            chartData: packetsData.packets
        };
    },
    computed: {
        viewBox() {
            return `0 0 ${this.width} ${this.height}`
        }
    },
    mounted() {
        console.log(this.chartData)
        const width = this.width
        const height = this.height

        // const width = 928;
        // const height = 500;
        // const marginTop = 20;
        // const marginRight = 30;
        // const marginBottom = 30;
        // const marginLeft = 40;

        const margin = 10;

        const data = this.chartData

        const svg = d3.select("svg").attr("width", width).attr("height", height)
        const g = svg.append("g")

        //2. Parse the dates
        // const parseTime = d3.timeParse("%d-%b-%y");

        const timestampFormat = "%y-%b-%d";
        const parseTime = d3.timeParse(timestampFormat)

        //3. Creating the Chart Axes
        const x = d3
        .scaleTime()
        .domain(
            d3.extent(data, function (d) {
                return new Date(d.Timestamp)
            })
        )
        .rangeRound([0, width]);

        const y = d3
        .scaleLinear()
        .domain(
            d3.extent(data, function (d) {
                return d.CaptureLength
            })
        )
        .rangeRound([height, 0])

        //4. Creating a Line
        const line = d3
        .line()
        .x(function (d) {
            return x(new Date(d.Timestamp))
        })
        .y(function (d) {
            return y(d.CaptureLength)
        });

        let xValues = this.chartData.map((item) => new Date(item.Timestamp))
        console.log(xValues)
   

        //5. Appending the Axes to the Chart
        g.append("g")
        .call(d3.axisBottom(x))
        .attr("transform", `translate(0,${height - 20})`)
        
        g.append("g")
        .call(d3.axisLeft(y))
        .attr("transform", `translate(${20},0)`)  
        .call(g => g.append("text")
          .attr("x", -20)
          .attr("y", 20)
          .attr("fill", "currentColor")
          .attr("text-anchor", "start")
          .text("Price"));

        //6. Appending a path to the Chart
        g.append("path")
        .datum(data)
        .attr("fill", "none")
        .attr("stroke", "red")
        .attr("stroke-width", 1.5)
        .attr("d", line);
  },
}
</script>

<style scoped>
    .chart-container {
        margin-top: 2rem;
        height: 100%;
        max-height: 100%;
    }

    /* .line-chart {
        width: 600;
        height: 400;
        max-width: 100%; 
        height: auto; 
        height: intrinsic;
    } */

    /* .line-chart-g {
        fill: #000;
    }
    
    .line-chart-line {
        fill: none;
        stroke: red;
        stroke-width: 1.5;
    } */
</style>