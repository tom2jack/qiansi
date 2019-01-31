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
            return [1, ['sessionKey' => $sessionKey]];
        }
        return 0;
    }

}