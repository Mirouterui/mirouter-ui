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

$("#doit").click(function() {
    new_name = $("#new_name").val();
    if (new_name != "") {
        url = '/api/xqsystem/set_device_nickname'
        postdata = {
            "mac": mac,
            "name": new_name
        }
        $.get(url, postdata, function(data) {
            if (data.code == 0) {
                mdui.snackbar({
                    message: '修改成功'
                });
            } else {
                mdui.snackbar({
                    message: data.msg
                });
            }
        });
    } else {
        mdui.snackbar({
            message: '请输入新的名称'
        });
    }
});
// 初次加载状态
getdeviceinfo();