window.onload = function () {
    let conn;
    let isOpen = false;
    const networkElement = document.getElementById("network-usage")
    const cpuUsageChart = new Chart(document.getElementById('cpu-usage'), {
        type: 'bar',
        data: {
            labels: ['Core 0', 'Core 1', 'Core 2', 'Core 3'],
            datasets: [{
                label: 'Usage in %',
                data: [0, 0, 0, 0],
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true,
                    min: 0,
                    max: 100,
                }
            }
        }
    });
    const memoryUsageChart = new Chart(document.getElementById('memory-usage'), {
        type: 'pie',
        data: {
            labels: [
                'Used',
                'Free',
                'Buff',
                'Shared'
            ],
            datasets: [{
                label: 'Memory Usage',
                data: [0, 0, 0, 0],
                backgroundColor: [
                    `#FA8F82`,
                    `#FAC869`,
                    `#AD8840`,
                    `#AD6B63`,
                ],
                hoverOffset: 4
            }]
        }
    });

    const cpuStat = {
        prevStat: [],
        calculatePercentageNow(stats) {
            return stats.map((v, i) => {
                let percentage = 0;
                if (cpuStat.prevStat[i]) {
                    // v.idle
                    // v.total
                    const prev = cpuStat.prevStat[i];
                    const totalD = v.total - prev.total;
                    const idleD = v.idle - prev.idle;
                    percentage = ((totalD - idleD) / totalD) * 100;
                }

                cpuStat.prevStat[i] = v;
                return percentage;
            })
        }
    };
    const networkStat = {
        prevStat: {tx: 0, rx: 0, t: 0},
        calculate: (tx, rx) => {
            const now = Math.ceil(Date.now() / 1000);
            const out = {rtx: 0, rrx: 0}
            if (networkStat.prevStat.t !== 0) {
                const td = now - networkStat.prevStat.t;
                out.rtx = Math.ceil((tx - networkStat.prevStat.tx) / td);
                out.rrx = Math.ceil((rx - networkStat.prevStat.rx) / td);
            }
            networkStat.prevStat.tx = tx;
            networkStat.prevStat.rx = rx;
            networkStat.prevStat.t = now;
            return out;
        }
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws/stats");
        conn.onclose = function (evt) {
            console.log("connection closed", evt);
            alert("Disconected")
        };
        conn.onmessage = function (evt) {
            const data = JSON.parse(evt.data)
            const netStat = networkStat.calculate(data.net_tx, data.net_rx)
            cpuUsageChart.data.labels = data.cpuStats.map((v, i) => `CPU #${i}`);
            cpuUsageChart.data.datasets[0].data = cpuStat.calculatePercentageNow(data.cpuStats);
            memoryUsageChart.data.datasets[0].data = [data.mem_used, data.mem_free, data.mem_buff, data.mem_shared];
            cpuUsageChart.update();
            memoryUsageChart.update();
            networkElement.textContent = `Download ${netStat.rrx} kB/s | Upload: ${netStat.rtx} kB/s`;
        };
        conn.onopen = function (evt) {
            console.log("connection opened", evt);
            isOpen = true;
        }
    } else {
        alert("browser doesn't support websocket");
    }
}
