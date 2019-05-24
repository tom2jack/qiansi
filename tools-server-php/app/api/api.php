<?php
/**
 * 客户端对接API
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019/5/19
 * Time: 下午1:49
 */

namespace app\api;
use app\Service\Utils;
use zhimiao\Data;
use zhimiao\Request;
use zhimiao\Response;

class api {
    private $db;

    /**
     * api constructor.
     */
    public function __construct()
    {
        $this->db = Data::pdo();
    }

    /**
     * 注册服务器
     * @param null $uid 用户ID
     * @param null $device 客户端设备号
     */
    public function regServer($uid = null, $device = null)
    {
        $uid > 0 || Response::json(-4, null, '用户UID非法');
        strlen($device) == 36 || Response::json(-4, null, '客户端唯一标识号非法');
        $uid = $this->db->quickPrepare('select uid from member where uid=:uid', [':uid' => $uid])->getOnce('uid');
        if ($uid === false) {
            return [-5, null, '用户不存在'];
        }
        $device_check = $this->db->quickPrepare('select exists(select 1 from `server` where device_id=:device_id) as ret', [':device_id' => $device])->getOnce('ret');
        if ($device_check == 1) {
            return [-5, null, '设备已存在，请勿重复注册'];
        }
        $api_secret = Utils::getRandomStr(mt_rand(50, 150));
        $statement = $this->db->quickPrepare(
            'insert into `server`(`uid`,`api_secret`,`device_id`,`domain`) values (:uid,:api_secret,:device_id,:domain)',
            [
                ':uid' => $uid,
                ':api_secret' => $api_secret,
                ':device_id' => $device,
                ':domain' => Request::getIp()
            ]
        );
        $result = $statement->rowCount() == 1 ? $this->db->lastInsertId('server') : 0;
        $statement->closeCursor();
        if ($result > 0) {
            return [1, [
                'server_id' => $result,
                'api_secret' => $api_secret
            ], '注册成功'];
        }
        return 0;
    }
}