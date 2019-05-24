<?php
/**
 * 验证码模块
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-01-29
 * Time: 11:56
 */
namespace app\Service;
use zhimiao\Request;
use zhimiao\Response;
use zhimiao\Utils;
use app\Service\Utils as APP_Utils;
use Respect\Validation\Validator as v;
class Verify {

    /**
     * 校验图片验证码正确性
     * @param null $token 生成验证码图片给的token
     * @param null $code 用户填写的验证码
     * @return bool|string
     */
    public function verifyImageCaptcha($token = null, $code = null)
    {
        if (!v::notEmpty()->validate($token)){
            return 'token不能为空';
        }
        if (!v::notEmpty()->validate($code)) {
            return 'code不能为空';
        }
        $token = Utils::decrypt($token);
        if (!v::json()->validate($token)) {
            return 'token无法解析';
        }
        $token = json_decode($token, true);
        $lock_id = 'verifyImageCaptcha:'. $token['key'];
        if (!APP_Utils::cacheNumLock($lock_id)) {
            return '操作过频';
        }
        if ($token['expire'] < time()) {
            return '验证码已经失效了';
        }
        if ($token['code'] == $code) {
            APP_Utils::cacheNumLock($lock_id, 0);
            return true;
        }
        return '验证失败';
    }

    /**
     * 校验手机号正确性
     * @param null $token 生成验证码图片给的token
     * @param null $code 用户填写的验证码
     * @return bool|string
     */
    public function verifySMSCode($phone = null, $token = null, $code = null)
    {
        if (!v::phone()->validate($phone)) {
            return [-5, null, '手机号无法识别'];
        }
        if (!v::notEmpty()->validate($token)){
            return 'token不能为空';
        }
        if (!v::notEmpty()->validate($code)) {
            return 'code不能为空';
        }
        $token = Utils::decrypt($token);
        if (!v::json()->validate($token)) {
            return 'token无法解析';
        }
        $token = json_decode($token, true);
        $lock_id = 'verifySMSCode:'. $token['key'];
        if (!APP_Utils::cacheNumLock($lock_id)) {
            return '操作过频';
        }
        if ($token['expire'] < time()) {
            return '验证码已经失效了';
        }
        if ($token['phone'] != $phone) {
            return '手机号不匹配';
        }
        if ($token['code'] == $code) {
            APP_Utils::cacheNumLock($lock_id, 0);
            return true;
        }
        return '验证失败';
    }

    /**
     * 判断用户登陆，并返回uid
     * @param bool $pass
     * @param bool $isExpire
     * @return int|string
     */
    public static function isLogin($pass = false, $isExpire = true)
    {
        $login_key = Request::header('LOGIN-KEY');
        if (v::notEmpty()->validate($login_key)) {
            $data = APP_Utils::parseSessionKey($login_key, $isExpire);
            if(!is_string($data)) {
                $data = (int) $data['uid'] ?? 0;
            }
        } else {
            $data = '头信息无法识别';
        }
        if (!$pass && is_string($data)) {
            Response::json(-1, null, $data);
        }
        return is_string($data) ? 0 : $data;
    }
}