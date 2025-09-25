const ctx = document.getElementById("networkChart").getContext("2d");
const networkChart = new Chart(ctx, {
  type: "line",
  data: {
    labels: Array(30).fill(""),
    datasets: [
      {
        label: "In (Mb/s)",
        data: Array(30).fill(0),
        borderColor: "rgba(0, 123, 255, 0.8)",
        backgroundColor: "rgba(0, 123, 255, 0.1)",
        borderWidth: 2,
        pointRadius: 0,
        tension: 0.4,
        fill: true,
      },
      {
        label: "Out (Mb/s)",
        data: Array(30).fill(0),
        borderColor: "rgba(40, 167, 69, 0.8)",
        backgroundColor: "rgba(40, 167, 69, 0.1)",
        borderWidth: 2,
        pointRadius: 0,
        tension: 0.4,
        fill: true,
      },
    ],
  },
  options: {
    maintainAspectRatio: false,
    scales: {
      y: {
        beginAtZero: true,
        ticks: { font: { size: 10 } },
      },
      x: {
        ticks: { display: false },
      },
    },
    plugins: {
      legend: { display: false },
    },
    animation: { duration: 250 },
  },
});

const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
const socket = new WebSocket(`${protocol}//${window.location.host}/ws`);

socket.onmessage = function (event) {
  const data = JSON.parse(event.data);

  function updateTextContent(selector, value) {
    const element = document.querySelector(selector);
    if (element) {
      element.textContent = value;
    }
  }

  function updateProgressBar(selector, value) {
    const element = document.querySelector(selector);
    if (element) {
      element.style.width = value;
    }
  }

  updateTextContent("#cpu-usage", data.cpu.usage);
  updateTextContent("#cpu-temp", data.cpu.temp);
  const cpuProgressBar = document.querySelector("#cpu-progress-bar");
  if (cpuProgressBar) {
    cpuProgressBar.style.width = data.cpu.usage;
  }

  updateTextContent("#ram-usage-gb", `${data.ram.used} / ${data.ram.total}`);
  updateTextContent("#ram-usage-percent", `(${data.ram.percent})`);
  updateProgressBar("#ram-progress-bar", data.ram.percent);

  updateTextContent("#docker-containers", data.docker.containers);
  updateTextContent("#docker-images", data.docker.images);
  updateTextContent("#docker-volumes", data.docker.volumes);

  updateTextContent(
    "#streaming-card .library-stat:nth-child(1) .count",
    data.streaming.films,
  );
  updateTextContent(
    "#streaming-card .library-stat:nth-child(2) .count",
    data.streaming.series,
  );
  updateTextContent(
    "#streaming-card .library-stat:nth-child(3) .count",
    data.streaming.animes,
  );

  const playingElement = document.querySelector(
    ".streaming-activity .activity-stat:nth-child(1)",
  );
  if (playingElement) {
    playingElement.innerHTML = `<i class="fas fa-play-circle"></i> ${data.streaming.playing} Media currently playing`;
  }

  const transcodingElement = document.querySelector(
    ".streaming-activity .activity-stat:nth-child(2)",
  );
  if (transcodingElement) {
    transcodingElement.innerHTML = `<i class="fas fa-cogs"></i> ${data.streaming.transcoding} Active transcoding`;
  }

  const arcCircle = document.querySelector(".cache-chart .chart-circle");
  if (arcCircle) {
    arcCircle.setAttribute(
      "stroke-dasharray",
      `${data.arcCache.arcHitRateNum}, 100`,
    );
  }
  updateTextContent(".cache-chart .chart-text", data.arcCache.arcHitRate);
  updateTextContent(
    ".cache-details .cache-stat:nth-child(1) .stat-value",
    data.arcCache.arcSize,
  );
  updateTextContent(
    ".cache-details .cache-stat:nth-child(2) .stat-value",
    data.arcCache.arcMaxSize,
  );
  updateTextContent(
    ".cache-details .cache-stat:nth-child(4) .stat-value",
    data.arcCache.l2arcSize,
  );
  updateTextContent(
    ".cache-details .cache-stat:nth-child(5) .stat-value",
    data.arcCache.l2arcHitRate,
  );

  updateTextContent(
    ".sysinfo-list li:nth-child(1) span:last-child",
    data.system.os,
  );
  updateTextContent(
    ".sysinfo-list li:nth-child(2) span:last-child",
    data.system.kernel,
  );
  updateTextContent(
    ".sysinfo-list li:nth-child(3) span:last-child",
    data.system.cpu,
  );

  if (data.zfsConfig && data.zfsConfig.poolName) {
    updateTextContent("#zfs-card .pool-name", data.zfsConfig.poolName);
    const statusElement = document.querySelector("#zfs-card .pool-status");
    if (statusElement) {
      statusElement.textContent = data.zfsConfig.poolStatus;
      statusElement.className =
        "pool-status " + data.zfsConfig.poolStatus.toLowerCase();
    }

    const dataVdevTree = document.querySelector(
      "#zfs-card .zfs-column:nth-child(1) .zfs-tree",
    );
    const cacheVdevTree = document.querySelector(
      "#zfs-card .zfs-column:nth-child(2) .zfs-tree",
    );

    if (dataVdevTree && cacheVdevTree) {
      dataVdevTree.innerHTML = "";
      cacheVdevTree.innerHTML = "";

      data.zfsConfig.dataVdevs.forEach((vdev) => {
        const vdevLi = document.createElement("li");
        let devicesHtml = "<ul>";
        vdev.devices.forEach((device) => {
          devicesHtml += `<li><i class="fas fa-hdd"></i> ${device}</li>`;
        });
        devicesHtml += "</ul>";
        vdevLi.innerHTML = `<i class="fas fa-layer-group"></i> ${vdev.name} <span class="device-status">${vdev.status}</span>${devicesHtml}`;
        dataVdevTree.appendChild(vdevLi);
      });

      if (data.zfsConfig.cacheVdev) {
        const vdev = data.zfsConfig.cacheVdev;
        const vdevLi = document.createElement("li");
        let devicesHtml = "<ul>";
        vdev.devices.forEach((device) => {
          devicesHtml += `<li><i class="fas fa-microchip"></i> ${device}</li>`;
        });
        devicesHtml += "</ul>";
        vdevLi.innerHTML = `<i class="fas fa-bolt"></i> ${vdev.name} <span class="device-status">${vdev.status}</span>${devicesHtml}`;
        cacheVdevTree.appendChild(vdevLi);
      }
    }
  }

  updateTextContent("#cpu-usage", data.cpu.usage);
  updateProgressBar(
    "#cpu-usage + .progress-bar-container .progress-bar",
    data.cpu.usage,
  );

  updateTextContent("#cpu-temp", data.cpu.temp);
  const tempGauge = document.querySelector(".gauge-fill");
  if (tempGauge) {
    const tempRotation = (data.cpu.tempDeg / 100) * 180;
    tempGauge.style.transform = `rotate(${tempRotation}deg)`;
  }

  updateTextContent("#net-in", data.net.in);
  updateTextContent("#net-out", data.net.out);

  const chartData = networkChart.data;

  chartData.datasets[0].data.push(parseFloat(data.net.in) || 0);
  chartData.datasets[0].data.shift();

  chartData.datasets[1].data.push(parseFloat(data.net.out) || 0);
  chartData.datasets[1].data.shift();

  networkChart.update();

  updateTextContent("#disk-mount-point", data.disk.mountPoint);
  updateTextContent("#disk-usage-tb", `${data.disk.used} / ${data.disk.total}`);
  updateTextContent("#disk-usage-percent", `(${data.disk.percent})`);
  updateProgressBar("#disk-progress-bar", data.disk.percent);
};

socket.onclose = function (event) {};

socket.onerror = function (error) {};
