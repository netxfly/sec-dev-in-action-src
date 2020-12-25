local core = require("apisix.core")
local exchange_owa = require("apisix.plugins.exchange-owa.exchange_owa")

local plugin_name = "exchange-owa"

local schema = {
    type = "object",
    properties = {
        mail_server = {
            type = "string",
        },
        debug = { type = "boolean",
                  enum = { true, false },
        },
    },
    required = { "mail_server" }
}

local _M = {
    version = 0.1,
    priority = 2001,
    name = plugin_name,
    schema = schema,
}

function _M.check_schema(conf)
    local ok, err = core.schema.check(schema, conf)

    if not ok then
        return false, err
    end

    return true
end

-- 在rewrite阶段，判断用户提交的动态口令是否正确
function _M.rewrite(conf, ctx)
    core.log.warn("plugin rewrite phase, conf: ", core.json.encode(conf))
    exchange_owa.auth_otp_token(conf.mail_server)
end

-- 用户访问登录页面时，在body_filter阶段修改返回的用户的html表单，增加动态口令输入框
function _M.body_filter(conf, ctx)
    exchange_owa.add_otp_token_form()
end

return _M
