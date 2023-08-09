// var download_traffic_data = [];
// var upload_traffic_data = [];

function updateStatus() {
    $.get('/api/misystem/status', function(data) {
        upspeed = convertSpeed(data.wan.upspeed)
        maxuploadspeed = convertSpeed(data.wan.maxuploadspeed)
        downspeed = convertSpeed(data.wan.downspeed)
        maxdownloadspeed = convertSpeed(data.wan.maxdownloadspeed)
        uploadtotal = convertbytes(data.wan.upload)
        downloadtotal = convertbytes(data.wan.download)
        cpuload = roundToOneDecimal(data.cpu.load * 100) + '%'
        memusage = roundToOneDecimal(data.mem.usage * 100) + '%'
        devicenum = data.count.all
        devicenum_now = data.count.online
        $('#platform').text("小米路由器" + data.hardware.platform);
        $('#cpu-used .mdui-progress-determinate').css('width', cpuload);
        $('#cpu-used-text').text(cpuload);
        $('#mem-used .mdui-progress-determinate').css('width', memusage);
        $('#mem-used-text').text(memusage);
        $('#uploadspeed').text(upspeed)
        $('#maxuploadspeed').text(maxuploadspeed)
        $('#downloadspeed').text(downspeed)
        $('#maxdownloadspeed').text(maxdownloadspeed)
        $('#uploadtotal').text(uploadtotal)
        $('#downloadtotal').text(downloadtotal)
        $("#devicenum").text(devicenum)
        $("#devicenum_now").text(devicenum_now)
        listdevices(data.dev)
    });
}

function get_messages() {
    $.get('/api/misystem/messages', function(data) {
        if (data.code != 0) {
            mdui.snackbar({
                message: '路由器有新信息，请登录路由器后台查看'
            });
        }
    });
}

function get_router_name() {
    $.get('/api/xqsystem/router_name', function(data) {
        if (data.code === 0) {
            router_name = data.routerName
            $("#router_name").text(router_name)
        }
    });
}

function check_internet_connect() {
    $.get('/api/xqsystem/internet_connect', function(data) {
        if (data.connect === 1) {
            mdui.snackbar({
                message: '路由器好像没联网呢😢'
            });

        }
    });
}
$(function() {
    // 初次加载状态
    updateStatus();
    check_internet_connect();
    get_router_name();
    get_messages();
    // 每5秒刷新状态
    setInterval(function() {
        updateStatus();
    }, 5000);
});



function listdevices(dev) {
    //获取已有的表格元素
    var table = document.querySelector("table");

    //获取表格内容区域
    var tbody = document.getElementById("device-list");

    //清空表格内容区域
    tbody.innerHTML = "";
    //遍历dev数组，创建表格内容行
    for (var i = 0; i < dev.length; i++) {
        //获取当前设备对象
        var device = dev[i];

        //创建内容行
        var tr = document.createElement("tr");

        //创建内容单元格，并添加到内容行中
        var td_devname = document.createElement("td");
        var detail_url = "/device/index.html?mac=" + device.mac;
        td_devname.innerHTML = "<a href='" + detail_url + "'>" + device.devname + "</a>";
        tr.appendChild(td_devname);

        // var td_detail = document.createElement("td");
        // var detail_url = "/device/?mac=" + device.mac;
        // td_detail.innerHTML = "<a href='" + detail_url + "'>点我</a>";
        // tr.appendChild(td_detail);

        var td_downspeed = document.createElement("td");
        td_downspeed.textContent = convertSpeed(device.downspeed);
        tr.appendChild(td_downspeed);

        var td_upspeed = document.createElement("td");
        td_upspeed.textContent = convertSpeed(device.upspeed);
        tr.appendChild(td_upspeed);

        var td_uptotal = document.createElement("td");
        td_uptotal.textContent = convertbytes(device.upload);
        tr.appendChild(td_uptotal);

        var td_downtotal = document.createElement("td");
        td_downtotal.textContent = convertbytes(device.download);
        tr.appendChild(td_downtotal);

        // var td_maxdownspeed = document.createElement("td");
        // td_maxdownspeed.textContent = convertSpeed(device.maxdownloadspeed);
        // tr.appendChild(td_maxdownspeed);

        // var td_maxupspeed = document.createElement("td");
        // td_maxupspeed.textContent = convertSpeed(device.maxuploadspeed);
        // tr.appendChild(td_maxupspeed);

        // var td_onlinetime = document.createElement("td");
        // td_onlinetime.textContent = convertSeconds(device.online);
        // tr.appendChild(td_onlinetime);

        // var td_mac = document.createElement("td");
        // td_mac.textContent = device.mac;
        // tr.appendChild(td_mac);

        //将内容行添加到表格内容区域中
        tbody.appendChild(tr);

        // pushuptrafficdata(device.devname, togb(device.upload));
        // pushdowntrafficdata(device.devname, togb(device.download));
    }

    //更新表格元素
    table.appendChild(tbody);
    // drawtrafficChart();
}

function pushuptrafficdata(name, value) {
    data = {
        value: value,
        name: name
    }
    upload_traffic_data.push(data);

}

function pushdowntrafficdata(name, value) {
    data = {
        value: value,
        name: name
    }
    download_traffic_data.push(data);

}

function drawtrafficChart() {
    // 获取div元素，用于放置图表
    var chart = document.getElementById("traffic-chart");
    // 初始化echarts实例
    var myChart = echarts.init(chart);
    // 定义图表的配置项和数据
    var option = {
        tooltip: {
            trigger: 'item',
            confine: true
        },
        series: [{
                name: '上传流量(GB)',
                type: 'pie',
                radius: '50%',
                center: ['25%', '50%'],
                data: upload_traffic_data,
            },
            {
                name: '下载流量(GB)',
                type: 'pie',
                radius: '50%',
                center: ['75%', '50%'],
                data: download_traffic_data,
            }
        ],
    };
    // 设置图表的配置项和数据
    myChart.setOption(option);
    //清空数据
    upload_traffic_data = [];
    download_traffic_data = [];
}