//token检查
if (document.cookie.indexOf('token') < 0) {
    document.getElementById("fuction").style.display = "none";
    mdui.snackbar({
        message: 'Cookies中没有token！'
    });
    document.getElementById("wtf").style.display = "unset";
} else {
    document.getElementById("token").value = getCookie('token')
    $.ajax({
        url: 'https://mcweb-api.hzchu.top/admin/checktoken',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            token: getCookie('token')
        }),
        success: function(response) {
            var code = response.code;
            var name = response.name
            console.log(code);
            //var msg = response.msg;
            //console.log(msg);

            // 处理返回的code和msg
            if (code === 0) {
                mdui.snackbar({
                    message: "欢迎：" + name + "<br>token过期时间：" + response.endtime
                });
                document.cookie = "name=" + name;
            } else {
                document.getElementById("fuction").style.display = "none";
                document.getElementById("wtf").style.display = "unset";
                mdui.snackbar({
                    message: response.msg
                });
            }

        }
    });



}

function getCookie(name) {
    var value = "; " + document.cookie;
    var parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
}

document.getElementById("token-enter").onclick = function() {
        var inputText = document.getElementById("token").value;
        document.cookie = "token=" + inputText;
        location.reload();
    }
    //ban部分
    // 获取选择框元素
var ban_choice = document.getElementById("ban_choice");
var ban_reason = document.getElementById("ban_reason");
var ban_player = document.getElementById("ban_gamename")
    // 添加事件监听器，在选择框的值改变时触发
ban_choice.addEventListener('change', function() {
    // 获取当前选择的选项值
    var selectedValue = ban_choice.value;

    // 根据选项值执行相应的操作
    switch (selectedValue) {
        case '1':
            ban_reason.value = "使用非法的修改器、外挂或其他作弊手段来获取优势，如飞行、无限物品、自动攻击等。"
            break;
        case '2':
            ban_reason.value = "故意破坏其他玩家的建筑、物品或服务器设施，包括恶意放火、爆炸、挖掘等。"
            break;
        case '3':
            ban_reason.value = "在聊天频道或其他地方发布未经许可的广告、垃圾信息或恶意链接"
            break;
        case '4':
            ban_reason.value = "对其他玩家进行辱骂、威胁、骚扰或恶意攻击"
            break;
        case '5':
            ban_reason.value = "故意进行大规模的红石电路、崩服机或其他可能导致服务器崩溃或卡顿的行为。"
            break;
        case '6':
            ban_reason.value = "使用多个账号进行作弊、恶意破坏或其他违规行为。"
            break;
    }
});
document.getElementById("ban-enter").onclick = function() {
    var inst = new mdui.Dialog('#dialog');
    var name = getCookie("name")
    $.ajax({
        url: 'https://mcweb-api.hzchu.top/admin/banplayer',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            playername: ban_player.value,
            token: getCookie('token'),
            reason: ban_reason.value + "（由管理员：" + name + "操作）"
        }),
        success: function(response) {
            var code = response.code;
            console.log(code);

            // 处理返回的code和msg
            if (code === 0) {
                document.getElementById("dialog-title").innerText = "返回结果："
                document.getElementById("dialog-content").innerText = response.response
                inst.open();
            } else {
                document.getElementById("dialog-title").innerText = "处理失败！返回结果："
                document.getElementById("dialog-content").innerText = response.msg
                inst.open();

            }

        }
    })
}
document.getElementById("deban-enter").onclick = function() {
        var inst = new mdui.Dialog('#dialog');

        $.ajax({
            url: 'https://mcweb-api.hzchu.top/admin/debanplayer',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                playername: ban_player.value,
                token: getCookie('token')
            }),
            success: function(response) {
                var code = response.code;
                console.log(code);

                // 处理返回的code和msg
                if (code === 0) {
                    document.getElementById("dialog-title").innerText = "返回结果："
                    document.getElementById("dialog-content").innerText = response.response
                    inst.open();
                } else {
                    document.getElementById("dialog-title").innerText = "处理失败！返回结果："
                    document.getElementById("dialog-content").innerText = response.msg
                    inst.open();

                }

            }
        })
    }
    //电源调节部分
var powerplan_choice = document.getElementById("powerplan_choice");

document.getElementById("power-enter").onclick = function() {
    var inst = new mdui.Dialog('#dialog');
    var selectedValue = powerplan_choice.value;

    // 根据选项值执行相应的操作
    switch (selectedValue) {
        case '1':
            plan = "BALANCE"
            break;
        case '2':
            plan = "HIGH_PERFORMANCE"
            break;
        case '3':
            plan = "ENERGY_SAVER"
            break;
    }
    $.ajax({
        url: 'https://mcweb-api.hzchu.top/admin/changepower',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            plan: plan,
            token: getCookie('token')
        }),
        success: function(response) {
            var code = response.code;
            console.log(code);

            // 处理返回的code和msg
            if (code === 0) {
                document.getElementById("dialog-title").innerText = "返回结果："
                document.getElementById("dialog-content").innerText = response.msg
                inst.open();
            } else {
                document.getElementById("dialog-title").innerText = "处理失败！返回结果："
                document.getElementById("dialog-content").innerText = response.msg
                inst.open();

            }

        }
    })
}
document.getElementById("unpower-enter").onclick = function() {
    var inst = new mdui.Dialog('#dialog');

    $.ajax({
        url: 'https://mcweb-api.hzchu.top/admin/unpower',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            token: getCookie('token')
        }),
        success: function(response) {
            var code = response.code;
            console.log(code);

            // 处理返回的code和msg
            if (code === 0) {
                document.getElementById("dialog-title").innerText = "返回结果："
                document.getElementById("dialog-content").innerText = response.msg
                inst.open();
            } else {
                document.getElementById("dialog-title").innerText = "处理失败！返回结果："
                document.getElementById("dialog-content").innerText = response.msg
                inst.open();

            }

        }
    })
}
document.getElementById("account-enter").onclick = function() {
    var inst = new mdui.Dialog('#dialog');
    var mode = document.getElementById("account_choice").value;
    $.ajax({
        url: 'https://mcweb-api.hzchu.top/admin/getaccount',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            accountinfo: document.getElementById("account").value,
            token: getCookie('token'),
            mode: parseInt(mode)
        }),
        success: function(response) {
            var code = response.code;
            // 处理返回的code和msg
            if (code === 0) {
                if (response.group == 0) {
                    var group = "管理组"

                } else if (response.group == 1) {
                    var group = "普通玩家"

                } else if (response.group == 2) {
                    var group = "非Q群玩家"

                } else {
                    var group = "未知"

                }
                if (response.status == 1) {
                    var status = "正常"

                } else if (response.group == 2) {
                    var status = "已被封禁"

                } else {
                    var status = "未知"

                }
                document.getElementById("dialog-title").innerText = "查询结果："
                var content = `
                <div class="mdui-table-fluid">
                <table class="mdui-table">
                  <thead>
                  <tr>
                    <th>项目</th>
                    <th>内容</th>
                  </tr>
                  </thead>
                  <tbody>
                  <tr>
                    <td>游戏名（JAVA）</td>
                    <td>` + xss(response.user_name) + `</td>
                  </tr>
                  <tr>
                    <td>游戏名（BE）</td>
                    <td>` + xss(response.be_name) + `</td>
                  </tr>
                  <tr>
                    <td>所属用户组</td>
                    <td>` + group + `</td>
                </tr>         
                <tr>
                    <td>状态</td>
                    <td>` + status + `</td>
                </tr>
                <tr>
                    <td>QQ</td>
                    <td>` + response.qq + `</td>
                </tr>
                <tr>
                    <td>绑定时间</td>
                    <td>` + response.bind_time + `</td>
                </tr>
                <tr>
                <td>IP</td>
                <td>` + response.ip + `</td>
                </tr>
                <tr>
                <td>IP归属地：</td>
                <td>` + response.country + `，` + response.city + `</td>
                </tr>  
                  </tbody>
                </table>
              </div>`
                document.getElementById("dialog-content").innerHTML = content
                inst.open();
            } else {
                document.getElementById("dialog-title").innerText = "处理失败！返回结果："
                document.getElementById("dialog-content").innerText = response.msg
                inst.open();

            }

        }
    })
}


function xss(str, kwargs) {
    return ('' + str)

    .replace(/&/g, '&amp;')

    .replace(/</g, '&lt;') // DEC=> &#60; HEX=> &#x3c; Entity=> &lt;

    .replace(/>/g, '&gt;')

    .replace(/"/g, '&quot;')

    .replace(/'/g, '&#x27;') // &apos; 不推荐，因为它不在HTML规范中

    .replace(/\//g, '&#x2F;');

};




























document.getElementById("titletext").innerHTML = "峰间云海|管理";