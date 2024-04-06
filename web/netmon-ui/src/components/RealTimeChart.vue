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
        },
        packets: {
            type: Array
        }
    },
    data() {
        return {
            padding: 20,
            chartData: packetsData.packets,
            marginTop: 20,
            marginRight: 30,
            marginBottom: 30,
            marginLeft: 40,
        };
    },
    computed: {
        viewBox() {
            return `0 0 ${this.width} ${this.height}`
        }
    },
    created() {
    },
    mounted() {
        const data = this.chartData
        const width = this.width
        const height = this.height

        const svg = d3.select("svg").attr("width", width).attr("height", height)
        const g = svg.append("g")

        const x = d3.scaleTime().domain(
            d3.extent(data, function (d) {
                return new Date(d.Timestamp)
            })
        ).rangeRound([0, width])
        const y = d3.scaleLinear().domain(
            d3.extent(data, function (d) {
                return d.CaptureLength
            })
        ).rangeRound([height, 0])

        const line = d3.line().x(function (d) {
            return x(new Date(d.Timestamp))
        }).y(function (d) {
            return y(d.CaptureLength)
        });

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
                .text("Price"))

        g.append("path")
            .datum(data)
            .attr("fill", "none")
            .attr("stroke", "red")
            .attr("stroke-width", 1.5)
            .attr("d", line)

        console.log(this.packets)
    },
}
</script>

<style scoped>
.chart-container {
    margin-top: 2rem;
}

.line-chart {
    width: 1200;
    height: 400;
    max-width: 100%;
    height: auto;
    height: intrinsic;
}
</style>