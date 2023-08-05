// 获取查询字符串
var queryString = window.location.search;
// 去掉第一个问号
queryString = queryString.substring(1);
// 用等号分割查询字符串
var queryArray = queryString.split("=");
// 获取MAC地址
var mac = queryArray[1];
if (mac) {
    $('#mac').text(mac);
} else {
    mdui.snackbar({
        message: '没有MAC地址😅'
    });
}
upspeed_data = []
downspeed_data = []
uptraffic_data = []
downtraffic_data = []

function updateStatus() {
    $.get('/api/misystem/status', function(data) {
        // upspeed = convertSpeed(data.wan.upspeed)
        // maxuploadspeed = convertSpeed(data.wan.maxuploadspeed)
        // downspeed = convertSpeed(data.wan.downspeed)
        // maxdownloadspeed = convertSpeed(data.wan.maxdownloadspeed)
        // uploadtotal = convertbytes(data.wan.upload)
        // downloadtotal = convertbytes(data.wan.download)
        // cpuload = roundToOneDecimal(data.cpu.load * 100) + '%'
        // memusage = roundToOneDecimal(data.mem.usage * 100) + '%'
        // $('#platform').text("小米路由器" + data.hardware.platform);
        // $('#mac').text(data.hardware.mac);
        // $('#cpu-used .mdui-progress-determinate').css('width', cpuload);
        // $('#cpu-used-text').text(cpuload);
        // $('#mem-used .mdui-progress-determinate').css('width', memusage);
        // $('#mem-used-text').text(memusage);
        dev = data.dev
        for (var i = 0; i < dev.length; i++) {
            //获取当前设备对象
            var device = dev[i];
            if (device.mac == mac) {
                upspeed = convertSpeed(device.upspeed)
                maxuploadspeed = convertSpeed(device.maxuploadspeed)
                downspeed = convertSpeed(device.downspeed)
                maxdownloadspeed = convertSpeed(device.maxdownloadspeed)
                uploadtotal = convertbytes(device.upload)
                downloadtotal = convertbytes(device.download)
                onlinetime = convertSeconds(device.online)
                $('#uploadspeed').text(upspeed)
                $('#maxuploadspeed').text(maxuploadspeed)
                $('#downloadspeed').text(downspeed)
                $('#maxdownloadspeed').text(maxdownloadspeed)
                $('#uploadtotal').text(uploadtotal)
                $('#downloadtotal').text(downloadtotal)
                $('#onlinetime').text(onlinetime)
                var upspeed = (device.upspeed / 1024 / 1024).toFixed(2);
                var downspeed = (device.downspeed / 1024 / 1024).toFixed(2);
                var uploadtotal = togb(device.upload)
                var downloadtotal = togb(device.download);
                upspeed_data.push(upspeed);
                downspeed_data.push(downspeed);
                uptraffic_data.push(uploadtotal);
                downtraffic_data.push(downloadtotal);
                // 调用drawChart函数，绘制图表
                drawspeedChart();
                drawtrafficChart();
            }
        }
    });
}

function getdeviceinfo() {
    $.get('/api/misystem/devicelist', function(data) {
        dev = data.list
        for (var i = 0; i < dev.length; i++) {
            //获取当前设备对象
            var device = dev[i];
            if (device.mac == mac) {
                if (device.icon != "") {
                    iconurl = "/img/" + device.icon
                } else {
                    iconurl = "/img/device_list_unknow.png"
                }
                $("#devicename").text(device.name);
                $("#deviceicon").attr("src", iconurl);
                $("#device_oname").text(device.oname);
                $("#ipaddress").text(device.ip[0].ip); //不管多少个ip地址，只显示第一个
                $("#authority_wan").text(getbooleantype(device.authority.wan));
                $("#authority_lan").text(getbooleantype(device.authority.lan));
                $("#authority_admin").text(getbooleantype(device.authority.admin));
                $("#authority_pridisk").text(getbooleantype(device.authority.pridisk));
                $("#connecttype").text(getconnecttype(device.type));
                $("#isap").text(getbooleantype(device.isap));
                $("#isonline").text(getbooleantype(device.online));

                var match = true
                if (device.mac == data.mac) {
                    $("#devicename").text($("#devicename").text() + " (页面所在地)");
                }
                console.log(device)
            }
        }
        if (match != true) {
            mdui.snackbar({
                message: '好像没有这个设备呢😢'
            });
        }


    });
}

function drawspeedChart() {
    // 获取div元素，用于放置图表
    var chart = document.getElementById("speed-chart");
    // 初始化echarts实例
    var myChart = echarts.init(chart);
    // 定义图表的配置项和数据
    var option = {
        tooltip: {
            trigger: "axis",
        },
        xAxis: {
            type: "category",
            data: upspeed_data.map(function(item, index) {
                return (index + 1) * 5 + "s"; // 返回请求次数作为横坐标
            }),
        },
        yAxis: {
            type: "value",
            name: "网络速度（MB/s）",
        },
        legend: {
            orient: 'vertical',
            left: 'right'
        },
        series: [{
                name: "上传速度（MB/s）",
                type: "line",
                data: upspeed_data, // 返回网络速度（MB/s）作为纵坐标
            },
            {
                name: "下载速度（MB/s）",
                type: "line",
                data: downspeed_data, // 返回网络速度（MB/s）作为纵坐标
            },
        ],
    };
    // 设置图表的配置项和数据
    myChart.setOption(option);
}

function drawtrafficChart() {
    // 获取div元素，用于放置图表
    var chart = document.getElementById("traffic-chart");
    // 初始化echarts实例
    var myChart = echarts.init(chart);
    // 定义图表的配置项和数据
    var option = {
        tooltip: {
            trigger: "axis",
        },
        xAxis: {
            type: "category",
            data: upspeed_data.map(function(item, index) {
                return (index + 1) * 5 + "s"; // 返回请求次数作为横坐标
            }),
        },
        legend: {
            orient: 'vertical',
            left: 'right'
        },
        yAxis: {
            type: "value",
            name: "上传/下载（GB）",
        },
        series: [{
                name: "上传总量（GB）",
                type: "line",
                data: uptraffic_data, // 返回网络速度（MB/s）作为纵坐标
            },
            {
                name: "下载总量（GB）",
                type: "line",
                data: downtraffic_data, // 返回网络速度（MB/s）作为纵坐标
            },
        ],
    };
    // 设置图表的配置项和数据
    myChart.setOption(option);
}

function getconnecttype(type) {
    // 0/1/2/3  有线 / 2.4G wifi / 5G wifi / guest wifi
    switch (type) {
        case 0:
            return "有线连接";
        case 1:
            return "2.4G wifi";
        case 2:
            return "5G wifi";
        case 3:
            return "guest wifi";
        default:
            return "未知";
    }
}

function get_router_name() {
    $.get('/api/xqsystem/router_name', function(data) {
        if (data.code === 0) {
            router_name = data.routerName
            $("#router_name").text(router_name)
        }
    });
}
$(function() {
    // 初次加载状态
    getdeviceinfo();
    get_router_name();
    updateStatus();
    // 每5秒刷新状态
    setInterval(function() {
        updateStatus();
    }, 5000);
});