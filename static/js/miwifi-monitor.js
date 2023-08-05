(function() {
    if (typeof window.MIWIFI_MONITOR !== "undefined") {
        return
    }
    var version = "v9.9.9 (2023.08.03)"
      , guidCookieDomains = [];
    var MIWIFI_MONITOR = (function(window, undefined) {
        var isLocal = false;
        var doc = document
          , nav = navigator
          , screen = window.screen
          , domain = isLocal ? "" : document.domain.toLowerCase()
          , ua = nav.userAgent.toLowerCase();
        var StringH = {
            trim: function(s) {
                return s.replace(/^[\s\xa0\u3000]+|[\u3000\xa0\s]+$/g, "")
            }
        };
        var NodeH = {
            on: function(el, type, fn) {
                if (el.addEventListener) {
                    el && el.addEventListener(type, fn, false)
                } else {
                    el && el.attachEvent("on" + type, fn)
                }
            },
            parentNode: function(el, tagName, deep) {
                deep = deep || 5;
                tagName = tagName.toUpperCase();
                while (el && deep-- > 0) {
                    if (el.tagName === tagName) {
                        return el
                    }
                    el = el.parentNode
                }
                return null
            }
        };
        var EventH = {
            fix: function(e) {
                if (!("target"in e)) {
                    var node = e.srcElement || e.target;
                    if (node && node.nodeType == 3) {
                        node = node.parentNode
                    }
                    e.target = node
                }
                return e
            }
        };
        var ObjectH = (function() {
            function getConstructorName(o) {
                if (o != null && o.constructor != null) {
                    return Object.prototype.toString.call(o).slice(8, -1)
                } else {
                    return ""
                }
            }
            return {
                isArray: function(obj) {
                    return getConstructorName(obj) == "Array"
                },
                isObject: function(obj) {
                    return obj !== null && typeof obj == "object"
                },
                mix: function(des, src, override) {
                    for (var i in src) {
                        if (override || !(des[i] || (i in des))) {
                            des[i] = src[i]
                        }
                    }
                    return des
                },
                encodeURIJson: function(obj) {
                    var result = [];
                    for (var p in obj) {
                        if (obj[p] == null) {
                            continue
                        }
                        result.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]))
                    }
                    return result.join("&")
                }
            }
        }
        )();
        var Cookie = {
            get: function(key) {
                try {
                    var a, reg = new RegExp("(^| )" + key + "=([^;]*)(;|$)");
                    if (a = doc.cookie.match(reg)) {
                        return unescape(a[2])
                    } else {
                        return ""
                    }
                } catch (e) {
                    return ""
                }
            },
            set: function(key, val, options) {
                options = options || {};
                var expires = options.expires;
                if (typeof (expires) === "number") {
                    expires = new Date();
                    expires.setTime(expires.getTime() + options.expires)
                }
                try {
                    doc.cookie = key + "=" + escape(val) + (expires ? ";expires=" + expires.toGMTString() : "") + (options.path ? ";path=" + options.path : "") + (options.domain ? "; domain=" + options.domain : "")
                } catch (e) {}
            }
        };
        var util = {
            getProject: function() {
                return ""
            },
            getReferrer: function() {
                var ref = doc.referrer || "";
                if (ref.indexOf("pass") > -1 || ref.indexOf("pwd") > -1) {
                    return "403"
                }
                return ref
            },
            getBrowser: function() {
                var browsers = {
                    "360se-ua": "360se",
                    TT: "tencenttraveler",
                    Maxthon: "maxthon",
                    GreenBrowser: "greenbrowser",
                    Sogou: "se 1.x / se 2.x",
                    TheWorld: "theworld"
                };
                for (var i in browsers) {
                    if (ua.indexOf(browsers[i]) > -1) {
                        return i
                    }
                }
                var result = ua.match(/(msie|chrome|safari|firefox|opera|trident)/);
                result = result ? result[0] : "";
                if (result == "msie") {
                    result = ua.match(/msie[^;]+/) + ""
                } else {
                    if (result == "trident") {
                        ua.replace(/trident\/[0-9].*rv[ :]([0-9.]+)/ig, function(a, c) {
                            result = "msie " + c
                        })
                    }
                }
                return result
            },
            getLocation: function() {
                var url = "";
                try {
                    url = location.href
                } catch (e) {
                    url = doc.createElement("a");
                    url.href = "";
                    url = url.href
                }
                url = url.replace(/[?#].*$/, "");
                url = /\.(s?htm|php)/.test(url) ? url : (url.replace(/\/$/, "") + "/");
                return url
            },
            getGuid: (function() {
                function hash(s) {
                    var h = 0
                      , g = 0
                      , i = s.length - 1;
                    for (i; i >= 0; i--) {
                        var code = parseInt(s.charCodeAt(i), 10);
                        h = ((h << 6) & 268435455) + code + (code << 14);
                        if ((g = h & 266338304) != 0) {
                            h = (h ^ (g >> 21))
                        }
                    }
                    return h
                }
                function guid() {
                    var s = [nav.appName, nav.version, nav.language || nav.browserLanguage, nav.platform, nav.userAgent, screen.width, "x", screen.height, screen.colorDepth, doc.referrer].join("")
                      , sLen = s.length
                      , hLen = window.history.length;
                    while (hLen) {
                        s += (hLen--) ^ (sLen++)
                    }
                    return (Math.round(Math.random() * 2147483647) ^ hash(s)) * 2147483647
                }
                var guidKey = "__guid"
                  , id = Cookie.get(guidKey);
                if (!id) {
                    id = [hash(isLocal ? "" : doc.domain), guid(), +new Date + Math.random() + Math.random()].join(".");
                    var config = {
                        expires: 24 * 3600 * 1000 * 300,
                        path: "/"
                    };
                    if (guidCookieDomains.length) {
                        for (var i = 0; i < guidCookieDomains.length; i++) {
                            var guidCookieDomain = guidCookieDomains[i]
                              , gDomain = "." + guidCookieDomain;
                            if ((domain.indexOf(gDomain) > 0 && domain.lastIndexOf(gDomain) == domain.length - gDomain.length) || domain == guidCookieDomain) {
                                config.domain = gDomain;
                                break
                            }
                        }
                    }
                    Cookie.set(guidKey, id, config)
                }
                return function() {
                    return id
                }
            }
            )(),
            getCount: (function() {
                var countKey = "monitor_count"
                  , count = Cookie.get(countKey);
                count = (parseInt(count) || 0) + 1;
                Cookie.set(countKey, count, {
                    expires: 24 * 3600 * 1000,
                    path: "/"
                });
                return function() {
                    return count
                }
            }
            )(),
            getFlashVer: function() {
                var ver = -1;
                if (nav.plugins && nav.mimeTypes.length) {
                    var plugin = nav.plugins["Shockwave Flash"];
                    if (plugin && plugin.description) {
                        ver = plugin.description.replace(/([a-zA-Z]|\s)+/, "").replace(/(\s)+r/, ".") + ".0"
                    }
                } else {
                    if (window.ActiveXObject && !window.opera) {
                        try {
                            var c = new ActiveXObject("ShockwaveFlash.ShockwaveFlash");
                            if (c) {
                                var version = c.GetVariable("$version");
                                ver = version.replace(/WIN/g, "").replace(/,/g, ".")
                            }
                        } catch (e) {}
                    }
                }
                ver = parseInt(ver, 10);
                return ver
            },
            getContainerId: function(el) {
                var areaStr, name, maxLength = 100;
                if (config.areaIds) {
                    areaStr = new RegExp("^(" + config.areaIds.join("|") + ")$","ig")
                }
                while (el) {
                    if (el.attributes && ("bk"in el.attributes || "data-bk"in el.attributes)) {
                        name = el.getAttribute("bk") || el.getAttribute("data-bk");
                        if (name) {
                            name = "bk:" + name;
                            return name.substr(0, maxLength)
                        }
                        if (el.id) {
                            name = el.getAttribute("data-desc") || el.id;
                            return name.substr(0, maxLength)
                        }
                    } else {
                        if (areaStr) {
                            if (el.id && areaStr.test(el.id)) {
                                name = el.getAttribute("data-desc") || el.id;
                                return name.substr(0, maxLength)
                            }
                        }
                    }
                    el = el.parentNode
                }
                return ""
            },
            getText: function(el) {
                var str = "";
                if (el.tagName.toLowerCase() == "input") {
                    str = el.getAttribute("text") || el.getAttribute("data-text") || el.value || el.title || ""
                } else {
                    str = el.getAttribute("text") || el.getAttribute("data-text") || el.innerText || el.textContent || el.title || ""
                }
                return StringH.trim(str).substr(0, 100)
            },
            getHref: function(el) {
                try {
                    return el.getAttribute("data-href") || el.href || ""
                } catch (e) {
                    return ""
                }
            },
            getDeviceId: function() {
                return "deviceId"
            },
            getAppVersion: function() {
                return "appVersion"
            },
            getRomVersion: function() {
                return "romVersion"
            },
            getHardwareVersion: function() {
                return "hardwareVersion"
            }
        };
        var data = {
            getBase: function() {
                return {
                    p: util.getProject(),
                    u: util.getLocation(),
                    id: util.getGuid(),
                    guid: util.getGuid(),
                    deviceId: util.getDeviceId(),
                    appVersion: util.getAppVersion(),
                    romVersion: util.getRomVersion(),
                    hardwareVersion: util.getHardwareVersion()
                }
            },
            getTrack: function() {
                return {
                    b: util.getBrowser(),
                    c: util.getCount(),
                    r: util.getReferrer(),
                    fl: util.getFlashVer()
                }
            },
            getClick: function(e) {
                e = EventH.fix(e || event);
                var target = e.target;
                if (target.attributes && ("data-log-element"in target.attributes)) {
                    return {
                        element: target.attributes["data-log-element"]
                    }
                }
                return false
            },
            getKeydown: function(e) {
                e = EventH.fix(e || event);
                if (e.keyCode != 13) {
                    return false
                }
                var target = e.target
                  , tagName = target.tagName
                  , containerId = util.getContainerId(target);
                if (tagName == "INPUT") {
                    var form = NodeH.parentNode(target, "FORM");
                    if (form) {
                        var formId = form.id || ""
                          , tId = target.id
                          , result = {
                            f: form.action,
                            c: "form:" + (form.name || formId),
                            cId: containerId
                        };
                        if (tId == "kw" || tId == "search-kw" || tId == "kw1") {
                            result.w = target.value
                        }
                        return result
                    }
                }
                return false
            }
        };
        var config = {
            trackUrl: null,
            clickUrl: null,
            areaIds: null
        };
        var $ = function(str) {
            return document.getElementById(str)
        };
        return {
            version: version,
            util: util,
            data: data,
            config: config,
            sendLog: (function() {
                window.__miwifi_monitor_imgs = {};
                return function(url) {
                    var id = "log_" + (+new Date)
                      , img = window.__miwifi_monitor_imgs[id] = new Image();
                    img.onload = img.onerror = function() {
                        if (window.__miwifi_monitor_imgs && window.__miwifi_monitor_imgs[id]) {
                            window.__miwifi_monitor_imgs[id] = null;
                            delete window.__miwifi_monitor_imgs[id]
                        }
                    }
                    ;
                    img.src = url
                }
            }
            )(),
            buildLog: (function() {
                var lastLogParams = "";
                return function(params, url) {
                    if (params === false) {
                        return
                    }
                    params = params || {};
                    var baseParams = data.getBase();
                    params = ObjectH.mix(baseParams, params, true);
                    var logParams = url + ObjectH.encodeURIJson(params);
                    if (logParams == lastLogParams) {
                        return
                    }
                    lastLogParams = logParams;
                    setTimeout(function() {
                        lastLogParams = ""
                    }, 100);
                    var sendParams = ObjectH.encodeURIJson(params);
                    sendParams += "&t=" + (+new Date);
                    url = url.indexOf("?") > -1 ? url + "&" + sendParams : url + "?" + sendParams;
                    this.sendLog(url)
                }
            }
            )(),
            log: function(params, type) {
                type = type || "click";
                var url = config[type + "Url"];
                if (!url) {
                    alert("Error : the " + type + "url does not exist!")
                }
                this.buildLog(params, url)
            },
            setConf: function(key, val) {
                var newConfig = {};
                if (!ObjectH.isObject(key)) {
                    newConfig[key] = val
                } else {
                    newConfig = key
                }
                this.config = ObjectH.mix(this.config, newConfig, true);
                return this
            },
            setUrl: function(url) {
                if (url) {
                    this.util.getLocation = function() {
                        return url
                    }
                }
                return this
            },
            setProject: function(prj) {
                if (prj) {
                    this.util.getProject = function() {
                        return prj
                    }
                }
                return this
            },
            setId: function() {
                var areaIds = [], i = 0, argument;
                while (argument = arguments[i++]) {
                    if (!ObjectH.isArray(argument)) {
                        areaIds.push(argument)
                    } else {
                        areaIds = areaIds.concat(argument)
                    }
                }
                this.setConf("areaIds", areaIds);
                return this
            },
            getTrack: function() {
                var params = this.data.getTrack();
                this.log(params, "track");
                return this
            },
            getClickAndKeydown: function() {
                var that = this;
                NodeH.on(doc, "mousedown", function(e) {
                    var params = that.data.getClick(e);
                    that.log(params, "click")
                });
                MIWIFI_MONITOR.getClickAndKeydown = function() {
                    return that
                }
                ;
                return this
            }
        }
    }
    )(window);
    window.MIWIFI_MONITOR = MIWIFI_MONITOR;
    if (typeof window.monitor === "undefined") {
        window.monitor = MIWIFI_MONITOR
    }
}
)();
