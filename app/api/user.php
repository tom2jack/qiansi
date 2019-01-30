<?php
/**
 *
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2018-11-27
 * Time: 09:54
 */

namespace app\api;

use Respect\Validation\Validator as v;
use zhimiao\Data;

class user
{
    public $db;

    public function __construct()
    {
        $this->db = Data::pdo();
    }

    public function signin($phone = null, $password = null)
    {

    }

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


    }

}