const scoreboardUrl = 'https://mcweb-api.hzchu.top/user/scoreboard';
const scoreboardBodyPickaxe = document.getElementById('scoreboard-body-pickaxe');
const scoreboardBodyKills = document.getElementById('scoreboard-body-kills');
const scoreboardBodyFish = document.getElementById('scoreboard-body-fish');
const scoreboardBodyGametime = document.getElementById('scoreboard-body-gametime');
const scoreboardBodyxp = document.getElementById('scoreboard-body-xp');
fetch(scoreboardUrl)
    .then(response => response.json())
    .then(scoreboard => {
        //定义一个比较函数，按照pickaxe_total的大小从高到低排序
        function compareByPickaxe(a, b) {
            return b.pickaxe_total - a.pickaxe_total;
        }

        //定义一个比较函数，按照Kills的大小从高到低排序
        function compareByKills(a, b) {
            return b.Kills - a.Kills;
        }

        //定义一个比较函数，按照fish的大小从高到低排序
        function compareByFish(a, b) {
            return b.fish - a.fish;
        }

        //定义一个比较函数，按照gametime的大小从高到低排序
        function compareByGametime(a, b) {
            return b.gametime - a.gametime;
        }

        //定义一个比较函数，按照xp的大小从高到低排序
        function compareByxp(a, b) {
            return b.xp - a.xp;
        }
        //定义一个过滤函数，排除掉pickaxe_total为0的玩家
        function filterByPickaxe(a) {
            return a.pickaxe_total > 0;
        }

        //定义一个过滤函数，排除掉Kills为null的玩家
        function filterByKills(a) {
            return a.Kills != null;
        }

        //定义一个过滤函数，排除掉fish为null的玩家
        function filterByFish(a) {
            return a.fish != null;
        }

        //定义一个过滤函数，排除掉gametime为null的玩家
        function filterByGametime(a) {
            return a.gametime != null;
        }
        //定义一个过滤函数，排除掉xp为null的玩家 
        function filterByxp(a) {
            return a.xp != null;
        }
        //对scoreboard数组进行过滤
        scoreboard = scoreboard.filter(filterByPickaxe);
        // 参照最初版重写！！！！
        //创建一个副本数组，按照pickaxe_total排序
        let scoreboardByPickaxe = [...scoreboard];
        scoreboardByPickaxe.sort(compareByPickaxe);
        console.log(scoreboardByPickaxe)

        //创建一个副本数组，按照Kills排序，并过滤掉null值
        let scoreboardByKills = [...scoreboard];
        scoreboardByKills.sort(compareByKills);
        scoreboardByKills = scoreboardByKills.filter(filterByKills);

        //创建一个副本数组，按照fish排序，并过滤掉null值
        let scoreboardByFish = [...scoreboard];
        scoreboardByFish.sort(compareByFish);
        scoreboardByFish = scoreboardByFish.filter(filterByFish);


        //创建一个副本数组，按照gametime排序，并过滤掉null值
        let scoreboardByGametime = [...scoreboard];
        scoreboardByGametime.sort(compareByGametime);
        scoreboardByGametime = scoreboardByGametime.filter(filterByGametime);
        //创建一个副本数组，按照xp排序，并过滤掉null值 
        let scoreboardByXp = [...scoreboard];
        scoreboardByXp.sort(compareByxp);
        scoreboardByXp = scoreboardByXp.filter(filterByxp);
        //遍历按照pickaxe_total排序的数组，创建表格行和单元格，并添加到对应的表格中
        scoreboardByPickaxe.forEach((player, index) => {
            const row = document.createElement('tr');
            const rank = document.createElement('td');
            const playerName = document.createElement('td');
            const pickaxeTotal = document.createElement('td');

            rank.textContent = index + 1;
            playerName.textContent = player.playername;
            pickaxeTotal.textContent = player.pickaxe_total;

            row.appendChild(rank);
            row.appendChild(playerName);
            row.appendChild(pickaxeTotal);

            scoreboardBodyPickaxe.appendChild(row);
        });

        //遍历按照Kills排序的数组，创建表格行和单元格，并添加到对应的表格中
        scoreboardByKills.forEach((player, index) => {
            const row = document.createElement('tr');
            const rank = document.createElement('td');
            const playerName = document.createElement('td');
            const kills = document.createElement('td');

            rank.textContent = index + 1;
            playerName.textContent = player.playername;
            kills.textContent = player.Kills;

            row.appendChild(rank);
            row.appendChild(playerName);
            row.appendChild(kills);

            scoreboardBodyKills.appendChild(row);
        });

        //遍历按照fish排序的数组，创建表格行和单元格，并添加到对应的表格中
        scoreboardByFish.forEach((player, index) => {
            const row = document.createElement('tr');
            const rank = document.createElement('td');
            const playerName = document.createElement('td');
            const fish = document.createElement('td');

            rank.textContent = index + 1;
            playerName.textContent = player.playername;
            fish.textContent = player.fish;

            row.appendChild(rank);
            row.appendChild(playerName);
            row.appendChild(fish);

            scoreboardBodyFish.appendChild(row);
        });

        //遍历按照gametime排序的数组，创建表格行和单元格，并添加到对应的表格中
        scoreboardByGametime.forEach((player, index) => {
            const row = document.createElement('tr');
            const rank = document.createElement('td');
            const playerName = document.createElement('td');
            const gametime = document.createElement('td');

            rank.textContent = index + 1;
            playerName.textContent = player.playername;
            gametime.textContent = Math.floor(player.gametime / 60 / 60) + "h";


            row.appendChild(rank);
            row.appendChild(playerName);
            row.appendChild(gametime);

            scoreboardBodyGametime.appendChild(row);
        });
        //遍历按照xp排序的数组，创建表格行和单元格，并添加到对应的表格中 
        scoreboardByXp.forEach((player, index) => {
            const row = document.createElement('tr');
            const rank = document.createElement('td');
            const playerName = document.createElement('td');
            const xp = document.createElement('td');

            rank.textContent = index + 1;
            playerName.textContent = player.playername;
            xp.textContent = getLevel(player.xp);

            row.appendChild(rank);
            row.appendChild(playerName);
            row.appendChild(xp);

            scoreboardBodyxp.appendChild(row);
        });
    });
// 定义一个函数，参数为原始经验值
function getLevel(exp) {
    // 定义一个变量，存储经验等级
    let level = 0;
    // 定义一个变量，存储升级所需的经验值
    let expToNext = 0;
    // 使用循环，不断更新经验等级和升级所需的经验值，直到原始经验值不足以升级为止
    while (true) {
        // 根据公式，计算升级所需的经验值
        if (level <= 15) {
            expToNext = 2 * level + 7;
        } else if (level <= 30) {
            expToNext = 5 * level - 38;
        } else {
            expToNext = 9 * level - 158;
        }
        // 判断原始经验值是否足够升级
        if (exp >= expToNext) {
            // 如果足够，减去升级所需的经验值，增加经验等级
            exp -= expToNext;
            level++;
        } else {
            // 如果不足，跳出循环
            break;
        }
    }
    // 返回经验等级
    return level;
}
document.getElementById("titletext").innerHTML = "峰间云海|排行榜";