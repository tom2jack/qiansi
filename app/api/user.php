<?php
/**
 * 用户信息处理接口模块
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2018-11-27
 * Time: 09:54
 */

namespace app\api;

use app\Service\Utils;
use Respect\Validation\Validator as v;
use zhimiao\Data;
use zhimiao\Request;

class user
{
    public $db;

    public function __construct()
    {
        $this->db = Data::pdo();
    }

    /**
     * 登陆
     * @param string $phone
     * @param string $password
     * @return array
     */
    public function signin($phone = null, $password = null)
    {
        if (!v::phone()->validate($phone)) {
            return [-5, null, '手机号无法识别'];
        }
        if (!v::length(6, 32)->validate($password)) {
            return [-4, null, '密码长度不合法(6位-32位)'];
        }
        $lock_key = [
            'signin:phone:'. $phone,
            'signin:ip:'. md5(Request::getIp())
        ];
        if (!Utils::cacheNumLock($lock_key[0]) || !Utils::cacheNumLock($lock_key[1], 15)) {
            return [-3, null, '尝试次数过频，请稍后重试'];
        }
        $user_info = $this->db->quickPrepare('select `uid`,`phone`,`password`,`status` from member where phone=:phone', [':phone' => $phone])->getOnce();
        if ($user_info === false) {
            return [-5, null, '用户不存在'];
        }
        if ($user_info['status'] != 1) {
            return [-5, null, '用户状态异常'];
        }
        if (!password_verify($password, $user_info['password'])) {
            return [-5, null, '密码输入错误'];
        }
        unset($user_info['password']);
        return [1, [
            'user_info' => $user_info,
            'sessionKey' => Utils::createSessionKey($user_info['uid'])
        ]];
    }

    /**
     * 注册账号
     * @param string $phone
     * @param string $password
     * @param string $code
     * @param string $token
     * @return array|int
     */
    public function signup($phone = null, $password = null, $code = null, $token = null)
    {
        $VerifyService = new \app\Service\Verify();
        $check_result = $VerifyService->verifySMSCode($phone, $token, $code);
        if (is_string($check_result)) {
            return [-5, null, $check_result];
        }
        if (!v::length(6, 32)->validate($password)) {
            return [-4, null, '密码长度不合法(6位-32位)'];
        }
        $statement = $this->db->quickPrepare('INSERT INTO `member`(`phone`, `password`) VALUES (:phone, :password)', [
            ':phone' => $phone,
            ':password' => password_hash($password, PASSWORD_DEFAULT)
        ]);
        $ret = $statement->rowCount();
        $statement->closeCursor();
        if (!$ret) {
            return [-6, null, '注册失败，系统繁忙或者此手机号已经注册了'];
        }
        $uid = $this->db->lastInsertId();
        if ($uid > 0) {
            $sessionKey = Utils::createSessionKey($uid);
            return [1, [
                'uid' => $uid,
                'sessionKey' => $sessionKey
            ]];
        }
        return 0;
    }

    /**
     * 修改密码
     * @param string $old_pwd 旧密码
     * @param string $new_pwd 新密码
     */
    public function resetPwd($old_pwd = null, $new_pwd = null)
    {
        $uid = \app\Service\Verify::isLogin();
        if (!v::length(6, 32)->validate($old_pwd) || !v::length(6, 32)->validate($new_pwd)) {
            return [-4, null, '密码长度不合法(6位-32位)'];
        }
        $old_pwd_hash = $this->db->quickPrepare('select `password` from member where uid=:uid', [':uid' => $uid])->getOnce('password');
        if ($old_pwd_hash === false) {
            return [-5, null, '用户不存在'];
        }
        if (!password_verify($old_pwd, $old_pwd_hash)) {
            return [-5, null, '旧密码输入错误'];
        }
        $statement = $this->db->quickPrepare('update member set `password`=:password where uid=:uid', [
                ':uid' => $uid,
                ':password' => password_hash($new_pwd, PASSWORD_DEFAULT)
            ]);
        $ret = $statement->rowCount();
        $statement->closeCursor();
        if (!$ret) {
            return [-6, null, '密码更新失败，请重新尝试'];
        }
        return 1;
    }
}