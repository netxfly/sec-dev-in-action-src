var txt = {
  cn: {
    "doc_title": "邮箱激活系统",
    "dear_user": "亲爱的<span class='email'></span>:",
    "new_device": "您的邮箱 <span class='email'></span>正在一台新的设备（<span class='device'></span>）上登录，详情如下",
    "device_name": "设备名称",
    "device_model": "设备型号",
    "device_type": "设备类型",
    "device_id": "设备ID",
    "imei_id": "IMEI",
    "phone_number":"手机号码",
    "ip_address": "新设备IP地址",
    "change_pwd": "如果这不是您本人操作，您的邮箱密码可能已经泄露，请修改密码并拒绝该请求。",
    "help_phone": "如需协助请联系IT部门",
    "outlook_tips": "您使用的是outlook客户端，用户名和密码会被保存到第三方微软公司的服务器中，存在安全风险，推荐使用手机自带的邮件客户端",
    "allow": " 允许",
    "reject": " 拒绝",
    "allow_confirm_tips": "确定后，将同步您的邮箱到该设备，如需禁用该设备，请登录内网帐号中心的“手机邮箱”中进行操作",
    "cancel": " 取消",
    "ok": "确定",
    "reject_manage_account": "拒绝后，将无法同步您的邮箱到该设备。如需重新激活该设备，请登录内网帐号中心的“手机邮箱”中进行操作",
    "device_actived": "设备已授权",
    "device_rejected": "设备已拒绝",
    "device_exceed": "设备数已超出限制",
    "page_not_found": "页面未找到",
    "link_invalid": "激活链接无效或已过期",
    "caution": "温馨提示",
    "auth_to_mis_manage": "已授权设备请登录到内网账号中心的“手机邮箱”中进行管理",
    "reject_to_mis_manage": "已拒绝设备请登录到内网账号中心的“手机邮箱”中进行管理",
    "to_mis_change_pwd": "如果这不是您本人操作，您的邮箱密码可能已泄露，并且有人正在尝试登陆您的邮箱，请立即点拒绝访问，并修改密码。",
    "devices_exceeds_the_limit": "已经激活的设备数超出了10台的限制，请中删除不再使用的设备",
    "more_sec_info": "",
    "new_ip": "您的邮箱 <span class='email'></span>正在使用（<span class='device'></span>）客户端进行同步邮件，详情如下",
    "client_type": "客户端类型",
    "device_ip": "移动设备IP地址",
    "client_ip": "客户端IP地址",
    "expire_tips": "该客户端 IP 有效期为 8 个小时，过期后需要重新进行授权",
    "allow_confirm_tips_ip": "确定后，将同步您的邮箱到该客户端，该客户端 IP 有效期为 8 个小时，如需禁用该 IP，请登录到“邮箱管理中心”中对 IP 进行管理",
    "reject_manage_account_ip": "拒绝后，将无法同步您的邮箱到该客户端。如需重新激活该 IP，请登录到“邮箱管理中心”中对 IP 进行管理",
    "ip_actived": "客户端 IP 已激活",
    "ip_rejected": "客户端 IP 已拒绝",
    "auth_to_mis_manage_ip": "已授权 IP 请登录到“邮箱管理中心”中对 IP 进行管理",
    "reject_to_mis_manage_ip": "已拒绝 IP 请登录到“邮箱管理中心”中对 IP 进行管理",
    "to_mis_change_pwd_ip": "如果该客户端 IP 不是您的 IP ，您的密码可能已经泄露，请立即修改密码"
  }
};
$(function() {
  var $i18n = $.i18n();
  $i18n.load(txt).done(function() {
    $('html').i18n();
    $("[data-i18n^='[html]']").html(function() {
      return $.i18n($(this).attr('data-i18n').replace('\[html\]', ''));
    });
  })
})