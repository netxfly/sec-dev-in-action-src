---

local ngx = ngx
local string = string
local core = require("apisix.core")
local http = require("resty.http")

local util = require("apisix.plugins.exchange-ews.util")
local active_code = require("apisix.plugins.exchange-ews.active_code")
local device_manager = require("apisix.plugins.exchange-ews.device_manager")

local fetch_local_conf = require("apisix.core.config_local").local_conf
local config = fetch_local_conf()

local zero_trust_api = config.api.zero_trust_api
local opa_url = string.format("%s%s", config.api.opa_host, config.api.opa_path)

local function ews_check(conf, ctx)
    local user_agent = core.request.header(ctx, "user-agent")
    local remote_ip = util.get_user_srcip(ctx)
    local client_type = util.get_client_type(ctx)
    local username = util.get_username(ctx, client_type)
    username = util.get_username_from_mail(username)

    if conf.debug then
        core.log.warn(string.format("username: %s, user_agent: %s, ip: %s", username, user_agent, remote_ip))
    end

    if username == "" then
        return
    end

    -- EWS访问控制
    -- get_ews_address是否注册过（ews_username中hmget ip有值）
    local result = device_manager.get_ews_address(username, remote_ip)
    -- 判断某个用户的IP是否被禁用，返回为true时表示未禁用，返回false表示禁用了，不再允许连接与激活了
    local ews_address_status = device_manager.get_ews_address_status(username, remote_ip)
    -- 判断状态，是否允许连接（只判断是否允许连接，不管禁用与未激活的状态）
    local ews_status = device_manager.chk_ews_address(username, remote_ip)
    core.log.warn(string.format("result: %s, ews_address_status: %s, ews_status: %s", result, ews_address_status, ews_status))

    -- 如果注册过
    if result then
        -- 该用户的EWS白名单是否存在且状态为非忽略的才进入激活流程
        if ews_status then
            -- 如果允许连接就续期
            device_manager.update_ews_address(username, remote_ip, 1, 0, client_type)
        elseif ews_address_status then
            -- 不允许连接时且账户未禁用时，就进入激活流程
            util.check_crack(username)
            local iplist = device_manager.get_ews_iplist(username)
            active_code.do_ews_active(user_agent, client_type, username, remote_ip, iplist)
        else
            util.check_crack(username)
            -- 否则直接阻断链接
            ngx.exit(ngx.HTTP_CLOSE)
        end
    else
        local iplist = device_manager.get_ews_iplist(username)
        active_code.do_ews_active(user_agent, client_type, username, remote_ip, iplist)
    end
    -- 判断逻辑结束
end

local function zero_trust_check(ctx)
    local result = false

    local user_agent = core.request.header(ctx, "user-agent")
    local remote_ip = util.get_user_srcip(ctx)
    local client_type = util.get_client_type(ctx)
    local username = util.get_username(ctx, client_type)
    username = util.get_username_from_mail(username)

    if conf.debug then
        core.log.debug(string.format("username: %s, user_agent: %s, ip: %s", username, user_agent, remote_ip))
    end

    if username == "" then
        return
    end

    local data = string.format("username=%s", username)
    local headers = { ["Content-Type"] = "application/x-www-form-urlencoded", ["Content-Length"] = #data }
    local trust_info = {}

    local http_client = http.new()
    local res, err = http_client:request_uri(zero_trust_api, {
        method = "POST",
        body = data,
        headers = headers,
    })

    if err == nil and res ~= nil and res.status == 200 then
        trust_info = core.json.decode(res.body)
    end

    result = util.check_opa_policy(opa_url, trust_info)

    return result

end

local function ews(conf, ctx)
    local remote_ip = util.get_user_srcip(ctx)
    local client_type = util.get_client_type(ctx)
    local email = util.get_username(ctx, client_type)
    local username = util.get_username_from_mail(email)

    if #username > 0 then
        -- 判断是否为办公网内网
        local is_office_vlan = core.strings.starts(remote_ip, "10.") or stringy.starts(remote_ip, "172.16.")
        -- local is_office_wlan = office_ip.chk_officeips(remote_ip)
        is_office_vlan = false
        local is_office_wlan = false

        if core.strings.startswith(ngx.var.uri, "/EWS/") then
            core.log.warn(string.format("is_office_vlan: %s, is_office_wlan: %s",
                    is_office_vlan, is_office_wlan))
            if is_office_vlan or is_office_wlan then
                -- 如果是内网地址或公司出口IP，直接跳过验证逻辑
            else
                core.log.warn(string.format("username: %s, run_mode: %s, remote_ip: %s, client_type: %s",
                        username, conf.run_mode, remote_ip, client_type))

                if conf.run_mode == "normal" then
                    ews_check(conf, ctx)
                else
                    zero_trust_check(conf, ctx)
                end
            end
        end
    else
        return
    end
end


-- 当用户未激活时，替换掉返回给邮件客户端的Body，只允许连接，但不允许收发邮件
local function replace_body()
    local is_not_activated = ngx.ctx.is_not_activated
    -- core.log.warn(string.format("is_not_activated: %s", is_not_activated))
    if is_not_activated then
        if core.strings.startswith(ngx.var.uri, "/EWS/") --[[or core.strings.startswith(ngx.var.uri, "/rpc/")]] then
            local chunk, eof = ngx.arg[1], ngx.arg[2]
            local buffered = ngx.ctx.buffered
            if not buffered then
                buffered = {}
                ngx.ctx.buffered = buffered
            end

            if chunk ~= "" then
                buffered[#buffered + 1] = chunk
                ngx.arg[1] = nil
            end

            if eof then
                -- local whole = table.concat(buffered)
                ngx.ctx.buffered = nil
                ngx.arg[1] = nil
            end
        end
    end
end

-- 当用户未激活时，替换掉邮件客户端的Header，只允许连接，但不允许收发邮件
local function replace_header()
    local is_not_activated = ngx.ctx.is_not_activated or false
    -- core.log.warn(string.format("is_not_activated: %s", is_not_activated))
    if is_not_activated then
        if core.strings.startswith(ngx.var.uri, "/EWS/") --[[or core.strings.startswith(ngx.var.uri, "/rpc/")]] then
            ngx.header.content_length = nil
            ngx.header.content_encoding = nil
        end
    end
end

local _M = {
    ews = ews,
    replace_body = replace_body,
    replace_header = replace_header,
}

return _M
