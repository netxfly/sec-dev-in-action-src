--[[

Copyright (c) 2016 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

]]

local rulematch = ngx.re.find
local unescape = ngx.unescape_uri

local config = require("config")
local util = require("util")

local _M = {
    RULES = {}
}

function _M.load_rules()
    _M.RULES = util.get_rules(config.config_rule_dir)
    for k, v in pairs(_M.RULES)
    do
        ngx.log(ngx.INFO, string.format("%s Rule Set", k))
        for kk, vv in pairs(v)
        do
            ngx.log(ngx.INFO, string.format("index:%s, Rule:%s", kk, vv))
        end
    end
    return _M.RULES
end

function _M.get_rule(rule_file_name)
    ngx.log(ngx.DEBUG, rule_file_name)
    return _M.RULES[rule_file_name]
end

-- white ip check
function _M.white_ip_check()
    if config.config_white_ip_check == "on" then
        local IP_WHITE_RULE = _M.get_rule('whiteip.rule')
        local WHITE_IP = util.get_client_ip()
        if IP_WHITE_RULE ~= nil then
            for _, rule in pairs(IP_WHITE_RULE) do
                if rule ~= "" and rulematch(WHITE_IP, rule, "jo") then
                    _M.log_record('White_IP', ngx.var_request_uri, "_", "_")
                    return true
                end
            end
        end
    end
end

-- Bad guys check
function _M.bad_guy_check()
    local client_ip = util.get_client_ip()
    local ret = false
    if client_ip ~= "" then
        ret = ngx.shared.badGuys:get(client_ip)
        if ret ~= nil and ret > 0 then
            ret = true
        end
    end
    return ret
end


-- deny black ip
function _M.black_ip_check()
    if config.config_black_ip_check == "on" then
        local IP_BLACK_RULE = _M.get_rule('blackip.rule')
        local BLACK_IP = util.get_client_ip()
        if IP_BLACK_RULE ~= nil then
            for _, rule in pairs(IP_BLACK_RULE) do
                if rule ~= "" and rulematch(BLACK_IP, rule, "jo") then
                    util.log_record('BlackList_IP', ngx.var_request_uri, "_", "_")
                    if config.config_waf_enable == "on" then
                        ngx.exit(403)
                        return true
                    end
                end
            end
        end
    end
end

-- allow white url
function _M.white_url_check()
    if config.config_white_url_check == "on" then
        local URL_WHITE_RULES = _M.get_rule('writeurl.rule')
        local REQ_URI = ngx.var.request_uri
        if URL_WHITE_RULES ~= nil then
            for _, rule in pairs(URL_WHITE_RULES) do
                if rule ~= "" and rulematch(REQ_URI, rule, "jo") then
                    return true
                end
            end
        end
    end
end

-- deny cc attack
function _M.cc_attack_check()
    if config.config_cc_check == "on" then
        local ATTACK_URI = ngx.var.uri
        local CC_TOKEN = util.get_client_ip() .. ATTACK_URI
        local limit = ngx.shared.limit
        local CCcount = tonumber(string.match(config.config_cc_rate, '(.*)/'))
        local CCseconds = tonumber(string.match(config.config_cc_rate, '/(.*)'))
        local req, _ = limit:get(CC_TOKEN)
        if req then
            if req > CCcount then
                util.log_record('CC_Attack', ngx.var.request_uri, "-", "-")
                if config.config_waf_enable == "on" then
                    ngx.log(ngx.DEBUG, string.format("CC URL:%s", CC_TOKEN))
                    ngx.header.content_type = "text/html"
                    ngx.status = ngx.HTTP_FORBIDDEN
                    ngx.say(string.format("您的IP是：%s，您的访问频率太快了", util.get_client_ip()))
                    ngx.exit(ngx.status)
                end
            else
                limit:incr(CC_TOKEN, 1)
            end
        else
            limit:set(CC_TOKEN, 1, CCseconds)
        end
    end
    return false
end

-- deny cookie
function _M.cookie_attack_check()
    if config.config_cookie_check == "on" then
        local COOKIE_RULES = _M.get_rule('cookie.rule')
        local USER_COOKIE = ngx.var.http_cookie
        if USER_COOKIE ~= nil then
            for _, rule in pairs(COOKIE_RULES) do
                if rule ~= "" and rulematch(USER_COOKIE, rule, "jo") then
                    util.log_record('Deny_Cookie', ngx.var.request_uri, "-", rule)
                    if config.config_waf_enable == "on" then
                        util.waf_output()
                        return true
                    end
                end
            end
        end
    end
    return false
end

-- deny url
function _M.url_attack_check()
    if config.config_url_check == "on" then
        local URL_RULES = _M.get_rule('url.rule')
        local REQ_URI = ngx.var.request_uri
        for _, rule in pairs(URL_RULES) do
            if rule ~= "" and rulematch(REQ_URI, rule, "jo") then
                util.log_record('Deny_URL', REQ_URI, "-", rule)
                if config.config_waf_enable == "on" then
                    util.waf_output()
                    return true
                end
            end
        end
    end
    return false
end

-- deny url args
function _M.url_args_attack_check()
    if config.config_url_args_check == "on" then
        local ARGS_RULES = _M.get_rule('args.rule')
        for _, rule in pairs(ARGS_RULES) do
            local REQ_ARGS = ngx.req.get_uri_args()
            for key, val in pairs(REQ_ARGS) do
                ngx.log(ngx.DEBUG, key)
                local ARGS_DATA = {}
                if type(val) == 'table' then
                    ARGS_DATA = table.concat(val, " ")
                    ngx.log(ngx.DEBUG, ARGS_DATA)
                else
                    ARGS_DATA = val
                end
                if ARGS_DATA and type(ARGS_DATA) ~= "boolean" and rule ~= "" and rulematch(unescape(ARGS_DATA), rule, "jo") then
                    util.log_record('Deny_URL_Args', ngx.var.request_uri, "-", rule)
                    if config.config_waf_enable == "on" then
                        util.waf_output()
                        return true
                    end
                end
            end
        end
    end
    return false
end

-- deny user agent
function _M.user_agent_attack_check()
    if config.config_user_agent_check == "on" then
        local USER_AGENT_RULES = _M.get_rule('useragent.rule')
        local USER_AGENT = ngx.var.http_user_agent
        if USER_AGENT ~= nil then
            for _, rule in pairs(USER_AGENT_RULES) do
                if rule ~= "" and rulematch(USER_AGENT, rule, "jo") then
                    util.log_record('Deny_USER_AGENT', ngx.var.request_uri, "-", rule)
                    if config.config_waf_enable == "on" then
                        util.waf_output()
                        return true
                    end
                end
            end
        end
    end
    return false
end

-- deny user agent
function _M.user_header_attack_check()
    if config.config_user_header_check == "on" then
        local USER_HEADER_RULES = _M.get_rule('headers.rule')
        local USER_HEADERS = util.get_headers() or {}
        if USER_HEADERS ~= nil and USER_HEADERS["Content-Type"] ~= nil then
            local content_type = USER_HEADERS["Content-Type"]
            local content_value = ""
            if content_value ~= nil then
                if type(content_type) == 'table' then
                    content_value = table.concat(content_type, " ")
                else
                    content_value = content_type
                end
            end

            for _, rule in pairs(USER_HEADER_RULES) do
                local result = rulematch(content_value, rule, "jo")
                -- ngx.log(ngx.DEBUG, string.format("content_value:%s, content_type:%s, rule:%s, result:%s", content_value, content_type, rule, result))
                if rule ~= "" and result then
                    ngx.log(ngx.DEBUG, string.format("Stuct2"))
                    util.log_record('Deny_USER_HEADERS', ngx.var.request_uri, "-", rule)
                    if config.config_waf_enable == "on" then
                        util.waf_output()
                        return false
                    end
                end
            end
        end
    end
    return true
end

-- deny post
function _M.post_attack_check()
    if config.config_post_check == "on" then
        ngx.req.read_body()
        local POST_RULES = _M.get_rule('post.rule')
        for _, rule in pairs(POST_RULES) do
            local POST_ARGS = ngx.req.get_post_args() or {}
            for _, v in pairs(POST_ARGS) do
                local post_data = ""
                if type(v) == "table" then
                    post_data = table.concat(v, ", ")
                else
                    post_data = v
                end
                -- ngx.log(ngx.INFO, post_data, type(post_data))
                if rule ~= "" and type(post_data) ~= "boolean" and rulematch(post_data, rule, "jo") then
                    util.log_record('Deny_USER_POST_DATA', post_data, "-", rule)
                    if config.config_waf_enable == "on" then
                        util.waf_output()
                        return true
                    end
                end
            end
        end
    end
    return false
end

-- start change to jinghuashuiyue mode, set in vhosts's location segument
function _M.start_jingshuishuiyue()
    local host = util.get_server_host()
    -- ngx.var.target = string.format("proxy_%s", host)
    ngx.var.target = "proxy_"
    if host and _M.bad_guy_check() then
        -- ngx.var.target = string.format("unreal_%s", host)
        ngx.var.target = "unreal_"
    end
end

-- waf start
function _M.check()
    if _M.white_ip_check() then
    elseif _M.black_ip_check() then
    elseif _M.user_agent_attack_check() then
    elseif _M.user_header_attack_check() then
    elseif _M.white_url_check() then
    elseif _M.url_attack_check() then
    elseif _M.cc_attack_check() then
    elseif _M.cookie_attack_check() then
    elseif _M.url_args_attack_check() then
    elseif _M.post_attack_check() then
    else
        return
    end
end

return _M